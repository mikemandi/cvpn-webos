# CVPN webOS — Licensing and publishing review

[English](LICENSE-PUBLISHING-REVIEW.md) | [Русский](LICENSE-PUBLISHING-REVIEW.ru.md)

Reviewed: **2026-10-28** · Published application: **v2.2.5**

## Current licensing status

**The CVPN-authored source and artwork have no declared project-wide license.** There is no root LICENSE file granting open-source reuse rights. Public GitHub source visibility and a free-to-use IPK **do not authorize redistribution, modification or reuse of proprietary CVPN code/artwork by default**.

The project owner must explicitly decide and document the intended licensing model. This review **does not grant new rights, select a license, or change the terms of third-party components**.

**Price and license are separate:**
- CVPN webOS is **free of charge** to install and use.
- Users can import **third-party VPN configurations**; a CVPN account/subscription is optional.
- The webOS Homebrew catalog uses **pool: non-free** for software not published under a declared open-source license, **even if it costs nothing**.
- Choose **pool: main** only **after** the owner adopts an appropriate open-source license and verifies licensing rights.

## Third-party components in the All-in-One IPK

| Component | Rights / included notices | Action |
|---|---|---|
| **Xray-core 26.10.21 ARMv7** | XTLS/Xray-core uses **Mozilla Public License 2.0 (MPL-2.0)**. The build fetches the pinned upstream binary and license, and packs payload/XRAY-LICENSE.txt. | Preserve license and copyright notices, keep an accessible link to the corresponding tagged source, and verify MPL-2.0 source availability obligations when distributing its binary. |
| **QR rendering library** | launcher/QR-NOTICE.txt identifies Kazuhiko Arase's original **MIT-licensed** QRCode for JavaScript and an adapted **Apache-2.0-licensed** qrcode-terminal distribution. launcher/QR-THIRD-PARTY-LICENSE.txt includes the relevant license text/notices. | Retain both applicable attributions/notices. Do not incorrectly describe the entire QR implementation as Apache-only. |
| **CVPN app/backend/icons** | Ownership and any upstream-derived assets or contributions require review. | Confirm every source/artwork contribution is licensed or owned and record the owner's final licensing decision. |

Upstream Xray project: https://github.com/XTLS/Xray-core/tree/v26.10.21

QR licensing source: [QR notice](launcher/QR-NOTICE.txt) and [bundled third-party license](launcher/QR-THIRD-PARTY-LICENSE.txt).

## What a final LICENSE would require

1. **Owner decision:** license CVPN-authored code under an identified OSI-approved license (if open source is intended), or publish explicit proprietary distribution/usage terms. Do not select a license automatically.
2. **Scope:** make clear which source code, UI artwork, logos, screenshots and names are covered; third-party code retains its own terms.
3. **Rights:** review contributors, generated/third-party assets, imported snippets and external dependencies before licensing them onward.
4. **Redistribution:** ensure any final policy agrees with public IPK releases, the Homebrew catalog, derivative works and third-party license notices.
5. **MPL-2.0:** provide a reliable way for recipients to obtain corresponding Xray-core source for the exact binary version distributed (and retain its license and required notices).
6. **Review before Homebrew submission:** check current [Homebrew rules](https://github.com/webosbrew/apps-repo), pricing text, catalog pool, owner review and screenshot hygiene.

**Do not simply copy Xray's MPL-2.0 into a root LICENSE file**: that would falsely imply the same license has been chosen for all CVPN-authored files.

## Public release privacy checklist

No IPK or source release should include live access.key, SSH keys, VLESS UUIDs, provider subscription URLs, account tokens, pairing codes, private network inventories or raw logs. README screenshots must be **empty-state captures** with nonfunctional TEST-NET demo data.

A passing build and a published GitHub Release do **not** certify independent security auditing, license compliance or acceptance to the Homebrew Channel catalog.

[Security](docs/SECURITY.md) · [Release limitations](docs/RELEASE_LIMITATIONS.md) · [README](README.md) · [Русский README](README.ru.md)
