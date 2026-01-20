# Virtual COM Port Setup for Windows

## Quick Start - Install com0com (FREE)

### Step 1: Download and Install
1. Download com0com from: https://sourceforge.net/projects/com0com/
2. Run the installer as **Administrator**
3. During installation, you may see Windows security warnings - click "Install" to proceed

### Step 2: Create Virtual Port Pair
After installation, the setup wizard will automatically create a pair. If not:

1. Open **Command Prompt as Administrator**
2. Navigate to com0com directory:
   ```cmd
   cd "C:\Program Files (x86)\com0com"
   ```

3. Run the setup command:
   ```cmd
   setupc install PortName=COM3 PortName=COM4
   ```

4. Verify the ports:
   ```cmd
   setupc list
   ```

### Step 3: Verify in Device Manager
1. Open **Device Manager** (Win + X → Device Manager)
2. Look for "com0com - serial port emulator" section
3. You should see: COM3 ↔ COM4 pair

### Step 4: Test the Setup
Run the loopback test first:
```powershell
go run .\cmd\scitek_loopback
```

Then test with actual ports:
```powershell
# Terminal 1: Run the simulator (sends to COM4)
go run .\cmd\scitek_sim

# Terminal 2: Run your driver (listens on COM3)
go run .\main.go
```

## Alternative: Virtual Serial Port Driver (Commercial)

If com0com doesn't work, try **VSPD** (trial available):
- Website: https://www.eltima.com/products/vspdxp/
- Easier GUI interface
- More stable on modern Windows

## Troubleshooting

### Port Already in Use
```powershell
# Check what's using COM3
Get-Process | Where-Object {$_.MainWindowTitle -like "*COM3*"}

# Or use PowerShell to release
[System.IO.Ports.SerialPort]::new("COM3").Close()
```

### Driver Not Signed Warning
- On Windows 10/11, you may need to disable driver signature enforcement temporarily:
  1. Hold Shift + Click Restart
  2. Troubleshoot → Advanced → Startup Settings → Restart
  3. Press F7 for "Disable driver signature enforcement"
  4. Install com0com

### Permission Denied
- Always run as **Administrator** when installing or configuring

## Without Virtual Ports (Testing on Real Hardware)

If you have actual lab equipment:
1. Connect Scitek analyzer to COM3
2. Modify `cmd/scitek_sim/main.go` to use a different port (or don't use simulator)
3. Run only your main driver application

## Quick Port Check

Run this in PowerShell to see available ports:
```powershell
[System.IO.Ports.SerialPort]::GetPortNames()
```
