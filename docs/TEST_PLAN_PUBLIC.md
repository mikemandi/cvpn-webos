# Public release acceptance checklist (not yet all executed)

- [x] Go 1.23.2 local unit tests and `go vet` pass (older toolchain).
- [x] Linux ARMv7 backend builds with Go 1.23.2 locally.
- [x] IPK is a valid Debian-format archive with `appinfo.json` and icons.
- [x] Local root package ZIP builds and opens.
- [ ] GitHub Actions Go 1.27.3 CI run for this exact repository.
- [ ] Test IPK 2.0.4 installation on physical LG.
- [ ] Test **new TV** fresh-root installer with a verified compatible Xray.
- [ ] Verify boot hook on cold reboot, no VPN auto-connect.
- [ ] Check 8765 local LAN access, QR and subscription secret redaction.
- [ ] Check abort/rollback on interrupted first install.
- [ ] Review code and licensing/branding rights for public release.
- [ ] Obtain current screenshots from actual device, free of private information.
- [ ] Publish GitHub Release and Homebrew listing PR.

The previously installed private Stable backend 2.0.0 on one LG has passed user tests, but **that does not satisfy new-device/bootstrap tests above**.
