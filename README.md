# CVPN webOS

**PRIVATE NETWORK · CVPN Connect** — VPN and private subscription management for
supported **rooted LG webOS televisions**. The interface is controlled from a
phone on the same trusted LAN via a QR pairing flow.

**Current candidate:** 2.2.1 All-in-One (feature branch; not yet Homebrew
published). **Backend:** 2.0.0. **Engine:** 0.1.4. **Xray:** 26.10.21 ARMv7.

[Website](https://cvpn.su) · [Security](docs/SECURITY.md) ·
[Hardware acceptance](docs/QA_ACCEPTANCE_2.2.1.md) ·
[Compatibility](docs/RELEASE_LIMITATIONS.md) ·
[Changelog](CHANGELOG.md)

## What is included in the All-in-One IPK

`su.cvpn.webos_2.2.1_all.ipk` includes:

- A TV Launcher: QR pairing and premium preparation/onboarding interface
- Native ARMv7 Go 2.0.0 control panel and its local mobile interface
- Xray-core ARMv7 26.10.21 (pinned upstream hash, bundled in build)
- VPN Engine 0.1.4 with safe TUN/DNS/route recovery
- Homebrew boot hook and installation/recovery scripts
- SHA256 checksum manifest for bundled runtime executables

The **TV does not download Xray or the backend** on first launch. Its
initialization screen prepares verified components already included in the IPK.

## Status of 2.2.1

The previous candidate **2.1.3** successfully connected and played YouTube on
the third LG, including recovery of video after temporary Wi-Fi interruptions.
On that TV, the initial installation had to be completed manually via root SSH;
this is an *unresolved first-run acceptance issue*, not a shipping workflow.

This 2.2.1 candidate addresses it by declaring `lunabus.query`, correctly
handling the Homebrew root-service callback, bounding the waiting UI, presenting
clearer error states and preventing parallel installers on retry. **The
improvement still requires validation on a fresh fourth LG** before a production
release is justified. Existing stable TVs are not modified by building this branch.

## Intended installation (for supported hardware)

1. Owner obtains working **root** and installs a compatible **Homebrew Channel**.
2. Install the single All-in-One IPK using supported webOS installation tools.
3. Open **CVPN webOS** on the TV. See the short, staged preparation screen.
4. Scan the TV's QR code to open **CVPN Connect** on a phone on the same LAN.
5. Add a VLESS/Reality configuration or HTTPS subscription and choose a server.
6. Connect, verify playback, and confirm the trial into persistent mode.

End users should not have to use SSH, PowerShell, or fetch Xray manually.
**That experience is a QA gate, not yet a hardware-certified claim.**

The preparation screen has three customer-facing stages (check device, prepare
components, launch CVPN), modest CSS animation and no fake download percentage.
It handles unexpected service responses, installer errors and timeouts with
Retry/Close. If an install command times out, it first checks existing status
instead of launching a concurrent installation.

### Requirements

- LG webOS with supported ARMv7 **32-bit userspace** (`armv7l` or compatible
  `aarch64` kernel plus 32-bit userspace)
- Working root, TUN device, Homebrew Channel root execution service
- The app's `lunabus.query` caller permission and the matching Homebrew service
  ACG mapping; mapping behavior differs by Homebrew version and webOS firmware
- Ethernet or Wi-Fi with a valid IPv4 default route
- Trusted private LAN for the unencrypted local mobile management interface

Not supported: stock/unrooted LG, LG Content Store, native 64-bit ARM target,
and systems with an exposed IPv6 default route the IPv4-only Engine cannot protect.

**Do not update or elevate Homebrew automatically** to work around a denied
service call; diagnose the root-service permission mapping on the actual TV.

## Keeping the network safe

- No VPN auto-start on reboot; no universal kill switch.
- Readiness checks validate a real HTTPS hostname over the current TUN routes
  with certificate verification enabled. They do not prove streaming or UDP/QUIC.
- The Engine monitors TUN, DNS and route consistency. On failure it restores
  the TV's original DNS/routes to ordinary connectivity.
- First connection uses a bounded trial and must be confirmed for persistent mode.
- Component upgrades **refuse** to overwrite active VPN files or a busy
  Engine lock. Disconnect VPN normally before updating.
- Existing subscription and session secrets stay local and survive upgrades.

See [security considerations](docs/SECURITY.md) and
[publication limitations](docs/RELEASE_LIMITATIONS.md).

## Developing and building

GitHub Actions `.github/workflows/ci.yml` builds the cross-compiled Go
control panel and produces the All-in-One IPK in the artifact
`cvpn-webos-preview`. The Xray binary is fetched and SHA256-verified in CI,
not installed from the internet at runtime. Published asset checksums
are **integrity evidence, not developer code signatures**.

On a Linux build environment with Go, Python, Node.js, `dpkg-deb` and the
repository's tooling:

```bash
bash scripts/build-release.sh
bash tests/run_offline.sh
```

The build output is `dist/su.cvpn.webos_2.2.1_all.ipk`. The offline suite
includes Go race/vet tests, fake-root installer integration/rollback,
headless Launcher JS tests, package checks and shell syntax tests.

Homebrew publication, manifest creation and third-party licensing must be
reviewed before tagging a public release. See
[LICENSE-PUBLISHING-REVIEW.md](LICENSE-PUBLISHING-REVIEW.md).

## Test the next fresh LG before release

Follow **[QA_ACCEPTANCE_2.2.1.md](docs/QA_ACCEPTANCE_2.2.1.md)**, particularly:

- Fresh TV and single IPK; **no command line or manual component installation**
- Normal/QR launch, retry after permission denial or delay, repeated launch
- QR pairing, profiles and subscriptions, VPN trial/persistent/stop/switch
- DNS/HTTPS/YouTube, Wi-Fi drop/recovery, Ethernet, standby/reboot
- Network restoration and preservation of access keys/secrets
- Uninstall/reinstall lifecycle, software permissions and LAN isolation

**A CI-green build is not a guarantee of production readiness on LG hardware.**
Do not merge the candidate into stable/main before the hardware gates pass.
