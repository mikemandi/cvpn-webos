# Publishing CVPN webOS to webOS Homebrew Channel

## 2.1.0 All-in-One prerelease gate

For 2.1.0, the single `su.cvpn.webos_2.1.0_all.ipk` is expected to include Xray (ARMv7), backend, Engine, recovery and the TV launcher. End users must not install a backend ZIP, Xray binary, SSH tools or any PC scripts.

Before a Homebrew PR, validate real-device first launch via `org.webosbrew.hbchannel.service/exec`, root failure states, working VPN, upgrade while disconnected, refusal while VPN is active, no credential loss, rollback on missing disk space, reboot/startup and uninstall cleanup behavior. The test builds from the feature branch are not yet considered production-ready.

Confirm bundled Xray MPL-2.0 license notices and the CVPN project's own publishing license/pool. Do not claim published until a reviewed `v2.1.0` tag exists.


Official instructions: https://github.com/webosbrew/apps-repo and https://repo.webosbrew.org/submit.

**The Homebrew reviewers decide whether an app is listed.** Uploading an IPK to GitHub does not automatically add it to Homebrew Channel.

## A. Before opening a PR

1. Verify copyright ownership and pick a source license if opting into `pool: main`. The generated template defaults to `pool: non-free` until ownership/licensing is resolved.
2. Launcher 2.0.4 and a fresh backend setup have been physically verified on LG 55UQ76003LD after the first-import and BusyBox timeout patches. Before PR, complete a fresh tagged release build check; separately verify cold boot, network restoration and uninstall behavior. See `docs/TESTED_DEVICES.md`.
3. Replace historical design-preview images under `assets/` with accurate fresh live screenshots that contain **no** IP, subscription URL, UUID, token, PIN or user information.
4. Confirm this is a separate **rooted TV product**, not an LG Content Store compatible application; IPK alone is a TV launcher without an installed backend.

## B. GitHub repo and release

1. Create a public GitHub repository e.g. `OWNER/cvpn-webos` and upload the source tree (do not upload historic `All.zip`).
2. In Actions, wait for `Verify build` to pass. It uses Go 1.27.3 and creates an IPK and backend/root zip.
3. Create a **new, verified release tag** (planned `v2.0.2`) and push it; do not retag or overwrite older releases. The `Publish verified release artifacts` workflow publishes:
   - `su.cvpn.webos_2.0.4_all.ipk`
   - `cvpn-webos-root-2.0.0.zip`
   - `su.cvpn.webos.manifest.json`
   - `SHA256SUMS`
4. Verify that the public manifest URL and `ipkUrl` return HTTP 200 and the SHA-256 equals the IPK published in the release. Also verify iconUri is public.

## C. Add Homebrew listing

1. Generate a Homebrew YAML with your **actual GitHub owner and repo**:

   ```bash
   python3 scripts/make-homebrew-entry.py --repository OWNER/REPO
   ```

2. Fork https://github.com/webosbrew/apps-repo, put the generated YAML at `packages/su.cvpn.webos.yml` and create a PR to the upstream `main` branch.
3. In PR include this (edit only claims you personally verified):

   **Title:** Add CVPN webOS — VPN companion for rooted LG TVs

   **Body:**

   > CVPN webOS consists of a TV QR/pairing launcher and a separately installed root VPN backend using Xray-core. The launcher by itself does not configure VPN. Target: rooted ARMv7 LG webOS with supported Homebrew startup hooks. IPv6 fails closed; local panel is trusted-LAN HTTP only. The repository includes source, repeatable Go/ARMv7 build, root installer, security limitations and screenshots. Tests performed on [TV model/firmware, completed test results]. This project used AI-assisted development and packaging; the owner has reviewed and tested the submission. Source license / third-party notices: [confirmed licensing decision].

4. Await maintainer review. They may require architectural changes, current screenshots, or direct testing before accepting.

## D. If rejected or not yet published

Release the `IPK` and root zip to GitHub first. Users can sideload the IPK manually with `ares-install`; their TV still needs root/backend setup. Do not claim it is already in Homebrew Channel.

## Verified publication state (2026-10-26)

The two TVs have demonstrated the working VPN system. On the fresh LG webOS 22 television, two fixes were installed as separate hotfixes before a successful connection. The **new combined GitHub Release is not yet verified** and must not be represented as already installed in one pass. Public repository status has been confirmed; a new tag and successful Release Assets are still necessary.

Homebrew index rules require honest descriptions, validated root-only compatibility, and disclosure of AI assistance. Generated Homebrew YAML omits historical design screenshots. Add actual screenshots before filing the PR, and sanitize all QR/PIN/server details.
