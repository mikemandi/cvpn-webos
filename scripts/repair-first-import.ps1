# Repairs Go control panel only on an existing rooted TV. Windows PowerShell 5.1+
param(
    [Parameter(Mandatory=$true)][ValidatePattern('^[a-zA-Z0-9][a-zA-Z0-9.-]*$')][string]$TV,
    [ValidateSet('check','install','rollback')][string]$Action='check'
)
$ErrorActionPreference='Stop'
$root=Split-Path -Parent $PSScriptRoot
$hostAddress='root@'+$TV
$scriptPath=Join-Path $PSScriptRoot 'repair-first-import.sh'
if (-not (Test-Path -LiteralPath $scriptPath)) { throw 'repair-first-import.sh missing' }
& scp.exe $scriptPath ($hostAddress+':/tmp/cvpn-first-import-fix.sh')
if ($LASTEXITCODE -ne 0) { throw 'SSH upload failed' }
if ($Action -eq 'install') {
    & ssh.exe $hostAddress 'sh /tmp/cvpn-first-import-fix.sh check'
    if ($LASTEXITCODE -ne 0) { throw 'Preflight failed; binary was not touched' }
    $binary=Join-Path $root 'dist/cvpn-webos-panel-armv7'
    if (-not (Test-Path -LiteralPath $binary -PathType Leaf)) { throw "Missing ARMv7 binary: $binary" }
    $hash=(Get-FileHash -Algorithm SHA256 -LiteralPath $binary).Hash.ToLowerInvariant()
    $file=Join-Path $env:TEMP 'cvpn-first-import-sha256.txt'
    [IO.File]::WriteAllText($file,"$hash  panel`n",[Text.Encoding]::ASCII)
    & ssh.exe $hostAddress 'mkdir -p /tmp/cvpn-first-import-fix && chmod 700 /tmp/cvpn-first-import-fix'
    if ($LASTEXITCODE -ne 0) { throw 'Cannot create remote staging' }
    & scp.exe $binary ($hostAddress+':/tmp/cvpn-first-import-fix/panel')
    if ($LASTEXITCODE -ne 0) { throw 'Binary upload failed' }
    & scp.exe $file ($hostAddress+':/tmp/cvpn-first-import-fix/SHA256SUMS')
    if ($LASTEXITCODE -ne 0) { throw 'SHA256 upload failed' }
}
& ssh.exe $hostAddress ('sh /tmp/cvpn-first-import-fix.sh '+$Action)
if ($LASTEXITCODE -ne 0) { throw "Remote repair $Action failed" }
