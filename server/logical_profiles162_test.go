package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

const logicalLinkBase = "vless://22222222-2222-2222-2222-222222222222@"
const logicalQuery = ":443?security=reality&type=xhttp&sni=www.google.com&pbk=test&mode=packet-up#"

func TestLogicalChoicesShareEndpointWithoutBeingDeduplicated(t *testing.T) {
	oldBase := base
	base = prepareFixture(t)
	defer func() { base = oldBase }()
	if err := bootstrap(); err != nil {
		t.Fatal(err)
	}
	oldFetcher := subscriptionFetcher
	defer func() { subscriptionFetcher = oldFetcher }()
	auto := logicalLinkBase + "1.0.0.1" + logicalQuery + "EU%20Automatic%20Selection"
	de := logicalLinkBase + "1.0.0.1" + logicalQuery + "DE%20Germany"
	nl := logicalLinkBase + "1.1.1.1" + logicalQuery + "NL%20Netherlands"
	us := logicalLinkBase + "193.233.131.83" + logicalQuery + "US%20America"
	links := []string{auto, nl, de, us}
	subscriptionFetcher = func(_ string) ([]string, string, string, *SubscriptionMeta, error) { return links, "CVPN", "cvpn.su", nil, nil }
	a := &app{key: strings.Repeat("x", 40)}
	source := "https://cvpn.su/sub/private-NEVER-SHOW"
	w := callAPI(a, "/api/import", map[string]string{"content": source})
	if w.Code != 200 {
		t.Fatalf("import %d %s", w.Code, w.Body.String())
	}
	var report struct {
		Imported int `json:"imported"`
		Skipped  int `json:"skipped"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.Imported != 4 || report.Skipped != 0 {
		t.Fatalf("report %+v", report)
	}
	p, err := readProfiles()
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Items) != 5 {
		t.Fatalf("want installed and 4 logical profiles, got %d", len(p.Items))
	}
	autoID, deID := p.Items[1].ID, p.Items[3].ID
	if autoID == deID || p.Items[1].Address != p.Items[3].Address {
		t.Fatal("different logical choices must share IP and have distinct IDs")
	}
	// A no-op refresh does not destroy either profile or change IDs.
	w = callAPI(a, "/api/groups/refresh", map[string]string{"id": subscriptionID(source)})
	if w.Code != 200 {
		t.Fatalf("refresh %d %s", w.Code, w.Body.String())
	}
	p, _ = readProfiles()
	if len(p.Items) != 5 || p.Items[1].ID != autoID || p.Items[3].ID != deID {
		t.Fatalf("lost logical profile after refresh: %+v", p.Items)
	}
	if strings.Contains(w.Body.String(), "NEVER-SHOW") || strings.Contains(w.Body.String(), "22222222-2222") {
		t.Fatal("sensitive data leaked")
	}
}

func TestMigratingPreviousFingerprintAndAutoEndpointChange(t *testing.T) {
	oldBase := base
	base = prepareFixture(t)
	defer func() { base = oldBase }()
	if err := bootstrap(); err != nil {
		t.Fatal(err)
	}
	originalFetcher := subscriptionFetcher
	defer func() { subscriptionFetcher = originalFetcher }()
	auto := logicalLinkBase + "1.0.0.1" + logicalQuery + "EU%20Automatic%20Selection"
	de := logicalLinkBase + "1.0.0.1" + logicalQuery + "DE%20Germany"
	links := []string{auto, de}
	subscriptionFetcher = func(_ string) ([]string, string, string, *SubscriptionMeta, error) { return links, "CVPN", "cvpn.su", nil, nil }
	a := &app{key: strings.Repeat("z", 40)}
	source := "https://cvpn.su/sub/private-test"
	w := callAPI(a, "/api/import", map[string]string{"content": source})
	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
	p, _ := readProfiles()
	id := p.Items[1].ID
	deID := p.Items[2].ID
	// Old backend stored the same hash without the fragment. Simulate upgrade.
	raw, err := os.ReadFile(path("webui", "secrets", id+".txt"))
	if err != nil {
		t.Fatal(err)
	}
	p.Items[1].Fingerprint = secretFingerprint(strings.Split(string(raw), "#")[0])
	if err := writeProfiles(p); err != nil {
		t.Fatal(err)
	}
	w = callAPI(a, "/api/select", map[string]string{"id": id})
	if w.Code != 200 {
		t.Fatalf("select %d %s", w.Code, w.Body.String())
	}
	// Automatic entry follows a different current endpoint; DE does not change.
	nextAuto := logicalLinkBase + "193.233.131.83" + logicalQuery + "EU%20Automatic%20Selection"
	links = []string{nextAuto, de}
	w = callAPI(a, "/api/groups/refresh", map[string]string{"id": subscriptionID(source)})
	if w.Code != 200 {
		t.Fatalf("refresh %d %s", w.Code, w.Body.String())
	}
	p, _ = readProfiles()
	if p.Active != id || len(p.Items) != 3 || p.Items[1].ID != id || p.Items[2].ID != deID || p.Items[1].Address != "193.233.131.83" {
		t.Fatalf("unstable logical IDs after switch: %+v", p)
	}
	newSecret, _ := os.ReadFile(path("webui", "secrets", id+".txt"))
	if string(newSecret) != nextAuto {
		t.Fatal("auto secret was not refreshed")
	}
	ip, _ := os.ReadFile(path("config", "server-ip.txt"))
	if strings.TrimSpace(string(ip)) != "193.233.131.83" {
		t.Fatalf("selected config IP not updated: %s", ip)
	}
	cfg, _ := os.ReadFile(path("config", "client.json"))
	if !strings.Contains(string(cfg), "193.233.131.83") {
		t.Fatal("selected Xray config was not updated")
	}
}

func TestOnlyIdenticalNamedLinkIsDuplicate(t *testing.T) {
	oldBase := base
	base = prepareFixture(t)
	defer func() { base = oldBase }()
	if err := bootstrap(); err != nil {
		t.Fatal(err)
	}
	orig := subscriptionFetcher
	defer func() { subscriptionFetcher = orig }()
	link := logicalLinkBase + "1.0.0.1" + logicalQuery + "DE%20Germany"
	subscriptionFetcher = func(_ string) ([]string, string, string, *SubscriptionMeta, error) { return []string{link, link}, "CVPN", "cvpn.su", nil, nil }
	a := &app{key: strings.Repeat("z", 40)}
	resp := callAPI(a, "/api/import", map[string]string{"content": "https://cvpn.su/sub/private"})
	if resp.Code != 200 {
		t.Fatal(resp.Body.String())
	}
	p, _ := readProfiles()
	if len(p.Items) != 2 {
		t.Fatalf("duplicated identical entry: %d", len(p.Items))
	}
	if !strings.Contains(resp.Body.String(), "Повтор") {
		t.Fatal("no duplicate warning")
	}
}
