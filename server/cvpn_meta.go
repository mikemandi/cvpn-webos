package main

import (
	"bufio"
	"encoding/base64"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Published subscription information only. No UUIDs, subscription URLs or secrets.
type SubscriptionMeta struct {
	ExpiresAt        string `json:"expiresAt,omitempty"`
	NeverExpires     bool   `json:"neverExpires,omitempty"`
	UpdateHours      int    `json:"updateHours,omitempty"` // provider recommendation, not an auto-update job
	SupportURL       string `json:"supportUrl,omitempty"`
	PortalURL        string `json:"portalUrl,omitempty"`
	Announcement     string `json:"announcement,omitempty"`
	UsedBytes        int64  `json:"usedBytes,omitempty"`
	TotalBytes       int64  `json:"totalBytes,omitempty"`
	UnlimitedTraffic bool   `json:"unlimitedTraffic,omitempty"`
}

func isCVPNSubscription(raw string) bool {
	u, e := url.Parse(strings.TrimSpace(raw))
	return e == nil && u.Scheme == "https" && u.User == nil && u.Host == "cvpn.su" && strings.HasPrefix(u.EscapedPath(), "/sub/") && len(strings.TrimPrefix(u.EscapedPath(), "/sub/")) > 0
}

func safePublishedURL(raw string) string {
	if len(raw) > 350 {
		return ""
	}
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Scheme != "https" || u.User != nil || u.Hostname() == "" || u.Port() != "" || u.RawQuery != "" || u.Fragment != "" {
		return ""
	}
	host := strings.ToLower(u.Hostname())
	if host != "cvpn.su" && host != "t.me" && host != "telegram.me" {
		return ""
	}
	if strings.HasPrefix(u.EscapedPath(), "/sub/") {
		return ""
	} // never expose a subscription token via public link
	return u.String()
}

func parseSubscriptionUserInfo(raw string, m *SubscriptionMeta) {
	// Missing values are unknown, not zero expiry / unlimited.
	var haveUpload, haveDownload bool
	var upload, download int64
	for _, item := range strings.Split(raw, ";") {
		pair := strings.SplitN(strings.TrimSpace(item), "=", 2)
		if len(pair) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(pair[0]))
		n, e := strconv.ParseInt(strings.TrimSpace(pair[1]), 10, 64)
		if e != nil || n < 0 {
			continue
		}
		switch key {
		case "upload":
			upload = n
			haveUpload = true
		case "download":
			download = n
			haveDownload = true
		case "total":
			m.TotalBytes = n
			m.UnlimitedTraffic = n == 0
		case "expire":
			if n == 0 {
				m.NeverExpires = true
			} else if n > 0 && n < 4102444800 {
				m.ExpiresAt = time.Unix(n, 0).UTC().Format(time.RFC3339)
			}
		}
	}
	if haveUpload || haveDownload {
		m.UsedBytes = upload + download
	}
}

// Provider metadata is accepted only from a verified cvpn.su/sub/* source.
// Header values take precedence over optional comment lines in plaintext subscriptions.
func parseCVPNSubscriptionMeta(source string, headers http.Header, body string) *SubscriptionMeta {
	if !isCVPNSubscription(source) {
		return nil
	}
	fields := map[string]string{}
	scanner := bufio.NewScanner(strings.NewReader(body))
	scanner.Buffer(make([]byte, 512), 8192)
	for i := 0; i < 10 && scanner.Scan(); i++ {
		row := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(row, "#") && !strings.HasPrefix(row, "//") {
			continue
		}
		row = strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(row, "#"), "//"))
		parts := strings.SplitN(row, ":", 2)
		if len(parts) != 2 {
			continue
		}
		fields[strings.ToLower(strings.TrimSpace(parts[0]))] = strings.TrimSpace(parts[1])
	}
	get := func(k string) string {
		v := strings.TrimSpace(headers.Get(k))
		if v == "" {
			v = fields[strings.ToLower(k)]
		}
		return v
	}
	m := &SubscriptionMeta{}
	parseSubscriptionUserInfo(get("subscription-userinfo"), m)
	if interval, e := strconv.Atoi(get("profile-update-interval")); e == nil && interval >= 1 && interval <= 720 {
		m.UpdateHours = interval
	}
	m.SupportURL = safePublishedURL(get("support-url"))
	m.PortalURL = safePublishedURL(get("profile-web-page-url"))
	announcement := get("announce")
	if strings.HasPrefix(strings.ToLower(announcement), "base64:") {
		if dec, err := base64.StdEncoding.DecodeString(announcement[7:]); err == nil {
			announcement = string(dec)
		}
	}
	m.Announcement = safeLabel(announcement, 150)
	if m.ExpiresAt == "" && !m.NeverExpires && m.UpdateHours == 0 && m.SupportURL == "" && m.PortalURL == "" && m.Announcement == "" && m.UsedBytes == 0 && m.TotalBytes == 0 && !m.UnlimitedTraffic {
		return nil
	}
	return m
}

// Source check is performed server side: the UI never receives the raw URL.
func officialGroupSource(g Group) bool {
	if g.Kind != "subscription" || !safeID.MatchString(g.ID) {
		return false
	}
	b, e := os.ReadFile(path("webui", "secrets", "subscription-"+g.ID+".txt"))
	return e == nil && isCVPNSubscription(string(b))
}
