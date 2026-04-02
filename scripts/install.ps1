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

# ZTVS Windows installer
# Downloads the latest release for your platform, installs the zt binary,
# and seeds first-party plugins into %USERPROFILE%\.ztvs\plugins\.
#
# One-liner (PowerShell):
#   irm https://raw.githubusercontent.com/mosesgameli/ztvs/main/install.ps1 | iex
#
# Pin a specific version:
#   $env:ZTVS_VERSION="v1.0.0"; irm https://raw.githubusercontent.com/mosesgameli/ztvs/main/install.ps1 | iex

[CmdletBinding()]
param()

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

# ── config ─────────────────────────────────────────────────────────────────
$Repo      = "mosesgameli/ztvs"
$InstallDir = if ($env:ZTVS_INSTALL_BIN) { $env:ZTVS_INSTALL_BIN } else { "$env:LOCALAPPDATA\Programs\ztvs" }
$ConfigDir  = if ($env:ZTVS_HOME) { $env:ZTVS_HOME } else { "$env:USERPROFILE\.ztvs" }
$PluginDir  = Join-Path $ConfigDir "plugins"

# ── helpers ─────────────────────────────────────────────────────────────────
function Write-Step  { param($msg) Write-Host "  -> $msg" -ForegroundColor Cyan   }
function Write-Ok    { param($msg) Write-Host "  v  $msg" -ForegroundColor Green  }
function Write-Fatal { param($msg) Write-Host "`nerror: $msg`n" -ForegroundColor Red; exit 1 }

function Get-LatestVersion {
    $url      = "https://api.github.com/repos/$Repo/releases/latest"
    $headers  = @{ "User-Agent" = "ztvs-installer" }
    $response = Invoke-RestMethod -Uri $url -Headers $headers -UseBasicParsing
    return $response.tag_name
}

function Get-Architecture {
    $hw = (Get-CimInstance -ClassName Win32_Processor).Architecture
    # 5 = ARM64, everything else we treat as amd64
    if ($hw -eq 5) { return "arm64" } else { return "amd64" }
}

function Add-ToUserPath {
    param([string]$Dir)
    $current = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($current -notlike "*$Dir*") {
        [Environment]::SetEnvironmentVariable("PATH", "$current;$Dir", "User")
        $env:PATH += ";$Dir"
        Write-Step "Added $Dir to your user PATH"
    }
}

# ── banner ──────────────────────────────────────────────────────────────────
Write-Host ""
Write-Host "  ███████╗████████╗██╗   ██╗███████╗" -ForegroundColor Magenta
Write-Host "     ███╔╝╚══██╔══╝██║   ██║██╔════╝" -ForegroundColor Magenta
Write-Host "    ███╔╝    ██║   ██║   ██║███████╗" -ForegroundColor Magenta
Write-Host "   ███╔╝     ██║   ╚██╗ ██╔╝╚════██║" -ForegroundColor Magenta
Write-Host "  ███████╗   ██║    ╚████╔╝ ███████║" -ForegroundColor Magenta
Write-Host "  ╚══════╝   ╚═╝     ╚═══╝  ╚══════╝" -ForegroundColor Magenta
Write-Host "  Zero Trust Vulnerability Scanner`n"

# ── resolve version & arch ───────────────────────────────────────────────────
$Version = if ($env:ZTVS_VERSION) { $env:ZTVS_VERSION } else { Get-LatestVersion }
if (-not $Version) { Write-Fatal "could not resolve latest release version" }

$Arch    = Get-Architecture
$OS      = "windows"
$Archive = "ztvs-$Version-$OS-$Arch.zip"
$DownloadUrl = "https://github.com/$Repo/releases/download/$Version/$Archive"

Write-Step "Detected platform : $OS/$Arch"
Write-Step "Installing version : $Version"

# ── download ─────────────────────────────────────────────────────────────────
$TmpDir = Join-Path $env:TEMP "ztvs-install-$(Get-Random)"
New-Item -ItemType Directory -Path $TmpDir | Out-Null

$ArchivePath = Join-Path $TmpDir $Archive
Write-Step "Downloading $Archive ..."
try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $ArchivePath -UseBasicParsing
} catch {
    Write-Fatal "download failed — check that release $Version exists for $OS/$Arch`n  URL: $DownloadUrl"
}

# ── extract ───────────────────────────────────────────────────────────────────
Write-Step "Extracting ..."
$ExtractDir = Join-Path $TmpDir "extracted"
Expand-Archive -Path $ArchivePath -DestinationPath $ExtractDir -Force

# ── install host binary ───────────────────────────────────────────────────────
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}
$BinSrc = Join-Path $ExtractDir "zt.exe"
if (-not (Test-Path $BinSrc)) {
    Write-Fatal "zt.exe not found in archive — the release may be malformed"
}
Write-Step "Installing zt.exe -> $InstallDir\"
Copy-Item $BinSrc -Destination "$InstallDir\zt.exe" -Force

# ── install plugins ───────────────────────────────────────────────────────────
if (-not (Test-Path $PluginDir)) {
    New-Item -ItemType Directory -Path $PluginDir | Out-Null
}

foreach ($plugin in @("plugin-os", "plugin-axios-mitigation")) {
    $pluginBin  = Join-Path $ExtractDir "$plugin.exe"
    $pluginYaml = Join-Path $ExtractDir "$plugin.yaml"
    if (Test-Path $pluginBin) {
        Write-Step "Installing plugin: $plugin"
        Copy-Item $pluginBin  -Destination "$PluginDir\$plugin.exe"  -Force
        if (Test-Path $pluginYaml) {
            Copy-Item $pluginYaml -Destination "$PluginDir\$plugin.yaml" -Force
        }
    }
}

# ── bootstrap config ──────────────────────────────────────────────────────────
if (-not (Test-Path "$ConfigDir\config.yaml")) {
    Write-Step "Bootstrapping $ConfigDir\config.yaml ..."
    try { & "$InstallDir\zt.exe" plugin init 2>$null } catch {}
}

# ── PATH ──────────────────────────────────────────────────────────────────────
Add-ToUserPath $InstallDir

# ── cleanup ───────────────────────────────────────────────────────────────────
Remove-Item $TmpDir -Recurse -Force -ErrorAction SilentlyContinue

# ── done ──────────────────────────────────────────────────────────────────────
Write-Host ""
Write-Ok "ZTVS $Version installed successfully!"
Write-Step "Run 'zt scan' to start your first audit."
Write-Host "  (You may need to restart your terminal for PATH changes to take effect.)`n"
