#!/bin/sh
# Stops CVPN control panel ONLY. No VPN cleanup, routes, or DNS modifications.
set -eu
BASE=/media/developer/cvpn-webos
PID="$BASE/webui/state/panel.pid"
BIN="$BASE/webui/bin/cvpn-webos-panel"
if [ -s "$PID" ]; then
 p=$(cat "$PID" 2>/dev/null || true)
 case "$p" in ''|*[!0-9]*) :;; *)
   if [ -r "/proc/$p/exe" ] && [ "$(readlink "/proc/$p/exe")" = "$BIN" ]; then kill "$p"; fi;;
 esac
 rm -f "$PID"
fi
echo PANEL_STOP_REQUESTED=YES
