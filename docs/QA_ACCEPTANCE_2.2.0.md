# CVPN webOS 2.2.0 — hardware acceptance and regression plan

## Release decision

CI passing means **release candidate only**. A shipping production build requires
the full **fresh-device** path with no SSH/manual commands and documented evidence
of network rollback. Do not experiment with boot/restart on a TV known to have
power-on difficulties. Do not change either previously working LG.

## Test matrix

| Scenario | Previously working LGs | Third LG (intermittent WLAN) | Fourth LG (fresh) |
| --- | --- | --- | --- |
| Existing persistent VPN remains unchanged | Preserve; no updates | Observed on 2.1.3 | N/A |
| Fresh single-IPK bootstrap without SSH | N/A | Upgrade only, manual bootstrap done | **Required** |
| Launcher transitions and retry/error states | Offline mocks only | Required after upgrade | **Required** |
| VLESS/Reality + subscription | Existing smoke | Observed YouTube success on 2.1.3 | **Required** |
| HTTPS DNS, 1080p/4K streaming | Leave alone | Observed YouTube success | **Required** |
| Disconnect and full network recovery | Previously used | Required | **Required** |
| Standby / safe restart | No changes | Postpone due device power issue | **Required** |
| Wi-Fi loss and restore / Ethernet | No changes | Inherent flaps observed | **Required** |

## 1. Fresh installation — no developer intervention

**Prerequisite:** rooted LG with Homebrew Channel already functioning, supported
ARMv7/TUN and stable LAN connection. No CVPN backend, binaries, private keys,
subscriptions, or prior CVPN app installed.

- Install **only** `su.cvpn.webos_2.2.0_all.ipk`.
- Open CVPN from TV home screen (no SSH). Expect the centered preparation screen
  to show: checking device -> preparing components -> launching CVPN -> ready.
  The screen must never show internal filenames, root commands, Xray, or fake percentages.
- Do not switch off the TV during first-run provisioning.
- Initialization may finish only when the installer has verified IPK SHA256 hashes
  and the local backend authenticated `/api/status` health.
- Relaunch after completion: setup must be idempotent and must not start another
  copy/restart of Xray or panel.
- Launch QR, pair a real smartphone on trusted LAN, import a subscription,
  select a server, then connect without command line.
- Record the complete UI and whether app-to-Homebrew `/exec` succeeds. A
  root SSH `luna-send -f` success **alone** is insufficient: Luna caller
  permissions differ on webOS. If denied, stop and inspect permission mapping;
  never patch service permissions or elevate components automatically.

## 2. Setup UI and fault handling

- Denied Homebrew permission: human-friendly message, Retry and Close; no spinner
  after a 20-second status response deadline.
- Missing root, Homebrew, TUN, or unsupported CPU: actionable message; no
  attempted modification to network routes, DNS or system partitions.
- Corrupt/modified IPK payload: installer refuses before copying any binaries.
- Panel start failure or component copy failure: original files restored;
  private configuration, profiles, secrets and access key retained.
- Long running install: bounded **90-second UI wait**; if timeout, user may
  **check status** only. No second install while original owns Engine lock.
- Close/reopen while install is running: no duplicate installer; stage resumes
  from manifest/stamp readiness check.
- Press Home/back during progress and relaunch. No stray overlays or visible
  secret/error stack traces; QR relaunch does not get closed by stale callbacks.
- Test existing VPN active while upgrading. Must refuse component replacement
  and instruct user to disconnect; existing VPN must remain online.
- Check old 720p/1080p TV viewport, readable 16px+ text, D-pad button focus,
  reduced motion and no broken typography, screenshot comparison against CVPN.

## 3. VPN/network functionality

- Valid VLESS/Reality subscription, multi-server selection, saved credentials
  retained after app update. Verify malformed, empty, duplicated and expired
  entries are handled without panic or secret leakage.
- First connection: trial -> ready -> confirm -> persistent, YouTube 1080p/4K
  playback, system HTTPS and DNS status; idle 15+ minutes; multiple reconnects.
- Server switching: preflight new node before stopping old; invalid new node
  cannot disconnect functioning VPN.
- Explicit disconnect: routes `0.0.0.0/1`, `128.0.0.0/1`, per-DNS route, bind
  mount and Xray process fully cleaned. Original system resolver and gateway return.
- Offline server, TLS failure, failed DNS/HTTPS, unsupported IPv6 default,
  Xray crash, blocked UDP, DNS delays and Wi-Fi flap: fail closed as designed
  and restore original network. Do not interpret initial HTTPS success as a
  streaming/UDP guarantee.
- Watchdog: repeated failures restore ordinary internet; no automatic boot
  VPN start; note that there is **no universal kill switch**.
- Wi-Fi -> Ethernet changes must be tested with VPN disconnected first, never
  manually delete routes/mounts under live VPN.
- Verify pairing cookies, CSRF, local-LAN access restrictions and that control
  port 8765 cannot be reached from the public internet.

## 4. Lifecycle and publication

- Clean LG standby and power-cycle: boots normally; root/Homebrew and local
  control panel return; VPN does not autostart; internet not trapped in stale
  TUN/DNS routes.
- Verify install, **upgrade, uninstall/reinstall**. Note persistent data may
  remain after uninstall; document approved cleanup flow before public release.
- Confirm that app permissions include `lunabus.query` and that the installed
  Homebrew version has a compatible service ACG mapping; old mapping may refuse
  calls on certain firmware.
- Inspect IPK contents for secrets/keys, third-party notices, source/version
  hashes, and compile/static reports. Complete owner testing screenshots.
- Log device model, webOS firmware, Homebrew version, CPU, network type,
  pass/fail, safe rollback observations, and reproduction steps; omit tokens,
  URLs carrying credentials, IPs from private endpoints and access keys.

## Stop conditions

**Do not publish** if fresh first launch needs SSH, root permission is denied,
an installer can run concurrently, the UI can hang indefinitely, any failure
leaves broken DNS/routes, or firmware/root boot becomes unstable.

### QA evidence template

| Field | Result |
| --- | --- |
| LG model / webOS / Homebrew | |
| Clean IPK first-run success without SSH | |
| Second start is idempotent | |
| Pairing / import / connect / persist | |
| YouTube 1080p and 4K | |
| Wi-Fi flaps / Ethernet | |
| Disconnect and route restore | |
| Standby, restart, recovery | |
| Root service permission | |
| Uninstall/reinstall | |
| Outstanding issues / release decision | |
