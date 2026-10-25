package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// RuntimeHealth is deliberately read-only. The safety-critical network cleanup
// remains owned by the existing, independently guarded CVPN Engine.
type RuntimeHealth struct {
	Healthy bool
	Reason  string
}

// Keep probe swappable in tests, so no test manipulates a real network interface.
var runtimeProbe = inspectEngineRuntime

func inspectEngineRuntime() RuntimeHealth {
	pidBytes, err := os.ReadFile(path("state", "xray.pid"))
	if err != nil {
		return RuntimeHealth{Reason: "Xray PID отсутствует"}
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
	if err != nil || pid < 2 {
		return RuntimeHealth{Reason: "Некорректный PID Xray"}
	}
	exe, err := os.Readlink(filepath.Join("/proc", strconv.Itoa(pid), "exe"))
	if err != nil || exe != path("bin", "xray") {
		return RuntimeHealth{Reason: "Не найден запущенный Xray CVPN"}
	}
	cmdline, err := os.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "cmdline"))
	if err != nil || (!strings.Contains(string(cmdline), path("config", "client.json")) && !strings.Contains(string(cmdline), path("state", "client.runtime.json"))) {
		return RuntimeHealth{Reason: "Xray работает с неизвестным конфигом"}
	}
	if _, err = os.Stat("/sys/class/net/cvpn0"); err != nil {
		return RuntimeHealth{Reason: "Сетевой интерфейс cvpn0 отсутствует"}
	}
	mounts, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return RuntimeHealth{Reason: "Невозможно проверить DNS"}
	}
	mounted := false
	for _, line := range strings.Split(string(mounts), "\n") {
		fields := strings.Fields(line)
		if len(fields) > 1 && fields[1] == "/var/lib/misc/resolv.conf" {
			mounted = true
			break
		}
	}
	if !mounted {
		return RuntimeHealth{Reason: "Защищённый DNS не смонтирован"}
	}
	// Directly inspect Linux's IPv4 route table rather than spawning shell commands.
	// Routes are to cvpn0 and cover the internet using two /1 ranges.
	routes, err := os.ReadFile("/proc/net/route")
	if err != nil {
		return RuntimeHealth{Reason: "Невозможно прочитать маршруты"}
	}
	lower, upper := false, false
	for _, line := range strings.Split(string(routes), "\n") {
		f := strings.Fields(line)
		if len(f) < 8 || f[0] != "cvpn0" {
			continue
		}
		if f[1] == "00000000" && f[7] == "00000080" {
			lower = true
		}
		if f[1] == "00000080" && f[7] == "00000080" {
			upper = true
		}
	}
	if !lower || !upper {
		return RuntimeHealth{Reason: "TUN-маршруты повреждены"}
	}
	return RuntimeHealth{Healthy: true}
}

func engineIsTransitioning() bool {
	if _, err := os.Stat(path("state", "lock")); err == nil {
		return true
	}
	return false
}

// Acquire the same lock directory used by cvpn.sh, but NEVER retain this lock
// while starting/stopping cvpn.sh itself (the shell manager needs its own lock).
func reserveEngineConfig() (func(), error) {
	dir := path("state", "lock")
	if err := os.Mkdir(dir, 0700); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(dir, "pid"), []byte(strconv.Itoa(os.Getpid())+"\n"), 0600); err != nil {
		_ = os.Remove(dir)
		return nil, err
	}
	return func() { _ = os.Remove(filepath.Join(dir, "pid")); _ = os.Remove(dir) }, nil
}

// If a trial is about to expire, callers must not start risky config writes.
func trialExpired() bool {
	b, err := os.ReadFile(path("state", "deadline"))
	if err != nil {
		return false
	}
	deadline, err := strconv.ParseInt(strings.TrimSpace(string(b)), 10, 64)
	return err == nil && time.Now().Unix() >= deadline
}
