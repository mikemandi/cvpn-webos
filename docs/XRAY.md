# Xray-core binary prerequisite

This repository does **not** redistribute Xray-core. Use the upstream [XTLS/Xray-core releases](https://github.com/XTLS/Xray-core/releases) and choose an ARM32 Linux build appropriate for the target webOS **userspace** (even where `uname -m` reports `aarch64` but `getconf LONG_BIT` reports `32`). Verify its upstream checksum and license.

For a fresh TV, pass the extracted executable file via `-XrayPath` (PowerShell) or the third positional argument to `scripts/setup.sh` (Linux). The installer transfers the local file via SSH, checks its uploaded SHA-256 in a separate payload manifest, and executes `xray version` on the TV to check basic compatibility.

Running `xray version` successfully does **not** guarantee full TUN support or protocol compatibility. Do not assume new Xray versions work on all LG firmware revisions without device testing. The previously installed and tested Xray binary is deliberately not replaced by the root installer on existing devices (which are rejected by this installer).
