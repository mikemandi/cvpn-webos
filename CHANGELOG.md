# Changelog

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
