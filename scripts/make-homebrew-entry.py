#!/usr/bin/env python3
"""Generate YAML for pull request to webosbrew/apps-repo; screenshots are deliberately omitted until verified."""
import argparse,pathlib,re
p=argparse.ArgumentParser()
p.add_argument('--repository',required=True,help='GitHub owner/repository')
p.add_argument('--output',default='dist/su.cvpn.webos.yml')
a=p.parse_args()
if not re.fullmatch(r'[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+',a.repository):p.error('expected owner/repo')
r=pathlib.Path(__file__).resolve().parent.parent
raw=f'https://raw.githubusercontent.com/{a.repository}/main'
s=f'''title: CVPN webOS
iconUri: {raw}/assets/icon160.png
manifestUrl: https://github.com/{a.repository}/releases/latest/download/su.cvpn.webos.manifest.json
category: utility
pool: non-free
description: |
  **CVPN Connect** is a VPN control experience for rooted LG webOS TVs.
  Import your own supported VLESS configurations and HTTPS subscriptions from a phone,
  pair securely with a temporary QR/PIN and control the VPN on the TV.

  **Important:** the CVPN webOS 2.1.0 IPK bundles its root backend, Xray-core,
  VPN Engine and recovery components. It configures them on first launch using
  Homebrew's root service. The app cannot run a system-wide VPN without root.
  IPv6 is fail-closed;
  do not expose the local HTTP control panel to untrusted networks.

  Installation and device requirements: https://github.com/{a.repository}#installation
# Add real, sanitized screenshots before submission; previews are not representative.
'''
output=r/a.output
output.parent.mkdir(parents=True,exist_ok=True)
output.write_text(s)
print('SUBMISSION_ENTRY_CREATED='+str(output))
