# Add CVPN webOS — a free VPN client for rooted LG TVs

## Application

CVPN webOS is a **free-to-use, provider-neutral VPN client** for compatible rooted LG webOS TVs. It supports VLESS/Reality profiles, compatible HTTPS subscriptions and supported Xray configurations from any provider. A CVPN account or paid CVPN subscription is **not required**.

- App ID: `su.cvpn.webos`
- Official stable release: https://github.com/mikemandi/cvpn-webos/releases/tag/v2.2.5
- Project and source: https://github.com/mikemandi/cvpn-webos
- International documentation: https://github.com/mikemandi/cvpn-webos/blob/main/README.md
- Russian documentation: https://github.com/mikemandi/cvpn-webos/blob/main/README.ru.md

## Implementation and requirements

A single All-in-One IPK contains the TV launcher, ARMv7 backend, pinned Xray-core, VPN Engine and recovery scripts; no additional executable needs to be downloaded during TV setup. Initial preparation is carried out through the existing Homebrew root service. A rooted TV, compatible 32-bit ARM userspace, working /dev/net/tun and a trusted LAN are required.

The app uses a local web UI paired via a temporary TV PIN. VPN connection is explicit, not automatic. During updates, the installer checks for active VPN routes/TUN and refuses unsafe binary replacement; rollback safeguards and tests are included.

## Tested hardware and limitations

The owner tested the All-in-One first-run, local pairing, Netherlands VPN exit IP, YouTube playback, disconnect/reconnect and ordinary power off/on on **LG 55UK6450PLC (webOS 4.4.3-22)**. Additional earlier tests were performed on LG 55UQ76003LD after compatibility fixes. The revised 2.2.5 TV dialogs and LG remote error-focus behavior were also tested on hardware. These reports are **not** a guarantee for every LG model/firmware.

The VPN currently handles **IPv4 only** and refuses connection when an active IPv6 default route would bypass it. There is **no universal kill switch**, no automatic VPN connection after boot and no assurance of leak-free behavior for all TV apps. The control panel uses **HTTP on the private LAN, port 8765**; it must never be exposed publicly. Persistent local profiles and secrets are stored with restrictive permissions but **not encrypted at rest**.

## Licensing / price

CVPN webOS is **free of charge**. `pool: non-free` is used because the CVPN-authored source does not currently have a declared open-source license; this **does not** indicate that a subscription or payment is necessary. Bundled Xray-core retains its **MPL-2.0** license; QR components contain **MIT / Apache-2.0** notices. See the repository's licensing review.

## AI-assisted development disclosure

AI tools assisted with implementation, interface design, documentation and testing. The project owner reviewed and iterated the application and personally tested it on LG hardware, including operational VPN connectivity and version 2.2.5 remote interactions. The project has not undergone an independent security audit. I remain responsible for maintenance, security fixes and compatibility issues.

The package description includes genuine screenshots of the application in empty/demo state, with **no real subscription data, server credentials, PINs or TV IP addresses**.

Thank you for reviewing this submission. I am happy to respond to compatibility, licensing and security feedback before it is listed.
