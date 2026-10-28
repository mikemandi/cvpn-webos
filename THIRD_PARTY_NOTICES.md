# CVPN webOS — Third-party licenses and source availability

[English](THIRD_PARTY_NOTICES.md) · [Русский](THIRD_PARTY_NOTICES.ru.md)

The [CVPN webOS proprietary license](LICENSE) applies **only** to material the
project's rightsholders own or are authorized to license. It does not replace
licenses for embedded third-party software or their copyright notices.

## Xray-core (MPL-2.0)

- Project: https://github.com/XTLS/Xray-core
- Version included in the official v2.2.5 All-in-One IPK: **v26.10.21**
- Exact-version source: https://github.com/XTLS/Xray-core/tree/v26.10.21
- Source archive: https://github.com/XTLS/Xray-core/archive/refs/tags/v26.10.21.tar.gz
- Original license: https://github.com/XTLS/Xray-core/blob/v26.10.21/LICENSE
- License: **Mozilla Public License Version 2.0 (MPL-2.0)**
- Downloaded executable: Xray-linux-arm32-v7a.zip, with pinned SHA-256
  validated by scripts/fetch-xray.sh.

The full Xray license is downloaded during the build and shipped inside the
All-in-One IPK at payload/XRAY-LICENSE.txt. Recipients may use, inspect, modify
and redistribute the MPL-covered portions **under MPL-2.0**, even though the
separate CVPN-authored components use a proprietary license. If a distributor
modifies MPL-covered Xray files and distributes executables, the
corresponding modified covered source must also be made available as required
by MPL-2.0. A license for CVPN-owned code cannot narrow these rights.

## Bundled QR generation code

- launcher/qr-bundle.js embeds/derives QR rendering code.
- Original **QRCode for JavaScript**, © 2009 Kazuhiko Arase:
  **MIT** (as identified in launcher/QR-NOTICE.txt).
- Adapted **qrcode-terminal 0.12.0**: **Apache License 2.0 (Apache-2.0)** with its retained
  original notices identifying the included MIT-licensed work.
- The QR attribution and license notices are stored in
  launcher/QR-NOTICE.txt and launcher/QR-THIRD-PARTY-LICENSE.txt, and
  both files are included in the IPK.

Recipients may exercise rights in the QR-derived parts under those licenses.
Neither the MIT/Apache notices nor attribution is replaced by the
CVPN proprietary license.

## Other possible upstream materials

The application backend is built with Go (including Go's standard library).
Any upstream standard-library components incorporated into binaries retain
the licensing terms applicable to those components; Go's licensing
information is available at https://go.dev/LICENSE.

Names such as LG, webOS, Homebrew Channel, Xray and QR Code belong to
their respective owners; their use is descriptive only.

## Publication checklist

Before changing any third-party version, verify its exact upstream terms,
notices and source availability, and include required notices in the shipped
binary package. The published v2.2.5 IPK was produced before the new
project-owned proprietary LICENSE was added to GitHub. Adding this document
does **not** silently alter or republish the already released IPK.

This file is an inventory based on the currently checked bundled components,
not a guarantee that every potential upstream right or dependency has been
independently audited.
