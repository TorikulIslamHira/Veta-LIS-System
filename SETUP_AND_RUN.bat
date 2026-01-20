@echo off
REM Veta LIS - Simple Setup Script (Batch)
REM This script will set up and run Veta LIS

setlocal enabledelayedexpansion

echo.
echo ========================================================
echo        VETA LIS - Setup and Run
echo ========================================================
echo.

REM Check if running as Administrator
net session >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: This script must be run as Administrator!
    echo Please right-click cmd.exe and select "Run as Administrator"
    pause
    exit /b 1
)

echo Checking system requirements...
echo.

REM Check Go
echo Checking Go installation...
if exist "C:\Program Files\Go\bin\go.exe" (
    echo [OK] Go is installed
    "C:\Program Files\Go\bin\go.exe" version
) else (
    echo [ERROR] Go is not installed
    echo Please download from: https://golang.org/dl/
    echo Then run this script again
    pause
    exit /b 1
)

echo.

REM Check Node.js
echo Checking Node.js installation...
if exist "C:\Program Files\nodejs\node.exe" (
    echo [OK] Node.js is installed
    "C:\Program Files\nodejs\node.exe" --version
) else (
    echo [ERROR] Node.js is not installed
    echo Please download from: https://nodejs.org/
    echo Then run this script again
    pause
    exit /b 1
)

echo.
echo Installing project dependencies...
echo.

REM Set paths
set PATH=%PATH%;C:\Program Files\Go\bin;C:\Program Files\nodejs

REM Download Go dependencies
echo [1/3] Downloading Go dependencies...
"C:\Program Files\Go\bin\go.exe" mod download
if %errorlevel% neq 0 (
    echo Warning: Go dependencies had issues
)

echo [2/3] Installing frontend dependencies...
cd frontend
"C:\Program Files\nodejs\npm.cmd" install --legacy-peer-deps
if %errorlevel% neq 0 (
    echo Warning: Frontend dependencies had issues
)
cd ..

echo.
echo ========================================================
echo        Starting Veta LIS Development Server
echo ========================================================
echo.
echo The application will open at http://localhost:3000
echo.

REM Try to find and run Wails
set WAILS_PATH=%GOPATH%\bin\wails.exe
if not exist "%WAILS_PATH%" (
    set WAILS_PATH=C:\Users\%USERNAME%\go\bin\wails.exe
)

if exist "%WAILS_PATH%" (
    "%WAILS_PATH%" dev
) else (
    echo Running via Go...
    "C:\Program Files\Go\bin\go.exe" run github.com/wailsapp/wails/v2/cmd/wails@latest dev
)

pause
