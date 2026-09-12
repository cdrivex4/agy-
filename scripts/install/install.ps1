<#
.SYNOPSIS
    Installs or uninstalls AGY++ on Windows.
.PARAMETER Uninstall
    Removes the AGY++ binary and optional PATH entry while preserving user configuration and sessions.
#>
[CmdletBinding()]
param(
    [switch]$Uninstall
)

$ErrorActionPreference = "Stop"

$installDir = Join-Path $env:LOCALAPPDATA "Programs\AgyPlusPlus"
$targetExe = Join-Path $installDir "agy++.exe"

if ($Uninstall) {
    Write-Host "Uninstalling AGY++..." -ForegroundColor Cyan
    if (Test-Path $targetExe) {
        Remove-Item $targetExe -Force
        Write-Host "Removed $targetExe" -ForegroundColor Green
    }
    if ((Test-Path $installDir) -and ((Get-ChildItem $installDir).Count -eq 0)) {
        Remove-Item $installDir -Force
    }
    Write-Host "Note: User data in %APPDATA%\AgyPlusPlus has been preserved." -ForegroundColor Yellow
    Write-Host "AGY++ uninstall complete." -ForegroundColor Green
    exit 0
}

Write-Host "=== AGY++ Windows Installer ===" -ForegroundColor Cyan

# 1. Verify AGY CLI presence
$agy = Get-Command agy.exe -ErrorAction SilentlyContinue
if (-not $agy) {
    Write-Warning "Antigravity CLI (agy.exe) was not found in PATH. AGY++ requires AGY to run."
} else {
    Write-Host "Detected upstream Antigravity CLI at $($agy.Source)" -ForegroundColor Green
}

# 2. Source binary
$sourceExe = Join-Path $PSScriptRoot "..\..\agy++.exe"
if (-not (Test-Path $sourceExe)) {
    Write-Host "Building agy++.exe from source first..." -ForegroundColor Yellow
    $buildScript = Join-Path $PSScriptRoot "..\build\build.ps1"
    & powershell -ExecutionPolicy Bypass -File $buildScript
}

if (-not (Test-Path $sourceExe)) {
    Write-Error "Could not find or build agy++.exe at $sourceExe"
    exit 1
}

# 3. Install binary
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
}

Copy-Item -Path $sourceExe -Destination $targetExe -Force
Write-Host "Installed executable to: $targetExe" -ForegroundColor Green

# 4. Add to User PATH if missing
$userPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User)
$pathEntries = $userPath -split ";"
if ($pathEntries -notcontains $installDir) {
    Write-Host "Adding $installDir to User PATH..." -ForegroundColor Cyan
    $newPath = "$userPath;$installDir"
    [Environment]::SetEnvironmentVariable("Path", $newPath, [EnvironmentVariableTarget]::User)
    Write-Host "Updated User PATH. Please restart your terminal for PATH changes to take effect." -ForegroundColor Green
} else {
    Write-Host "$installDir is already in User PATH." -ForegroundColor Gray
}

Write-Host "`nInstallation successful! Run 'agy++ --version' or 'agy++' to get started." -ForegroundColor Green
