# Publishing CVPN webOS to webOS Homebrew Channel

Official instructions: https://github.com/webosbrew/apps-repo and https://repo.webosbrew.org/submit.

**The Homebrew reviewers decide whether an app is listed.** Uploading an IPK to GitHub does not automatically add it to Homebrew Channel.

## A. Before opening a PR

1. Verify copyright ownership and pick a source license if opting into `pool: main`. The generated template defaults to `pool: non-free` until ownership/licensing is resolved.
2. Test the 2.0.4 launcher and **new TV bootstrap** on actual LG ARMv7. Test QR, cold boot, a fresh subscription, DNS/IPv6 safety and uninstall/recovery. The user's previous LG validated backend 2.0.0 but not this new installation pipeline.
3. Replace historical design-preview images under `assets/` with accurate fresh live screenshots that contain **no** IP, subscription URL, UUID, token, PIN or user information.
4. Confirm this is a separate **rooted TV product**, not an LG Content Store compatible application; IPK alone is a TV launcher without an installed backend.

## B. GitHub repo and release

1. Create a public GitHub repository e.g. `OWNER/cvpn-webos` and upload the source tree (do not upload historic `All.zip`).
2. In Actions, wait for `Verify build` to pass. It uses Go 1.27.3 and creates an IPK and backend/root zip.
3. Create tag **`v2.0.0`** and push it. The `Publish verified release artifacts` workflow publishes:
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
