# Veta LIS - Automated Setup and Run Script
# This script will:
# 1. Download and install Go
# 2. Download and install Node.js
# 3. Install Wails CLI
# 4. Install project dependencies
# 5. Start the development server

Write-Host "╔════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║        VETA LIS - Automated Setup and Run                 ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# Check if running as Administrator
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")

if (-not $isAdmin) {
    Write-Host "⚠️  This script needs to be run as Administrator to install software." -ForegroundColor Yellow
    Write-Host "Please right-click on PowerShell and select 'Run as Administrator'" -ForegroundColor Yellow
    exit 1
}

Write-Host "✓ Running with Administrator privileges" -ForegroundColor Green
Write-Host ""

# Step 1: Check and install Go
Write-Host "Step 1: Checking Go installation..." -ForegroundColor Yellow
if (Test-Path "C:\Program Files\Go\bin\go.exe") {
    Write-Host "✓ Go is already installed" -ForegroundColor Green
    & "C:\Program Files\Go\bin\go.exe" version
} else {
    Write-Host "⬇ Downloading Go 1.21..." -ForegroundColor Yellow
    $goUrl = "https://go.dev/dl/go1.21.0.windows-amd64.msi"
    $goMsi = "$env:TEMP\go1.21.0.windows-amd64.msi"
    
    Invoke-WebRequest -Uri $goUrl -OutFile $goMsi -ErrorAction SilentlyContinue
    
    if (Test-Path $goMsi) {
        Write-Host "📦 Installing Go..." -ForegroundColor Yellow
        Start-Process msiexec.exe -ArgumentList "/i `"$goMsi`" /quiet" -Wait
        Write-Host "✓ Go installed successfully" -ForegroundColor Green
        
        # Add Go to PATH
        $env:Path += ";C:\Program Files\Go\bin"
        [Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\Go\bin", "Machine")
        
        & "C:\Program Files\Go\bin\go.exe" version
    } else {
        Write-Host "✗ Failed to download Go" -ForegroundColor Red
        Write-Host "Please download from https://golang.org/dl/" -ForegroundColor Yellow
        exit 1
    }
}

Write-Host ""

# Step 2: Check and install Node.js
Write-Host "Step 2: Checking Node.js installation..." -ForegroundColor Yellow
if (Test-Path "C:\Program Files\nodejs\node.exe") {
    Write-Host "✓ Node.js is already installed" -ForegroundColor Green
    & "C:\Program Files\nodejs\node.exe" --version
    & "C:\Program Files\nodejs\npm.cmd" --version
} else {
    Write-Host "⬇ Downloading Node.js LTS..." -ForegroundColor Yellow
    $nodeUrl = "https://nodejs.org/dist/v20.10.0/node-v20.10.0-x64.msi"
    $nodeMsi = "$env:TEMP\node-v20.10.0-x64.msi"
    
    Invoke-WebRequest -Uri $nodeUrl -OutFile $nodeMsi -ErrorAction SilentlyContinue
    
    if (Test-Path $nodeMsi) {
        Write-Host "📦 Installing Node.js..." -ForegroundColor Yellow
        Start-Process msiexec.exe -ArgumentList "/i `"$nodeMsi`" /quiet" -Wait
        Write-Host "✓ Node.js installed successfully" -ForegroundColor Green
        
        # Add Node to PATH
        $env:Path += ";C:\Program Files\nodejs"
        [Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\nodejs", "Machine")
        
        & "C:\Program Files\nodejs\node.exe" --version
        & "C:\Program Files\nodejs\npm.cmd" --version
    } else {
        Write-Host "✗ Failed to download Node.js" -ForegroundColor Red
        Write-Host "Please download from https://nodejs.org/" -ForegroundColor Yellow
        exit 1
    }
}

Write-Host ""

# Step 3: Install Wails CLI
Write-Host "Step 3: Installing Wails CLI..." -ForegroundColor Yellow
Write-Host "📦 This may take a minute..." -ForegroundColor Yellow

$env:Path = $env:Path + ";C:\Program Files\Go\bin;C:\Program Files\nodejs"

& "C:\Program Files\Go\bin\go.exe" install github.com/wailsapp/wails/v2/cmd/wails@latest
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Wails CLI installed successfully" -ForegroundColor Green
} else {
    Write-Host "⚠ Wails CLI installation had issues, but continuing..." -ForegroundColor Yellow
}

Write-Host ""

# Step 4: Install project dependencies
Write-Host "Step 4: Installing project dependencies..." -ForegroundColor Yellow
Write-Host "Current directory: $(Get-Location)" -ForegroundColor Cyan

Write-Host "📦 Downloading Go dependencies..." -ForegroundColor Yellow
& "C:\Program Files\Go\bin\go.exe" mod download
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Go dependencies installed" -ForegroundColor Green
} else {
    Write-Host "⚠ Go dependencies had issues" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "📦 Installing frontend dependencies..." -ForegroundColor Yellow
cd frontend
& "C:\Program Files\nodejs\npm.cmd" install --legacy-peer-deps
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Frontend dependencies installed" -ForegroundColor Green
} else {
    Write-Host "⚠ Frontend dependencies had issues" -ForegroundColor Yellow
}
cd ..

Write-Host ""

# Step 5: Start Wails development server
Write-Host "╔════════════════════════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║         Starting Veta LIS Development Server              ║" -ForegroundColor Green
Write-Host "╚════════════════════════════════════════════════════════════╝" -ForegroundColor Green
Write-Host ""
Write-Host "🚀 Starting Wails development server..." -ForegroundColor Cyan
Write-Host "   The application will open at http://localhost:3000" -ForegroundColor Cyan
Write-Host ""

$wailsPath = "$env:GOPATH\bin\wails.exe"
if (-not (Test-Path $wailsPath)) {
    $wailsPath = "C:\Users\$env:USERNAME\go\bin\wails.exe"
}

if (Test-Path $wailsPath) {
    & $wailsPath dev
} else {
    Write-Host "⚠ Wails executable not found, trying alternative method..." -ForegroundColor Yellow
    & "C:\Program Files\Go\bin\go.exe" run github.com/wailsapp/wails/v2/cmd/wails@latest dev
}
