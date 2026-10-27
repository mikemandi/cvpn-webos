# CVPN webOS 2.2.2 — supported devices and release gates

**Status: Release candidate on the `feature/all-in-one-v2.1.0` branch, NOT a proven production release.**
Do not merge into stable/main, publish to Homebrew, or promise one-click compatibility
until a clean TV passes the hardware acceptance plan in `docs/QA_ACCEPTANCE_2.2.2.md`.

## Supported platform

- LG webOS with **working root and Homebrew Channel**. The LG Content Store is not supported.
- ARMv7 32-bit userspace: `armv7l` or compatible `aarch64` kernel with `getconf LONG_BIT=32`.
- Working `/dev/net/tun`, IPv4 default gateway, supported WAN interface `eth0`, `wlan0` or `wlan1`.
- Homebrew root execution service must permit `su.cvpn.webos` to call
  `org.webosbrew.hbchannel.service/exec`. The app requests `lunabus.query`;
  Homebrew service mappings vary by release. Test this from the **app itself**,
  not only from a root SSH shell. **Do not automatically elevate/update Homebrew**
  or patch system Luna permission files.
- Native Xray v26.10.21 ARMv7 and the Go panel are shipped inside the same IPK;
  no first-run internet download, remote install script, or SSH is intended.
- IPK package version 2.2.2 / backend 2.0.0 / VPN Engine 0.1.4.
  The existing proven Engine routing/recovery code is retained.

## Network and security trade-offs

- IPv4-only tunneling. If IPv6 has a default route, the Engine refuses activation.
- No system-wide kill switch. On disconnect/error, **ordinary connectivity is restored**.
- No VPN autostart at boot. A currently active VPN is not automatically resumed after restart.
- Temporary DNS redirect to 8.8.8.8 through TUN; HTTPS readiness by hostname.
  The startup check does **not** prove DoH, UDP/QUIC, streaming, or absence of every DNS leak.
- If Wi-Fi fails, watchdog/recovery may restore ordinary internet rather than reconnect;
  automatic retry could hide real network issues and is deliberately not enabled.
- Local control over HTTP port 8765 must be limited to a **trusted private LAN**;
  never port-forward. Profile secrets are protected by filesystem permissions,
  **not encrypted at rest**.
- Updates refuse to replace Engine/Xray while an active VPN, its route, its TUN,
  the DNS bind mount or Engine lock is present. Users should disconnect before upgrading.
- Browser lifecycle, root-service permissions and power management vary by LG model.
  We cannot validate those scenarios from CI alone.

## Required evidence before publication

1. CI succeeds: Go tests, race detector, Go vet, installer rollback, offline
   launcher behavior, integrity and IPK inspection.
2. Fresh installation on a fourth untouched rooted/Homebrew TV with **no SSH or
   manual Engine installation**, including proper `lunabus.query` mapping.
3. First-launch/second-launch, QR pairing, profile import, connection (trial then
   persistent), actual video playback, disconnect and rollback.
4. Wi-Fi flaps, standby/return, unplanned process death, corrupted file scenarios,
   no hanging preparation UI, and recovery after reboot on that test TV.
5. No personal secrets/TV IPs in the IPK; review permissions, CVE/dependencies
   and third-party licensing; complete real screenshots and uninstall policy.

**Stop-release:** any unrecovered routes/DNS, loss of power/root, exposed secrets,
no-response first-run without user-visible recovery, or need for SSH on a fresh TV.
