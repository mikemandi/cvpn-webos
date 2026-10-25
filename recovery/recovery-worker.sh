#!/bin/sh
# CVPN webOS RC4.9.1. Safe recovery status reconciliation. NEVER VPN autostart.
# Never alters system partitions, DNS, routes, tunnel interfaces or Xray.
set -u
umask 077
BASE=/media/developer/cvpn-webos
BOOTID=/proc/sys/kernel/random/boot_id
MOUNTS=/proc/mounts
SYS_NET=/sys/class/net
# The alternate paths are for offline unit tests, never for the production hook.
if [ "${CVPN_RECOVERY_TEST_MODE:-}" = 1 ]; then
    BASE=${CVPN_TEST_BASE:?}
    BOOTID=${CVPN_TEST_BOOTID:?}
    MOUNTS=${CVPN_TEST_MOUNTS:?}
    SYS_NET=${CVPN_TEST_SYS_NET:?}
fi
STATE="$BASE/state"
REC="$STATE/recovery"
LOCK="$STATE/lock"
PANEL_START="$BASE/webui/start-panel.sh"
MODE=${1:-inspect}

log() {
    [ -d "$REC" ] || return 0
    printf '%s %s\n' "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" "$*" >> "$REC/recovery.log" 2>/dev/null || true
}
status() {
    printf '%s\n' "$1" > "$REC/status.new.$$" || return 1
    mv -f "$REC/status.new.$$" "$REC/status" || return 1
    log "status=$1"
}
# Record the boot identifier only AFTER a conclusive result is safely stored.
# ERROR and RUNNING are intentionally retryable on the same Linux boot.
finish() {
    result=$1
    status "$result" || return 1
    case "$result" in ERROR|RUNNING) return 0 ;; esac
    [ -n "${current:-}" ] || return 1
    printf '%s\n' "$current" > "$REC/last-boot-id.new.$$" || return 1
    mv -f "$REC/last-boot-id.new.$$" "$REC/last-boot-id" || return 1
    log "boot complete=$result"
}
split_routes_present() {
    r1=$(ip -4 route show 0.0.0.0/1 2>/dev/null) || return 0
    r2=$(ip -4 route show 128.0.0.0/1 2>/dev/null) || return 0
    case "$r1 $r2" in *'dev cvpn0'*) return 0 ;; esac
    return 1
}
resolver_mounted() {
    awk '$2=="/var/lib/misc/resolv.conf"{x=1} END{exit(!x)}' "$MOUNTS" 2>/dev/null
}
xray_running() {
    # Fail closed if a saved PID exists and is in use, even when reused.
    if [ -s "$STATE/xray.pid" ]; then
        p=$(cat "$STATE/xray.pid" 2>/dev/null || true)
        case "$p" in
            ''|*[!0-9]*) ;;
            *)
                if kill -0 "$p" 2>/dev/null; then
                    exe=$(readlink "/proc/$p/exe" 2>/dev/null || true)
                    if [ "$exe" = "$BASE/bin/xray" ]; then return 0; fi
                fi
                ;;
        esac
    fi
    if command -v pidof >/dev/null 2>&1; then
        [ -n "$(pidof xray 2>/dev/null || true)" ] && return 0
    fi
    return 1
}
resources_present() {
    # Unknown route state is a reason to refuse cleanup, not to guess.
    [ -d "$SYS_NET/cvpn0" ] && return 0
    [ -r "$MOUNTS" ] || return 0
    ip -4 route show 0.0.0.0/1 >/dev/null 2>&1 || return 0
    ip -4 route show 128.0.0.0/1 >/dev/null 2>&1 || return 0
    split_routes_present && return 0
    resolver_mounted && return 0
    xray_running && return 0
    return 1
}
lock_is_busy() {
    [ -d "$LOCK" ] || return 1
    if [ ! -s "$LOCK/pid" ]; then return 1; fi
    holder=$(cat "$LOCK/pid" 2>/dev/null || true)
    case "$holder" in ''|*[!0-9]*) return 1 ;; esac
    kill -0 "$holder" 2>/dev/null || return 1
    # Ignore a recycled PID owned by an unrelated new-boot process.
    cmd=$(tr '\000' ' ' < "/proc/$holder/cmdline" 2>/dev/null || true)
    case "$cmd" in
        *cvpn.sh*|*cvpn-webos-panel*|*upgrade-1.8.0.sh*) return 0 ;;
    esac
    return 1
}
# Repair an informational INSTALL_PENDING_REBOOT status only when every
# read-only check proves there is no active CVPN session or network residue.
# A matching boot ID is not evidence that a pending status was processed.
# This function NEVER changes TUN, DNS, routes, Xray or engine state files.
reconcile_pending() {
    [ -s "$REC/status" ] && [ "$(cat "$REC/status" 2>/dev/null)" = INSTALL_PENDING_REBOOT ] || return 2
    [ -s "$REC/last-boot-id" ] && [ -r "$BOOTID" ] || { log 'pending remains: boot identifier unavailable'; return 1; }
    current=$(cat "$BOOTID" 2>/dev/null || true)
    saved=$(cat "$REC/last-boot-id" 2>/dev/null || true)
    [ -n "$current" ] && [ "$saved" = "$current" ] || { log 'pending remains: boot identifiers differ'; return 1; }
    [ ! -e "$STATE/active" ] || { log 'pending remains: active session marker'; return 1; }
    [ ! -e "$LOCK" ] || { log 'pending remains: engine lock present'; return 1; }
    [ ! -e "$REC/worker.lock" ] || { log 'pending remains: recovery worker lock present'; return 1; }
    if resources_present; then
        log 'pending remains: CVPN resources present or network checks inconclusive'
        return 1
    fi
    # Final guard against concurrent VPN start. No engine locks are created.
    [ ! -e "$STATE/active" ] && [ ! -e "$LOCK" ] || { log 'pending remains: concurrent session start'; return 1; }
    [ "$(cat "$REC/status" 2>/dev/null)" = INSTALL_PENDING_REBOOT ] || return 1
    status CLEAN || return 1
    log 'reconciled pending status after safe read-only network verification; no network changes'
    return 0
}
start_panel() {
    [ -r "$PANEL_START" ] || { log 'panel script absent'; return 0; }
    sh "$PANEL_START" >> "$REC/panel-start.log" 2>&1 || log 'panel start failed'
}

