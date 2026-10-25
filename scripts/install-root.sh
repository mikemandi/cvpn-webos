#!/bin/sh
# Fresh ROOT-only bootstrap on rooted LG ARMv7, never an updater.
# Expected /tmp/cvpn-root-setup/: panel, engine, worker, boot-hook, start-panel,
# stop-panel, payload.sha256; xray optional only if preinstalled in BASE/bin/xray.
set -eu
umask 077
MODE=${1:-check}
BASE=/media/developer/cvpn-webos
WEB="$BASE/webui"
STAGE=/tmp/cvpn-root-setup
HOOK=/var/lib/webosbrew/init.d/90-cvpn-webos
refuse() { echo "REFUSED=$1" >&2; exit 1; }
idle() {
  [ ! -e "$BASE/state/active" ] || refuse vpn_active
  [ ! -e "$BASE/state/lock" ] || refuse vpn_engine_locked
  if command -v ip >/dev/null 2>&1; then
    if ip link show dev cvpn0 >/dev/null 2>&1; then refuse cvpn_tun_exists; fi
    ra=$(ip -4 route show 0.0.0.0/1) || refuse route_unreadable
    rb=$(ip -4 route show 128.0.0.0/1) || refuse route_unreadable
    case "$ra $rb" in *'dev cvpn0'*) refuse cvpn_routes_exist;; esac
  else
    refuse ip_command_missing
  fi
  [ -r /proc/mounts ] || refuse proc_mounts_unreadable
  if awk '$2=="/var/lib/misc/resolv.conf"{m=1} END{exit(m?0:1)}' /proc/mounts; then
    refuse resolver_bind_active
  fi
}
preflight() {
  [ "$(id -u)" = 0 ] || refuse root_required
  [ "$(uname -m)" = armv7l ] || refuse armv7l_required
  for c in sh ip curl awk od tr sha256sum readlink; do
    command -v "$c" >/dev/null 2>&1 || refuse "missing_tool:$c"
  done
  [ -f /var/lib/webosbrew/startup.sh ] || refuse homebrew_startup_missing
  idle
  # Strictly refuse overwriting any CVPN pre-existing files or keys.
  for f in "$WEB/bin/cvpn-webos-panel" "$WEB/access.key" "$WEB/profiles.json" "$BASE/bin/cvpn.sh"; do
    [ ! -e "$f" ] || refuse "existing_install_use_legacy_upgrade:$f"
  done
  if [ -e "$HOOK" ] || [ -L "$HOOK" ]; then refuse preexisting_boot_hook; fi
  if [ -x "$BASE/bin/xray" ]; then
    echo XRAY_SOURCE=ALREADY_ON_TV
  else
    echo XRAY_SOURCE=LOCAL_PROVIDED_BINARY_REQUIRED
  fi
  echo FRESH_PREFLIGHT_OK=YES
  echo NETWORK_MUTATIONS=NONE
}
check_stage() {
  [ -d "$STAGE" ] || refuse staging_missing
  for f in panel engine worker boot-hook start-panel stop-panel payload.sha256; do
    [ -s "$STAGE/$f" ] || refuse "staging_missing:$f"
  done
  (cd "$STAGE" && sha256sum -c payload.sha256) >/dev/null || refuse staged_hash_mismatch
  for f in engine worker boot-hook start-panel stop-panel; do sh -n "$STAGE/$f" || refuse "invalid_shell:$f"; done
  [ -x "$BASE/bin/xray" ] || [ -s "$STAGE/xray" ] || refuse xray_binary_required
}
install() {
  preflight
  check_stage
  mkdir -p "$BASE/bin" "$BASE/config" "$BASE/state" "$WEB/bin" "$WEB/secrets" "$WEB/state" "$WEB/recovery" /var/lib/webosbrew/init.d
  chmod 700 "$BASE/config" "$BASE/state" "$WEB/secrets" "$WEB/state"
  if [ ! -x "$BASE/bin/xray" ]; then
    cp "$STAGE/xray" "$BASE/bin/xray.tmp.$$"
    chmod 700 "$BASE/bin/xray.tmp.$$"
    mv -f "$BASE/bin/xray.tmp.$$" "$BASE/bin/xray"
  fi
  # Check actual ABI and executable on TV before proceeding.
  "$BASE/bin/xray" version >/dev/null 2>&1 || refuse incompatible_xray_binary
  cp "$STAGE/panel" "$WEB/bin/cvpn-webos-panel"
  cp "$STAGE/engine" "$BASE/bin/cvpn.sh"
  cp "$STAGE/worker" "$WEB/recovery/recovery-worker.sh"
  cp "$STAGE/boot-hook" "$WEB/recovery/boot-hook.sh"
  cp "$STAGE/start-panel" "$WEB/start-panel.sh"
  cp "$STAGE/stop-panel" "$WEB/stop-panel.sh"
  chmod 700 "$WEB/bin/cvpn-webos-panel" "$BASE/bin/cvpn.sh" "$WEB/recovery/recovery-worker.sh" "$WEB/recovery/boot-hook.sh" "$WEB/start-panel.sh" "$WEB/stop-panel.sh"
  if [ ! -e "$WEB/access.key" ]; then
    # /dev/urandom is required; never create a predictable or reused key.
    key=$(od -An -v -tx1 -N32 /dev/urandom | tr -d ' \n')
    [ "${#key}" = 64 ] || refuse random_key_generation_failed
    printf '%s\n' "$key" > "$WEB/access.key"
    chmod 600 "$WEB/access.key"
  fi
  # Note: no VPN configured and no VPN auto-connect. Empty profile state is created by backend.
  ln -s "$WEB/recovery/boot-hook.sh" "$HOOK" || refuse create_boot_hook_failed
  sh "$WEB/start-panel.sh" || refuse panel_failed_to_start
  key=$(cat "$WEB/access.key")
  resp=$(curl -fsS --noproxy '*' --connect-timeout 2 --max-time 5 -H "X-CVPN-Key: $key" http://127.0.0.1:8765/api/status 2>/dev/null || true)
  case "$resp" in *'"version":"2.0.0"'*) :;; *) refuse panel_api_unavailable;; esac
  echo FRESH_INSTALLED=YES
  echo PANEL_VERSION=2.0.0
  echo VPN_AUTOSTART=NO
  echo USER_PROFILES=DEVICE_LOCAL
  echo 'NOTE=Install launcher IPK separately via Homebrew Channel or ares-install'
}
doctor() {
  [ -x "$WEB/bin/cvpn-webos-panel" ] && echo PANEL_BINARY=YES || echo PANEL_BINARY=NO
  [ -x "$BASE/bin/xray" ] && echo XRAY_BINARY=YES || echo XRAY_BINARY=NO
  [ -s "$WEB/access.key" ] && echo UNIQUE_KEY=YES || echo UNIQUE_KEY=NO
  [ -L "$HOOK" ] && echo BOOT_HOOK=YES || echo BOOT_HOOK=NO
  [ -s "$BASE/state/active" ] && echo VPN_ACTIVE=YES || echo VPN_ACTIVE=NO
  [ -f "$WEB/profiles.json" ] && echo PROFILE_STORE=YES || echo PROFILE_STORE=NO
}
case "$MODE" in check) preflight;; install) install;; doctor) doctor;; *) echo 'Usage: check|install|doctor' >&2; exit 2;; esac
