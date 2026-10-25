package main

import (
    "os"
    "strings"
)

// Read only a fixed list of neutral diagnostic status strings. Do not expose
// boot IDs, raw logs or paths through the HTTP API.
func safeBootRecoverySummary() string {
    b, err := os.ReadFile(path("state", "recovery", "status"))
    if err != nil { return "пока не запускалось" }
    switch strings.TrimSpace(string(b)) {
    case "CLEANED_STALE": return "старая сессия очищена, VPN не запускался"
    case "CLEAN": return "сеть не требовала восстановления; VPN выключен"
    case "SKIPPED_ACTIVE": return "активный VPN не затронут"
    case "SKIPPED_UNSAFE", "SKIPPED_LOCK", "ERROR": return "требуется ручная проверка; сеть не изменена"
    case "ALREADY_HANDLED": return "уже проверено при текущей загрузке"
    case "INSTALL_PENDING_REBOOT": return "обновлено; восстановление проверится при следующей полной загрузке"
    case "RUNNING": return "проверка восстановления выполняется"
    default: return "нет достоверных данных"
    }
}
