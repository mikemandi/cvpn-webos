<div align="center">
  <img src="assets/icon160.png" width="112" height="112" alt="CVPN webOS" />

  # CVPN webOS

  **PRIVATE NETWORK · VPN for LG webOS**

  A free, independent VPN client for rooted LG televisions.
  Import compatible VLESS/Reality profiles or HTTPS subscriptions **from any provider**.

  **No CVPN subscription or account required.** CVPN is an optional recommended service.

  [**Download IPK**](https://github.com/mikemandi/cvpn-webos/releases/latest) ·
  [Installation](#getting-started) ·
  [Compatibility](#compatibility) ·
  [Security](docs/SECURITY.md) ·
  [Changelog](CHANGELOG.md)
</div>

---

## Simple setup, full control

<table>
<tr>
<td width="50%" align="center"><img src="assets/ru-dashboard.png" alt="CVPN Connect home and server selection" width="100%"></td>
<td width="50%" align="center"><img src="assets/ru-settings.png" alt="CVPN Connect settings" width="100%"></td>
</tr>
<tr><td align="center"><sub><b>Servers & VPN controls</b></sub></td><td align="center"><sub><b>Settings & diagnostics</b></sub></td></tr>
</table>

CVPN webOS installs from a **single All-in-One IPK**. On supported rooted TVs,
the app prepares its bundled backend and Xray automatically. Pair a phone
using the TV's QR code to manage connections, profiles, subscriptions and
server selection from the trusted local network.

| | |
|---|---|
| **App version** | `2.2.4` — Unified TV UI & RU/EN localization |
| **Backend / Engine** | 2.0.0 / 0.1.4 |
| **Xray-core** | 26.10.21 ARMv7, included and checksum-verified |
| **Language** | Russian on Russian-language LG; English otherwise; manual RU/EN choice in CVPN Connect |
| **Price** | Free to use; no mandatory subscription |
| **Supported VPN input** | VLESS/Reality and compatible HTTPS subscriptions (including third-party services) |
| **Installation** | Root + Homebrew Channel required |

## Getting started

1. Ensure your LG TV has **root**, **Homebrew Channel** and TUN support.
2. Download and install `su.cvpn.webos_2.2.4_all.ipk` from the [release page](https://github.com/mikemandi/cvpn-webos/releases/latest) **once v2.2.4 is published**. Until then, v2.2.3 remains the stable public release.
3. Open CVPN webOS and let it prepare the **locally bundled** components.
4. Scan the TV's QR code with a phone connected to the same trusted LAN.
5. Import a supported provider configuration or subscription and connect.
6. Confirm the initial trial-style connection to maintain a persistent VPN session.

**No separate backend or Xray download on the TV** is necessary. Existing
private configurations are preserved during an app upgrade. Disconnect VPN
normally before installing an updated IPK; the installer refuses to replace
active Engine/Xray binaries.

## Why CVPN webOS?

- **Universal client:** compatible providers and self-hosted endpoints work without a CVPN account.
- **TV-native onboarding:** compact matching QR and preparation screens with clear progress and nonblocking recovery.
- **RU / EN:** automatic locale choice plus persistent manual preference in CVPN Connect.
- **Recovery-focused:** integrity checks, automatic local setup, engine preflight,
  fail-safe route and DNS restoration, and no automatic VPN connection at boot.
- **LAN pairing:** temporary QR/PIN access and per-device revocation.

## Compatibility

CVPN webOS is built for **rooted LG webOS TVs with 32-bit ARMv7 userspace**,
a working `/dev/net/tun` and a compatible Homebrew root execution service.
It is **not** an LG Content Store app and will not provide system-wide VPN
connectivity on unrooted TVs.

**Real hardware acceptance:** the owner tested All-in-One v2.2.2 on
**LG 55UK6450PLC (webOS 4.4.3-22, firmware 05.50.70)**, including clean
installation, external Netherlands IP, YouTube playback, disconnect/reconnect
and normal power off/on. v2.2.3 changed app icons only. **v2.2.4 UI and
localization require separate physical LG verification before broad claims.**
See the [hardware QA record](docs/QA_ACCEPTANCE_2.2.2.md).

## Security & limitations

| Feature | Current behavior |
|---|---|
| VPN routing | IPv4 only; connection is refused when an active IPv6 default route is present |
| Kill switch | No universal kill switch; errors restore ordinary internet |
| Autostart | No automatic VPN connection after reboot |
| Local panel | HTTP port 8765 — **trusted private LAN only**, never forward publicly |
| Certificates | Readiness checks verify TLS certificates |
| Subscriptions | Stored locally; never include token URLs or keys in public bug reports |

[Security considerations](docs/SECURITY.md) ·
[Compatibility details](docs/RELEASE_LIMITATIONS.md) ·
[Hardware testing](docs/QA_ACCEPTANCE_2.2.2.md)

## Development

The project's CI builds the ARMv7 Go backend, verifies the pinned Xray asset,
creates the All-in-One IPK, then runs offline Go, launcher, installer,
packaging and rollback tests.

```bash
bash scripts/build-release.sh
bash tests/run_offline.sh
```

The public source is **not automatically an open-source license**. Until a
license is explicitly adopted, the Homebrew catalog classifies the listing as
`pool: non-free`; this is a licensing category, **not a purchase price**.
Bundled third-party components retain their own licenses. AI assistance was
used in development and testing; the project owner reviewed and tested real LG
hardware.

[Release & Homebrew submission instructions](docs/HOMEBREW_SUBMISSION.md) ·
[Third-party licensing review](LICENSE-PUBLISHING-REVIEW.md)
