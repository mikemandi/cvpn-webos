# CVPN webOS 2.2.5 — Homebrew Channel submission guide

[English](HOMEBREW_SUBMISSION.md) | [Русский](HOMEBREW_SUBMISSION.ru.md)

**Current status, 2026-10-28:** [v2.2.5 Stable](https://github.com/mikemandi/cvpn-webos/releases/tag/v2.2.5) is published on GitHub with IPK, SHA256SUMS, manifest and Homebrew YAML. **A GitHub Release does not mean inclusion in the official webOS Homebrew Channel catalog.** No maintainer approval is asserted here.

## Before requesting catalog inclusion

1. **Rights:** verify CVPN-authored code/artwork ownership and third-party license compliance. No open-source license is currently assigned to CVPN-authored code. Keep **pool: non-free** rather than main until an appropriate source license is explicitly adopted. This refers to source redistribution terms, **not the price**; the app is free.
2. **Third parties:** Xray-core **MPL-2.0** and QR **MIT / Apache-2.0** license notices must remain bundled. Verify corresponding-source availability and attribution.
3. **Privacy:** use real-UI empty-state screenshots with no subscriptions, UUIDs, real PINs, private LAN addresses or credentials. Sanitized English images: [screenshots](../assets/screenshots). Russian screenshots are also available.
4. **Safety:** document root/Homebrew requirements, IPv4-only routing, active IPv6 route refusal, absence of universal kill switch and autostart, and HTTP-only trusted-LAN control on TCP/8765.
5. **Review:** independent hardware acceptance across all models has not been established. Complete meaningful human review of system-level effects and rollback/standby risks before catalog distribution.
6. **AI disclosure:** current [webOS Homebrew rules](https://github.com/webosbrew/apps-repo) require meaningful human development, testing and review when AI tools assisted work; disclose use in the pull request. No piracy/DRM bypass claims or content.

## Submission workflow

1. Check all release assets at [v2.2.5](https://github.com/mikemandi/cvpn-webos/releases/tag/v2.2.5): IPK, SHA256SUMS, su.cvpn.webos.manifest.json and su.cvpn.webos.yml. Verify the file hashes, download addresses, version and icon.
2. Generate/review the catalog package YAML from scripts/make-homebrew-entry.py if needed. For this published version, use the **2.2.5 artifact**, not an older 2.2.4 or a manually patched IPK.
3. Fork [webosbrew/apps-repo](https://github.com/webosbrew/apps-repo), add the appropriate **packages/su.cvpn.webos.yml** entry, and keep pool: non-free unless licensing changes.
4. Use accurate actual screenshots, feature summary and limitations. Explain the root installer, local HTTP risk and exactly which television models were owner-tested.
5. Submit a PR to upstream maintainers **only after owner review and approval**. Explicitly disclose AI assistance, what was personally tested, remaining risks and who will maintain the app.
6. Address maintainer feedback. **Do not claim Homebrew catalog availability until the maintainers accept it**.

Useful references: [Homebrew submission rules](https://repo.webosbrew.org/submit) · [apps-repo](https://github.com/webosbrew/apps-repo) · [proposed PR text](HOMEBREW_PR_DRAFT.md).

## What users install

CVPN webOS installs as a single All-in-One IPK **after** the TV already has root/Homebrew. Its Xray, Go backend and VPN Engine are inside the package and prepared automatically. No account or paid CVPN subscription is required to import compatible third-party configurations.

**Do not publish a root bypass or alter firmware ACLs as a quick compatibility fix.** For safety details read [Security](SECURITY.md), [Release limitations](RELEASE_LIMITATIONS.md) and [Licensing review](../LICENSE-PUBLISHING-REVIEW.md).
