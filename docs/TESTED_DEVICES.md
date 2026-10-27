# CVPN webOS — tested hardware & compatibility

Updated: 2026-10-28 · Current TV app: **2.2.5**

This is an evidence log, **not a list of guaranteed supported LG TVs**. Rooting, Homebrew versions, firmware updates, userspace architecture and TUN availability can affect compatibility.

## Hardware evidence

| Television | Software / architecture | Verified | Important limits |
|---|---|---|---|
| **LG 55UK6450PLC** | webOS 4.4.3-22, firmware 05.50.70; rooted Homebrew, ARMv7-compatible system | All-in-One **2.2.2** first-run preparation, QR and local panel; Netherlands VPN exit IP; YouTube, disconnect/reconnect and normal power off/on | Verified functional baseline; 2.2.3 only changed the app icon |
| **LG 55UQ76003LD (2022)** | webOS 22 / 7.6.0-1901, rooted; aarch64 kernel with **32-bit userspace** | Fresh bootstrap and ARMv7 Xray 26.10.21, first import and VPN connection after sequential first-import/BusyBox fixes | This test used hotfixes; it does **not** prove every firmware supports a clean 2.2.5 install |
| **Owner's LG test device for 2.2.5** | Exact model not separately recorded against the 2.2.5 visual acceptance | Owner confirmed the refined English/Russian setup and QR UI and remote-focused error dialog work on real LG hardware | A UI acceptance report; not a new cross-model network certification |

Earlier TV tests used Ethernet and Wi-Fi; no universal compatibility claim is made for either.

## What changed since the earlier release

- **2.2.2:** self-contained All-in-One IPK with locally bundled backend, Xray, launcher, rollback and root service safeguards.
- **2.2.3:** transparent rounded app icon only.
- **2.2.4:** matched floating TV preparation/QR screens, RU/EN locale selection, bounded component checks and automatic completion polling.
- **2.2.5:** corrected Flexbox spacing for older LG Chromium, centered text-only TV wordmark and explicit remote navigation between **Retry** and **Close**. VPN Engine and Xray transport remained unchanged.

The consolidated package also contains previous repairs for ARMv7-compatible 32-bit userspace running under aarch64 kernels, the first validated import and portable BusyBox DNS-probe timeout.

## Compatibility prerequisites

A testable installation requires:

1. Rooted LG webOS TV with a functioning Homebrew root service.
2. ARMv7-compatible **32-bit** userspace and working **/dev/net/tun**.
3. Sufficient app storage, a trusted LAN and a phone/browser to pair with the TV.
4. A supported VPN provider profile or configuration for an actual connection.

Unrooted TVs and standard LG Content Store installations are not supported. VPN traffic is IPv4-only; connection is refused when an active IPv6 default route would bypass it. There is no universal kill switch or automatic VPN connection after a reboot.

## What still needs testing

- One-pass installation of the current IPK on additional untouched rooted LG models.
- Firmware-specific Homebrew root-service compatibility and rollback behavior.
- Repeated standby/cold-boot cycles and uninstall/reinstall scenarios on new TV hardware.

Keep bug reports anonymized. Never include subscription URLs, VLESS UUIDs, tokens, pairing codes, local access keys or private TV IP addresses.

[README](../README.md) · [Security](SECURITY.md) · [Release limitations](RELEASE_LIMITATIONS.md)
