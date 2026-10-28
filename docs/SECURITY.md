# CVPN webOS — Security

[English](SECURITY.md) | [Русский](SECURITY.ru.md)

Last reviewed: **2026-10-28** · Published app: **2.2.5 Stable** · [Known limitations](RELEASE_LIMITATIONS.md)

> CVPN webOS is an experimental root-dependent VPN client for supported LG TVs, **not a security-audited firewall or a general-purpose privacy guarantee**. Read the boundaries below before relying on it.

## Threat model and prerequisites

- **Trusted LG with root and Homebrew:** the owner enables root; CVPN neither roots the TV nor grants itself privileges. The launcher calls a Homebrew privileged Luna service for **package-owned, constant** status/install commands. Homebrew ACL and firmware compatibility vary by device.
- **Trusted home LAN:** the local CVPN Connect control panel listens on **HTTP TCP/8765**. The connection between phone and TV is **not TLS-encrypted**. Other devices/users with access to an untrusted LAN may observe or attack traffic; pairing/session checks do not replace transport encryption. Never forward 8765 to the internet or expose it through a public proxy.
- **Trusted VPN provider:** importing a provider configuration permits its servers to transport traffic covered by the VPN route. It does not establish provider trust. CVPN is an optional provider; no account or purchase is required to use another compatible service.
- **No independent audit:** automated tests and owner hardware checks reduce risk but do not constitute a professional security audit or promise perfect isolation.

## Authentication, credentials and private files

The local panel has device pairing with a **short-lived, single-use TV PIN**, browser sessions and CSRF checks; panel access also has LAN/host/request validation. Do not assume these measures can secure traffic on a hostile LAN. Revoke old browser devices in **Settings → Devices and access**.

Access keys, subscriptions and profile data are stored **locally on the TV**, with restrictive filesystem permissions (typically 0600/0700). They are **not encrypted at rest**: privileged TV users, root malware or offline access to storage can read them. Keep full subscription URLs, VLESS UUIDs, token-bearing QR codes, session cookies, local IP addresses, access keys and diagnostic logs out of public issues/screenshots.

Release screenshots use only a reserved **192.0.2.0/24 TEST-NET** demo IP and a **nonfunctional example PIN**, never actual personal pairing details.

## VPN, DNS and failure behavior

| Topic | Security boundary |
|---|---|
| IPv4 | Engine 0.1.4 protects supported **IPv4 routes only** |
| IPv6 | VPN activation is refused when an active IPv6 default route would bypass the tunnel; no IPv6 VPN tunneling |
| Kill switch | **No universal kill switch**. VPN failure/disconnect may restore ordinary internet access |
| Startup | VPN **does not automatically connect** when the TV starts |
| DNS | Route/DNS recovery exists, but **encrypted DNS and zero leaks are not guaranteed** |
| HTTPS readiness | Startup HTTPS probes validate certificates/hostnames; this does **not** prove all apps/protocols use the VPN |
| TCP/UDP/apps | No guarantee for every app, UDP/QUIC transport, streaming provider or standby behavior |

The watchdog and recovery logic attempt to restore the previous DNS and network routes. An interrupted session or failed update still requires user verification; do not assume a power loss or firmware update is always recoverable.

## Root installation and updates

The single **All-in-One IPK** contains the launcher, ARMv7 backend, VPN Engine and a pinned Xray binary. Component checksums detect corruption of bundled payloads; **SHA-256 is not a code-signing certificate or independent proof of the publisher's identity**. Download only from the project's official GitHub Releases and verify the published release checksums.

The installer refuses to replace live Engine/Xray files when a VPN session, its network routes, TUN, DNS mount or lock is active. It backs up installed binaries and attempts rollback on partial failure; user profiles are not deliberately overwritten. A timed-out visual preparation screen does **not** necessarily mean root installation stopped; read-only polls should not launch overlapping installers.

Uninstalling the IPK can leave app-owned local backend or profile files. Do not delete backend paths automatically unless you explicitly intend to erase locally stored configurations and paired devices.

## Vulnerability reporting

For a potentially exploitable security issue, **do not post secrets, live credentials or attack details in a public GitHub issue**. First contact the maintainer privately via [CVPN support](https://t.me/insightdev0); share a minimal, redacted description, affected app version and reproducible steps. Use a secure agreed channel for any sensitive evidence. General, nonsensitive problems can be reported via [GitHub Issues](https://github.com/mikemandi/cvpn-webos/issues).

No guaranteed response time, bounty or formal coordinated-disclosure program is currently offered.

## Before broad device/catalog distribution

- Repeat first-run, standby, power-loss, network-loss, rollback and uninstall tests on additional LG firmware/Homebrew combinations.
- Review LAN authentication, session/CSRF behavior, source-code ownership, QR notices and the Xray **MPL-2.0** source-availability obligations.
- Keep private data out of the IPK, CI logs, README galleries and GitHub Releases.
- Disclose AI-assisted development and perform meaningful owner review before submitting to the [Homebrew application catalog](https://github.com/webosbrew/apps-repo).

[Project license](../LICENSE) · [Third-party notices](../THIRD_PARTY_NOTICES.md) · [Licensing review](../docs/legal/LICENSING_REVIEW.md) · [Device evidence](TESTED_DEVICES.md) · [English README](../README.md) · [Русский README](../README.ru.md)
