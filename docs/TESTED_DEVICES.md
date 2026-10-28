# CVPN webOS — tested LG devices and compatibility

[English](TESTED_DEVICES.md) | [Русский](TESTED_DEVICES.ru.md)

Updated: **2026-10-28** · Current release: **2.2.6 Stable**

The project owner reports **CVPN webOS working on four LG TVs**. The new 2.2.6 release is a package-layout and documentation update; it has not been independently re-tested on every device. Older device logs have different levels of detail: a working installation reported by the owner is not necessarily a documented full clean-install and recovery test of the latest IPK on that model.

## Owner-tested devices

| TV model | webOS version | Evidence and results | Scope |
|---|---|---|---|
| **LG 55UK6450PLC** | **4.4.3-22**, firmware **05.50.70** | All-in-One **2.2.2** clean setup, QR panel, VPN connection with Netherlands exit IP, YouTube playback, disconnect/reconnect, normal off/on; earlier sessions also verified VPN routes and HTTPS | Detailed device logs for previous builds; owner confirms continued operation |
| **LG 65UK6710PLB** | **4.4.3-22** | Model and webOS version recorded during earlier root/Homebrew discussions; **owner now confirms CVPN works** | No retained step-by-step clean-install, transport or DNS evidence for this device; do not claim one |
| **LG 86NANO906NA** | **5.6.2** | Root/SSH, ARMv7 userspace, TUN, Xray/VLESS/Reality/**XHTTP** and changed exit IP were tested; YouTube/HTTPS connectivity and route restoration were checked during development | One later protected-DNS experiment rolled back safely; the owner now reports working CVPN. Not a documented clean install of 2.2.5 |
| **LG 55UQ76003LD** | **webOS 22 / 7.6.0-1901**, LG software **04.06.91** | Root/Homebrew, aarch64 kernel with **32-bit ARMv7-compatible userspace**, Xray 26.10.21, TUN, VLESS import and VPN connection after first-import/BusyBox/DNS fixes | Functional test on this firmware after iterative fixes, not a verified one-pass 2.2.5 clean install |

**Version 2.2.5**: the owner separately confirmed updated QR/preparation dialog layout, RU/EN interface and remote navigation between **Retry** and **Close** on real LG hardware. The model for that specific UI acceptance run was not independently recorded.

## Minimum version for Homebrew metadata

- **Oldest directly recorded successful webOS build:** **4.4.3-22** (on two UK-series models, with differing levels of evidence).
- **Catalog requirement:** `webosRelease: '>=4.4'` to avoid claiming all of webOS 4.0–4.3 were tested.
- This requirement is a **conservative catalog filter, not a guarantee that every webOS 4.4.x build works**. The owner-tested baseline is **4.4.3-22**. Later versions also require compatible root, Homebrew root service, 32-bit ARMv7-compatible userspace and TUN.
- No successful test is documented for webOS **4.0–4.3**, unrooted TVs or ordinary LG Content Store installations.

## App features and limits

**CVPN webOS** imports compatible **Xray/VLESS** profiles and subscriptions using **TCP, XHTTP and gRPC** transport options. gRPC/TLS and compatible Reality configurations are supported; these transports are not separate VPN protocols. The All-in-One IPK bundles Xray-core, the local backend and recovery scripts.

The VPN routes supported **IPv4** traffic. Active default IPv6 routes prevent activation; there is **no universal kill switch** and no automatic VPN connection after boot. CVPN Connect uses **unencrypted local HTTP** on TCP/8765 and must be kept on a trusted private LAN. A CVPN subscription is optional.

## Evidence still needed for wider support claims

- Repeat a fully clean **2.2.5** install, reboot, VPN tests and uninstall/reinstall on every named model if claiming model-specific production certification.
- Record release, firmware, root/Homebrew, 32-bit userspace, TUN and result for any new LG model.
- Verify sleep/power loss, network loss, rollback and uninstallation on more firmware versions.

Do not post real subscription links, VLESS UUIDs, access tokens, pairing PINs or local TV IPs in public issues.

[README](../README.md) · [Security](SECURITY.md) · [Release limitations](RELEASE_LIMITATIONS.md)
