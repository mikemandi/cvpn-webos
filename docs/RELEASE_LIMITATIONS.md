# Supported configuration and publication gates

- TV must be rooted with Homebrew Channel startup hooks, Linux ARMv7/armv7l, `ip`, `curl`,
  `sha256sum`, BusyBox/ash, and working TUN support. Not LG Content Store compatible.
- IPK is the television launcher / pairing UI **only**. A separate root backend install is
  required. Users MUST NOT expect VPN to work by installing the IPK alone.
- New TV needs a compatible Xray ARMv7 binary supplied by the owner and checked on that TV.
  None is bundled here (no hidden downloads or arbitrary remote scripts).
- Web-panel over HTTP port 8765: trusted LAN ONLY; do not port-forward; cookie interception
  is possible on hostile Wi-Fi networks.
- IPv6 does NOT pass through tunnel: the engine refuses activation on IPv6 default route.
- There is no generic Kill Switch, and the engine restores normal connectivity on failure.
- Existing TV's Stable 2.0.0 backend is not changed by this publication package.
- TV launcher 2.0.4 is a separate component version; backend/engine are 2.0.0/0.1.4.
- GitHub workflow must pass and root bootstrap must be tested on an independent TV before
  reviewers should consider final acceptance. Root/boot lifecycle is hardware-dependent.
- Material screenshots under assets are *historical interface previews*, not certified
  live screenshots of the final release. Replace with real captures before PR submission.
