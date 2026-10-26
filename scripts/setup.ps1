# PowerShell 5.1+ / Windows 11 — new rooted TV with ARMv7 userspace (armv7l or compatible aarch64 kernel) only. Run in repo directory.
param(
  [Parameter(Mandatory=$true)][ValidatePattern('^[A-Za-z0-9][A-Za-z0-9.-]*$')][string]$TV,
  [ValidateSet('check','install','doctor')][string]$Action='check',
  [string]$XrayPath=''
)
$ErrorActionPreference = 'Stop'
$Root = Split-Path -Parent $PSScriptRoot
$HostName = 'root@' + $TV
function Remote([string]$command) {
  & ssh.exe $HostName $command
  if ($LASTEXITCODE -ne 0) { throw "SSH failed: $command" }
}
function Upload([string]$src,[string]$target) {
  & scp.exe $src ($HostName + ':' + $target)
  if ($LASTEXITCODE -ne 0) { throw "Upload failed: $src" }
}
Upload (Join-Path $PSScriptRoot 'install-root.sh') '/tmp/cvpn-root-setup.sh'
if ($Action -eq 'install') {
  Remote 'sh /tmp/cvpn-root-setup.sh check'
  Remote 'mkdir -p /tmp/cvpn-root-setup && chmod 700 /tmp/cvpn-root-setup'
  $items = @(
    @('dist/cvpn-webos-panel-armv7','panel'), @('engine/cvpn.sh','engine'),
    @('recovery/recovery-worker.sh','worker'), @('recovery/boot-hook.sh','boot-hook'),
    @('scripts/start-panel.sh','start-panel'), @('scripts/stop-panel.sh','stop-panel')
  )
  $hashLines = @()
  foreach ($pair in $items) {
    $path = Join-Path $Root $pair[0]
    if (-not (Test-Path -LiteralPath $path)) { throw "Missing payload: $path" }
    $hash = (Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash.ToLowerInvariant()
    $hashLines += "$hash  $($pair[1])"
    Upload $path ('/tmp/cvpn-root-setup/' + $pair[1])
  }
  if ($XrayPath) {
    if (-not (Test-Path -LiteralPath $XrayPath)) { throw 'XrayPath does not exist' }
    $xrayFullPath = (Resolve-Path $XrayPath).Path
    $xrayHash = (Get-FileHash -LiteralPath $xrayFullPath -Algorithm SHA256).Hash.ToLowerInvariant()
    $hashLines += "$xrayHash  xray"
    Upload $xrayFullPath '/tmp/cvpn-root-setup/xray'
  }
  $tmp = Join-Path $env:TEMP 'cvpn-payload.sha256'
  [IO.File]::WriteAllText($tmp,($hashLines -join "`n")+"`n",[Text.Encoding]::ASCII)
  Upload $tmp '/tmp/cvpn-root-setup/payload.sha256'
}
Remote ('sh /tmp/cvpn-root-setup.sh ' + $Action)
