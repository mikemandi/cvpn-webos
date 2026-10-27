package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// The launcher is deliberately UI-less: a loopback request schedules a system
// alert after the webOS launcher has had time to close. Presenting the alert
// before closing the launcher caused it to disappear along with the app.
const nativeLaunchDelay = 2300 * time.Millisecond

var nativeLaunchState struct {
	sync.Mutex
	generation  uint64
	lastRequest time.Time
	language    string
}

// Overridden in tests. Requests never invoke a system command synchronously.
var nativeSchedule = func(f func()) { go f() }
var nativeDelay = func(d time.Duration) { time.Sleep(d) }
var nativeAlert = showNativePairingAlert

var nativeGif = []byte{71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 0, 0, 0, 255, 255, 255, 33, 249, 4, 1, 0, 0, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 2, 68, 1, 0, 59}

func (a *app) tvNativeNotify(w http.ResponseWriter, r *http.Request) {
	if !localTVRequest(r) {
		http.Error(w, "TV only", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "GET only", http.StatusMethodNotAllowed)
		return
	}
	ip := tvLANAddress()
	if ip == "" {
		http.Error(w, "TV LAN address unavailable", http.StatusServiceUnavailable)
		return
	}

	nativeLaunchState.Lock()
	// Duplicate load/relaunch events should not produce multiple alerts.
	if time.Since(nativeLaunchState.lastRequest) > 1100*time.Millisecond {
		nativeLaunchState.generation++
		nativeLaunchState.lastRequest = time.Now()
		lang := r.URL.Query().Get("lang")
		if lang != "ru" && lang != "en" {
			lang = "en"
		}
		nativeLaunchState.language = lang
		sequence := nativeLaunchState.generation
		nativeSchedule(func() { showAfterLauncherExit(sequence, ip) })
	}
	nativeLaunchState.Unlock()

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "image/gif")
	_, _ = w.Write(nativeGif)
}

func showAfterLauncherExit(sequence uint64, ip string) {
	nativeDelay(nativeLaunchDelay)
	nativeLaunchState.Lock()
	latest := sequence == nativeLaunchState.generation
	nativeLaunchState.Unlock()
	if !latest {
		return
	}

	code, err := issueTVPin()
	if err != nil {
		log.Printf("native pairing PIN creation failed: %v", err)
		return
	}
	nativeLaunchState.Lock()
	tvLanguage := nativeLaunchState.language
	nativeLaunchState.Unlock()
	tvLanguage = languageForTV(tvLanguage)
	nativeLaunchState.Lock()
	nativeLaunchState.language = tvLanguage
	nativeLaunchState.Unlock()
	message := fmt.Sprintf("Open on your phone: http://%s:8765\nPairing code: %s %s (3 minutes)", ip, code[:3], code[3:])
	if tvLanguage == "ru" {
		message = fmt.Sprintf("Откройте на телефоне: http://%s:8765\nКод сопряжения: %s %s (3 минуты)", ip, code[:3], code[3:])
	}
	if err := nativeAlert(message); err != nil {
		log.Printf("native notification failed: %v", err)
	}
}

func nativePairingAlertPayload(message string) ([]byte, error) {
	return nativePairingAlertPayloadLocalized(message, "ru")
}

func nativePairingAlertPayloadLocalized(message string, lang string) ([]byte, error) {
	qrLabel, closeLabel := "Показать QR", "Закрыть"
	if lang == "en" {
		qrLabel, closeLabel = "Show QR code", "Close"
	}
	return json.Marshal(map[string]any{
		"title":   "CVPN webOS",
		"message": message,
		"modal":   true,
		"buttons": []map[string]any{
			{"label": qrLabel, "onclick": "luna://com.webos.applicationManager/launch", "params": map[string]any{"id": "su.cvpn.webos", "params": map[string]any{"mode": "qr"}}},
			{"label": closeLabel, "focus": true},
		},
	})
}

func showNativePairingAlert(message string) error {
	nativeLaunchState.Lock()
	lang := nativeLaunchState.language
	nativeLaunchState.Unlock()
	if pref := readLanguageSetting(); pref == "ru" || pref == "en" {
		lang = pref
	}
	payload, err := nativePairingAlertPayloadLocalized(message, lang)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	binary, err := lunaSendExecutable()
	if err != nil {
		return err
	}
	output, err := exec.CommandContext(ctx, binary, "-n", "1", "-f", "-a", "com.webos.service.secondscreen.gateway", "luna://com.webos.notification/createAlert", string(payload)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("luna-send: %w", err)
	}
	var response struct {
		ReturnValue bool   `json:"returnValue"`
		ErrorText   string `json:"errorText"`
	}
	if err := json.Unmarshal(output, &response); err != nil {
		return fmt.Errorf("bad notification response: %w", err)
	}
	if !response.ReturnValue {
		return fmt.Errorf("notification rejected: %s", strings.TrimSpace(response.ErrorText))
	}
	return nil
}

// TV launcher locale is read-only and returns no credentials or pairing PIN.
// Used during first-run setup to honor the selected CVPN Connect language even
// when the LG web app's navigator.language does not reflect webOS menu locale.
func (a *app) tvLanguageData(w http.ResponseWriter, r *http.Request) {
    if !localTVRequest(r) {
        http.Error(w, "TV only", http.StatusForbidden)
        return
    }
    if r.Method != http.MethodGet {
        http.Error(w, "GET only", http.StatusMethodNotAllowed)
        return
    }
    pref := readLanguageSetting()
    resolved := ""
    if pref == "auto" {
        resolved = detectLGUILanguage()
    }
    payload, _ := json.Marshal(map[string]string{
        "localePreference": pref,
        "resolvedTVLocale": resolved,
    })
    w.Header().Set("Cache-Control", "no-store")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
    _, _ = fmt.Fprintf(w, "window.cvpnLanguageInfo(%s);", payload)
}

// tvQRData is private to the television loopback. It never issues another PIN
// and does not expose the secret to phones on the LAN.
func (a *app) tvQRData(w http.ResponseWriter, r *http.Request) {
	if !localTVRequest(r) {
		http.Error(w, "TV only", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "GET only", http.StatusMethodNotAllowed)
		return
	}
	ip := tvLANAddress()
	pairLock.Lock()
	code, expiry := pairState.Code, pairState.Expires
	pairLock.Unlock()
	pref := readLanguageSetting()
	resolved := ""
	if pref == "auto" {
		resolved = detectLGUILanguage()
	}
	response := map[string]any{"ok": false, "localePreference": pref, "resolvedTVLocale": resolved}
	if ip != "" && len(code) == 6 && time.Now().Before(expiry) {
		response = map[string]any{"ok": true, "address": "http://" + ip + ":8765", "pairURL": (&url.URL{Scheme: "http", Host: ip + ":8765", Fragment: "pair=" + code}).String(), "code": code, "expires": expiry.UTC().Format(time.RFC3339), "localePreference": pref, "resolvedTVLocale": resolved}
	}
	b, _ := json.Marshal(response)
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	_, _ = fmt.Fprintf(w, "window.cvpnQRInfo(%s);", b)
}
