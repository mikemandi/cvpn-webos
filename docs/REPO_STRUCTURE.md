# Repository structure

- `server/` Go backend, embedded web/mobile interface, source and tests.
- `launcher/` webOS launcher IPK source including third-party QR library license.
- `engine/` Xray wrapper, network route/DNS modifications and recovery logic.
- `recovery/` boot-only recovery worker + hook.
- `scripts/` offline IPK builder, root first-time installer, Windows/Linux clients, manifest generators.
- `assets/` icon and older historical design previews (replace with real screenshots).
- `.github/workflows/` modern Go build, tests and release manifest publishing.
- `dist/` generated IPK and ARMv7 binary (NOT tracked by git).
- `docs/` safety, first-run, publication instructions and acceptance checklist.
