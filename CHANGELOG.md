# Changelog

## 2.2.1 RC — UI reference, installer resilience, code audit (2026-10-27)

- Revised the entire first-run TV layout to match the supplied CVPN webOS
  design reference: cobalt/navy gradient, glass panel, branded app mark,
  animated TV connection ring, connected progress timeline and clear safety note.
- Extended first-run status timeout to 35 seconds to cover the four
  5-second panel health retries on slow ARMv7 televisions.
- Differentiate live vs stale installer locks. Safely recover a stale PID lock
  **only after** confirming there is no active CVPN state, TUN, split-route or
  DNS bind mount. Never discard a lock with ambiguous ownership.
- Runtime/UI test and installer regression coverage updated.
- Engine 0.1.4, Go backend 2.0.0 and native Xray unchanged.
- Release candidate only: fresh-LG first-run/permissions/rollback/standby
  and uninstall lifecycle still require actual device acceptance.

## 2.2.0 RC — premium onboarding and first-install hardening (2026-10-27)

- Reworked the installation overlay to a legible, understated three-step
  CVPN preparation experience, with reduced-motion support, accessible retry
  focus, success and safe recovery states. All binaries are included in IPK:
  there is **no** runtime download, fake percentage, or network installer.
- Declare `lunabus.query` application permission for Homebrew root-service
  calls. Older Homebrew service mapping may require a supported version;
  no automatic elevation or firmware changes are performed.
- Finite Homebrew callback deadlines; reject malformed permission responses;
  prevent concurrent install retries and stale callbacks closing QR.
- In-progress installation state reported from Engine lock; preserve existing
  network and private profiles on update or failure.
- Runtime Launcher smoke tests, offline installer lock/rollback tests,
  refreshed release safety policy, and fresh-LG acceptance checklist.
- VPN Engine 0.1.4, backend 2.0.0, Xray 26.10.21 unchanged.
- **Not yet production certified:** fresh-TV root-service/ACG permission
  and no-SSH installation remain a required hardware acceptance test.

## 2.0.0 Stable — public source/packaging preparation (2026-10-25)
- Backend 2.0.0 and VPN Engine 0.1.4 unchanged from the TV-tested Stable.
- Versioned launcher IPK 2.0.4 (separate from VPN engine version), historical hard-coded
  sample LAN address removed from unused launcher helper.
- Standalone launcher IPK, clean-root bootstrap scripts, boot hook, public-source layout,
  release integrity metadata builder, GitHub Actions, Homebrew submission draft.
- No personal TV profiles, authentication keys, subscriptions or Xray binary distributed.
- Separate root backend bootstrap is a NEW workflow and has NOT been tested on real LG.


### Unreleased — fresh installation import repair

- Built-in secret-free TUN template for clean installs without a pre-existing client configuration.
- First selected and Xray-validated server creates active configuration files with private permissions.
- Existing installations keep their original VPN profiles and configuration.
- Recoverable panel-only repair scripts; no changes to Xray, networking or TV Launcher.

### Unreleased — portable DNS timeout on webOS 22

- Use `timeout 6 nslookup ...` instead of `timeout -t 6 nslookup ...` in the Engine's startup DNS probe; the latter fails immediately with BusyBox 1.31.1.
- Preserve bounded DNS probing and the existing safe rollback; do not disable TLS checks or extend the VPN trial timeout.
- No tagged GitHub Release contains this update yet; existing TVs require a reviewed Engine-only update.
