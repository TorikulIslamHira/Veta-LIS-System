# Veta LIS - Veterinary Laboratory Information System

## Project Overview

Veta LIS is a production-grade Wails (v2) application for managing veterinary laboratory operations. It follows **Clean Architecture** principles with a focus on modularity, testability, and extensibility.

**Tech Stack:**
- **Frontend:** React 18 (Vite) + Tailwind CSS + TypeScript
- **Backend:** Go 1.21 (standard library) + pgx for PostgreSQL
- **Hardware:** Serial (COM Port) & TCP/IP communication
- **Supported Machines:** HA3100, Scitek, Mindray (extensible via Adapter Pattern)

---

## Project File Structure

```
veta-lis/
│
├── main.go                           # Application entry point (Wails bootstrap)
├── app.go                            # (Optional) Wails app configuration
├── wails.json                        # Wails framework configuration
├── go.mod                            # Go module definition
├── config.yaml                       # System configuration (machines, DB, server)
├── config.sqlite                     # (Optional) Local SQLite cache for settings
├── README.md                         # This file
│
├── internal/                         # Private Go packages (not importable externally)
│   │
│   ├── app/
│   │   ├── app.go                   # Main App struct (Wails bindings)
│   │   └── handlers.go              # (Optional) Wails event handlers
│   │
│   ├── core/                        # Pure business logic (no DB, no hardware)
│   │   ├── business_logic.go        # ResultValidator, PetService, OwnerService
│   │   ├── models.go                # (Optional) Domain models
│   │   └── services.go              # (Optional) Core service interfaces
│   │
│   ├── drivers/                     # ⭐ MOST IMPORTANT - Machine communication adapters
│   │   ├── drivers.go               # MachineDriver interface & DriverFactory
│   │   ├── errors.go                # Driver-specific errors
│   │   │
│   │   ├── ha3100/
│   │   │   └── ha3100.go            # HA3100 analyzer implementation
│   │   │
│   │   ├── scitek/
│   │   │   └── scitek.go            # Scitek analyzer implementation (HL7 + MLLP)
│   │   │
│   │   └── mindray/
│   │       └── mindray.go           # Mindray analyzer implementation
│   │
│   ├── parser/                      # Protocol parsing (HL7, ASTM)
│   │   └── parser.go                # HL7Parser, ASTMParser
│   │
│   ├── repository/                  # PostgreSQL data access layer
│   │   ├── db.go                    # Database connection & pool
│   │   │
│   │   ├── owners/
│   │   │   └── repository.go        # Owner CRUD operations
│   │   │
│   │   ├── pets/
│   │   │   └── repository.go        # Pet CRUD operations
│   │   │
│   │   └── results/
│   │       └── repository.go        # Lab result storage & retrieval
│   │
│   └── config/                      # Configuration management
│       └── config.go                # ConfigManager, YAML loading/saving
│
├── frontend/                         # React application (Vite + Tailwind)
│   ├── index.html                   # Entry HTML
│   ├── package.json                 # NPM dependencies
│   ├── vite.config.ts               # Vite build configuration
│   ├── tailwind.config.js           # Tailwind CSS configuration
│   ├── tsconfig.json                # TypeScript configuration
│   │
│   ├── src/
│   │   ├── index.css                # Global styles
│   │   ├── main.tsx                 # React entry point
│   │   ├── App.tsx                  # Root component
│   │   │
│   │   ├── features/                # Feature-based organization
│   │   │   ├── dashboard/
│   │   │   │   ├── Dashboard.tsx    # Machine status overview
│   │   │   │   └── index.ts         # Feature exports
│   │   │   │
│   │   │   ├── machines/
│   │   │   │   ├── MachineList.tsx
│   │   │   │   ├── MachineDetail.tsx
│   │   │   │   └── index.ts
│   │   │   │
│   │   │   ├── settings/
│   │   │   │   ├── MachineSettings.tsx   # Configure machines
│   │   │   │   ├── DatabaseSettings.tsx  # DB config
│   │   │   │   └── index.ts
│   │   │   │
│   │   │   ├── pets/
│   │   │   │   ├── PetList.tsx
│   │   │   │   ├── PetForm.tsx
│   │   │   │   └── index.ts
│   │   │   │
│   │   │   └── results/
│   │   │       ├── ResultsView.tsx
│   │   │       └── index.ts
│   │   │
│   │   ├── hooks/                   # Custom React hooks (Wails integration)
│   │   │   ├── useBackend.ts        # useMachineControl(), useMachineConfig()
│   │   │   ├── useResults.ts        # (Optional) Live result streaming
│   │   │   └── useAuth.ts           # (Optional) Authentication
│   │   │
│   │   ├── context/                 # Global state management
│   │   │   ├── MachineContext.tsx   # Machine status context
│   │   │   ├── ResultsContext.tsx   # (Optional) Live results context
│   │   │   └── UserContext.tsx      # (Optional) User session context
│   │   │
│   │   ├── components/              # Shared UI components
│   │   │   ├── common/
│   │   │   │   ├── Header.tsx
│   │   │   │   ├── Sidebar.tsx
│   │   │   │   ├── LoadingSpinner.tsx
│   │   │   │   └── ErrorBoundary.tsx
│   │   │   │
│   │   │   ├── MachineCard.tsx
│   │   │   └── StatusBadge.tsx
│   │   │
│   │   ├── types/                   # TypeScript type definitions
│   │   │   ├── index.ts
│   │   │   ├── machine.ts
│   │   │   ├── result.ts
│   │   │   └── api.ts
│   │   │
│   │   └── utils/                   # (Optional) Utility functions
│   │       ├── formatters.ts
│   │       └── validators.ts
│   │
│   ├── public/                       # Static assets
│   │   ├── favicon.ico
│   │   └── images/
│   │
│   ├── dist/                        # Build output (generated)
│   └── wailsjs/                     # Wails auto-generated bindings (generated)
│
├── scripts/                         # (Optional) Build and deployment scripts
│   ├── build.sh                     # Cross-platform build
│   └── docker/                      # Docker containerization
│
└── docs/                            # (Optional) Architecture documentation
    ├── ARCHITECTURE.md              # Detailed design patterns
    ├── API.md                       # Go functions exposed to frontend
    └── SETUP.md                     # Development environment setup
```