[ "$MODE" = inspect ] || [ "$MODE" = boot ] || [ "$MODE" = reconcile ] || { echo 'Usage: recovery-worker.sh inspect|boot|reconcile'; exit 2; }
if [ "$MODE" = inspect ]; then
    echo 'CVPN webOS recovery inspect (read only)'
    echo 'engine_state_active='$(test -s "$STATE/active" && echo yes || echo no)
    echo 'cvpn_resources_present='$(resources_present && echo yes_or_unknown || echo no)
    echo 'lock_present='$(test -d "$LOCK" && echo yes || echo no)
    if [ -r "$REC/status" ]; then echo 'last_status='$(cat "$REC/status"); fi
    if [ -r "$REC/last-boot-id" ] && [ -r "$BOOTID" ]; then
        if [ "$(cat "$REC/last-boot-id")" = "$(cat "$BOOTID")" ]; then
            echo 'boot_record=matches_current_boot'
        else
            echo 'boot_record=previous_boot'
        fi
    else
        echo 'boot_record=not_available'
    fi
    exit 0
fi

# Boot worker is never invoked by opening the TV launcher or polling the panel.
[ "$(id -u)" = 0 ] || [ "${CVPN_RECOVERY_TEST_MODE:-}" = 1 ] || exit 1
[ -d "$STATE" ] || exit 0
mkdir -p "$REC" || exit 0
chmod 700 "$REC" 2>/dev/null || true
if [ "$MODE" = reconcile ]; then
    if reconcile_pending; then
        echo 'RECOVERY_PENDING_RECONCILED=YES'
    else
        echo 'RECOVERY_PENDING_RECONCILED=NO'
    fi
    echo 'NETWORK_CHANGES=NEVER'
    exit 0
