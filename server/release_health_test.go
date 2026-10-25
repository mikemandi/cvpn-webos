package main

import (
	"encoding/json"
	"net"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReleaseHealthNoAutoConnect(t *testing.T) {
	oldBase := base
	base = t.TempDir()
	defer func() { base = oldBase }()
	out := releaseHealth()
	if out["autoConnect"] != false || out["failOpen"] != true || out["encryptedPanel"] != false {
		t.Fatal(out)
	}
	b, _ := json.Marshal(out)
	if strings.Contains(string(b), "vless://") {
		t.Fatal("private subscription leaked")
	}
}
func TestLegacyHeaderIsLoopbackOnly(t *testing.T) {
	// API guard exists on actual HTTP entrypoint, while loopback utility access
	// used by the installer and offline tests remains supported.
	a := &app{key: strings.Repeat("k", 48)}
	r := httptest.NewRequest("GET", "/api/status", nil)
	r.RemoteAddr = "192.168.0.42:1234"
	r.Header.Set("X-CVPN-Key", a.key)
	w := httptest.NewRecorder()
	a.api(w, r)
	if w.Code != 401 {
		t.Fatalf("LAN legacy secret accepted: %d", w.Code)
	}
	r.RemoteAddr = "127.0.0.1:1234"
	w = httptest.NewRecorder()
	a.api(w, r)
	if w.Code != 200 {
		t.Fatalf("loopback installer rejected: %d", w.Code)
	}
}
func TestLegacyPairingAuthHeaderCannotBeUsedFromLAN(t *testing.T) {
	a := &app{key: strings.Repeat("k", 48)}
	r := httptest.NewRequest("GET", "/api/auth", nil)
	r.RemoteAddr = "192.168.0.42:1234"
	r.Header.Set("X-CVPN-Key", a.key)
	w := httptest.NewRecorder()
	a.api(w, r)
	if w.Code != 401 {
		t.Fatalf("LAN legacy key accepted by pairing endpoint: %d", w.Code)
	}
}
func TestLANFilteringIsNotBroadPrivateRange(t *testing.T) {
	for _, s := range []string{"8.8.8.8", "203.0.113.5", "100.64.0.1"} {
		if localClientAllowed(net.ParseIP(s)) {
			t.Fatal(s)
		}
	}
	if !localClientAllowed(net.ParseIP("127.0.0.1")) {
		t.Fatal("loopback rejected")
	}
}
