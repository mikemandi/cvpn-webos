package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFingerprintPreservesDifferentIPv4Endpoints(t *testing.T) {
	base := "vless://11111111-1111-1111-1111-111111111111@"
	q := ":443?security=reality&type=xhttp&sni=example.com&pbk=test#"
	nl := base + "1.1.1.1" + q + "CVPN%20NL"
	de := base + "1.0.0.1" + q + "CVPN%20DE"
	if secretFingerprint(nl) == secretFingerprint(de) {
		t.Fatal("different IP endpoints were incorrectly deduplicated")
	}
	duplicateByLabel := base + "1.1.1.1" + q + "CVPN%20DE"
	if secretFingerprint(nl) == secretFingerprint(duplicateByLabel) {
		t.Fatal("distinct labels represent separate logical servers")
	}
}

func TestDuplicateReportIdentifiesMatchWithoutSecrets(t *testing.T) {
	p := Profile{Name: "CVPN · NL Нидерланды", Address: "1.1.1.1"}
	result := duplicateProfileReason(p)
	if !strings.Contains(result, "Нидерланды") || !strings.Contains(result, "1.1.1.1") {
		t.Fatalf("missing match detail: %s", result)
	}
	bad := duplicateProfileReason(Profile{Name: "vless://some-secret@1.2.3.4:443?uuid=secret", Address: "1.0.0.1"})
	if strings.Contains(bad, "some-secret") || strings.Contains(bad, "uuid") {
		t.Fatalf("leaked secret %s", bad)
	}
}

func TestImportReportDifferentEndpointsAndSameLinkLabel(t *testing.T) {
	oldBase := base
	base = prepareFixture(t)
	defer func() { base = oldBase }()
	if err := bootstrap(); err != nil {
		t.Fatal(err)
	}
	prev := subscriptionFetcher
	defer func() { subscriptionFetcher = prev }()
	q := ":443?security=reality&type=xhttp&sni=www.google.com&pbk=test#CVPN%20"
	u1 := "vless://22222222-2222-2222-2222-222222222222@1.1.1.1" + q + "NL"
	u2 := "vless://22222222-2222-2222-2222-222222222222@1.0.0.1" + q + "DE"
	u3 := "vless://22222222-2222-2222-2222-222222222222@1.1.1.1" + q + "Germany"
	subscriptionFetcher = func(_ string) ([]string, string, string, *SubscriptionMeta, error) { return []string{u1, u2, u3}, "CVPN", "cvpn.su", nil, nil }
	a := &app{key: strings.Repeat("x", 40)}
	resp := callAPI(a, "/api/import", map[string]string{"content": "https://cvpn.su/sub/test-secret"})
	if resp.Code != 200 {
		t.Fatalf("status %d: %s", resp.Code, resp.Body.String())
	}
	var result struct {
		Imported int           `json:"imported"`
		Issues   []ImportIssue `json:"issues"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Imported != 3 || len(result.Issues) != 0 {
		t.Fatalf("unexpected %+v", result)
	}

	if strings.Contains(resp.Body.String(), "22222222-2222") || strings.Contains(resp.Body.String(), "sub/test-secret") {
		t.Fatal("secret leaked")
	}
}

func TestDuplicateRefreshReasonShowsEndpointWithoutUUID(t *testing.T) {
	source := "vless://22222222-2222-2222-2222-222222222222@1.1.1.1:443?security=reality&type=xhttp#CVPN%20NL"
	got := duplicateInSourceReason(0, source)
	if !strings.Contains(got, "1.1.1.1") || !strings.Contains(got, "CVPN NL") {
		t.Fatalf("missing endpoint detail: %q", got)
	}
	if strings.Contains(got, "22222222") || strings.Contains(got, "security=reality") {
		t.Fatalf("secret leaked: %q", got)
	}
}
