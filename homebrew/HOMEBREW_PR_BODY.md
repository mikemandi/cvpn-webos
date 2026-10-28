# Add CVPN webOS — Xray/VLESS VPN for rooted LG TVs (TCP, XHTTP, gRPC)

## Application

**CVPN webOS is an Xray/VLESS VPN client for rooted LG webOS TVs**, with support for **TCP, XHTTP and gRPC transports**. Import compatible `vless://` links, HTTPS subscriptions or Xray JSON from any provider and manage the TV's VPN from your phone. VLESS/Reality is supported; gRPC can also use TLS. These are transport options rather than different VPN protocols. The app is free to use, and a CVPN subscription is **not required**.

- App ID: `su.cvpn.webos`
- Official stable release: https://github.com/mikemandi/cvpn-webos/releases/tag/v2.2.5
- Project and source: https://github.com/mikemandi/cvpn-webos
- International documentation: https://github.com/mikemandi/cvpn-webos/blob/main/README.md
- Russian documentation: https://github.com/mikemandi/cvpn-webos/blob/main/README.ru.md

## Implementation and requirements

A single All-in-One IPK contains the TV launcher, ARMv7 backend, pinned Xray-core, VPN Engine and recovery scripts; no additional executable needs to be downloaded during TV setup. Initial preparation is carried out through the existing Homebrew root service. A rooted TV, compatible 32-bit ARM userspace, working /dev/net/tun and a trusted LAN are required.

The app uses a local web UI paired via a temporary TV PIN. VPN connection is explicit, not automatic. During updates, the installer checks for active VPN routes/TUN and refuses unsafe binary replacement; rollback safeguards and tests are included.

## Tested hardware and limitations

The owner confirms CVPN webOS works on **four LG models**: **55UK6450PLC (webOS 4.4.3-22)**, **65UK6710PLB (4.4.3-22)**, **86NANO906NA (5.6.2)** and **55UQ76003LD (webOS 22 / 7.6.0-1901)**. Earlier logs document All-in-One setup, VPN exit-IP, YouTube and reconnect tests on 55UK6450PLC; VLESS/Reality/XHTTP, TUN and routing tests on 86NANO906NA; and functional VPN after first-import/BusyBox/DNS fixes on 55UQ76003LD. For 65UK6710PLB, operation is owner-reported without retained detailed network logs. Version 2.2.5 dialogs and remote focus were also accepted on LG hardware. **The oldest verified webOS build is 4.4.3-22; the catalog requests 4.4+, not untested 4.0–4.3.** These results are not a universal LG compatibility guarantee.

The VPN currently handles **IPv4 only** and refuses connection when an active IPv6 default route would bypass it. There is **no universal kill switch**, no automatic VPN connection after boot and no assurance of leak-free behavior for all TV apps. The control panel uses **HTTP on the private LAN, port 8765**; it must never be exposed publicly. Persistent local profiles and secrets are stored with restrictive permissions but **not encrypted at rest**.

## Licensing / price

CVPN webOS is **free of charge**. `pool: non-free` is used because CVPN-authored portions use a **custom proprietary license with free end-user use**; this **does not** indicate that a subscription or payment is necessary. Bundled Xray-core retains its **MPL-2.0** license; QR components contain **MIT / Apache-2.0** notices. See [LICENSE](https://github.com/mikemandi/cvpn-webos/blob/main/LICENSE), [Brand policy](https://github.com/mikemandi/cvpn-webos/blob/main/BRAND_POLICY.md) and [third-party notices](https://github.com/mikemandi/cvpn-webos/blob/main/THIRD_PARTY_NOTICES.md). The v2.2.5 binary predates the new project LICENSE file; the next package should include the updated notices.

## AI-assisted development disclosure

AI tools assisted with implementation, interface design, documentation and testing. The project owner reviewed and iterated the application and personally tested it on LG hardware, including operational VPN connectivity and version 2.2.5 remote interactions. The project has not undergone an independent security audit. I remain responsible for maintenance, security fixes and compatibility issues.

The package description includes genuine screenshots of the application in empty/demo state, with **no real subscription data, server credentials, PINs or TV IP addresses**.

Thank you for reviewing this submission. I am happy to respond to compatibility, licensing and security feedback before it is listed.
