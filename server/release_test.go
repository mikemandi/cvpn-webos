package main

import (
	"net"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseSecurityGuards(t *testing.T) {
	// Even fully formed Reality links cannot target loopback, metadata or LAN endpoints.
	for _, h := range []string{"127.0.0.1", "192.168.0.1", "169.254.169.254"} {
		link := "vless://11111111-1111-1111-1111-111111111111@" + h + ":443?security=reality&type=xhttp&sni=www.google.com&pbk=test"
		if _, _, _, err := makeVLESS(link); err == nil {
			t.Fatalf("unsafe VLESS endpoint accepted: %s", h)
		}
	}
	for _, v := range []string{"127.0.0.1", "192.168.1.1", "100.64.0.42", "198.18.0.1", "169.254.169.254", "203.0.113.2", "::1", "fc00::1", "2001:db8::1"} {
		if publicDestination(net.ParseIP(v)) {
			t.Fatalf("SSRF guard allowed restricted IP %s", v)
		}
	}
	if !publicDestination(net.ParseIP("1.1.1.1")) {
		t.Fatal("public DoH IP was wrongly rejected")
	}
	for _, v := range []string{"malicious.example:8765", "127.0.0.1.example:8765", "192.0.2.5:8765"} {
		if allowedHost(v) {
			t.Errorf("host validation allowed %s", v)
		}
	}
	if !allowedHost("192.168.0.2:8765") || !allowedHost("127.0.0.1:8765") {
		t.Fatal("local panel host was denied")
	}
	masked := redactSecrets("failed vless://user@host?security=reality https://cvpn.su/sub/SECRET 11111111-1111-1111-1111-111111111111")
	for _, token := range []string{"SECRET", "vless://", "11111111-1111-1111-1111-111111111111"} {
		if strings.Contains(masked, token) {
			t.Fatalf("secret leaked: %q", masked)
		}
	}
}
func TestRateLimitPasswordAttempts(t *testing.T) {
	a := &app{key: strings.Repeat("A", 48)}
	remote := "192.168.0.88:44321"
	defer clearFailures(remote)
	for i := 0; i < 12; i++ {
		r := httptest.NewRequest("GET", "/api/status", nil)
		r.RemoteAddr = remote
		w := httptest.NewRecorder()
		a.api(w, r)
		if w.Code != 401 {
			t.Fatalf("expected auth rejection %d => %d", i, w.Code)
		}
	}
	r := httptest.NewRequest("GET", "/api/status", nil)
	r.RemoteAddr = remote
	r.Header.Set("X-CVPN-Key", a.key)
	w := httptest.NewRecorder()
	a.api(w, r)
	if w.Code != 429 {
		t.Fatalf("expected rate limit, got %d", w.Code)
	}
}
func TestFailedSwitchRestoresPreviousProfile(t *testing.T) {
	oldBase := base
	base = prepareFixture(t)
	defer func() { base = oldBase }()
	if err := bootstrap(); err != nil {
		t.Fatal(err)
	}
	a := &app{key: strings.Repeat("Q", 48)}
	link := "vless://22222222-2222-2222-2222-222222222222@1.1.1.1:443?security=reality&type=xhttp&sni=www.google.com&pbk=test&mode=packet-up#CVPN%20NL"
	w := callAPI(a, "/api/import", map[string]string{"content": link})
	if w.Code != 200 {
		t.Fatalf("import %s", w.Body.String())
	}
	p, _ := readProfiles()
	// The start operation is simulated to fail after selecting the new profile.
	script := "#!/bin/sh\ncase \"$1\" in start) echo 'PROBE_FAILED'; exit 1;; stop) rm -f \"" + path("state", "active") + "\"; echo 'NETWORK_RESTORED=YES';; esac\n"
	if err := os.WriteFile(filepath.Join(base, "bin", "cvpn.sh"), []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	w = callAPI(a, "/api/connect", map[string]string{"id": p.Items[1].ID})
	if w.Code != 409 {
		t.Fatalf("failed start incorrectly accepted: %d %s", w.Code, w.Body.String())
	}
	after, _ := readProfiles()
	if after.Active != "installed" {
		t.Fatalf("selected profile lost on failed connection: %s", after.Active)
	}
	if active() {
		t.Fatal("fake engine remained connected")
	}
}
func TestConcurrentSSHStyleEngineLock(t *testing.T) {
	oldBase := base
	base = prepareFixture(t)
	defer func() { base = oldBase }()
	if err := bootstrap(); err != nil {
		t.Fatal(err)
	}
	lock := path("state", "lock")
	if err := os.Mkdir(lock, 0700); err != nil {
		t.Fatal(err)
	}
	p, _ := readProfiles()
	if err := selectProfileLocked(&p, "installed"); err == nil {
		t.Fatal("ignored engine lock")
	}
	if err := os.Remove(lock); err != nil {
		t.Fatal(err)
	}
	if err := selectProfileLocked(&p, "installed"); err != nil {
		t.Fatal(err)
	}
}
func TestDegradedStatusNeverClaimsConnected(t *testing.T) {
	oldBase := base
	base = prepareFixture(t)
	defer func() { base = oldBase }()
	if err := bootstrap(); err != nil {
		t.Fatal(err)
	}
	for k, v := range map[string]string{"active": "session", "phase": "ready", "mode": "persistent"} {
		if err := os.WriteFile(path("state", k), []byte(v), 0600); err != nil {
			t.Fatal(err)
		}
	}
	runtimeProbe = func() RuntimeHealth { return RuntimeHealth{Healthy: false, Reason: "Missing TUN"} }
	s := status()
	if s["connected"] != false || s["needsRecovery"] != true || s["phase"] != "degraded" {
		t.Fatalf("false positive: %#v", s)
	}
}
