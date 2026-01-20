# Veta LIS - Dependency Installer
# Run as Administrator

Write-Host "Installing dependencies..." -ForegroundColor Cyan

# Check Admin
if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Host "ERROR: Must run as Administrator" -ForegroundColor Red
    exit 1
}

Write-Host "Running as Administrator - OK" -ForegroundColor Green
Write-Host ""

# Paths
$goPath = "C:\Program Files\Go"
$nodePath = "C:\Program Files\nodejs"
$tempDir = "$env:TEMP\veta_install"

if (-not (Test-Path $tempDir)) {
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
}

# Step 1: Install Go
Write-Host "Step 1: Installing Go 1.21.5..." -ForegroundColor Yellow

if (Test-Path "$goPath\bin\go.exe") {
    Write-Host "Go already installed" -ForegroundColor Green
    & "$goPath\bin\go.exe" version
} else {
    Write-Host "Downloading Go..." -ForegroundColor Cyan
    $ProgressPreference = 'SilentlyContinue'
    Invoke-WebRequest -Uri "https://dl.google.com/go/go1.21.5.windows-amd64.zip" -OutFile "$tempDir\go.zip"
    
    Write-Host "Extracting Go..." -ForegroundColor Cyan
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    
    if (Test-Path $goPath) {
        Remove-Item -Path $goPath -Recurse -Force
    }
    
    [System.IO.Compression.ZipFile]::ExtractToDirectory("$tempDir\go.zip", "C:\Program Files")
    
    # Update PATH
    $currentPath = [System.Environment]::GetEnvironmentVariable("Path", "Machine")
    if ($currentPath -notlike "*$goPath\bin*") {
        $newPath = "$goPath\bin;$currentPath"
        [System.Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")
        $env:Path = $newPath
    }
    
    Write-Host "Go installed" -ForegroundColor Green
}

Write-Host ""

# Step 2: Install Node.js
Write-Host "Step 2: Installing Node.js 20.10.0..." -ForegroundColor Yellow

if ((Test-Path "$nodePath\node.exe") -or (Get-Command node -ErrorAction SilentlyContinue)) {
    Write-Host "Node.js already installed" -ForegroundColor Green
    node --version
} else {
    Write-Host "Downloading Node.js..." -ForegroundColor Cyan
    $ProgressPreference = 'SilentlyContinue'
    Invoke-WebRequest -Uri "https://nodejs.org/dist/v20.10.0/node-v20.10.0-win-x64.zip" -OutFile "$tempDir\node.zip"
    
    Write-Host "Extracting Node.js..." -ForegroundColor Cyan
    
    $extractPath = "$tempDir\node_extract"
    if (Test-Path $extractPath) {
        Remove-Item -Path $extractPath -Recurse -Force
    }
    
    [System.IO.Compression.ZipFile]::ExtractToDirectory("$tempDir\node.zip", $extractPath)
    
    $nodeFolder = Get-ChildItem -Path $extractPath -Directory | Select-Object -First 1
    
    if (Test-Path $nodePath) {
        Remove-Item -Path $nodePath -Recurse -Force
    }
    
    Move-Item -Path "$extractPath\$($nodeFolder.Name)" -Destination $nodePath -Force
    
    $currentPath = [System.Environment]::GetEnvironmentVariable("Path", "Machine")
    if ($currentPath -notlike "*$nodePath*") {
        $newPath = "$nodePath;$currentPath"
        [System.Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")
        $env:Path = $newPath
    }
    
    Write-Host "Node.js installed" -ForegroundColor Green
}

Write-Host ""

# Refresh PATH
$env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")

# Step 3: Verify
Write-Host "Step 3: Verifying installations..." -ForegroundColor Yellow
Write-Host ""

go version
node --version
npm --version

Write-Host ""
Write-Host "Step 4: Installing Wails CLI..." -ForegroundColor Yellow

go install github.com/wailsapp/wails/v2/cmd/wails@latest

Write-Host ""
Write-Host "SUCCESS! All dependencies installed." -ForegroundColor Green
Write-Host ""

# Cleanup
Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue

Read-Host "Press Enter to exit"
