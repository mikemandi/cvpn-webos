package main

import (
	"context"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// LAN restriction is based on the actual local TV subnet, never arbitrary
// RFC1918 addresses. Do not admit the VPN TUN interface as a controller LAN.
func localClientAllowed(client net.IP) bool {
	if client == nil {
		return false
	}
	if client.IsLoopback() {
		return true
	}
	for _, name := range []string{"eth0", "wlan0", "wlan1"} {
		nic, err := net.InterfaceByName(name)
		if err != nil || nic.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, err := nic.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			_, subnet, err := net.ParseCIDR(a.String())
			if err == nil && subnet.Contains(client) {
				return true
			}
		}
	}
	return false
}

// Network diagnostics are read-only: the independently guarded Engine alone
// may change DNS, routes or process lifetimes. This monitor does NOT start VPN.
func networkSnapshot() (bool, string) {
	routes, err := os.ReadFile("/proc/net/route")
	if err != nil {
		return false, "Невозможно прочитать маршруты"
	}
	anyUplink := false
	for _, line := range strings.Split(string(routes), "\n") {
		f := strings.Fields(line)
		if len(f) < 8 || f[1] != "00000000" || f[7] != "00000000" {
			continue
		}
		name := f[0]
		if name != "eth0" && name != "wlan0" && name != "wlan1" {
			continue
		}
		anyUplink = true
		nic, err := net.InterfaceByName(name)
		if err != nil || nic.Flags&net.FlagUp == 0 {
			continue
		}
		if b, err := os.ReadFile("/sys/class/net/" + name + "/carrier"); err == nil && strings.TrimSpace(string(b)) == "0" {
			continue
		}
		addrs, err := nic.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			ip, _, err := net.ParseCIDR(a.String())
			if err == nil && ip.To4() != nil && ip.IsPrivate() {
				return true, "Локальная сеть и маршрут по умолчанию доступны: " + name
			}
		}
	}
	if anyUplink {
		return false, "Интерфейс Ethernet/Wi-Fi не готов или не имеет локального IPv4"
	}
	return false, "Маршрут Ethernet/Wi-Fi по умолчанию не найден"
}

var recentNetwork struct {
	sync.RWMutex
	checked time.Time
	ok      bool
	detail  string
}

func releaseHealth() map[string]any {
	up, detail := networkSnapshot()
	p := inspectEngineRuntime
	vpn := active()
	runtime := "off"
	if vpn {
		health := p()
		if health.Healthy {
			runtime = "running"
		} else {
			runtime = "needs_attention"
		}
	}
	return map[string]any{
		"version": version, "network": map[string]any{"ethernetReady": up, "uplinkReady": up, "detail": detail},
		"vpn": runtime, "autoConnect": false, "failOpen": true,
		"recovery":        safeBootRecoverySummary(),
		"supportedUplink": "eth0 / wlan0 / wlan1", "encryptedPanel": false,
		"notes": "VPN requires manual connection; local HTTP requires a trusted LAN",
	}
}

// A debounced read-only observer. The existing engine watchdog remains the
// single writer and recovery owner; this observer records only transitions.
func watchNetworkHealth(ctx context.Context) {
	tick := time.NewTicker(15 * time.Second)
	defer tick.Stop()
	first, _ := networkSnapshot()
	last := first
	candidate := first
	count := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
		}
		now, reason := networkSnapshot()
		recentNetwork.Lock()
		recentNetwork.checked, recentNetwork.ok, recentNetwork.detail = time.Now(), now, reason
		recentNetwork.Unlock()
		if now == last {
			count = 0
			candidate = now
			continue
		}
		if now != candidate {
			candidate = now
			count = 1
			continue
		}
		count++
		if count < 2 {
			continue
		}
		last = now
		count = 0
		if now {
			recordEvent("ethernet_recovered", "Соединение с сетью восстановлено; VPN не включался автоматически")
		} else {
			recordEvent("ethernet_lost", "Локальное подключение к сети отсутствует; безопасность маршрутов контролирует Engine")
		}
	}
}
