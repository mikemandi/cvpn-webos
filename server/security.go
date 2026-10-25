package main

import (
	"net"
	"net/netip"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Fail closed for subscription origin IPs, even after DNS rebinding.
var nonPublicRanges = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"), netip.MustParsePrefix("10.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"), netip.MustParsePrefix("127.0.0.0/8"),
	netip.MustParsePrefix("169.254.0.0/16"), netip.MustParsePrefix("172.16.0.0/12"),
	netip.MustParsePrefix("192.0.0.0/24"), netip.MustParsePrefix("192.0.2.0/24"),
	netip.MustParsePrefix("192.168.0.0/16"), netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("198.51.100.0/24"), netip.MustParsePrefix("203.0.113.0/24"),
	netip.MustParsePrefix("224.0.0.0/4"), netip.MustParsePrefix("240.0.0.0/4"),
	netip.MustParsePrefix("::/128"), netip.MustParsePrefix("::1/128"),
	netip.MustParsePrefix("fc00::/7"), netip.MustParsePrefix("fe80::/10"),
	netip.MustParsePrefix("2001:db8::/32"),
}

func publicDestination(ip net.IP) bool {
	addr, ok := netip.AddrFromSlice(ip)
	if !ok {
		return false
	}
	addr = addr.Unmap()
	if !addr.IsGlobalUnicast() || addr.IsPrivate() {
		return false
	}
	for _, r := range nonPublicRanges {
		if r.Contains(addr) {
			return false
		}
	}
	return true
}
func allowedHost(hostport string) bool {
	host := hostport
	if strings.Contains(hostport, ":") {
		h, _, err := net.SplitHostPort(hostport)
		if err != nil {
			return false
		}
		host = h
	}
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	if ip.IsLoopback() {
		return true
	}
	// A control panel is only reachable via numeric, private local address.
	return ip.IsPrivate()
}

var sensitiveURL = regexp.MustCompile(`(?i)(?:vless://|https?://)[^\s"'<>]+`)
var uuidLike = regexp.MustCompile(`(?i)\b[0-9a-f]{8}-(?:[0-9a-f]{4}-){3}[0-9a-f]{12}\b`)

func redactSecrets(s string) string {
	s = sensitiveURL.ReplaceAllString(s, "[hidden URL]")
	return uuidLike.ReplaceAllString(s, "[hidden UUID]")
}

type failureBucket struct {
	Count      int
	First      time.Time
	BlockUntil time.Time
}

var failureMu sync.Mutex
var failedAuth = map[string]failureBucket{}

func clientIP(remote string) string {
	ip, _, err := net.SplitHostPort(remote)
	if err != nil {
		return remote
	}
	return ip
}
func exceededFailures(remote string) bool {
	failureMu.Lock()
	defer failureMu.Unlock()
	b, ok := failedAuth[clientIP(remote)]
	return ok && time.Now().Before(b.BlockUntil)
}
func recordFailure(remote string) {
	failureMu.Lock()
	defer failureMu.Unlock()
	now := time.Now()
	id := clientIP(remote)
	b := failedAuth[id]
	if now.Sub(b.First) > 5*time.Minute {
		b = failureBucket{First: now}
	}
	if b.First.IsZero() {
		b.First = now
	}
	b.Count++
	if b.Count >= 12 {
		b.BlockUntil = now.Add(3 * time.Minute)
	}
	failedAuth[id] = b
	// Prevent unbounded memory under rotating IP attacks.
	if len(failedAuth) > 2048 {
		failedAuth = map[string]failureBucket{id: b}
	}
}
func clearFailures(remote string) {
	failureMu.Lock()
	delete(failedAuth, clientIP(remote))
	failureMu.Unlock()
}
