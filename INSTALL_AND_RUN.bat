@echo off
echo ============================================================
echo VETA LIS - Complete Installation and Run Script
echo ============================================================
echo.

cd /d "%~dp0"

REM Refresh PATH to include Go and Node.js
set "PATH=%PATH%;C:\Program Files\Go\bin;C:\Program Files\nodejs;%USERPROFILE%\go\bin"

echo Step 1: Verifying Go installation...
where go >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Go is not installed or not in PATH
    echo Please ensure Go 1.21+ is installed from https://go.dev/dl/
    pause
    exit /b 1
)
go version
echo.

echo Step 2: Verifying Node.js installation...
where node >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Node.js is not installed or not in PATH
    echo Please ensure Node.js 18+ is installed from https://nodejs.org/
    pause
    exit /b 1
)
node --version
npm --version
echo.

echo Step 3: Installing Wails CLI...
go install github.com/wailsapp/wails/v2/cmd/wails@latest
echo Wails CLI installed
echo.

echo Step 4: Installing Go dependencies...
go mod download
echo Go dependencies installed
echo.

echo Step 5: Installing frontend dependencies...
cd frontend
call npm install
cd ..
echo Frontend dependencies installed
echo.

echo Step 6: Generating Wails bindings...
%USERPROFILE%\go\bin\wails generate module
echo.

echo ============================================================
echo Starting Veta LIS Development Server...
echo ============================================================
echo.
echo The application will open in your default browser
echo Press Ctrl+C to stop the server
echo.

%USERPROFILE%\go\bin\wails dev

pause
