#!/usr/bin/env python3
"""Generate root bootstrap zip without personal secrets or unverified Xray-core."""
import pathlib,zipfile,hashlib
root=pathlib.Path(__file__).resolve().parent.parent
files=["dist/cvpn-webos-panel-armv7","engine/cvpn.sh","recovery/recovery-worker.sh", "recovery/boot-hook.sh", "scripts/start-panel.sh", "scripts/stop-panel.sh", "scripts/install-root.sh", "scripts/setup.ps1","scripts/setup.sh","docs/RELEASE_LIMITATIONS.md"]
for f in files:
 if not (root/f).is_file(): raise SystemExit('Missing '+f)
dest=root/'dist/cvpn-webos-root-2.0.0.zip'
with zipfile.ZipFile(dest,'w',compression=zipfile.ZIP_DEFLATED,compresslevel=9) as z:
 for f in files: z.write(root/f,f)
 print('ROOT_BUNDLE_CREATED='+str(dest))
