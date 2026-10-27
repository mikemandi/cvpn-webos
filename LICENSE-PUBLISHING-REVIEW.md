# Publishing / licensing review required

This repository contains CVPN-authored code, artwork, and a bundled third-party QR library.
The operator must confirm ownership and choose the license for CVPN source before assigning
`pool: main` in the Homebrew repo. Public availability of source without a license is NOT
an open-source license grant.

The QR implementation's Apache-2.0 license is shipped as `launcher/QR-THIRD-PARTY-LICENSE.txt`.
The 2.2.3 All-in-One **IPK includes an Xray-core ARMv7 binary** downloaded from
upstream and SHA256-verified during CI, together with `XRAY-LICENSE.txt`.
The upstream code is not vendored into this repository, but binary redistribution
rights and all required upstream copyright/license notices must be reviewed
before public release.

Do not upload old historical project archives containing `private/access.key` or real subscriptions.
