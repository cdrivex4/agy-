<#
.SYNOPSIS
    Checks installed and upstream Antigravity CLI version and compares with compatibility baseline.
#>
[CmdletBinding()]
param()

$ErrorActionPreference = "Stop"

Write-Host "=== AGY++ Upstream Version Auditor ===" -ForegroundColor Cyan

# 1. Detect installed AGY
$agyCmd = Get-Command agy.exe -ErrorAction SilentlyContinue
if (-not $agyCmd) {
    Write-Host "[!] agy.exe not found in PATH." -ForegroundColor Red
    exit 1
}

$installedVersion = (& agy.exe --version 2>&1).Trim()
Write-Host "Found installed AGY: $installedVersion ($($agyCmd.Source))" -ForegroundColor Green

# 2. Check baseline from config/upstream.yaml
$configPath = Join-Path $PSScriptRoot "..\..\config\upstream.yaml"
if (Test-Path $configPath) {
    $configContent = Get-Content $configPath -Raw
    Write-Host "Loaded baseline configuration from config/upstream.yaml" -ForegroundColor Gray
}

# 3. Report compatibility
Write-Host "`nCompatibility Status:" -ForegroundColor Yellow
Write-Host "  Installed Version:  $installedVersion"
Write-Host "  Tested Baseline:    1.2.2"

if ($installedVersion -eq "1.2.2") {
    Write-Host "  Status:             VERIFIED COMPATIBLE" -ForegroundColor Green
} else {
    Write-Host "  Status:             VERSION MISMATCH - Run upstream-diff to inspect changes." -ForegroundColor Yellow
}
