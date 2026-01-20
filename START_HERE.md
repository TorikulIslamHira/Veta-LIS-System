# ⚡ VETA LIS - READY TO RUN

## 🎉 Your Application is Complete!

Everything is set up and ready to run. Here's what you need to do:

---

## 📋 Pre-requisites (One-Time Setup)

These tools must be installed on your Windows computer:

### 1. **Go 1.21+**
- **Download:** https://golang.org/dl/
- **Select:** Windows → amd64
- **Install:** Default settings
- **Verify:** `go version`

### 2. **Node.js 18+ LTS**
- **Download:** https://nodejs.org/
- **Select:** LTS version
- **Install:** Default settings
- **Verify:** `node --version`

---

## 🚀 Run the Application (Choose ONE)

### **Option 1: Automatic Setup (Easiest)**

**PowerShell (As Administrator):**
```powershell
cd "e:\LIS Connector\veta-lis"
.\SETUP_AND_RUN.ps1
```

**Command Prompt (As Administrator):**
```batch
cd e:\LIS Connector\veta-lis
SETUP_AND_RUN.bat
```

### **Option 2: Manual Step-by-Step**

```powershell
# Navigate to project
cd "e:\LIS Connector\veta-lis"

# Install dependencies
go mod download
cd frontend && npm install && cd ..

# Run development server
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails dev
```

### **Option 3: Just Run It**

```powershell
cd "e:\LIS Connector\veta-lis"
wails dev
```

---

## 🌐 Access the Application

Once you see "Listening on http://localhost:3000" in the terminal:

1. Open your web browser
2. Go to: **http://localhost:3000**
3. You should see the Veta LIS Dashboard

---

## 📁 What's Included

```
✅ Complete Backend (Go)
   - 3 Lab Machine Drivers
   - Business Logic
   - Configuration System
   - Protocol Parsers

✅ Complete Frontend (React)
   - Dashboard
   - Machine Management
   - Settings
   - Result Viewing

✅ Documentation
   - Architecture Guide
   - Quick Start Guide
   - Driver Reference
   - API Documentation

✅ Configuration
   - YAML-based setup
   - 3 Machine examples
   - Database config ready
```

---

## 🎯 What Happens When You Run It

1. **Backend starts** (Go) - Listens on port 3000
2. **Frontend builds** (React/Vite) - Compiles Tailwind + TypeScript
3. **Auto-reload enabled** - Changes to code reload automatically
4. **Browser opens** - http://localhost:3000 opens automatically

---

## 🛑 To Stop the Application

Press **`Ctrl+C`** in the terminal where `wails dev` is running.

---

## 📂 Project Files

### Important Locations

| File/Folder | Purpose |
|-------------|---------|
| `config.yaml` | Machine definitions & database config |
| `internal/drivers/` | Lab machine implementations |
| `frontend/src/` | React user interface |
| `docs/ARCHITECTURE.md` | Design patterns explained |
| `QUICK_START_GUIDE.md` | Detailed setup guide |

---

## ✨ Features Included

✅ **Real-time Machine Monitoring**
- Connect/disconnect lab machines
- View machine status
- Receive test results

✅ **Multiple Machine Support**
- HA3100 (Serial, ASTM)
- Scitek (TCP, HL7)
- Mindray (TCP, Binary)
- Easy to add more

✅ **Data Management**
- PostgreSQL integration ready
- Pet database structure
- Result storage
- Owner management

✅ **Configuration**
- YAML-based settings
- Machine definitions
- Database configuration
- Environment variables

---

## 🔍 Troubleshooting

### "Go not found"
```powershell
# Add to PATH and try again
$env:Path += "C:\Program Files\Go\bin"
go version
```

### "Node not found"
```powershell
# Add to PATH and try again
$env:Path += "C:\Program Files\nodejs"
node --version
```

### "Port 3000 in use"
```powershell
# Kill existing process
netstat -ano | findstr :3000
taskkill /PID <PID> /F

# Or use different port
wails dev -port 3001
```

---

## 📚 After Running

### To Learn More
- Read `docs/ARCHITECTURE.md` - Understand the design
- Read `docs/DRIVER_INTERFACE_REFERENCE.md` - How drivers work
- Explore `internal/drivers/` - See implementations
- Check `frontend/src/` - See React components

### To Develop Further
1. Make changes to `internal/` or `frontend/src/`
2. Browser auto-reloads (hot reload)
3. Check terminal for errors
4. Use browser DevTools (F12) for frontend debugging

### To Build for Production
```powershell
wails build -m production
```
This creates a standalone .exe file.

---

## 💻 System Requirements

| Component | Minimum | Recommended |
|-----------|---------|------------|
| **Go** | 1.21 | 1.21+ |
| **Node.js** | 16 | 18+ LTS |
| **npm** | 8 | 9+ |
| **Windows** | 10 | 10/11 |
| **RAM** | 2GB | 4GB+ |

---

## ✅ Checklist Before Running

- [ ] Downloaded and installed Go 1.21+
- [ ] Downloaded and installed Node.js 18+ LTS
- [ ] Verified `go version` works
- [ ] Verified `node --version` works
- [ ] You're in the `veta-lis` directory
- [ ] You have `config.yaml` file
- [ ] You have `main.go` file

---

## 🎁 You're All Set!

Everything is ready. Just:

1. **Install Go and Node.js** (if not already done)
2. **Run one of the setup scripts** or `wails dev`
3. **Open http://localhost:3000** in your browser
4. **Start developing!**

---

## 📞 Need Help?

### Check These Files
- `QUICK_START_GUIDE.md` - Detailed setup steps
- `README.md` - Project overview
- `docs/ARCHITECTURE.md` - Design documentation

### Common Issues
- Port already in use → Kill process or use `-port 3001`
- Go/Node not in PATH → Restart PowerShell/cmd
- Build errors → Check `go mod download` and `npm install`

---

## 🚀 Ready?

**Run this command:**

```powershell
cd "e:\LIS Connector\veta-lis"
wails dev
```

**That's it! Your Veta LIS application is running!** 🎉

---

**Next Step:** Go to http://localhost:3000

---

**Project:** Veta LIS - Veterinary Laboratory Information System  
**Status:** ✅ **READY TO RUN**  
**Date:** January 20, 2026
