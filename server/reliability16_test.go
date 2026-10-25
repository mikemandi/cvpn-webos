package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestXHTTPExtraPreservedAndInvalidRejected(t *testing.T) {
	uri := "vless://11111111-1111-1111-1111-111111111111@1.1.1.1:443?security=reality&type=xhttp&sni=www.google.com&pbk=test&extra=%7B%22xPaddingBytes%22%3A%22100-200%22%2C%22xmux%22%3A%7B%22maxConcurrency%22%3A%223-5%22%7D%7D#CVPN-DE"
	outbound, _, _, err := makeVLESS(uri)
	if err != nil {
		t.Fatal(err)
	}
	xh := outbound["streamSettings"].(map[string]any)["xhttpSettings"].(map[string]any)
	extra := xh["extra"].(map[string]any)
	if extra["xPaddingBytes"] != "100-200" {
		t.Fatalf("lost padding: %#v", extra)
	}
	if extra["xmux"].(map[string]any)["maxConcurrency"] != "3-5" {
		t.Fatalf("lost nested: %#v", extra)
	}
	b, err := json.Marshal(outbound)
	if err != nil || !strings.Contains(string(b), "xPaddingBytes") {
		t.Fatalf("marshal: %v %s", err, b)
	}
	bad := strings.Replace(uri, "100-200", "100-200", 1) // preserves the valid URL before replacing extra
	bad = strings.Replace(bad, "%7B%22xPaddingBytes%22%3A%22100-200%22%2C%22xmux%22%3A%7B%22maxConcurrency%22%3A%223-5%22%7D%7D", "%5B1%2C2%5D", 1)
	if _, _, _, err := makeVLESS(bad); err == nil {
		t.Fatal("accepted extra array")
	}
}

func TestPartialRefreshPreservesOldAndReportsGermany(t *testing.T) {
	prev := base
	base = prepareFixture(t)
	defer func() { base = prev }()
	if err := bootstrap(); err != nil {
		t.Fatal(err)
	}
	oldFetcher := subscriptionFetcher
	defer func() { subscriptionFetcher = oldFetcher }()
	nl := "vless://22222222-2222-2222-2222-222222222222@1.1.1.1:443?security=reality&type=xhttp&sni=google.com&pbk=test#CVPN%20%C2%B7%20NL"
	de := "vless://33333333-3333-3333-3333-333333333333@1.1.1.1:443?security=reality&type=xhttp&sni=google.com&pbk=test#CVPN%20%C2%B7%20DE"
	us := "vless://44444444-4444-4444-4444-444444444444@1.1.1.1:443?security=reality&type=tcp&sni=google.com&pbk=test#CVPN%20%C2%B7%20US"
	content := []string{nl, de}
	subscriptionFetcher = func(_ string) ([]string, string, string, *SubscriptionMeta, error) { return content, "CVPN", "cvpn.su", nil, nil }
	a := &app{key: strings.Repeat("k", 40)}
	source := "https://cvpn.su/sub/test-secret-key-123"
	resp := callAPI(a, "/api/import", map[string]string{"content": source})
	if resp.Code != 200 {
		t.Fatalf("import %d %s", resp.Code, resp.Body.String())
	}
	p, _ := readProfiles()
	gid := subscriptionID(source)
	if len(p.Items) != 3 {
		t.Fatalf("profiles %d", len(p.Items))
	}
	oldIDs := map[string]bool{}
	for _, item := range p.Items {
		if item.GroupID == gid {
			oldIDs[item.ID] = true
		}
	}
	invalidDE := strings.Replace(de, "type=xhttp", "type=unknown", 1)
	content = []string{nl, invalidDE, us}
	resp = callAPI(a, "/api/groups/refresh", map[string]string{"id": gid})
	if resp.Code != 200 {
		t.Fatalf("refresh %d %s", resp.Code, resp.Body.String())
	}
	var report struct {
		Partial  bool          `json:"partial"`
		Retained int           `json:"retained"`
		Removed  int           `json:"removed"`
		Issues   []ImportIssue `json:"issues"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if !report.Partial || report.Retained != 1 || report.Removed != 0 || len(report.Issues) != 1 {
		t.Fatalf("bad report %+v", report)
	}
	if !strings.Contains(report.Issues[0].Server, "DE") && !strings.Contains(report.Issues[0].Server, "Сервер") {
		t.Fatalf("missing label %+v", report)
	}
	if strings.Contains(resp.Body.String(), "test-secret-key-123") || strings.Contains(resp.Body.String(), "33333333-3333") {
		t.Fatal("secret disclosed in report")
	}
	p, _ = readProfiles()
	if len(p.Items) != 4 {
		t.Fatalf("expected old NL+DE and new US, got %d", len(p.Items))
	}
	for id := range oldIDs {
		if !containsID(p.Items, id) {
			t.Fatalf("deleted old ID %s", id)
		}
		if _, err := os.Stat(filepath.Join(base, "webui", "secrets", id+".txt")); err != nil {
			t.Fatal(err)
		}
	}
	// A completely unrecognized refresh must leave the saved list intact.
	content = []string{invalidDE}
	resp = callAPI(a, "/api/groups/refresh", map[string]string{"id": gid})
	if resp.Code != 422 {
		t.Fatalf("expected 422, got %d: %s", resp.Code, resp.Body.String())
	}
	q, _ := readProfiles()
	if len(q.Items) != 4 {
		t.Fatal("refresh with zero valid profiles deleted saved entries")
	}
}

func TestSafeReportNeverLeaksCredentialInFragment(t *testing.T) {
	issue := newImportIssue(0, "vless://secret@1.1.1.1:443#token%3Athis-secret-value", "Bad")
	b, _ := json.Marshal(issue)
	if strings.Contains(string(b), "secret-value") || strings.Contains(string(b), "vless://") {
		t.Fatalf("unsafe report: %s", b)
	}
	if strings.Contains(importFailureReason(os.ErrNotExist), "/") {
		t.Fatal("unexpected raw error")
	}
}
