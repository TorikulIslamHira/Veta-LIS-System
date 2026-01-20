# 🚀 Veta LIS - Complete Setup & Run Guide

## Quick Start (3 Steps)

### **Step 1: Install Required Software**

You need to install these tools (one-time only):

#### **Go 1.21+**
- Download: https://golang.org/dl/
- Select: **Windows** → **amd64**
- Install with default settings
- Verify: Open PowerShell and run `go version`

#### **Node.js 18+ LTS**
- Download: https://nodejs.org/
- Install with default settings
- Verify: Open PowerShell and run `node --version`

#### **Git** (Optional but recommended)
- Download: https://git-scm.com/
- Install with default settings

---

### **Step 2: Run the Setup Script**

#### **Option A: PowerShell (Recommended)**

```powershell
# Right-click PowerShell → Run as Administrator
# Navigate to project folder
cd "e:\LIS Connector\veta-lis"

# Run the setup script
.\SETUP_AND_RUN.ps1
```

#### **Option B: Command Prompt**

```batch
# Right-click cmd.exe → Run as Administrator
# Navigate to project folder
cd e:\LIS Connector\veta-lis

# Run the setup script
SETUP_AND_RUN.bat
```

#### **Option C: Manual Setup**

```powershell
# In PowerShell (as Administrator)
cd "e:\LIS Connector\veta-lis"

# Install Go dependencies
go mod download

# Install frontend dependencies
cd frontend
npm install
cd ..

# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Run development server
wails dev
```

---

### **Step 3: Access the Application**

Once `wails dev` is running:
- Open your browser
- Go to: **http://localhost:3000**
- You should see the Veta LIS dashboard

---

## 🔍 Verification Checklist

Before running, verify everything is installed:

```powershell
# Check Go
go version
# Expected: go version go1.21.X windows/amd64

# Check Node.js
node --version
# Expected: v18.X.X or higher

# Check npm
npm --version
# Expected: v9.X.X or higher

# Check Wails (after go mod download)
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails version
# Expected: Wails version output
```

---

## 🎯 Troubleshooting

### **"Go not found" error**
- Go is installed but not in PATH
- **Solution:** Restart PowerShell or restart your computer
- Or add Go manually: `$env:Path += "C:\Program Files\Go\bin"`

### **"Node not found" error**
- Node.js is installed but not in PATH
- **Solution:** Restart PowerShell or restart your computer
- Or add Node manually: `$env:Path += "C:\Program Files\nodejs"`

### **"wails: command not found"**
- Wails CLI not installed or not in PATH
- **Solution:** Run `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- Then restart PowerShell

### **Port 3000 already in use**
```powershell
# Find process using port 3000
netstat -ano | findstr :3000

# Kill the process (replace PID with actual process ID)
taskkill /PID <PID> /F

# Or use different port
wails dev -port 3001
```

### **"Cannot find path" error**
- Make sure you're in the correct directory
- Should see config.yaml, main.go, wails.json
```powershell
cd "e:\LIS Connector\veta-lis"
dir
```

### **Frontend dependencies error**
```powershell
cd frontend
npm cache clean --force
rm -r node_modules
npm install --legacy-peer-deps
cd ..
```

---

## 📋 Project Structure After Setup

```
veta-lis/
├── main.go                 ← Go entry point
├── config.yaml             ← Configuration
├── go.mod                  ← Go dependencies
│
├── internal/               ← Backend code
│   ├── app/app.go
│   ├── drivers/            ← Lab machine drivers
│   └── ...
│
├── frontend/               ← React frontend
│   └── src/
│       ├── App.tsx
│       ├── main.tsx
│       └── ...
│
└── docs/                   ← Documentation
    ├── ARCHITECTURE.md
    └── ...
```

---

## 🚀 Development Workflow

### **Running the App**

```powershell
cd "e:\LIS Connector\veta-lis"
wails dev
```

This will:
- Start Go backend (port 3000)
- Serve React frontend (port 3000)
- Enable hot-reload for both
- Open browser automatically

### **Making Changes**

**Backend changes (Go files):**
- Edit `internal/` files
- Backend recompiles automatically
- Browser refreshes automatically

**Frontend changes (React files):**
- Edit `frontend/src/` files
- Frontend rebuilds automatically
- Browser refreshes automatically

### **Stopping the Server**

Press `Ctrl+C` in the terminal where `wails dev` is running.

---

## 🔧 Production Build

When ready to build for production:

```powershell
cd "e:\LIS Connector\veta-lis"
wails build -m production
```

This creates a standalone `.exe` file that can be distributed.

---

## 📚 Next Steps After Setup

1. **Read Documentation:**
   - `README.md` - Project overview
   - `docs/ARCHITECTURE.md` - Design patterns
   - `docs/DRIVER_INTERFACE_REFERENCE.md` - Machine drivers

2. **Explore the Code:**
   - `internal/drivers/` - Lab machine implementations
   - `frontend/src/features/` - UI features
   - `config.yaml` - Configuration

3. **Start Development:**
   - Set up PostgreSQL database
   - Implement repository layer
   - Connect to real lab machines
   - Build custom UI

---

## 🎁 What You'll See

Once running at http://localhost:3000, you'll see:

- **Dashboard:** Machine status overview
- **Machines:** Lab machine list and management
- **Settings:** Configure machine connections
- **Results:** View lab test results
- **Pets:** Pet database (when implemented)

---

## 💡 Tips

1. **Keep terminal open:** The `wails dev` terminal must stay open while developing
2. **Check console:** Browser DevTools (F12) shows frontend errors
3. **Monitor output:** Terminal shows backend logs and errors
4. **Save files:** Changes auto-reload, but some changes may require restart
5. **Use config.yaml:** Manage machines and database settings there

---

## 📞 Support Resources

- **Wails Documentation:** https://wails.io/docs/
- **Go Documentation:** https://golang.org/doc/
- **React Documentation:** https://react.dev/
- **Project Documentation:** See `docs/` folder

---

**Status:** ✅ Ready to run!

Once you have Go and Node.js installed, just run one of the setup scripts and you're good to go! 🚀
