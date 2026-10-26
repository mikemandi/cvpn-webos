#!/usr/bin/env bash
# Reproducible Debian-format webOS IPK containing the TV launcher, backend, Xray and recovery.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
command -v dpkg-deb >/dev/null || { echo 'dpkg-deb required (Linux/Ubuntu or GitHub Actions)' >&2; exit 1; }
python3 - "$ROOT/launcher/appinfo.json" <<'PY'
import json,sys,pathlib
p=pathlib.Path(sys.argv[1]); d=json.loads(p.read_text())
for k in ('id','title','version','type','main','icon'):
    if not d.get(k): raise SystemExit('Missing appinfo key: '+k)
assert d['id']=='su.cvpn.webos', 'Unexpected app ID'
assert d['type']=='web'
for k in ('main','icon','largeIcon'):
    if k in d and not (p.parent/d[k]).is_file(): raise SystemExit('Missing file: '+d[k])
PY
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
APP="$TMP/usr/palm/applications/su.cvpn.webos"
PKG="$TMP/usr/palm/packages/su.cvpn.webos"
mkdir -p "$APP" "$PKG" "$TMP/DEBIAN" "$ROOT/dist"
[[ -s "$ROOT/dist/payload/SHA256SUMS" && -s "$ROOT/dist/payload/bin/xray" && -s "$ROOT/dist/payload/XRAY-LICENSE.txt" ]] || { echo "Missing verified All-in-One payload; run prepare-bundle.sh" >&2; exit 1; }
cp -a "$ROOT/dist/payload" "$APP/payload"
for f in appinfo.json index.html icon.png largeicon.png icon-400.png qr-bundle.js QR-NOTICE.txt QR-THIRD-PARTY-LICENSE.txt; do
  cp "$ROOT/launcher/$f" "$APP/$f"
done
python3 - "$ROOT/launcher/appinfo.json" "$PKG/packageinfo.json" <<'PY'
import json,sys
p=json.load(open(sys.argv[1]))
with open(sys.argv[2],'w') as f: json.dump({'id':p['id'],'version':p['version'],'app':p['id']},f,indent=2)
PY
VERSION="$(python3 -c 'import json,sys;print(json.load(open(sys.argv[1]))["version"])' "$ROOT/launcher/appinfo.json")"
cat > "$TMP/DEBIAN/control" <<EOF
Package: su.cvpn.webos
Version: $VERSION
Section: misc
Priority: optional
Architecture: all
Maintainer: CVPN Project
Description: CVPN webOS All-in-One: launcher, root VPN backend and Xray-core.
webOS-Package-Format-Version: 2
webOS-Packager-Version: custom-debian2
EOF
chmod -R go-w "$TMP"
OUT="$ROOT/dist/su.cvpn.webos_${VERSION}_all.ipk"
dpkg-deb -Zgzip --build --root-owner-group "$TMP" "$OUT" >/dev/null
printf 'IPK_CREATED=%s\n' "$OUT"
sha256sum "$OUT"
dpkg-deb --info "$OUT" | head -n 20
