#!/bin/sh
# Repairs ONLY the CVPN Go panel on an already-installed rooted LG.
# Does not modify Xray, VPN Engine, TUN, routes, DNS, profiles, or subscriptions.
set -eu
umask 077
MODE=${1:-check}
BASE=/media/developer/cvpn-webos
WEB="$BASE/webui"
BIN="$WEB/bin/cvpn-webos-panel"
STAGE=/tmp/cvpn-first-import-fix
BACKUPS="$WEB/backups"
LAST="$BACKUPS/last-first-import-fix.txt"
refuse() { echo "REFUSED=$1" >&2; exit 1; }
check_idle() {
  [ "$(id -u)" = 0 ] || refuse root_required
  [ -x "$BIN" ] || refuse installed_panel_missing
  [ -s "$WEB/access.key" ] || refuse access_key_missing
  [ ! -e "$BASE/state/active" ] || refuse vpn_active
  [ ! -e "$BASE/state/lock" ] || refuse engine_locked
  if command -v ip >/dev/null 2>&1; then
    if ip link show dev cvpn0 >/dev/null 2>&1; then refuse tun_present; fi
    if ip -4 route show 0.0.0.0/1 | grep -q 'dev cvpn0'; then refuse vpn_route_present; fi
    if ip -4 route show 128.0.0.0/1 | grep -q 'dev cvpn0'; then refuse vpn_route_present; fi
  fi
}
panel_health() {
  key=$(cat "$WEB/access.key")
  attempt=0
  while [ "$attempt" -lt 4 ]; do
    resp=$(curl -fsS --noproxy '*' --connect-timeout 2 --max-time 5 -H "X-CVPN-Key: $key" 'http://127.0.0.1:8765/api/status' 2>/dev/null || true)
    case "$resp" in *'"version":"2.0.0"'*) return 0;; esac
    sleep 1
    attempt=$((attempt + 1))
  done
  return 1
}
install_fix() {
  [ -s "$STAGE/panel" ] || refuse staged_binary_missing
  [ -s "$STAGE/SHA256SUMS" ] || refuse staged_hash_missing
  (cd "$STAGE" && sha256sum -c SHA256SUMS) >/dev/null || refuse staged_hash_mismatch
  mkdir -p "$BACKUPS"
  chmod 700 "$BACKUPS"
  previous="$BACKUPS/pre-first-import-$(date -u +%Y%m%dT%H%M%SZ)-$$"
  mkdir -m 700 "$previous"
  cp -p "$BIN" "$previous/cvpn-webos-panel"
  [ -s "$previous/cvpn-webos-panel" ] || refuse backup_failed
  sh "$WEB/stop-panel.sh"
  sleep 2
  cp "$STAGE/panel" "$WEB/bin/cvpn-webos-panel.next.$$"
  chmod 700 "$WEB/bin/cvpn-webos-panel.next.$$"
  mv -f "$WEB/bin/cvpn-webos-panel.next.$$" "$BIN"
  if sh "$WEB/start-panel.sh" && panel_health; then
    printf '%s\n' "$previous" > "$LAST"
    chmod 600 "$LAST"
    echo PANEL_PATCHED=YES
    echo FIRST_IMPORT_TEMPLATE=EMBEDDED
    echo BACKUP_PATH="$previous"
    echo VPN_NETWORK_UNCHANGED=YES
    return 0
  fi
  echo 'Panel health check failed; restoring previous binary' >&2
  sh "$WEB/stop-panel.sh" || true
  cp -p "$previous/cvpn-webos-panel" "$BIN"
  if ! sh "$WEB/start-panel.sh" || ! panel_health; then
    refuse rollback_panel_unhealthy_manual_recovery_required
  fi
  refuse patch_failed_original_restored
}
rollback_fix() {
  [ -s "$LAST" ] || refuse backup_pointer_missing
  previous=$(cat "$LAST")
  case "$previous" in "$BACKUPS"/pre-first-import-*) :;; *) refuse unsafe_backup_path;; esac
  [ -s "$previous/cvpn-webos-panel" ] || refuse backup_missing
  sh "$WEB/stop-panel.sh"
  sleep 2
  cp -p "$previous/cvpn-webos-panel" "$BIN"
  sh "$WEB/start-panel.sh" && panel_health || refuse rollback_panel_unhealthy
  echo ROLLBACK_COMPLETE=YES
  echo VPN_NETWORK_UNCHANGED=YES
}
check_idle
case "$MODE" in
  check)
    echo FIRST_IMPORT_FIX_PREFLIGHT_OK=YES
    if [ -s "$BASE/config/client.json" ]; then echo CLIENT_CONFIG_PRESENT=YES; else echo CLIENT_CONFIG_PRESENT=NO; fi
    if [ -s "$BASE/config/server-ip.txt" ]; then echo SERVER_IP_PRESENT=YES; else echo SERVER_IP_PRESENT=NO; fi
    ;;
  install) install_fix;;
  rollback) rollback_fix;;
  *) echo 'Usage: sh repair-first-import.sh check|install|rollback' >&2; exit 2;;
esac
