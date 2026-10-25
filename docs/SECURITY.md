# Security notes — CVPN webOS

**TV backend:** root-required. **Frontend:** local HTTP port 8765. **Distribution:** TV launcher IPK and separate root install.

- The control panel checks local-network sources, client subnet, Host, short PIN/QR pairing, sessions and CSRF. These do not provide transport encryption: malicious clients on shared or hostile Wi-Fi may intercept unencrypted traffic. Use only a trusted private LAN and NEVER forward TCP/8765 to the internet.
- Profiles, imported secrets and subscription URLs are stored on the local TV and protected primarily by filesystem permissions (mode 0600/0700), **not** encrypted at rest.
- Source and IPK contain **no** personal subscriptions, access keys, sessions, developer SSH keys or vendor credentials.
- The fresh install generates a device-specific random access key from `/dev/urandom`, and refuses to overwrite an existing CVPN installation.
- No automatic VPN at boot. IPv6 fails closed. No system-wide network kill switch: a VPN failure restores regular internet.
- Root shell scripts modify files in `/media/developer/cvpn-webos` and a symlink in Homebrew `init.d`; they do not edit the kernel/rootfs or install root access.
- Hashes in `SHA256SUMS` and the Homebrew manifest are integrity checks, **not an independent publisher signature**.
- Do not include `All.zip`, historic private keys, access logs or user profiles in public releases.

## Release review gates

1. Review permissions and source licensing with the owner.
2. Build backend with up-to-date supported Go (GitHub Actions pins Go 1.27.3), then inspect the SHA-256 digests of generated artifacts.
3. Test TV launcher IPK and root bootstrap on at least one fresh supported LG; verify pairing, importing a real subscription, VPN stop and cold reboot.
4. Obtain current screenshots from a real TV and phone, remove real addresses, tokens, PINs and identities.
5. Disclose AI-assistance in Homebrew PR as required by `webosbrew/apps-repo` policy.

Historical `RC5 -> Stable` scripts are not necessary for fresh installs and are intentionally excluded from this clean public-source repository.
