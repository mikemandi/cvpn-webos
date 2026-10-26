#!/usr/bin/env bash
# Run from the repository root on Linux or macOS, rooted TV with ARMv7 userspace only.
set -euo pipefail
if [[ $# -lt 2 || $# -gt 3 ]]; then echo 'Usage: bash scripts/setup.sh TV_IP check|install|doctor [LOCAL_XRAY_ARMV7_BINARY]'; exit 2; fi
TV="$1"; MODE="$2"; XRAY="${3:-}"
[[ "$TV" =~ ^[a-zA-Z0-9][a-zA-Z0-9.-]*$ ]] || { echo 'Invalid TV name/IP'; exit 2; }
[[ "$MODE" =~ ^(check|install|doctor)$ ]] || exit 2
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
command -v ssh >/dev/null && command -v scp >/dev/null && command -v sha256sum >/dev/null
TVHOST="root@$TV"
scp -o StrictHostKeyChecking=ask "$ROOT/scripts/install-root.sh" "$TVHOST:/tmp/cvpn-root-setup.sh"
if [[ "$MODE" == install ]]; then
  ssh "$TVHOST" 'sh /tmp/cvpn-root-setup.sh check'
  ssh "$TVHOST" 'mkdir -p /tmp/cvpn-root-setup && chmod 700 /tmp/cvpn-root-setup'
  scp "$ROOT/dist/cvpn-webos-panel-armv7" "$TVHOST:/tmp/cvpn-root-setup/panel"
  scp "$ROOT/engine/cvpn.sh" "$TVHOST:/tmp/cvpn-root-setup/engine"
  scp "$ROOT/recovery/recovery-worker.sh" "$TVHOST:/tmp/cvpn-root-setup/worker"
  scp "$ROOT/recovery/boot-hook.sh" "$TVHOST:/tmp/cvpn-root-setup/boot-hook"
  scp "$ROOT/scripts/start-panel.sh" "$TVHOST:/tmp/cvpn-root-setup/start-panel"
  scp "$ROOT/scripts/stop-panel.sh" "$TVHOST:/tmp/cvpn-root-setup/stop-panel"
  if [[ -n "$XRAY" ]]; then [[ -f "$XRAY" ]] || exit 1; scp "$XRAY" "$TVHOST:/tmp/cvpn-root-setup/xray"; fi
  (cd "$ROOT" && sha256sum     dist/cvpn-webos-panel-armv7 engine/cvpn.sh recovery/recovery-worker.sh     recovery/boot-hook.sh scripts/start-panel.sh scripts/stop-panel.sh     | sed 's|dist/cvpn-webos-panel-armv7|panel|;s|engine/cvpn.sh|engine|;s|recovery/recovery-worker.sh|worker|;s|recovery/boot-hook.sh|boot-hook|;s|scripts/start-panel.sh|start-panel|;s|scripts/stop-panel.sh|stop-panel|') > "$ROOT/dist/payload.sha256"
  if [[ -n "$XRAY" ]]; then sha256sum "$XRAY" | awk '{print $1 "  xray"}' >> "$ROOT/dist/payload.sha256"; fi
  scp "$ROOT/dist/payload.sha256" "$TVHOST:/tmp/cvpn-root-setup/payload.sha256"
fi
ssh "$TVHOST" "sh /tmp/cvpn-root-setup.sh $MODE"
