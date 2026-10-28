# CVPN webOS — Product description / Описание продукта

Short, accurate product copy for GitHub, Homebrew and release pages. The supported transports are confirmed by the VLESS importer in `server/main.go` and the gRPC validation in `server/grpc.go`.

## English

**One-line description (GitHub About)**

Xray/VLESS VPN for rooted LG webOS TVs · TCP, XHTTP & gRPC · import your configs and subscriptions · free, provider-independent.

**Short description (app listings)**

Xray/VLESS for rooted LG TVs — TCP, XHTTP and gRPC.

**Product introduction**

**CVPN webOS** brings **Xray/VLESS configurations to LG TVs with root**. Import VLESS links, compatible HTTPS subscriptions or supported Xray JSON, and connect over **TCP, XHTTP or gRPC**. Control the VPN from your phone — free to use, with no CVPN subscription required.

## Русский

**Одной строкой (описание GitHub)**

Xray/VLESS для LG webOS с root · TCP, XHTTP, gRPC · импорт своих конфигураций и подписок · бесплатно, без привязки к CVPN.

**Короткое описание (каталоги)**

Xray/VLESS для LG с root — TCP, XHTTP и gRPC.

**Описание продукта**

**CVPN webOS** — **Xray/VLESS-клиент для телевизоров LG с root**. Импортируй VLESS-ссылки, совместимые HTTPS-подписки или Xray JSON и подключайся через **TCP, XHTTP или gRPC**. Управляй VPN с телефона — бесплатно и без обязательной подписки CVPN.

## Accuracy and compatibility

- TCP, XHTTP and gRPC are **transport modes**, not three separate VPN protocols.
- This release supports **VLESS** configurations with compatible Reality settings; **gRPC also supports TLS**. It does **not** claim support for every Xray protocol or security configuration.
- Only compatible links, HTTPS subscription formats and supported Xray JSON fields can be imported.
- A TV with root, Homebrew root service, compatible 32-bit ARM userspace and TUN is required.
- No universal kill switch, VPN autoconnect or IPv6 tunneling is promised.
- These descriptions refer to code in the repository. The installed v2.2.5 IPK still contains its previously published app metadata until a new release is issued.

### Where this copy belongs

- `README.md` and `README.ru.md`: short, prominent headline and first paragraph.
- `homebrew/su.cvpn.webos.yml`: the Homebrew card and shortDescription.
- `homebrew/HOMEBREW_PR_BODY.md` and Homebrew submission drafts.
- `launcher/appinfo.json`: the bundled app description for **the next IPK build**.
- **GitHub About**: edit manually under the repository's About settings; GitHub's repository metadata is not a file and cannot be changed by editing README.
