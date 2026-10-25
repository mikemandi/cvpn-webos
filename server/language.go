package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Preference shared by authenticated mobile settings and the TV QR screen.
func readLanguageSetting() string {
	b, err := os.ReadFile(path("webui", "language.txt"))
	if err != nil {
		return "auto"
	}
	v := strings.TrimSpace(string(b))
	if v != "en" && v != "ru" {
		return "auto"
	}
	return v
}

func (a *app) languageAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		jsonResponse(w, http.StatusOK, map[string]any{"language": readLanguageSetting()})
		return
	}
	if !requirePost(w, r) {
		return
	}
	var req struct {
		Language string `json:"language"`
	}
	if err := decodeJSON(r, &req); err != nil || (req.Language != "auto" && req.Language != "ru" && req.Language != "en") {
		errorOut(w, http.StatusBadRequest, "invalid language")
		return
	}
	if err := atomic(path("webui", "language.txt"), []byte(req.Language+"\n"), 0600); err != nil {
		errorOut(w, http.StatusInternalServerError, "language preference could not be saved")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{"ok": true, "language": req.Language})
}

// LG Settings Service exposes the actual menu language. A denied/unavailable
// service must never prevent showing the pairing QR or a native alert.
func detectLGUILanguage() string {
	binary, err := lunaSendExecutable()
	if err != nil {
		return ""
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	out, err := exec.CommandContext(ctx, binary, "-n", "1", "-f", "-a", "com.webos.service.secondscreen.gateway", "luna://com.webos.settingsservice/getSystemSettings", `{"keys":["localeInfo"]}`).Output()
	if err != nil {
		return ""
	}
	var r struct {
		ReturnValue bool `json:"returnValue"`
		Settings    struct {
			LocaleInfo struct {
				Locales struct {
					UI string `json:"UI"`
				} `json:"locales"`
			} `json:"localeInfo"`
		} `json:"settings"`
	}
	if json.Unmarshal(out, &r) != nil || !r.ReturnValue {
		return ""
	}
	lang := strings.ToLower(r.Settings.LocaleInfo.Locales.UI)
	if lang == "ru" || strings.HasPrefix(lang, "ru-") || strings.HasPrefix(lang, "ru_") {
		return "ru"
	}
	if lang != "" {
		return "en"
	}
	return ""
}

func languageForTV(fallback string) string {
	if pref := readLanguageSetting(); pref == "ru" || pref == "en" {
		return pref
	}
	if detected := detectLGUILanguage(); detected != "" {
		return detected
	}
	if fallback == "ru" {
		return "ru"
	}
	return "en"
}
