# Draft Homebrew catalog PR — CVPN webOS v2.2.5

[English](HOMEBREW_PR_DRAFT.md) | [Русский](HOMEBREW_PR_DRAFT.ru.md)

**Suggested title:** Add CVPN webOS — free provider-neutral VPN client for rooted LG televisions

**Suggested PR description (review and adapt before submitting):**

> CVPN webOS is a **free-to-use VPN client for supported rooted LG webOS TVs** with ARMv7-compatible 32-bit userspace, working TUN and Homebrew Channel.
>
> Users can import supported VLESS/Reality configurations and compatible HTTPS subscriptions **from any provider**. No paid CVPN account or CVPN subscription is required.
>
> A single All-in-One IPK includes launcher, Xray-core, Go local panel and network recovery scripts. The root service prepares bundled components on first launch; it does not download extra executables on the television.
>
> **Network/security boundaries:** VPN routing is **IPv4-only**; an active IPv6 default route blocks VPN activation. There is **no universal kill switch and no VPN autostart**. The web control panel on **HTTP port 8765** is for **trusted LAN access only** and must never be publicly exposed.
>
> **Owner-reported hardware evidence:** LG 55UK6450PLC (webOS 4.4.3-22, firmware 05.50.70) was used for All-in-One baseline installation, QR, Netherlands VPN exit-IP, YouTube playback, reconnect and normal off/on. Version 2.2.5 TV layout and remote-control improvements were also confirmed on LG hardware. Other devices/firmware are not universally supported.
>
> **AI disclosure:** AI tools assisted implementation, UI, documentation and testing. The owner reviewed the code and tested the application on real devices; this is **not an independent security audit**. The submitter remains responsible for updates and bug reports.
>
> **Price/licensing:** free of charge. CVPN-authored source is publicly visible but has no declared open-source license, so catalog **pool: non-free** is requested (not a paid app). Bundled Xray-core uses **MPL-2.0**, and QR notices identify **MIT and Apache-2.0** components.
>
> Official release: https://github.com/mikemandi/cvpn-webos/releases/tag/v2.2.5
>
> Documentation: https://github.com/mikemandi/cvpn-webos

**Before sending:** check current [webOS Homebrew rules](https://github.com/webosbrew/apps-repo), approve the final text with the project owner, inspect packages/su.cvpn.webos.yml, verify official release file hashes, and attach sanitized real-UI screenshots already in assets/screenshots. Do not claim catalog inclusion before upstream merges the PR.

[Submission guide](HOMEBREW_SUBMISSION.md) · [Licensing review](../LICENSE-PUBLISHING-REVIEW.md)
