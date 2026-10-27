# CVPN webOS — security notes (2.2.0 candidate)

**Installation:** All-in-One IPK (TV Launcher + root backend + VPN Engine +
pinned ARMv7 Xray). **Root:** obtained by the TV owner, never by CVPN.
**Mobile control:** local HTTP on TCP/8765, not encrypted in transit.

## Trust boundaries

- TV with root/Homebrew is trusted by design. The `/exec` Luna endpoint
  can run privileged code; the app requests `lunabus.query` and invokes
  only hard-coded, package-owned `install.sh status` or `install.sh install`.
  Do not pass user-originated filenames, URLs or other strings to root execution.
- Homebrew root-service ACLs vary by firmware and Homebrew release. A request
  made as root from SSH does not prove the browser application can call it.
  Do not use system permission changes/elevation as an automatic workaround.
- Bundled component SHA256 hashes are checked before installation. These
  prevent corrupted payload usage but are **not** an independent signature
  authenticating the publisher.
- Existing access keys, profiles, subscription URLs and sessions remain local,
  protected by 0600/0700 permissions and **not encrypted at rest**. No
  credentials, unique TV IPs, private URLs or SSH keys may be included in IPK.
- Local web-panel controls check LAN source, host, pairing PIN/QR, sessions and
  CSRF, but client-server HTTP traffic remains observable on hostile networks.
  Use only a **trusted private LAN**; NEVER port-forward TCP 8765 to the internet.
- Do not open local-control URLs containing sensitive tokens on a shared or
  untrusted network; redact QR/PIN/private server addresses before publishing logs.

## Network behavior and release risk

- Engine 0.1.4 protects IPv4 only. It rejects an existing IPv6 default route.
- No system-wide kill switch; failures can restore ordinary internet rather
  than block every connection. No automatic VPN startup at boot.
- Startup verifies HTTPS by hostname through the active TUN with normal TLS
  certificate validation, never `curl -k`. This does not guarantee encrypted
  DNS transport, leak-proof behavior for all applications, UDP/QUIC or video.
- Watchdog/expiry restore original DNS bind mount and routes. A release
  candidate must pass physical rollback, Wi-Fi failure and reboot/standby
  tests before being presented as production-ready.
- Installer refuses network-active upgrades, keeps user data, backs up
  binaries before replacement and attempts rollback after partial failure.
  A timeout in the **UI** does not stop a potentially active installer;
  repeated checking must not create overlapping root install processes.
- Removing the IPK can leave persistent local backend/profile data. Agree on
  an owner-controlled uninstall/removal policy prior to Homebrew publication.

## Publication requirements

1. Reconcile code/artwork ownership and distribution licenses, Xray-core
   license notices and QR library notices, and choose a valid CVPN license.
2. Inspect artifact SHA256 digest and embedded files for secrets.
3. Prove first-run bootstrap via IPK **without SSH** on a separate clean LG.
4. Confirm frontend permission ACG and root-service behavior on the target webOS.
5. Validate normal shutdown, standby, boot, DNS/route recovery and user-data
   preservation; use `docs/QA_ACCEPTANCE_2.2.0.md`.
6. Capture real TV and phone screenshots, removing identifiers.
7. Satisfy Homebrew submission rules including disclosure requirements.

Do not include `All.zip`, historic profiles, exported subscriptions,
private keys or home-network diagnostics in public releases.