---

## Internal Drivers Folder - Detailed Responsibility

### 🎯 Purpose

The `internal/drivers/` folder is the **cornerstone of hardware abstraction** in Veta LIS. It implements the **Adapter Pattern** to provide a unified interface for multiple lab machine types, enabling seamless device switching and future extensibility.

### 📋 Responsibilities

#### 1. **MachineDriver Interface** (`drivers.go`)
   - **Definition:** Core contract that all machine implementations must satisfy
   - **Methods:**
     - `Connect()` - Establish connection (serial/TCP)
     - `Disconnect()` - Gracefully close connection
     - `IsConnected()` - Query connection status
     - `Listen()` - Continuously receive and parse data
     - `SendCommand()` - Transmit commands to machine
     - `GetMachineInfo()` - Retrieve device metadata
     - `ParseData()` - Convert raw bytes to LabResult
     - `GetDriverType()` - Return driver identifier

#### 2. **Machine-Specific Implementations**
   - **`ha3100/ha3100.go`**
     - Serial communication (typically 9600-19200 baud)
     - ASTM E1381-04 protocol parsing
     - Handles HA-3100 biochemical analyzer responses
   
   - **`scitek/scitek.go`**
     - TCP/IP network communication
     - HL7 v2.5+ message parsing
     - MLLP framing (0x0B start, 0x1C 0x0D end)
     - Ideal for larger laboratory setups
   
   - **`mindray/mindray.go`**
     - TCP-based hematology analyzer
     - Support for multiple result types (CBC, etc.)
     - Real-time data streaming

#### 3. **DriverFactory** (`drivers.go`)
   - **Pattern:** Registry pattern for dynamic driver instantiation
   - **Responsibility:** Create driver instances without hard-coded dependencies
   - **Benefits:**
     - Runtime machine type selection
     - Easy addition of new machine types
     - Decouples machine type string from implementation

