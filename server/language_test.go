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