fi
# Small delay runs in detached worker, not in Homebrew run-parts startup.
[ "${CVPN_RECOVERY_TEST_MODE:-}" = 1 ] || sleep 5
current=$(cat "$BOOTID" 2>/dev/null || true)
[ -n "$current" ] || { status ERROR; exit 0; }
if [ -r "$REC/last-boot-id" ] && [ "$(cat "$REC/last-boot-id")" = "$current" ]; then
    # v1.8.0 installer deliberately pre-marked the installation boot without
    # writing a status. Never reinterpret that as a successful recovery.
    if [ ! -s "$REC/status" ]; then
        status INSTALL_PENDING_REBOOT || true
        log 'current boot was pre-marked by installer; no recovery cleanup performed'
    elif [ "$(cat "$REC/status" 2>/dev/null)" = INSTALL_PENDING_REBOOT ]; then
        # Legacy installs can keep a matching boot ID but a stale pending status.
        # Reconcile ONLY after all network/resource checks are conclusively clean.
        reconcile_pending || true
    else
        log 'already processed current boot'
    fi
    start_panel
    exit 0
fi
# Atomic one-worker-at-a-time lock, separate from the VPN engine lock.
if ! mkdir "$REC/worker.lock" 2>/dev/null; then
    finish SKIPPED_LOCK || true
    exit 0
fi
trap 'rmdir "$REC/worker.lock" 2>/dev/null || true' EXIT HUP INT TERM
status RUNNING || { log 'cannot write running status'; exit 0; }

if resources_present; then
    finish SKIPPED_ACTIVE
    start_panel
    exit 0
fi
if lock_is_busy; then
    finish SKIPPED_LOCK
    start_panel
    exit 0
fi
# Remove only clearly orphaned lock entries from a *previous* boot.
if [ -d "$LOCK" ]; then
    if [ -e "$LOCK/pid" ]; then rm -f "$LOCK/pid" || { status ERROR; start_panel; exit 0; }; fi
    rmdir "$LOCK" 2>/dev/null || { finish SKIPPED_LOCK; start_panel; exit 0; }
fi
# Acquire the SAME lock directory the engine uses for configuration changes.
if ! mkdir "$LOCK" 2>/dev/null; then
    finish SKIPPED_LOCK
    start_panel
    exit 0
fi
printf '%s\n' "$$" > "$LOCK/pid"
# Recheck because another service may have begun connecting during startup.
if resources_present; then
    rm -f "$LOCK/pid"; rmdir "$LOCK" 2>/dev/null || true
    finish SKIPPED_UNSAFE
    start_panel
    exit 0
fi
stale=0
for f in active mode phase deadline xray.pid active-server-ip dns.route server.route gateway; do
    if [ -e "$STATE/$f" ]; then stale=1; break; fi
done
if [ "$stale" = 1 ]; then
    archive="$REC/stale-$(date -u '+%Y%m%dT%H%M%SZ')-$$"
    mkdir -p "$archive" || { status ERROR; rm -f "$LOCK/pid"; rmdir "$LOCK" 2>/dev/null || true; start_panel; exit 0; }
    chmod 700 "$archive"
    for f in active mode phase deadline xray.pid active-server-ip dns.route server.route gateway; do
        [ -e "$STATE/$f" ] || continue
        mv "$STATE/$f" "$archive/$f" || log "could not archive state field $f"
    done
    if [ -e "$STATE/active" ]; then status ERROR; else finish CLEANED_STALE; fi
else
    finish CLEAN
fi
rm -f "$LOCK/pid"
rmdir "$LOCK" 2>/dev/null || true
# Starting the web CONTROL PANEL is not starting the VPN.
start_panel
exit 0
