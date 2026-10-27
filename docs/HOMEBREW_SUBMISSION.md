# Publish CVPN webOS 2.2.3 to webOS Homebrew Channel

**As of 2026-10-28: Homebrew listing not yet approved.** The 2.2.2 All-in-One IPK has
passed offline CI and the core fresh-TV functional flow on one LG
55UK6450PLC / webOS 4.4.3-22 / firmware 05.50.70 (owner-reported).
Version 2.2.3 changes only the rounded app icon PNG assets; the tested VPN components are unchanged. The user did **not** manually install Xray/backend: the Homebrew root
service prepared bundled components on first app launch. The owner also
confirmed a Netherlands exit address, YouTube playback, VPN disconnect,
reconnect and ordinary television off/on.

These results are **not evidence that uninstall, power-loss or
fault-injection work on every model**. See
[QA acceptance](QA_ACCEPTANCE_2.2.2.md) and
[limitations](RELEASE_LIMITATIONS.md).

## Review gates before tagging

1. Confirm CVPN code and artwork ownership and distribution rights. The app is
   **free of charge**, supports **other providers' subscriptions** and has
   **no mandatory CVPN account or payment**. Homebrew `pool: non-free` means
   the CVPN-authored code currently has **no declared open-source license**;
   it does **not** mean paid software. Do not invent a source license.
   Keep the MPL-2.0 and QR library notices for bundled third-party code. Refer to
   [licensing review](../LICENSE-PUBLISHING-REVIEW.md).
2. Review the live Homebrew [submission rules](https://github.com/webosbrew/apps-repo):
   **NO AI SLOP** — meaningful owner code review and testing, plus explicit
   disclosure of AI-assisted development in the PR; safe system networking;
   no DRM bypass/piracy.
3. Finish remaining acceptance: handling of network-loss/failure conditions,
   correct persistent-data/uninstall policy, release-mode source/security audit,
   no system DNS/routes left behind, and accurate sanitized screenshots
   captured from the **actual LG application** rather than AI mockups.
4. Record tested platforms precisely, don't claim broad model support,
   a universal kill switch, or IPv6 tunneling. No Content Store compatibility.
5. Do not include personal television addresses, real subscription URLs,
   client UUIDs, Homebrew SSH keys, tokens or other private data in screenshots.

## Release process (after review approval)

1. Merge the verified candidate PR into `main` through normal review;
   do not change working deployed LG televisions in the process.
2. Create an **immutable new** release tag `v2.2.3` on the reviewed commit.
   `.github/workflows/release.yml` builds and tests the IPK and creates
   GitHub Release assets.
3. Inspect the release's
   `su.cvpn.webos_2.2.3_all.ipk`, `su.cvpn.webos.manifest.json`, and
   `SHA256SUMS`. Verify SHA-256, manifest URL, icon URL and release download
   all resolve correctly; verify no credentials and no separate root ZIP.
4. Generate `dist/su.cvpn.webos.yml` with
   `python3 scripts/make-homebrew-entry.py --repository mikemandi/cvpn-webos`.
   Its root requirement, version and limitations must remain accurate.
   Use `pool: non-free` until the project owner explicitly adopts a valid
   open-source license (even though the code is publicly visible).
5. Fork [webosbrew/apps-repo](https://github.com/webosbrew/apps-repo), add
   `packages/su.cvpn.webos.yml` with **real** sanitized screenshots, and
   submit a PR to its maintainers. Clearly disclose AI-assisted code/design
   and summarize your review and hardware testing. The connected GitHub
   integration currently has read-only permissions to upstream apps-repo;
   submit from a personal fork instead.
6. Wait for maintainer approval before claiming Homebrew Channel availability.

## What users install

Root and Homebrew are prerequisites and **not supplied by CVPN**.
Users sideload the **single All-in-One IPK**; backend, Xray and Engine are
bundled and configured automatically by the root service on supported TVs.
Afterwards users pair via QR on a trusted LAN and control VPN from CVPN Connect.

**Known boundaries:** rooted ARMv7 32-bit user space only, active IPv6
default route blocks connection, no universal kill switch, no VPN autostart,
and local HTTP management only on trusted LAN. App update must not modify
Engine while VPN is connected. See [security](SECURITY.md).

## Homebrew maintainers' current rules

- https://repo.webosbrew.org/submit
- https://github.com/webosbrew/apps-repo
- https://github.com/webosbrew/webos-homebrew-channel

Do not publish, open a submission PR or automatically merge before the
owner has approved the public-facing description, licensing and risk policy.
