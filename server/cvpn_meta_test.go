package main

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestCVPNMetaFromHeadersAndBody(t *testing.T) {
	h := http.Header{}
	h.Set("Subscription-Userinfo", "upload=1024; download=2048; total=4096; expire=1800000000")
	h.Set("Profile-Update-Interval", "1")
	h.Set("Support-Url", "https://t.me/official_cvpn_bot")
	body := "#profile-update-interval: 8\n#support-url: https://t.me/someone\nvless://x\n"
	m := parseCVPNSubscriptionMeta("https://cvpn.su/sub/TEST_TOKEN", h, body)
	if m == nil || m.UpdateHours != 1 || m.UsedBytes != 3072 || m.TotalBytes != 4096 || m.SupportURL != "https://t.me/official_cvpn_bot" {
		t.Fatalf("incorrect metadata: %+v", m)
	}
	if got := time.Unix(1800000000, 0).UTC().Format(time.RFC3339); m.ExpiresAt != got {
		t.Fatalf("expiry=%s want %s", m.ExpiresAt, got)
	}
}
func TestCVPNMetadataOfficialOnlyAndNoTokenLeaks(t *testing.T) {
	h := http.Header{"Support-Url": []string{"https://cvpn.su/sub/SECRETTOKEN"}}
	for _, source := range []string{
		"https://evil.test/sub/key", "https://cvpn.su.evil.test/sub/key", "https://cvpn.su/other/key",
		"http://cvpn.su/sub/key", "https://cvpn.su:8443/sub/key", "https://cvpn.su/sub/",
	} {
		if isCVPNSubscription(source) {
			t.Fatalf("untrusted source accepted: %q", source)
		}
		if meta := parseCVPNSubscriptionMeta(source, h, "#support-url: https://t.me/bot"); meta != nil {
			t.Fatalf("metadata leaked: %+v", meta)
		}
	}
	if m := parseCVPNSubscriptionMeta("https://cvpn.su/sub/key", h, ""); m != nil {
		t.Fatalf("secret URL leaked: %+v", m)
	}
	for _, raw := range []string{"javascript:alert(1)", "http://t.me/unsafe", "https://t.me/user?token=secret", "https://cvpn.su/sub/abcdefg", "https://google.com"} {
		if safePublishedURL(raw) != "" {
			t.Fatalf("allowed untrusted support URL %q", raw)
		}
	}
}
func TestCVPNMetadataMissingIsUnknown(t *testing.T) {
	if m := parseCVPNSubscriptionMeta("https://cvpn.su/sub/key", http.Header{}, "vless://example"); m != nil {
		t.Fatalf("invented provider fields: %+v", m)
	}
	if m := parseCVPNSubscriptionMeta("https://cvpn.su/sub/key", http.Header{}, "#profile-update-interval: 4\n#announce: Hello\n"); m == nil || m.UpdateHours != 4 || m.Announcement != "Hello" {
		t.Fatalf("comment metadata missing: %+v", m)
	}
	h := http.Header{"Subscription-Userinfo": []string{"total=0; expire=0"}}
	m := parseCVPNSubscriptionMeta("https://cvpn.su/sub/key", h, "")
	if m == nil || !m.UnlimitedTraffic || !m.NeverExpires {
		t.Fatalf("explicit unlimited ignored %+v", m)
	}
}
func TestValidateTCPProbeTarget(t *testing.T) {
	if _, _, err := validProbeTarget("vless://" + strings.Repeat("1", 36) + "@example.org:443?type=grpc#NL"); err != nil {
		t.Fatal(err)
	}
	for _, raw := range []string{"http://example.org:443", "vless://id@127.0.0.1", "vless://id@example.org:0", "vless://id@example.org:99999"} {
		if _, _, err := validProbeTarget(raw); err == nil {
			t.Fatalf("invalid accepted %q", raw)
		}
	}
}

func TestOfficialSubscriptionPublicListDoesNotTrustDomain(t *testing.T) {
	old := base
	base = prepareFixture(t)
	defer func() { base = old }()
	official := Group{ID: "0123456789abcdef01234567", Kind: "subscription", Domain: "cvpn.su"}
	if officialGroupSource(official) {
		t.Fatal("verified domain without source secret")
	}
	if e := atomic(path("webui", "secrets", "subscription-"+official.ID+".txt"), []byte("https://cvpn.su/other/key"), 0600); e != nil {
		t.Fatal(e)
	}
	if officialGroupSource(official) {
		t.Fatal("untrusted CVPN path accepted")
	}
	if e := atomic(path("webui", "secrets", "subscription-"+official.ID+".txt"), []byte("https://cvpn.su/sub/valid"), 0600); e != nil {
		t.Fatal(e)
	}
	if !officialGroupSource(official) {
		t.Fatal("real CVPN URL not accepted")
	}
}
