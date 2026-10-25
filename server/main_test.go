package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestVLESSXHTTP(t *testing.T) {
	link := "vless://11111111-1111-1111-1111-111111111111@1.1.1.1:443?security=reality&type=xhttp&sni=www.google.com&pbk=ABCdef123&fp=firefox&path=%2Fassets%2Fapi&mode=packet-up#CVPN-NL"
	out, ip, name, e := makeVLESS(link)
	if e != nil {
		t.Fatal(e)
	}
	if ip != "1.1.1.1" || name != "CVPN-NL" {
		t.Fatalf("metadata %q %q", ip, name)
	}
	_, tr, sec, e := metaFromOutbound(out)
	if e != nil || tr != "XHTTP" || sec != "REALITY" {
		t.Fatalf("meta %q %q %v", tr, sec, e)
	}
	stream := out["streamSettings"].(map[string]any)
	xh := stream["xhttpSettings"].(map[string]any)
	if xh["mode"] != "packet-up" || xh["path"] != "/assets/api" {
		t.Fatalf("xhttp settings %#v", xh)
	}
}
func TestVLESSTCP(t *testing.T) {
	link := "vless://11111111-1111-1111-1111-111111111111@1.1.1.1:443?security=reality&type=tcp&sni=www.google.com&pbk=ABCdef123&flow=xtls-rprx-vision#Germany"
	out, _, _, e := makeVLESS(link)
	if e != nil {
		t.Fatal(e)
	}
	_, tr, _, e := metaFromOutbound(out)
	if e != nil || tr != "TCP" {
		t.Fatalf("%q %v", tr, e)
	}
}
func TestRejectUnsafeLinks(t *testing.T) {
	links := []string{"vless://fake@127.0.0.1:443?security=reality", "trojan://x@example.com:443", "vless://11111111-1111-1111-1111-111111111111@127.0.0.1:443?security=tls&type=tcp"}
	for _, v := range links {
		if _, _, _, e := makeVLESS(v); e == nil {
			t.Fatalf("accepted unsafe/unsupported link %q", v)
		}
	}
}
func TestAuthorization(t *testing.T) {
	a := &app{key: strings.Repeat("a", 48)}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/status", nil)
	a.api(w, r)
	if w.Code != 401 {
		t.Fatalf("unauthenticated expected 401 got %d", w.Code)
	}
	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/api/status", nil)
	r.Header.Set("X-CVPN-Key", strings.Repeat("a", 48))
	r.RemoteAddr = "127.0.0.1:12345"
	a.api(w, r)
	if w.Code != 200 {
		t.Fatalf("auth expected 200 got %d", w.Code)
	}
	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/engine", strings.NewReader(`{"action":"rm -rf /"}`))
	r.Header.Set("X-CVPN-Key", strings.Repeat("a", 48))
	r.RemoteAddr = "127.0.0.1:12345"
	r.Header.Set("Origin", "http://evil.com")
	a.api(w, r)
	if w.Code != 403 {
		t.Fatalf("foreign origin expected 403 got %d", w.Code)
	}
	_ = http.MethodPost
}

func TestStatusAfterTrialRollback(t *testing.T) {
	oldProbe := runtimeProbe
	runtimeProbe = func() RuntimeHealth { return RuntimeHealth{Healthy: true} }
	defer func() { runtimeProbe = oldProbe }()
	oldBase := base
	tmp := t.TempDir()
	base = tmp
	defer func() { base = oldBase }()
	if err := os.MkdirAll(filepath.Join(tmp, "state"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "state", "active"), []byte("test\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "state", "mode"), []byte("trial\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "state", "phase"), []byte("ready\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "state", "deadline"), []byte(strconv.FormatInt(time.Now().Unix()+60, 10)), 0600); err != nil {
		t.Fatal(err)
	}
	s := status()
	if s["connected"] != true {
		t.Fatalf("active marker should be connected: %#v", s)
	}
	rem, ok := s["trialRemainingSec"].(int64)
	if !ok || rem < 55 || rem > 61 {
		t.Fatalf("remaining trial seconds: %#v", s)
	}
	if err := os.Remove(filepath.Join(tmp, "state", "active")); err != nil {
		t.Fatal(err)
	}
	s = status()
	if s["connected"] != false {
		t.Fatalf("after rollback should be disconnected: %#v", s)
	}
	if _, ok := s["trialRemainingSec"]; ok {
		t.Fatal("stale countdown after rollback")
	}
}
