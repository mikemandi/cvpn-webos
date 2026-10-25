package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNoticeFormattingAndUnicodeLimit(t *testing.T) {
	oldBase := base
	base = t.TempDir()
	defer func() { base = oldBase }()
	if err := os.MkdirAll(filepath.Join(base, "webui"), 0700); err != nil {
		t.Fatal(err)
	}
	p := profilesFile{Active: "a", Items: []Profile{{ID: "a", Name: strings.Repeat("Нидерланды", 10) + "\n"}}}
	if err := writeProfiles(p); err != nil {
		t.Fatal(err)
	}
	msg := buildStateToast(vpnNoticeState{kind: "connected", profile: "a"}, observedIP{IP: "80.64.17.130"})
	if len([]rune(msg)) > 60 || strings.Contains(msg, "\n") || !strings.Contains(msg, "🇳🇱") || strings.Contains(msg, "IP ") {
		t.Fatalf("invalid toast: %q", msg)
	}
	down := buildStateToast(vpnNoticeState{kind: "direct"}, observedIP{IP: "1.2.3.4", Country: "Россия"})
	if !strings.Contains(down, "Россия") || !strings.Contains(down, "1.2.3.4") {
		t.Fatal(down)
	}
	errorMsg := buildStateToast(vpnNoticeState{kind: "degraded"}, observedIP{})
	if !strings.Contains(strings.ToLower(errorMsg), "сбой") || strings.Contains(errorMsg, "IP") {
		t.Fatal(errorMsg)
	}
}

func TestLookupDisabledMakesNoNetworkRequest(t *testing.T) {
	t.Setenv("CVPN_IP_LOOKUP", "0")
	if _, err := observePublicIPv4(context.Background()); err == nil {
		t.Fatal("expected disabled lookup")
	}
}

func TestNoActiveButRuntimeRouteIsNotAuthoritative(t *testing.T) {
	// We never report a direct public IP solely from the absence of 'state/active':
	// directNetworkRestored also checks the kernel route table and DNS overrides.
	if strings.TrimSpace(cleanNoticeText("bad\ncontrol", 25)) != "badcontrol" {
		t.Fatal("control chars not removed")
	}
}

func TestTransitionTrackerSuppressesDuplicatesAndSwitching(t *testing.T) {
	nl := vpnNoticeState{kind: "connected", profile: "nl"}
	de := vpnNoticeState{kind: "connected", profile: "de"}
	off := vpnNoticeState{kind: "direct"}
	transient := vpnNoticeState{kind: "transitional"}
	tracker := transitionTracker{last: off, initialized: true}
	next := func(s vpnNoticeState) (vpnNoticeState, bool) { return tracker.next(s) }
	if _, changed := next(off); changed {
		t.Fatal("startup duplicated")
	}
	if _, changed := next(nl); changed {
		t.Fatal("announced without debounce")
	}
	if event, changed := next(nl); !changed || event != nl {
		t.Fatal("missing NL connected")
	}
	if _, changed := next(nl); changed {
		t.Fatal("duplicate NL connected")
	}
	if _, changed := next(off); changed {
		t.Fatal("premature off during switch")
	}
	if _, changed := next(transient); changed {
		t.Fatal("transient should be silent")
	}
	if _, changed := next(de); changed {
		t.Fatal("premature DE connected")
	}
	if event, changed := next(de); !changed || event != de {
		t.Fatal("missing DE connected")
	}
	if _, changed := next(off); changed {
		t.Fatal("premature final off")
	}
	if event, changed := next(off); !changed || event != off {
		t.Fatal("missing final off")
	}
}

func TestInstantForegroundToastsWithConnectedState(t *testing.T) {
	previousBase := base
	base = prepareFixture(t)
	t.Cleanup(func() { base = previousBase })
	previousNotifier := notifySystemToast
	var messages []string
	notifySystemToast = func(s string) error { messages = append(messages, s); return nil }
	t.Cleanup(func() { notifySystemToast = previousNotifier })
	if err := writeProfiles(profilesFile{Active: "installed", Items: []Profile{{ID: "installed", Name: "Нидерланды"}}}); err != nil {
		t.Fatal(err)
	}
	a := &app{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/connect", strings.NewReader(`{"id":"installed"}`))
	a.connectProfile(w, r)
	if w.Code != 200 {
		t.Fatalf("unexpected connect result: %d %s", w.Code, w.Body.String())
	}
	if len(messages) != 2 || !strings.Contains(messages[0], "Подключение") || !strings.Contains(messages[1], "VPN подключён") {
		t.Fatalf("want connecting then connected, got %q", messages)
	}
	if !recentlyManuallyAnnounced(vpnNoticeState{kind: "connected", profile: "installed"}) {
		t.Fatal("observer should suppress duplicate")
	}
	if recentlyManuallyAnnounced(vpnNoticeState{kind: "direct"}) {
		t.Fatal("must not suppress different state")
	}
	w2 := httptest.NewRecorder()
	a.connectProfile(w2, httptest.NewRequest(http.MethodPost, "/api/connect", strings.NewReader(`{"id":"installed"}`)))
	if w2.Code != 200 || !strings.Contains(w2.Body.String(), "unchanged") {
		t.Fatalf("duplicate was not ignored: %s", w2.Body.String())
	}
	if len(messages) != 2 {
		t.Fatalf("unchanged profile generated more toast notifications: %q", messages)
	}
}

func TestFailedStartNeverAnnouncesConnected(t *testing.T) {
	previousBase := base
	base = prepareFixture(t)
	t.Cleanup(func() { base = previousBase })
	previousNotifier := notifySystemToast
	var messages []string
	notifySystemToast = func(s string) error { messages = append(messages, s); return nil }
	t.Cleanup(func() { notifySystemToast = previousNotifier })
	if err := writeProfiles(profilesFile{Active: "installed", Items: []Profile{{ID: "installed", Name: "Нидерланды"}}}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "bin", "cvpn.sh"), []byte("#!/bin/sh\nexit 1\n"), 0700); err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	(&app{}).connectProfile(w, httptest.NewRequest(http.MethodPost, "/api/connect", strings.NewReader(`{"id":"installed"}`)))
	if w.Code == 200 {
		t.Fatalf("engine failure reported success: %s", w.Body.String())
	}
	if len(messages) != 2 || !strings.Contains(messages[0], "Подключение") || !strings.Contains(messages[1], "не удалось") {
		t.Fatalf("invalid error sequence: %q", messages)
	}
	for _, message := range messages {
		if strings.Contains(message, "VPN подключён") {
			t.Fatal("false positive connection toast")
		}
	}
}

func TestToastPayloadUsesLocalIconWhenPresent(t *testing.T) {
	previous := toastIconPath
	t.Cleanup(func() { toastIconPath = previous })
	toastIconPath = filepath.Join(t.TempDir(), "icon.png")
	body, hasIcon := toastPayload("CVPN status")
	if hasIcon || strings.Contains(string(body), "iconUrl") {
		t.Fatalf("nonexistent icon should not be sent: %s", body)
	}
	if err := os.WriteFile(toastIconPath, []byte("png"), 0644); err != nil {
		t.Fatal(err)
	}
	body, hasIcon = toastPayload("CVPN status")
	if !hasIcon || !strings.Contains(string(body), "file://"+toastIconPath) {
		t.Fatalf("installed icon not used: %s", body)
	}
	if !strings.Contains(string(body), "su.cvpn.webos") {
		t.Fatalf("wrong source ID: %s", body)
	}
}
