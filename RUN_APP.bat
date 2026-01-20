@echo off
REM Veta LIS - Run with Administrator
REM This batch file will check for admin and re-run if needed

setlocal
cd /d "%~dp0"

REM Check for administrator privileges
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo Requesting administrator privileges...
    powershell -Command "Start-Process cmd -ArgumentList '/c cd %~dp0 && %~0' -Verb RunAs"
    exit /b
)

echo.
echo ============================================================
echo VETA LIS - Setup and Installation
echo ============================================================
echo.

REM Step 1: Check Go
echo Step 1: Checking Go...
if exist "C:\Program Files\Go\bin\go.exe" (
    echo Go is installed
    "C:\Program Files\Go\bin\go.exe" version
) else (
    echo Installing Go 1.21.5...
    powershell -Command "$ProgressPreference = 'SilentlyContinue'; Invoke-WebRequest -Uri 'https://dl.google.com/go/go1.21.5.windows-amd64.zip' -OutFile '%TEMP%\go.zip'; Add-Type -AssemblyName System.IO.Compression.FileSystem; [System.IO.Compression.ZipFile]::ExtractToDirectory('%TEMP%\go.zip', 'C:\Program Files'); setx Path 'C:\Program Files\Go\bin;!Path!'"
    echo Go installed
)

echo.

REM Step 2: Check Node.js
echo Step 2: Checking Node.js...
if exist "C:\Program Files\nodejs\node.exe" (
    echo Node.js is installed
    "C:\Program Files\nodejs\node.exe" --version
) else (
    echo Installing Node.js 20.10.0...
    powershell -Command "$ProgressPreference = 'SilentlyContinue'; Invoke-WebRequest -Uri 'https://nodejs.org/dist/v20.10.0/node-v20.10.0-win-x64.zip' -OutFile '%TEMP%\node.zip'; Add-Type -AssemblyName System.IO.Compression.FileSystem; [System.IO.Compression.ZipFile]::ExtractToDirectory('%TEMP%\node.zip', '%TEMP%'); Move-Item -Path '%TEMP%\node-v20.10.0-win-x64' -Destination 'C:\Program Files\nodejs' -Force"
    echo Node.js installed
)

echo.

REM Refresh PATH
echo Step 3: Refreshing PATH...
set "PATH=C:\Program Files\Go\bin;C:\Program Files\nodejs;%PATH%"

REM Step 4: Install Wails
echo Step 4: Installing Wails CLI...
call go install github.com/wailsapp/wails/v2/cmd/wails@latest

echo.
echo Step 5: Installing project dependencies...
cd /d "%~dp0"

REM Install frontend dependencies
echo Installing frontend packages...
cd frontend
call npm install
cd ..

echo.
echo ============================================================
echo Setup complete! Starting Veta LIS...
echo ============================================================
echo.

REM Run the application
call wails dev

pause
