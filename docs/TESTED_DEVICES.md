# Device validation — CVPN webOS Stable

Last updated: 2026-10-26.

## Tested installations

| Device | Environment | Evidence | Outcome |
|---|---|---|---|
| Original development LG TV (model not recorded in this report) | Rooted webOS, established CVPN backend; Ethernet and Wi-Fi tested | Successful RC5 to Stable upgrade, Launcher 2.0.4 installed and launched, VPN exit-IP tests performed | Passed |
| LG 55UQ76003LD (2022) | webOS 22 / 7.6.0-1901; rooted Homebrew; `uname -m=aarch64`; `getconf LONG_BIT=32` | Fresh bootstrap, ARMv7 Xray 26.10.21 executed, Launcher installed, QR control and subscription import, VPN working after first-import and BusyBox timeout fixes | Passed with **post-install hotfixes** |

The second television reported firmware `04.06.91` in conversation; verify the exact on-device revision before treating that number as tested compatibility. No claim is made for all UQ7600 models or all webOS 22 updates.

## Why this matters for the next release

The old first-install bundle lacks a trusted client template and rejects `aarch64` kernels even when the userspace is ARMv7-compatible. The originally packaged Engine also used `timeout -t 6`, which BusyBox 1.31.1 rejects. Fixes are now in `main`:

- Kernel `aarch64` + **32-bit userspace** accepted by bootstrap; 64-bit only rejected.
- Secret-free embedded TUN template supports the first validated import and profile selection.
- Startup DNS probe invokes the portable `timeout 6 nslookup …` form.
- 90-second trial deadline, recovery, DNS routing safeguards and fail-closed IPv6 checks remain in place.

**Important validation boundary:** The second TV successfully reached working VPN after sequential hotfixes. A one-pass install from the consolidated next tagged release still needs CI and package inspection; cold boot, uninstall, and reinstall are not separately certified here.

## Release gate

1. Pass full `Verify build` workflow on final `main`.
2. Create a **new tag** (planned `v2.0.2`) rather than rewriting `v2.0.1`.
3. Verify release assets, SHA-256, public URLs and Homebrew manifest fields.
4. Add accurate, sanitized screenshots of the actual final TV launcher and mobile management UI.
5. Review source licensing, third-party notices and secrets in git history. Disclose AI-assisted development in the Homebrew PR.
6. Document that the IPK is only the TV launcher and that the root backend and external Xray binary must be installed separately.

Never include or publish subscriber links, VLESS UUIDs, auth cookies, login PINs or private IP addresses in screenshots or sample configurations.