#### 4. **Error Handling** (`errors.go`)
   - Centralized, typed error definitions
   - Specific error context (e.g., `ErrNotConnected`, `ErrParseError`)
   - Enables frontend to provide meaningful error messages

#### 5. **Data Models** (`drivers.go`)
   - **`LabResult`** - Standardized result structure across all machines
   - **`MachineInfo`** - Device metadata and status
   - **`MachineConfig`** - Connection parameters (port, baud rate, timeout)
   - Ensures consistency regardless of machine type

### 🔄 Data Flow

```
Raw Data from Machine (Serial/TCP)
          ↓
    [Listen goroutine]
          ↓
    ParseData(rawBytes)
          ↓
    LabResult struct
          ↓
    Send via dataChan
          ↓
    [Core/Repository layer receives]
```

### 🏗️ How to Add a New Machine Type

1. Create a new folder: `internal/drivers/newmachine/`
2. Implement `MachineDriver` interface in `newmachine.go`
3. Register in `main.go`:
   ```go
   factory.Register("NewMachine", newmachine.NewNewMachineDriver)
   ```
4. Update `config.yaml` with new machine configuration
5. Test connection and data parsing

---

## MachineDriver Interface - Code Snippet

Here's a detailed look at how the interface handles both **Serial** (HA3100) and **TCP** (Scitek) machines:

```go
package drivers

import "context"

// MachineDriver is the core interface for all lab machines
type MachineDriver interface {
	// Connect establishes connection - handles both Serial and TCP
	// For HA3100: Opens COM port with specified baud rate
	// For Scitek: Establishes TCP socket to IP:port
	Connect(ctx context.Context) error

	// Disconnect closes connection gracefully
	Disconnect(ctx context.Context) error

	// IsConnected returns current connection state
	IsConnected() bool

	// Listen runs in a goroutine, continuously receives data
	// Machine-specific parsing happens here
	// Results flow through the channel to the App layer
	Listen(ctx context.Context, dataChan chan<- *LabResult) error

	// SendCommand sends control commands to the machine
	// HA3100: "*IDN?" queries device info
	// Scitek: HL7 queries wrapped in MLLP frames
	SendCommand(ctx context.Context, cmd string) (string, error)

	// GetMachineInfo returns metadata
	GetMachineInfo(ctx context.Context) (*MachineInfo, error)

	// ParseData converts raw bytes to LabResult
	// HA3100: Parses ASTM format (|^~\& delimiters)
	// Scitek: Parses HL7 format (\r segment delimiters)
	ParseData(rawData []byte) (*LabResult, error)

	// GetDriverType returns machine type identifier
	GetDriverType() string
}

// LabResult is standardized across ALL machines
type LabResult struct {
	ID                string      // UUID
	MachineType       string      // "HA3100", "Scitek", etc.
	TestID            string      // From machine
	PetID             string      // Links to pet in DB
	TestName          string      // "Glucose", "BUN", etc.
	ResultValue       string      // "85" (could be numeric or text)
	ResultUnit        string      // "mg/dL", "cells/µL"
	ReferenceRange    string      // "60-110" for validation
	Status            string      // "NORMAL", "ABNORMAL", "CRITICAL"
	Timestamp         string      // ISO 8601
	RawData           []byte      // Original bytes (for debugging)
}

// MachineInfo standardizes device metadata
type MachineInfo struct {
	Type              string      // "HA3100", "Scitek"
	ConnectionType    string      // "SERIAL", "TCP"
	Port              string      // "COM3" or "192.168.1.100:5000"
	BaudRate          int         // For serial (9600, 19200, etc.)
	Status            string      // "READY", "BUSY", "ERROR"
	FirmwareVersion   string
	LastTestTime      string
	TotalTestsRun     int
}
```

### Key Design Benefits

