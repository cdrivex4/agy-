<#
.SYNOPSIS
    Builds agy++.exe with embedded version and commit information.
#>
[CmdletBinding()]
param(
    [string]$Output = "agy++.exe"
)

$ErrorActionPreference = "Stop"

$version = (Get-Content (Join-Path $PSScriptRoot "..\..\VERSION") -Raw).Trim()
$commit = (& git rev-parse --short HEAD 2>$null)
if (-not $commit) { $commit = "unknown" }
$date = (Get-Date).ToString("yyyy-MM-dd")

$ldflags = "-s -w -X 'github.com/cdrivex4/agy-plus-plus/internal/version.Version=$version' " +
           "-X 'github.com/cdrivex4/agy-plus-plus/internal/version.GitCommit=$commit' " +
           "-X 'github.com/cdrivex4/agy-plus-plus/internal/version.BuildDate=$date' " +
           "-X 'github.com/cdrivex4/agy-plus-plus/internal/version.BaselineAgy=1.2.2'"

Write-Host "Building AGY++ v$version ($commit)..." -ForegroundColor Cyan
go build -ldflags $ldflags -o $Output ./cmd/agy-plus-plus
Write-Host "Build complete: $Output" -ForegroundColor Green
