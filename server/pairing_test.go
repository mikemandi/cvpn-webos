package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func pairReq(a *app, t *testing.T, method, endpoint, payload, ip string, cookies []*http.Cookie, csrf string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, endpoint, strings.NewReader(payload))
	req.RemoteAddr = ip + ":32142"
	req.Host = "192.168.0.2:8765"
	for _, c := range cookies {
		req.AddCookie(c)
	}
	if csrf != "" {
		req.Header.Set("X-CVPN-CSRF", csrf)
	}
	rec := httptest.NewRecorder()
	if endpoint == "/tv" {
		a.tvPairingPage(rec, req)
	} else {
		a.api(rec, req)
	}
	return rec
}
func readJSON(t *testing.T, r *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var data map[string]any
	if e := json.Unmarshal(r.Body.Bytes(), &data); e != nil {
		t.Fatal(e)
	}
	return data
}
func TestPairingEndToEnd(t *testing.T) {
	oldBase := base
	base = t.TempDir()
	defer func() { base = oldBase }()
	pairLock.Lock()
	pairState = tvPairState{}
	pairLock.Unlock()
	a := &app{key: strings.Repeat("a", 64)}
	denied := pairReq(a, t, "GET", "/tv", "", "192.168.0.12", nil, "")
	if denied.Code != 403 {
		t.Fatalf("remote TV page must be denied: %d", denied.Code)
	}
	tv := pairReq(a, t, "GET", "/tv", "", "127.0.0.1", nil, "")
	if tv.Code != 200 {
		t.Fatalf("TV page: %d", tv.Code)
	}
	pairLock.Lock()
	code := pairState.Code
	pairLock.Unlock()
	if len(code) != 6 || !strings.Contains(tv.Body.String(), code[:3]+" "+code[3:]) {
		t.Fatal("valid code not shown on TV")
	}
	wrong := pairReq(a, t, "POST", "/api/pairing/request", `{"code":"999999","name":"Phone"}`, "192.168.0.23", nil, "")
	if wrong.Code != 401 && wrong.Code != 429 {
		t.Fatalf("wrong code accepted: %d", wrong.Code)
	}
	request := pairReq(a, t, "POST", "/api/pairing/request", `{"code":"`+code+`","name":"Samsung"}`, "192.168.0.23", nil, "")
	if request.Code != 200 {
		t.Fatalf("request failed: %d %s", request.Code, request.Body)
	}
	r := readJSON(t, request)
	if r["state"] != "approved" {
		t.Fatal("valid PIN did not grant session")
	}
	approved := request
	replay := pairReq(a, t, "POST", "/api/pairing/request", `{"code":"`+code+`","name":"Attacker"}`, "192.168.0.24", nil, "")
	if replay.Code == 200 {
		t.Fatal("single-use PIN was accepted twice")
	}
	cookies := approved.Result().Cookies()
	if len(cookies) != 1 || !cookies[0].HttpOnly || cookies[0].SameSite != http.SameSiteStrictMode {
		t.Fatal("session cookie flags missing")
	}
	me := pairReq(a, t, "GET", "/api/auth", "", "192.168.0.23", cookies, "")
	if me.Code != 200 {
		t.Fatalf("session not authorized: %d", me.Code)
	}
	csrf := readJSON(t, me)["csrf"].(string)
	noCSRF := pairReq(a, t, "POST", "/api/auth/logout", `{}`, "192.168.0.23", cookies, "")
	if noCSRF.Code != 403 {
		t.Fatal("CSRF protection not enforced")
	}
	devices := pairReq(a, t, "GET", "/api/devices", "", "192.168.0.23", cookies, "")
	if devices.Code != 200 || len(readJSON(t, devices)["items"].([]any)) != 1 {
		t.Fatal("paired device not persisted")
	}
	logout := pairReq(a, t, "POST", "/api/auth/logout", `{}`, "192.168.0.23", cookies, csrf)
	if logout.Code != 200 {
		t.Fatal("logout failed")
	}
	after := pairReq(a, t, "GET", "/api/auth", "", "192.168.0.23", cookies, "")
	if after.Code != 401 {
		t.Fatal("revoked device still authorized")
	}
}
func TestPairingExpiresAndRateLimits(t *testing.T) {
	oldBase := base
	base = t.TempDir()
	defer func() { base = oldBase }()
	a := &app{key: strings.Repeat("a", 64)}
	tv := pairReq(a, t, "GET", "/tv", "", "127.0.0.1", nil, "")
	if tv.Code != 200 {
		t.Fatal("cannot start pairing")
	}
	for i := 0; i < 6; i++ {
		wrong := pairReq(a, t, "POST", "/api/pairing/request", `{"code":"not-a-code","name":"Unknown"}`, "192.168.0.40", nil, "")
		if wrong.Code != 401 {
			t.Fatalf("bad-code attempt %d returned %d", i, wrong.Code)
		}
	}
	pairLock.Lock()
	code := pairState.Code
	pairLock.Unlock()
	denied := pairReq(a, t, "POST", "/api/pairing/request", `{"code":"`+code+`"}`, "192.168.0.41", nil, "")
	if denied.Code != 429 {
		t.Fatal("global attempt limit can be bypassed with another IP")
	}
	pairLock.Lock()
	pairState.Expires = time.Now().Add(-time.Second)
	pairLock.Unlock()
	expired := pairReq(a, t, "POST", "/api/pairing/request", `{"code":"`+code+`"}`, "192.168.0.42", nil, "")
	if expired.Code != 429 {
		t.Fatal("expired PIN accepted")
	}
}
