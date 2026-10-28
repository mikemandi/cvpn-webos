# Draft Homebrew catalog PR — CVPN webOS v2.2.5

[English](HOMEBREW_PR_DRAFT.md) | [Русский](HOMEBREW_PR_DRAFT.ru.md)

**Suggested title:** Add CVPN webOS — Xray/VLESS client with TCP, XHTTP and gRPC for rooted LG TVs

**Suggested PR description (review and adapt before submitting):**

> CVPN webOS is a **free Xray/VLESS VPN client for rooted LG webOS TVs**, featuring **TCP, XHTTP and gRPC transports**. It requires ARMv7-compatible 32-bit userspace, working TUN and Homebrew Channel.
>
> Users can import supported `vless://` links, VLESS/Reality configurations, compatible HTTPS subscriptions or Xray JSON **from any provider**. gRPC also supports TLS. TCP, XHTTP and gRPC are transports (not independent VPN protocols). No CVPN account or subscription is required.
>
> A single All-in-One IPK includes launcher, Xray-core, Go local panel and network recovery scripts. The root service prepares bundled components on first launch; it does not download extra executables on the television.
>
> **Network/security boundaries:** VPN routing is **IPv4-only**; an active IPv6 default route blocks VPN activation. There is **no universal kill switch and no VPN autostart**. The web control panel on **HTTP port 8765** is for **trusted LAN access only** and must never be publicly exposed.
>
> **Owner-tested LG models:** **55UK6450PLC (webOS 4.4.3-22)**, **65UK6710PLB (4.4.3-22)**, **86NANO906NA (5.6.2)** and **55UQ76003LD (webOS 22 / 7.6.0-1901)**. Logs document clean All-in-One/YouTube/reconnect on 55UK6450PLC, VLESS/Reality/XHTTP and TUN on 86NANO906NA, and functional VPN after initial compatibility fixes on 55UQ76003LD. The 65UK6710PLB is owner-confirmed without retained detailed logs. Oldest documented build: **4.4.3-22**; catalog requires **4.4+**. No guarantee for all firmware.
>
> **AI disclosure:** AI tools assisted implementation, UI, documentation and testing. The owner reviewed the code and tested the application on real devices; this is **not an independent security audit**. The submitter remains responsible for updates and bug reports.
>
> **Price/licensing:** free of charge. CVPN-authored code is offered under a **custom proprietary license with free end-user use**, so catalog **pool: non-free** is requested (not a paid app). Bundled Xray-core uses **MPL-2.0**, and QR notices identify **MIT and Apache-2.0** components. The proprietary terms do not restrict those third-party rights. License: https://github.com/mikemandi/cvpn-webos/blob/main/LICENSE
>
> Official release: https://github.com/mikemandi/cvpn-webos/releases/tag/v2.2.5
>
> Documentation: https://github.com/mikemandi/cvpn-webos

**Before sending:** check current [webOS Homebrew rules](https://github.com/webosbrew/apps-repo), approve the final text with the project owner, inspect packages/su.cvpn.webos.yml, verify official release file hashes, and attach sanitized real-UI screenshots already in assets/screenshots. Do not claim catalog inclusion before upstream merges the PR.

[Submission guide](HOMEBREW_SUBMISSION.md) · [Licensing review](../docs/legal/LICENSING_REVIEW.md)
