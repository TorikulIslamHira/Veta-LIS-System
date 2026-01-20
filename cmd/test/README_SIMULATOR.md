# Scitek Simulator - Testing Instructions

## Overview
This simulator mimics the **Scitek Biochemistry Analyzer** by sending HL7 messages over RS232 serial port with MLLP framing.

## Setup

### 1. Install Virtual Serial Port Software (for testing without real hardware)
- **Windows**: Use [com0com](https://sourceforge.net/projects/com0com/) or [Virtual Serial Port Driver](https://www.eltima.com/products/vspdxp/)
- Create a pair: COM3 ↔ COM4

### 2. Install Dependencies
```powershell
go get go.bug.st/serial
```

### 3. Build the Simulator
```powershell
cd e:\Github\Veta-LIS-System
go build -o scitek_sim.exe .\cmd\scitek_sim
```

## Running the Simulator

### Option 1: Direct Run
```powershell
go run .\cmd\scitek_sim
```

### Option 2: Execute Built Binary
```powershell
.\scitek_sim.exe
```

## Expected Output

```
🔬 Scitek Biochemistry Analyzer Simulator
==========================================
Port: COM4 | Baud: 9600

✅ Serial port COM4 opened successfully

📤 Sending HL7 Message:
------------------------
MSH|^~\&|Scitek||||20240120100000||ORU^R01|1|P|2.3.1
PID|1||101||Test^Patient^^^||19900101|M|||123 Main St^^Cityville^^12345
OBR|1||1|CHEM^Chemistry Panel|||20240120100000
OBX|1|NM|GLU^Glucose||105|mg/dL|70-110|N|||F
OBX|2|NM|UREA^Blood Urea Nitrogen||40|mg/dL|7-20|H|||F
OBX|3|NM|CREA^Creatinine||1.2|mg/dL|0.6-1.2|N|||F
------------------------

📦 Total bytes to send: 357 (including MLLP framing)
✅ Sent 357 bytes to COM4

⏳ Waiting for ACK from driver...
✅ ACK Received (6 bytes)
📥 ACK Data: [11 65 67 75 28 13]
📥 ACK String: "\vACK\x1c\r"
✅ Valid MLLP ACK received - transmission successful!

🏁 Simulator finished
```

## Testing Workflow

1. **Start your main Veta LIS application** (listening on COM3)
2. **Run the simulator** (sends to COM4)
3. The simulator sends HL7 data → Your driver receives it on COM3
4. Your driver sends ACK back → Simulator receives it on COM4

## Customization

### Change Port or Baud Rate
Edit constants in `cmd/scitek_sim/main.go`:
```go
const (
    PORT_NAME = "COM4"  // Change to your simulator port
    BAUD_RATE = 9600    // Change if needed
)
```

### Modify Test Data
Edit the `buildHL7Message()` function to add/modify test results.

## Troubleshooting

### "Failed to open port COM4"
- Ensure COM4 exists (check Device Manager on Windows)
- Close any other programs using the port
- Run with administrator privileges if needed

### "Timeout - no ACK received"
- Verify your main application is running and listening on COM3
- Check that COM3 and COM4 are paired correctly
- Increase timeout in code if needed: `port.SetReadTimeout(10 * time.Second)`

### Serial Port Access Denied
```powershell
# Run as Administrator
runas /user:Administrator "go run .\cmd\test\scitek_sim.go"
```

## MLLP Framing Details

The simulator wraps HL7 messages with MLLP (Minimal Lower Layer Protocol):
- **Start Byte**: `0x0B` (Vertical Tab)
- **Message**: HL7 content
- **End Bytes**: `0x1C` (File Separator) + `0x0D` (Carriage Return)

Example:
```
[0x0B] MSH|^~\&|Scitek|... [0x1C][0x0D]
```
