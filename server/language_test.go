package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLanguageSettings(t *testing.T) {
	original := base
	base = t.TempDir()
	defer func() { base = original }()
	a := &app{}
	for _, test := range []struct {
		val  string
		code int
	}{{"en", 200}, {"ru", 200}, {"auto", 200}, {"fr", 400}, {"", 400}} {
		r := httptest.NewRequest(http.MethodPost, "/api/language", strings.NewReader(`{"language":"`+test.val+`"}`))
		w := httptest.NewRecorder()
		a.languageAPI(w, r)
		if w.Code != test.code {
			t.Fatalf("%q: status %d", test.val, w.Code)
		}
		if test.code == 200 && readLanguageSetting() != test.val {
			t.Fatalf("saved language wrong: %q", readLanguageSetting())
		}
	}
	b, err := os.ReadFile(filepath.Join(base, "webui", "language.txt"))
	if err != nil || strings.TrimSpace(string(b)) != "auto" {
		t.Fatal("locale not stored safely")
	}
	r := httptest.NewRequest(http.MethodGet, "/api/language", nil)
	w := httptest.NewRecorder()
	a.languageAPI(w, r)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"language":"auto"`) {
		t.Fatalf("GET wrong: %s", w.Body.String())
	}
}

func TestNativePairingAlertEnglishButtons(t *testing.T) {
	data, err := nativePairingAlertPayloadLocalized("test", "en")
	if err != nil || !strings.Contains(string(data), "Show QR code") || !strings.Contains(string(data), "Close") {
		t.Fatalf("english buttons: %s", data)
	}
}

func TestTVLauncherLanguageEndpointLocalAndExplicitPreference(t *testing.T) {
 original := base
 base = t.TempDir()
 defer func(){base=original}()
 a:= &app{}
 // Manual preference must be available to Launcher without exposing secrets.
 post:=httptest.NewRequest(http.MethodPost,"/api/language",strings.NewReader(`{"language":"en"}`))
 a.languageAPI(httptest.NewRecorder(),post)
 remote:=httptest.NewRequest(http.MethodGet,"/tv/language.js",nil)
 remote.RemoteAddr="192.168.0.12:4567"
 w:=httptest.NewRecorder()
 a.tvLanguageData(w,remote)
 if w.Code!=http.StatusForbidden {t.Fatalf("non-TV could request locale endpoint: %d",w.Code)}
 local:=httptest.NewRequest(http.MethodGet,"/tv/language.js",nil)
 local.RemoteAddr="127.0.0.1:1234"
 w=httptest.NewRecorder()
 a.tvLanguageData(w,local)
 if w.Code!=200 {t.Fatalf("local locale rejected: %d",w.Code)}
 result:=w.Body.String()
 if !strings.Contains(result,`"localePreference":"en"`) ||
    !strings.Contains(result,"window.cvpnLanguageInfo(") ||
    strings.Contains(result,"access.key") || strings.Contains(result,"pairURL") {
      t.Fatalf("unexpected localized response: %s",result)
 }
 w=httptest.NewRecorder()
 local.Method=http.MethodPost
 a.tvLanguageData(w,local)
 if w.Code!=http.StatusMethodNotAllowed {t.Fatalf("POST permitted: %d",w.Code)}
}
