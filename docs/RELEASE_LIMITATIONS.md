# CVPN webOS 2.2.6 — compatibility and release limitations

[English](RELEASE_LIMITATIONS.md) | [Русский](RELEASE_LIMITATIONS.ru.md)

**Current public release:** [v2.2.6 Stable](https://github.com/mikemandi/cvpn-webos/releases/tag/v2.2.6) (2026-10-28). Stable means a publicly published, owner-accepted build; it **does not** assert comprehensive firmware support, an independent security audit or Homebrew Channel catalog approval.

## Supported environment

- Rooted **LG webOS** television with compatible Homebrew Channel root-service permissions.
- ARMv7-compatible **32-bit userspace**: native armv7l or a compatible aarch64 kernel reporting 32-bit userspace.
- Working **/dev/net/tun**, IPv4 network gateway and supported network interface (Ethernet/Wi-Fi as provided by the target platform).
- Trusted local network and a phone/computer for browser pairing. The local API is HTTP on **TCP/8765**.
- A compatible VPN profile/subscription from any provider. **No account or payment to CVPN is necessary** to run the client.

Not compatible with unrooted TVs or official LG Content Store installation. Firmware, Homebrew service permissions and available Linux tools vary by LG model.

The Homebrew-compatible 2.2.6 IPK includes the CVPN license and third-party notices; no VPN behavior changed from 2.2.5. The All-in-One IPK contains the **TV launcher 2.2.6, backend 2.0.0, VPN Engine 0.1.4 and Xray-core 26.10.21 ARMv7**. The initial setup does **not** require separate downloads of Xray/backend on the TV.

## VPN and security trade-offs

- **IPv4-only:** VPN applies to supported IPv4 routing. If an active IPv6 default route exists, the engine refuses to connect to avoid an obvious bypass.
- **No universal kill switch.** If the VPN disconnects or encounters an error, ordinary internet access can be restored. Never treat the product as guaranteed leak prevention.
- **No VPN autoconnect after reboot.** Reopen and connect manually.
- **DNS:** network recovery attempts to restore routes and DNS. HTTPS connectivity tests validate certificates/hostnames, but they do **not** prove that DoH, all UDP/QUIC traffic, streaming applications or every DNS request is tunneled.
- **LAN management:** pairing codes, sessions and CSRF provide application checks, but **HTTP over LAN is not end-to-end encrypted**. Never expose port 8765 to the WAN.
- **Secret storage:** profiles and access keys stay on the TV with restrictive file permissions; they are **not encrypted at rest**.

## Installation and failure behavior

- The root installer refuses to overwrite engine/Xray binaries if it detects an active VPN process, TUN, VPN route, DNS mount or lock.
- The installer keeps existing profiles, backs up app binaries and attempts rollback after partial failure. Unexpected power loss or unsupported firmware can still require recovery.
- A launcher Homebrew callback timeout is not proof that installation stopped; the app performs read-only status polling rather than attempting a concurrent installer.
- Removing the IPK may **leave persistent application profile/backend data**. Do not assume uninstall securely erases keys.
- First-run Homebrew root permission changes are not made automatically, and users must not silently patch system ACLs to bypass compatibility checks.

## Actual acceptance evidence

- **LG 55UK6450PLC, webOS 4.4.3-22 / firmware 05.50.70:** owner-reported All-in-One 2.2.2 first run, QR/pairing, Netherlands VPN exit IP, YouTube playback, disconnect/reconnect and normal off/on.
- **LG 55UQ76003LD, webOS 22 / 7.6.0-1901, LG software 04.06.91:** VPN and Xray after first-import/BusyBox/DNS fixes, with aarch64 kernel and compatible 32-bit ARM userspace; not a one-pass 2.2.5 certification.
- **LG 86NANO906NA, webOS 5.6.2:** root/SSH, ARMv7, TUN, VLESS/Reality/XHTTP and connectivity/route restoration checked in development; an isolated protected-DNS test rolled back safely. Owner currently confirms working CVPN.
- **LG 65UK6710PLB, webOS 4.4.3-22:** owner currently confirms working CVPN; retained detailed clean-install/VPN logs unavailable.
- **v2.2.5 UI on LG:** the owner confirmed that refined TV dialogs, Russian/English interface and error-dialog remote focus work on hardware; do not infer additional model-specific network testing.

**The oldest directly tested build is webOS 4.4.3-22.** Homebrew metadata uses **4.4+** conservatively; it does not certify all firmware from that range. See [tested hardware](TESTED_DEVICES.md) and [2.2.2 QA evidence](QA_ACCEPTANCE_2.2.2.md).

## Known gaps and future validation

Physical tests across additional LG models and firmware, fault injection, standby, network flaps, clean reinstall/uninstall, power loss during upgrade, cross-app DNS/IPv6 behavior, security review, component/asset licensing and upstream Homebrew approval remain limited or incomplete.

**Stop and report immediately** if routes/DNS cannot be restored, existing private profiles are unexpectedly erased, the TV loses normal network access, or an installation requires privilege bypasses or exposed secrets.

[Security](SECURITY.md) · [Licensing review](../docs/legal/LICENSING_REVIEW.md) · [English README](../README.md) · [Русский README](../README.ru.md)
