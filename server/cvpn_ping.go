package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type tcpProbeResult struct {
	ID        string `json:"id"`
	MS        int64  `json:"ms,omitempty"`
	Reachable bool   `json:"reachable"`
}

var lastSubscriptionProbe time.Time // protected by mx

func validProbeTarget(raw string) (string, string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Scheme != "vless" || u.User == nil || u.Hostname() == "" {
		return "", "", errors.New("invalid VLESS URI")
	}
	host := u.Hostname()
	port := u.Port()
	pn, e := strconv.Atoi(port)
	if e != nil || pn < 1 || pn > 65535 {
		return "", "", errors.New("invalid destination port")
	}
	// URI includes credentials; they must never be logged or returned to the browser.
	return host, port, nil
}

// Only public destinations are allowed, including after DNS resolution.
// This is a TCP-connect RTT, not an ICMP ping or a full VPN health check.
func tcpEndpointProbe(ctx context.Context, host, port string) (int64, bool) {
	if ip := net.ParseIP(host); ip != nil && !publicDestination(ip) {
		return 0, false
	}
	ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return 0, false
	}
	for _, ip := range ips {
		if ip.IP.To4() == nil || !publicDestination(ip.IP) {
			continue
		}
		start := time.Now()
		c, err := (&net.Dialer{Timeout: 1600 * time.Millisecond}).DialContext(ctx, "tcp", net.JoinHostPort(ip.IP.String(), port))
		if err == nil {
			_ = c.Close()
			ms := time.Since(start).Milliseconds()
			if ms < 1 {
				ms = 1
			}
			return ms, true
		}
	}
	return 0, false
}

func probeBatch(items []Profile) []tcpProbeResult {
	result := make([]tcpProbeResult, len(items))
	sem := make(chan struct{}, 4)
	var wg sync.WaitGroup
	for i, item := range items {
		wg.Add(1)
		go func(i int, p Profile) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			result[i].ID = p.ID
			if !safeID.MatchString(p.ID) {
				return
			}
			secret, err := os.ReadFile(path("webui", "secrets", p.ID+".txt"))
			if err != nil {
				return
			}
			host, port, err := validProbeTarget(string(secret))
			if err != nil {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), 2400*time.Millisecond)
			defer cancel()
			result[i].MS, result[i].Reachable = tcpEndpointProbe(ctx, host, port)
		}(i, item)
	}
	wg.Wait()
	return result
}

func (a *app) subscriptionPing(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string `json:"id"`
	}
	if err := decodeJSON(r, &req); err != nil || !safeID.MatchString(req.ID) {
		errorOut(w, 400, "Некорректная подписка")
		return
	}
	mx.Lock()
	if active() || isEngineLocked() {
		mx.Unlock()
		errorOut(w, 409, "Отключите VPN перед проверкой задержки")
		return
	}
	if time.Since(lastSubscriptionProbe) < 15*time.Second {
		mx.Unlock()
		errorOut(w, 429, "Повторите проверку через несколько секунд")
		return
	}
	p, err := readProfiles()
	if err != nil {
		mx.Unlock()
		errorOut(w, 500, "Список профилей недоступен")
		return
	}
	g := getGroup(p, req.ID)
	if g == nil || g.Kind != "subscription" {
		mx.Unlock()
		errorOut(w, 404, "Подписка не найдена")
		return
	}
	// Only subscriptions that were actually imported via HTTPS can be probed.
	// A provider name or a source domain stored in profiles.json is not proof.
	source, err := os.ReadFile(path("webui", "secrets", "subscription-"+req.ID+".txt"))
	src, parseErr := url.Parse(strings.TrimSpace(string(source)))
	if err != nil || parseErr != nil || src.Scheme != "https" || src.Hostname() == "" || src.User != nil || src.Fragment != "" {
		mx.Unlock()
		errorOut(w, 403, "Проверка доступна для HTTPS-подписок")
		return
	}
	items := make([]Profile, 0, 60)
	for _, item := range p.Items {
		if item.GroupID == g.ID && safeID.MatchString(item.ID) {
			items = append(items, item)
		}
	}
	if len(items) == 0 {
		mx.Unlock()
		errorOut(w, 409, "В подписке нет серверов")
		return
	}
	if len(items) > 60 {
		mx.Unlock()
		errorOut(w, 400, "Слишком много серверов")
		return
	}
	lastSubscriptionProbe = time.Now()
	mx.Unlock() // never hold the panel mutex while making network probes
	jsonResponse(w, 200, map[string]any{"ok": true, "kind": "tcp", "results": probeBatch(items)})
}

func isEngineLocked() bool { _, e := os.Stat(path("state", "lock")); return e == nil }
