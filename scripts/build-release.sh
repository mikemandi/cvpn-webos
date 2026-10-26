#!/usr/bin/env bash
# Builds Go ARMv7 backend and TV launcher IPK. Bundle pinned Xray-core inside IPK.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
mkdir -p "$ROOT/dist"
cd "$ROOT/server"
go test ./...
go vet ./...
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -ldflags='-s -w' -o "$ROOT/dist/cvpn-webos-panel-armv7" .
cd "$ROOT"
bash scripts/prepare-bundle.sh
bash scripts/build-ipk.sh
printf 'BACKEND_SHA256='; sha256sum dist/cvpn-webos-panel-armv7
printf 'BACKEND_BUILD_GO='; go version
python3 scripts/package-root.py
