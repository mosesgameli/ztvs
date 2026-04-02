# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# ZTVS Windows uninstaller
# Removes the zt binary and %USERPROFILE%\.ztvs configuration directory.
#
# One-liner (PowerShell):
#   irm https://raw.githubusercontent.com/mosesgameli/ztvs/main/uninstall.ps1 | iex

[CmdletBinding()]
param()

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

# ── config ─────────────────────────────────────────────────────────────────
$InstallDir = "$env:LOCALAPPDATA\Programs\ztvs"
$ConfigDir  = "$env:USERPROFILE\.ztvs"

# ── helpers ─────────────────────────────────────────────────────────────────
function Write-Step { param($msg) Write-Host "  -> $msg" -ForegroundColor Cyan }
function Write-Ok   { param($msg) Write-Host "  v  $msg" -ForegroundColor Green }

function Remove-FromUserPath {
    param([string]$Dir)
    $current = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($current -like "*$Dir*") {
        # Filter out the ztvs directory from PATH
        $paths = $current -split ";"
        $newPaths = @()
        foreach ($p in $paths) {
            if ($p -and $p -ne $Dir) {
                $newPaths += $p
            }
        }
        $newPathStr = $newPaths -join ";"
        [Environment]::SetEnvironmentVariable("PATH", $newPathStr, "User")
        Write-Step "Removed $Dir from your user PATH"
    }
}

# ── main ───────────────────────────────────────────────────────────────────
Write-Host ""
Write-Host "  Uninstalling Zero Trust Vulnerability Scanner..." -ForegroundColor Magenta
Write-Host ""

# ── remove binaries ────────────────────────────────────────────────────────
if (Test-Path $InstallDir) {
    Write-Step "Removing installation directory $InstallDir ..."
    Remove-Item -Path $InstallDir -Recurse -Force
} else {
    Write-Step "Installation directory $InstallDir not found. Skipping."
}

# ── remove config & plugins ────────────────────────────────────────────────
if (Test-Path $ConfigDir) {
    Write-Step "Removing config directory $ConfigDir ..."
    Remove-Item -Path $ConfigDir -Recurse -Force
} else {
    Write-Step "Config directory $ConfigDir not found. Skipping."
}

# ── remove PATH entry ──────────────────────────────────────────────────────
Remove-FromUserPath $InstallDir

Write-Host ""
Write-Ok "ZTVS has been completely uninstalled from your system."
Write-Host "  (You may need to restart your terminal for PATH changes to take effect.)`n"
