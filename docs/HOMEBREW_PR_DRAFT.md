# Homebrew listing PR — CVPN webOS v2.2.4

**Title:** Add CVPN webOS — free VPN client for rooted LG televisions

**Suggested PR body:**

> CVPN webOS is a **free-of-charge** VPN client for rooted LG webOS TVs
> with 32-bit ARMv7 userspace, working TUN and Homebrew Channel.
> Users can import compatible VLESS/Reality configurations and HTTPS
> subscriptions from **any provider**. A CVPN (cvpn.su) subscription is
> optional, not a prerequisite and not a paywall.
>
> The All-in-One IPK includes native Xray-core, Go web backend, TV
> pairing UI and safety/recovery scripts. It prepares the components on
> first launch using Homebrew's root execution service; it does **not**
> download unknown executable payloads at runtime.
>
> **Safety and limitations:** no universal kill switch; no VPN autostart;
> the IPv4-only engine refuses connection with an IPv6 default route.
> Admin/pairing UI listens on port 8765 over HTTP for trusted local LAN
> use and must not be publicly exposed. Automatic recovery restores
> ordinary routes/DNS instead of leaving the TV offline.
>
> **Real hardware tests:** LG 55UK6450PLC, firmware 05.50.70,
> webOS 4.4.3-22: clean single-IPK automatic first launch, QR/pairing,
> VPN connection, a Netherlands exit IP via 2ip.io, YouTube playback,
> disconnect/reconnect, ordinary TV off/on. Additional tests were
> confirmed by the owner. Other models and firmware not guaranteed.
>
> **Development disclosure:** This project used AI assistance in code,
> interface design, packaging and tests. The project owner tested the
> installation and VPN on real LG devices, reviewed the actual behavior,
> and is responsible for releases and maintenance. This does not claim
> an independent professional code audit.
>
> **License and pricing:** The app is free of charge. CVPN-authored
> source is publicly visible but has no OSI-approved redistribution
> license, so this listing uses Homebrew `pool: non-free`. That does
> **not** describe an app price. Bundled Xray-core retains its
> MPL-2.0 license; QR library license and notices are bundled.
>
> See https://github.com/mikemandi/cvpn-webos for installation and
> security limitations. Public release:
> https://github.com/mikemandi/cvpn-webos/releases/tag/v2.2.4

**Required submission steps once release assets are visible:**

1. Create a personal fork of `webosbrew/apps-repo` on GitHub.
2. Put the generated `su.cvpn.webos.yml` in the fork under
   `packages/su.cvpn.webos.yml`.
3. If possible, add **actual screenshots** from the app, scrubbed of
   private pairing PIN, subscription tokens, UUIDs and LAN addresses.
   Do not substitute AI concept renders for real product screenshots.
4. Open a pull request into `webosbrew/apps-repo:main` with the draft above.
5. Respect maintainers' requests before claiming Homebrew listing approval.
