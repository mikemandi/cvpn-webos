#!/bin/sh
# Control-panel only. Never starts VPN or changes system networking.
set -eu
umask 077
BASE=/media/developer/cvpn-webos
ROOT="$BASE/webui"
BIN="$ROOT/bin/cvpn-webos-panel"
PID="$ROOT/state/panel.pid"
LOG="$ROOT/state/panel.log"
mkdir -p "$ROOT/state"
if [ -s "$PID" ]; then
 p=$(cat "$PID" 2>/dev/null || true)
 case "$p" in ''|*[!0-9]*) :;; *)
   if [ -r "/proc/$p/exe" ] && [ "$(readlink "/proc/$p/exe")" = "$BIN" ]; then
     echo "PANEL_RUNNING=YES PID=$p"; exit 0
   fi;;
 esac
 rm -f "$PID"
fi
[ -x "$BIN" ] || { echo PANEL_BINARY_MISSING >&2; exit 1; }
[ -s "$ROOT/access.key" ] || { echo ACCESS_KEY_MISSING >&2; exit 1; }
nohup "$BIN" -base "$BASE" -listen ':8765' > "$LOG" 2>&1 </dev/null &
p=$!
printf '%s\n' "$p" > "$PID"
sleep 2
if ! kill -0 "$p" 2>/dev/null; then echo PANEL_START_FAILED >&2; exit 1; fi
echo "PANEL_RUNNING=YES PID=$p"
