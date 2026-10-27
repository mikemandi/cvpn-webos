#!/usr/bin/env python3
"""Generate YAML for pull request to webosbrew/apps-repo; screenshots are deliberately omitted until verified."""
import argparse,json,pathlib,re
p=argparse.ArgumentParser()
p.add_argument('--repository',required=True,help='GitHub owner/repository')
p.add_argument('--output',default='dist/su.cvpn.webos.yml')
a=p.parse_args()
if not re.fullmatch(r'[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+',a.repository):p.error('expected owner/repo')
r=pathlib.Path(__file__).resolve().parent.parent
version=json.loads((r/'launcher/appinfo.json').read_text())['version']
raw=f'https://raw.githubusercontent.com/{a.repository}/main'
s=f'''title: CVPN webOS
iconUri: {raw}/assets/icon160.png
manifestUrl: https://github.com/{a.repository}/releases/latest/download/su.cvpn.webos.manifest.json
category: utility
# Homebrew pool classification: no OSI-approved license has been declared
# for CVPN-authored source. This does NOT mean the application costs money.
pool: non-free
description: |
  **CVPN webOS** is a **free-to-use VPN client** for supported rooted LG TVs.
  Import your own VLESS/Reality configuration or compatible HTTPS subscription
  from **any provider**. No CVPN subscription, payment or CVPN account is needed
  to install or use its ordinary VPN functionality.

  **CVPN (cvpn.su)** is an **optional recommended provider** with additional
  subscription conveniences. Using a different compatible provider is supported.

  **Compatibility:** root and Homebrew Channel are prerequisites. The {version}
  All-in-One IPK includes the ARMv7 backend, Xray-core, VPN Engine and recovery
  scripts; installation prepares them locally through the Homebrew root service.
  A supported 32-bit ARM user space and TUN are required. Not for the LG Content Store.

  **Safety:** an active IPv6 default route prevents connection; there is no
  universal kill switch or automatic VPN start at boot. The local management
  interface uses HTTP: access it only on a trusted private network; never
  expose the TV's port 8765 to the public internet.

  First-run component setup was tested on LG 55UK6450PLC / webOS 4.4.3-22.
  Support varies with firmware and Homebrew version.
  Installation and requirements: https://github.com/{a.repository}#installation-for-supported-hardware
# Add screenshots captured from a real TV after sanitizing private values.
'''
output=r/a.output
output.parent.mkdir(parents=True,exist_ok=True)
output.write_text(s)
print('SUBMISSION_ENTRY_CREATED='+str(output))
