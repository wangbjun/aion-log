$ErrorActionPreference = "Stop"

$root = Resolve-Path (Join-Path $PSScriptRoot "..")
$desktop = Join-Path $root "desktop"
$bin = Join-Path $desktop "bin"
$frontendDist = Join-Path $root "frontend\dist"

function Invoke-Native {
  param(
    [Parameter(Mandatory = $true)]
    [string]$FilePath,
    [Parameter(Mandatory = $true)]
    [string[]]$Arguments
  )

  & $FilePath @Arguments
  if ($LASTEXITCODE -ne 0) {
    throw "$FilePath exited with code $LASTEXITCODE"
  }
}

New-Item -ItemType Directory -Force -Path $bin | Out-Null

Push-Location $root
Invoke-Native -FilePath "go" -Arguments @("build", "-o", (Join-Path $bin "aion.exe"), "main.go")
Pop-Location

Push-Location (Join-Path $root "frontend")
$env:NODE_OPTIONS = "--openssl-legacy-provider"
Invoke-Native -FilePath "npm" -Arguments @("install", "--legacy-peer-deps")
Invoke-Native -FilePath "npm" -Arguments @("run", "build")
Pop-Location

$targetFrontend = Join-Path $bin "frontend"
if (Test-Path $targetFrontend) {
  Remove-Item $targetFrontend -Recurse -Force
}
New-Item -ItemType Directory -Force -Path $targetFrontend | Out-Null
Copy-Item $frontendDist (Join-Path $targetFrontend "dist") -Recurse -Force

$appIni = Join-Path $root "app.ini"
if (Test-Path $appIni) {
  Copy-Item $appIni (Join-Path $bin "app.ini") -Force
}
Copy-Item (Join-Path $root "aion.db") (Join-Path $bin "aion.db") -Force

$targetStorage = Join-Path $bin "storage"
if (Test-Path $targetStorage) {
  Remove-Item $targetStorage -Recurse -Force
}
Copy-Item (Join-Path $root "storage") $targetStorage -Recurse -Force

Push-Location $desktop
Invoke-Native -FilePath "npm" -Arguments @("install")
$env:ELECTRON_MIRROR = "https://npmmirror.com/mirrors/electron/"
$env:ELECTRON_BUILDER_BINARIES_MIRROR = "https://npmmirror.com/mirrors/electron-builder-binaries/"
Invoke-Native -FilePath "npm" -Arguments @("run", "dist")
Pop-Location

Write-Host "Desktop installer has been generated under desktop\release."
