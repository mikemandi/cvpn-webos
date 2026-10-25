package main

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFailedPreflightKeepsCurrentConnection(t *testing.T) {
	prev := base
	base = prepareFixture(t)
	defer func() { base = prev }()
	if err := bootstrap(); err != nil {
		t.Fatal(err)
	}
	p, err := readProfiles()
	if err != nil {
		t.Fatal(err)
	}
	newID := strings.Repeat("a", 24)
	p.Items = append(p.Items, Profile{ID: newID, Name: "Unusable target", Address: "198.51.100.4", Transport: "XHTTP", Security: "REALITY"})
	if err := writeProfiles(p); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path("webui", "secrets", newID+".txt"), []byte("this is not a VLESS link"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path("state", "active"), []byte("session-123"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path("state", "phase"), []byte("ready"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path("state", "mode"), []byte("persistent"), 0600); err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/connect", strings.NewReader(`{"id":"`+newID+`"}`))
	a := &app{key: strings.Repeat("z", 40)}
	r.Header.Set("X-CVPN-Key", a.key)
	r.RemoteAddr = "127.0.0.1:12345"
	a.api(w, r)
	if w.Code != 409 {
		t.Fatalf("want 409 preflight rejection got %d: %s", w.Code, w.Body.String())
	}
	if !active() {
		t.Fatal("preflight error disconnected the previously active VPN")
	}
	next, err := readProfiles()
	if err != nil {
		t.Fatal(err)
	}
	if next.Active != "installed" {
		t.Fatalf("profile switched despite preflight error: %q", next.Active)
	}
	if _, err := os.Stat(path("state", "lock")); !os.IsNotExist(err) {
		t.Fatalf("lock retained after rejection: %v", err)
	}
	events := readEvents()
	if len(events) == 0 || events[len(events)-1].Kind != "preflight_failed" {
		t.Fatalf("missing failure event: %+v", events)
	}
}

func TestEventTimelineBoundedAndPrivate(t *testing.T) {
	previous := base
	base = t.TempDir()
	defer func() { base = previous }()
	for i := 0; i < 75; i++ {
		recordEvent("connected", "VPN подключён")
	}
	b, e := os.ReadFile(filepath.Join(base, "webui", "state", "connection-events.json"))
	if e != nil {
		t.Fatal(e)
	}
	var ev []auditEvent
	if e = json.Unmarshal(b, &ev); e != nil {
		t.Fatal(e)
	}
	if len(ev) != eventLimit {
		t.Fatalf("event count =%d want %d", len(ev), eventLimit)
	}
	if strings.Contains(string(b), "vless://") {
		t.Fatal("sensitive URI leaked")
	}
	info, e := os.Stat(filepath.Join(base, "webui", "state", "connection-events.json"))
	if e != nil {
		t.Fatal(e)
	}
	if info.Mode().Perm()&0077 != 0 {
		t.Fatalf("events permissions too open: %v", info.Mode().Perm())
	}
	if !strings.Contains(safeEventSummary(), "VPN подключён") {
		t.Fatal("missing readable UI event")
	}
}