| Aspect | Benefit |
|--------|---------|
| **Unified Interface** | Frontend doesn't need to know machine type details |
| **Easy Testing** | Mock driver can be injected for unit tests |
| **Extensible** | New machines added without modifying existing code |
| **Separation of Concerns** | Machine communication isolated from business logic |
| **Error Handling** | Consistent error types across all drivers |
| **Thread-Safe** | Goroutine-based Listen() pattern with channels |

---

## Configuration Management

### `config.yaml` Structure

```yaml
database:
  host: localhost
  port: 5432
  user: veta_user
  password: "${DB_PASSWORD}"  # Environment variable
  database: veta_lis
  ssl_mode: disable

machines:
  - id: machine_1
    type: HA3100                     # Serial analyzer
    connection_type: SERIAL
    port: COM3
    baud_rate: 9600
    parse_format: ASTM
    
  - id: machine_2
    type: Scitek                     # Network analyzer
    connection_type: TCP
    port: "192.168.1.100:5000"
    parse_format: HL7
```

### Configuration Persistence

- **YAML** for static configuration (machines, database)
- **SQLite** (optional) for runtime settings (last used port, connection history)
- **Environment variables** for sensitive data (passwords, API keys)

---

## Frontend Architecture

### Context API for Global State

**MachineContext:**
```typescript
- machines: Machine[]           // All configured machines
- activeMachine: string | null  // Currently selected machine
- machineStatus: MachineStatus  // Real-time status per machine
```

### Custom Hooks for Wails Integration

**useBackend.ts:**
```typescript
const { connectMachine, disconnectMachine, getMachineStatus } = useMachineControl();
// Calls Go methods via Wails runtime
// Updates context automatically
```

### Feature-Based Organization

```
features/
├── dashboard/        # Real-time machine overview
├── machines/         # Machine list and details
├── settings/         # Configuration UI
├── pets/            # Pet database management
└── results/         # Lab result viewer/exporter
```

---

## Getting Started

### Prerequisites
- Go 1.21+
- Node.js 18+ (for React build)
- PostgreSQL 13+ (production)
- Wails CLI v2.7+

### Installation

```bash
# Clone repository
git clone <repo>
cd veta-lis

# Install Go dependencies
go mod download

# Install frontend dependencies
cd frontend
npm install

# Build and run (development)
cd ..
wails dev

# Build for production
wails build
```

### Database Setup

```bash
# Create PostgreSQL user and database
createuser -P veta_user
createdb -O veta_user veta_lis

# Run migrations (to be added)
psql -U veta_user -d veta_lis -f migrations/001_initial.sql
```

---

## Development Guidelines

### Adding a New Lab Machine

1. **Create Driver:**
   ```go
   // internal/drivers/newmachine/newmachine.go
   type NewMachineDriver struct { ... }
   func (d *NewMachineDriver) Connect(ctx context.Context) error { ... }
   ```

2. **Register in App:**
   ```go
   factory.Register("NewMachine", newmachine.NewNewMachineDriver)
   ```

3. **Update Config:**
   ```yaml
   machines:
     - id: machine_3
       type: NewMachine
       connection_type: TCP
       port: "192.168.1.101:5000"
   ```

### Testing

```bash
# Run Go tests
go test ./...

# Run with coverage
go test -cover ./...

# Frontend tests
cd frontend && npm test
```

---

## Production Deployment

### Docker Deployment

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN wails build -m production

FROM alpine:latest
COPY --from=builder /app/build/bin /usr/local/bin
CMD ["veta-lis"]
```

### Security Considerations

- Use **TLS** for Scitek TCP connections (use `tls:// protocol`)
- Store credentials in **environment variables**, never in config files
- Implement **role-based access control** (admin, technician, viewer)
- Add **audit logging** for result modifications
- Use **parameterized SQL queries** (pgx already handles this)

---

## Support & Documentation

- **API Reference:** See `docs/API.md`
- **Architecture Deep Dive:** See `docs/ARCHITECTURE.md`
- **Setup Guide:** See `docs/SETUP.md`

---

**Version:** 1.0.0  
**Last Updated:** January 2026  
**License:** MIT
