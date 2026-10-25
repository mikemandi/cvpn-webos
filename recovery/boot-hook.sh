#!/bin/sh
# CVPN webOS 1.8.0: Homebrew boot entrypoint. No VPN autostart.
# Return immediately: do not delay or block webOS startup.
WORKER=/media/developer/cvpn-webos/webui/recovery/recovery-worker.sh
[ -r "$WORKER" ] || exit 0
( nohup /bin/sh "$WORKER" boot >/dev/null 2>&1 </dev/null & ) >/dev/null 2>&1
exit 0
