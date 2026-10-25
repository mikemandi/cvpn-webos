#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
bash -n scripts/build-ipk.sh scripts/build-release.sh scripts/setup.sh
sh -n engine/cvpn.sh recovery/recovery-worker.sh recovery/boot-hook.sh scripts/start-panel.sh scripts/stop-panel.sh scripts/install-root.sh
(cd server && go test -race ./... && go vet ./... && go test -cover ./...)
python3 -m unittest discover -s tests -v
printf 'PUBLIC_OFFLINE_CHECKS_OK=YES\n'
