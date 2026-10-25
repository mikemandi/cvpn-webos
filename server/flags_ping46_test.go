package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFlagSelectionAndNoticeLabels46(t *testing.T) {
	cases := []struct{ name, flag string }{
		{"🇳🇱 Нидерланды", "🇳🇱"},
		{"🇫🇷 Париж", "🇫🇷"},
		{"Париж", "🇫🇷"},
		{"🇹🇷 Стамбул (через РФ)", "🇹🇷"},
		{"Santa Clara", "🇺🇸"},
		{"[DE] Germany", "🇩🇪"},
		{"(GB) London", "🇬🇧"},
		{"JP Tokyo", "🇯🇵"},
		{"🇦🇷 Buenos Aires", "🇦🇷"},
		{"Custom Private", ""},
	}
	for _, tc := range cases {
		if got := countryFlagForNotice(tc.name); got != tc.flag {
			t.Errorf("flag %q: got %q want %q", tc.name, got, tc.flag)
		}
		label := noticeProfileLabel(tc.name)
		if tc.flag != "" && strings.Count(label, tc.flag) != 1 {
			t.Errorf("duplicate emoji in %q: %q", tc.name, label)
		}
	}
	if got := limitToast(strings.Repeat("a", 60) + "🇳🇱"); len([]rune(got)) > 60 {
		t.Fatal("toast too long")
	}
}

func TestOtherHTTPSProviderCanPing46(t *testing.T) {
	oldBase := base
	oldProbe := lastSubscriptionProbe
	base = prepareFixture(t)
	defer func() { base = oldBase; lastSubscriptionProbe = oldProbe }()
	groupID := "0123456789abcdef01234567"
	profileID := "fedcba9876543210abcdef12"
	data := profilesFile{Groups: []Group{{ID: groupID, Kind: "subscription", Name: "SVR Telecom", Domain: "example.net"}}, Items: []Profile{{ID: profileID, Name: "🇫🇷 Paris", GroupID: groupID}}}
	if err := writeProfiles(data); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(base, "webui", "secrets"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "webui", "secrets", "subscription-"+groupID+".txt"), []byte("https://subscription.example.net/subscription/key"), 0600); err != nil {
		t.Fatal(err)
	}
	// Invalid secret: the probe must not send any network traffic.
	if err := os.WriteFile(filepath.Join(base, "webui", "secrets", profileID+".txt"), []byte("vless://abc@127.0.0.1:443"), 0600); err != nil {
		t.Fatal(err)
	}
	lastSubscriptionProbe = time.Time{}
	request := func() *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/groups/ping", strings.NewReader(`{"id":"`+groupID+`"}`))
		(&app{}).subscriptionPing(w, r)
		return w
	}
	response := request()
	if response.Code != 200 || !strings.Contains(response.Body.String(), `"reachable":false`) {
		t.Fatalf("non-CVPN subscription probe: %d %s", response.Code, response.Body.String())
	}
	if retry := request(); retry.Code != 429 {
		t.Fatalf("missing rate limit: %d %s", retry.Code, retry.Body.String())
	}
	lastSubscriptionProbe = time.Time{}
	if err := os.WriteFile(filepath.Join(base, "webui", "secrets", "subscription-"+groupID+".txt"), []byte("http://example.net/subscription"), 0600); err != nil {
		t.Fatal(err)
	}
	if res := request(); res.Code != 403 {
		t.Fatalf("non-HTTPS source allowed: %d", res.Code)
	}
}
