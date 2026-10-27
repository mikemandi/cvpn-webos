# Changelog

## 2.2.4 RC — unified onboarding, internationalization and automatic recovery (2026-10-28)

- Preparation and QR now use the same floating TV overlay, card and branded
  app icon: no full-screen blue setup illustration, no cramped stepper layout.
- Added RU/EN to **all** TV setup, error, QR and expiration countdown states.
  Russian TV settings render Russian; other languages fall back to English.
  A manual selection in CVPN Connect overrides the TV language on both screens.
- First launch retrieves the persisted preference via a local-only JSONP endpoint.
  Existing TV QR locale metadata is honored on display.
- Short, bounded panel readiness probe on repeated app launches; only a genuine
  install waits for the slower health/restart path. No fake loading percent.
- Removed the spurious intermediate **Continue** UI from interrupted installs:
  automatically poll READ-ONLY state until the existing installation completes,
  then close normally; only actionable errors show Retry/Close.
- The main VPN Engine, Xray, subscription storage and network routes are unchanged.
- Comprehensive offline UI localization, retry, timeout, and Go language API tests.
- **Hardware rollout:** v2.2.4 is a candidate until a TV confirms the new UI.
  Existing tested 2.2.2/2.2.3 VPN connection behavior is unchanged.

## 2.2.3 — transparent rounded app icon (2026-10-28)

- Fixed opaque square corners on the LG launcher and GitHub/Homebrew artwork.
- Preserved the original blue rounded icon tile, white TV/cast symbol,
  lighting and color treatment; only the pixels beyond the rounded square
  are transparent (true 8-bit alpha, antialiased borders).
- Updated all four shipped icon asset sizes from one 400px source
  (80, 130, 160 and 400px).
- No changes to Xray, backend, subscription support, root installer,
  recovery scripts or VPN Engine 0.1.4.
- 2.2.2 physical clean-LG test remains applicable to VPN functionality;
  2.2.3 is an app-artwork-only package revision.

## 2.2.2 RC — flash-friendly updates and orphan safety (2026-10-27)

- Bounded app-owned binary backup retention to two archives after successful
  upgrade, reducing flash/disk exhaustion on older LG TVs.
- Installer now refuses replacing files while an orphaned Xray process is active,
  even when Engine's session marker is absent.
- Regression tests cover backup bounds, rollback safety and live Xray refusal.
- 2.2.1 UI and timeout compatibility fixes unchanged. Root/TV startup, networking
  and fresh-install still require acceptance on clean LG hardware.

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
