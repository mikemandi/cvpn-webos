package main

import (
    "net/http"
    "net"
    "os"
    "strings"
)

// Recovery is deliberately restricted to rollback of a CVPN-owned live session.
// We NEVER flush routes, reset ConnMan or unmount arbitrary system resources.
// The existing Engine owns all network mutations, and remains unchanged.
var recoveryStop = func() (string, error) { return runEngine("stop") }

func ownedSessionSnapshotsExist() bool {
    for _, name := range []string{"active", "active-server-ip", "dns.route", "server.route"} {
        if _, err := os.Stat(path("state", name)); err != nil { return false }
    }
    // The active endpoint must be a valid public IPv4 as required by Engine start.
    b, err := os.ReadFile(path("state", "active-server-ip"))
    if err != nil { return false }
    ip := net.ParseIP(strings.TrimSpace(string(b)))
    if ip == nil || ip.To4() == nil || !publicDestination(ip) { return false }
    return true
}

type recoveryInfo struct {
    OK bool `json:"ok"`
    Active bool `json:"active"`
    Direct bool `json:"direct"`
    CanRestore bool `json:"canRestore"`
    State string `json:"state"`
}

func networkRecoverySnapshot() recoveryInfo {
    live := active()
    direct := directNetworkRestored()
    outcome := recoveryInfo{OK: true, Active: live, Direct: direct, CanRestore:false}
    switch {
    case engineIsTransitioning(): outcome.State="busy"
    case live && ownedSessionSnapshotsExist(): outcome.CanRestore=true; outcome.State="owned_session"
    case live: outcome.State="incomplete_session"
    case direct: outcome.State="clean"
    default: outcome.State="orphaned_resources"
    }
    return outcome
}

func (a *app) restoreNetwork(w http.ResponseWriter, r *http.Request) {
    // POST+CSRF/session checks are done centrally by app.api.
    var body struct { Confirm bool `json:"confirm"` }
    if err := decodeJSON(r, &body); err != nil || !body.Confirm {
        errorOut(w, 400, "Требуется подтверждение восстановления")
        return
    }
    mx.Lock()
    defer mx.Unlock()
    before := networkRecoverySnapshot()
    if before.State == "clean" {
        jsonResponse(w, 200, map[string]any{"ok":true,"alreadyClean":true,"state":before})
        return
    }
    if !before.CanRestore {
        // Do not guess the correct route or DNS state from an incomplete snapshot.
        errorOut(w, 409, "Автоматическое восстановление небезопасно для текущего состояния. Нужна диагностика по SSH; не удаляйте маршруты вручную.")
        return
    }
    out, err := recoveryStop()
    if err != nil || !strings.Contains(out, "NETWORK_RESTORED=YES") || !directNetworkRestored() {
        recordEvent("recovery_needs_attention", "Восстановление не подтвердилось; требуется ручная диагностика")
        errorOut(w, 409, "Восстановление не подтверждено. Проверьте состояние сети в диагностике CVPN.")
        return
    }
    recordEvent("network_restored", "Обычная сеть восстановлена штатным механизмом Engine")
    announceVPNDisconnected()
    jsonResponse(w, 200, map[string]any{"ok":true,"restored":true,"state":networkRecoverySnapshot()})
}
