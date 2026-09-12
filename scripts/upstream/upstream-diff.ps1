<#
.SYNOPSIS
    Dumps upstream changelog and help to compare with fixtures.
#>
[CmdletBinding()]
param()

$ErrorActionPreference = "Stop"

Write-Host "Fetching upstream changelog..." -ForegroundColor Cyan
$changelog = & agy.exe changelog 2>&1
$reportDate = (Get-Date).ToString("yyyy-MM-dd")
$reportDir = Join-Path $PSScriptRoot "..\..\docs\upstream\reports"
$reportFile = Join-Path $reportDir "$reportDate-agy-update.md"

if (-not (Test-Path $reportDir)) {
    New-Item -ItemType Directory -Force -Path $reportDir | Out-Null
}

$reportContent = @"
# Upstream Audit Report - $reportDate

## Installed AGY Version
$(& agy.exe --version 2>&1)

## Upstream Changelog Summary
$changelog
"@

Set-Content -Path $reportFile -Value $reportContent
Write-Host "Saved upstream report to: $reportFile" -ForegroundColor Green
