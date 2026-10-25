package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNativeLaunchNeedsLoopback(t *testing.T) {
	a := &app{}
	r := httptest.NewRequest(http.MethodGet, "/tv/notify", nil)
	r.RemoteAddr = "192.168.0.50:1234"
	w := httptest.NewRecorder()
	a.tvNativeNotify(w, r)
	if w.Code != http.StatusForbidden {
		t.Fatalf("remote control was allowed: %d", w.Code)
	}
}

func TestNativeLaunchIsAsynchronous(t *testing.T) {
	oldSchedule, oldDelay, oldAlert := nativeSchedule, nativeDelay, nativeAlert
	defer func() { nativeSchedule, nativeDelay, nativeAlert = oldSchedule, oldDelay, oldAlert }()
	nativeLaunchState.Lock()
	nativeLaunchState.generation = 0
	nativeLaunchState.lastRequest = time.Time{}
	nativeLaunchState.Unlock()
	pairLock.Lock()
	oldPair := pairState
	pairState = tvPairState{}
	pairLock.Unlock()
	defer func() { pairLock.Lock(); pairState = oldPair; pairLock.Unlock() }()
	var job func()
	nativeSchedule = func(f func()) { job = f }
	var slept time.Duration
	nativeDelay = func(d time.Duration) { slept = d }
	var message string
	nativeAlert = func(m string) error { message = m; return nil }
	a := &app{}
	r := httptest.NewRequest(http.MethodGet, "/tv/notify", nil)
	r.RemoteAddr = "127.0.0.1:4444"
	w := httptest.NewRecorder()
	a.tvNativeNotify(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected HTTP status %d: %s", w.Code, w.Body.String())
	}
	if w.Header().Get("Content-Type") != "image/gif" {
		t.Fatal("image handoff changed")
	}
	if job == nil {
		t.Fatal("notification was not queued")
	}
	pairLock.Lock()
	before := pairState.Code
	pairLock.Unlock()
	if before != "" {
		t.Fatal("PIN generated before launcher was closed")
	}
	job()
	if slept != nativeLaunchDelay {
		t.Fatalf("notification delay %v", slept)
	}
	if !strings.Contains(message, ":8765") || !strings.Contains(message, "Pairing code:") {
		t.Fatalf("missing pairing information: %q", message)
	}
	pairLock.Lock()
	code := pairState.Code
	pairLock.Unlock()
	if len(code) != 6 || !strings.Contains(message, code[:3]+" "+code[3:]) {
		t.Fatal("incorrect PIN on notification")
	}
}

func TestQRLocalOnlyAndPINSharedWithAlert(t *testing.T) {
	old := pairState
	defer func() { pairLock.Lock(); pairState = old; pairLock.Unlock() }()
	code, err := issueTVPin()
	if err != nil {
		t.Fatal(err)
	}
	a := &app{}
	remote := httptest.NewRequest(http.MethodGet, "/tv/qr-data.js", nil)
	remote.RemoteAddr = "192.168.0.31:8888"
	w := httptest.NewRecorder()
	a.tvQRData(w, remote)
	if w.Code != http.StatusForbidden {
		t.Fatalf("remote PIN leaked: %d", w.Code)
	}
	r := httptest.NewRequest(http.MethodGet, "/tv/qr-data.js", nil)
	r.RemoteAddr = "127.0.0.1:3333"
	w = httptest.NewRecorder()
	a.tvQRData(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("local request failed: %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "window.cvpnQRInfo(") {
		t.Fatal("wrong JSONP output")
	}
	// On CI / computers without eth0 private LAN address, no QR is issued.
	// Crucially, viewing QR never replaces the PIN in the active pairing state.
	pairLock.Lock()
	same := pairState.Code
	pairLock.Unlock()
	if same != code {
		t.Fatal("QR handler unexpectedly replaced PIN")
	}
	payload, err := nativePairingAlertPayload("Connect")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(payload), "Показать QR") || !strings.Contains(string(payload), "su.cvpn.webos") || !strings.Contains(string(payload), "CVPN webOS") {
		t.Fatal("QR button not wired")
	}
}
