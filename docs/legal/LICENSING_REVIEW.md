# CVPN webOS — Licensing and publishing review

[English](LICENSING_REVIEW.md) | [Русский](LICENSING_REVIEW.ru.md)

Reviewed: **2026-10-28** · Published binary: **v2.2.5**

## Chosen license

CVPN webOS now uses a **custom proprietary license with free end-user use**, documented in [LICENSE](../../LICENSE) ([Russian translation](CVPN_TERMS_RU.md)).

Users may download, install, update and run official distributions at no cost, keep private backups, and inspect/build the published CVPN-owned code for their own devices. They may use **any compatible VPN provider** without purchasing a CVPN subscription. Public redistribution, rebranding or resale of CVPN-owned portions requires written authorization. Valid GitHub forks, accurate reviews and official catalog imagery are allowed within the stated limits.

**The project's own code is not OSI open source. Homebrew should retain `pool: non-free`, which does not mean the app costs money.** See [brand policy](../../BRAND_POLICY.md) and [third-party notices](../../THIRD_PARTY_NOTICES.md).

This grant applies **only to content that the relevant CVPN rightsholders lawfully control**. GitHub visibility and choosing a license do not establish exclusive rights in contributions, third-party material, generated content or all images. Contributor and asset provenance still need confirmation.

## Components with independent licenses

| Component | Terms | Required preservation |
|---|---|---|
| Xray-core **v26.10.21** | **MPL-2.0** — https://github.com/XTLS/Xray-core/tree/v26.10.21 | Keep source access, license and notices for the MPL-covered Xray code; CVPN proprietary terms cannot limit MPL rights |
| QRCode for JavaScript © 2009 Kazuhiko Arase | **MIT** | Copyright and complete MIT permission text |
| qrcode-terminal 0.12.0 adaptation | **Apache-2.0** and MIT upstream notices | Apache and MIT notices in launcher/QR-THIRD-PARTY-LICENSE.txt |
| Go standard library in compiled backend | Go standard-library license: https://go.dev/LICENSE | Verify required notices in redistributed binaries |
| CVPN-owned original files/icons | **Custom proprietary LICENSE** and brand policy | Free use and narrow catalog permission; public rebranding/redistribution restricted |

The IPK's existing payload includes Xray's XRAY-LICENSE.txt and both QR notice files. The repository QR notice now also includes the full MIT permission text. This licensing review does **not** assert a complete independent code/asset audit. Mozilla's [MPL FAQ](https://www.mozilla.org/en-US/MPL/2.0/FAQ/) explains corresponding-source obligations and combination with separate proprietary code.

## Catalog, releases and remaining action

The official **v2.2.5 IPK was published before this project LICENSE and updated QR notice were added**. Updating the repository does **not** replace or alter that historical binary. A subsequent IPK should visibly bundle the applicable proprietary and third-party notices; review Go standard-library notices as part of that packaging change.

Homebrew can display accurate CVPN icons/screenshots and link to official releases under [LICENSE](../../LICENSE), while source/code rebranding remains restricted. The brand policy does not claim a registered trademark. Public GitHub forks remain allowed through the service under GitHub's Terms.

Before any broader distribution: verify who owns each authored contribution, artwork, and any AI-assisted/third-party material; maintain exact-version Xray source availability; keep Homebrew at non-free; prevent secrets from entering builds or screenshots; seek qualified legal review when necessary.

[README](../../README.md) · [Security](../SECURITY.md) · [Homebrew submission](../HOMEBREW_SUBMISSION.md)
