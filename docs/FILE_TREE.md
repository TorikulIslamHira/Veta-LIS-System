# Veta LIS - Complete File Tree

## Project Structure Visualization

```
veta-lis/
├── 📄 main.go                             ← Wails application entry point
├── 📄 wails.json                          ← Wails framework config
├── 📄 go.mod                              ← Go module dependencies
├── 📄 config.yaml                         ← System configuration
├── 📄 README.md                           ← Main documentation
│
├── 📁 internal/                           ← Private packages
│   │
│   ├── 📁 app/                            ← Application layer (Wails bindings)
│   │   └── 📄 app.go                      ├─ App struct
│   │                                      ├─ ConnectMachine()
│   │                                      ├─ DisconnectMachine()
│   │                                      ├─ GetMachineStatus()
│   │                                      └─ Result channel management
│   │
│   ├── 📁 core/                           ← Pure business logic
│   │   └── 📄 business_logic.go           ├─ ResultValidator
│   │                                      ├─ PetService
│   │                                      ├─ OwnerService
│   │                                      └─ Domain models
│   │
│   ├── 📁 drivers/                        ⭐ ADAPTER PATTERN (Most Important)
│   │   ├── 📄 drivers.go                  ├─ MachineDriver interface
│   │   │                                  ├─ LabResult struct
│   │   │                                  ├─ MachineInfo struct
│   │   │                                  └─ DriverFactory
│   │   │
│   │   ├── 📄 errors.go                   └─ Typed error definitions
│   │   │
│   │   ├── 📁 ha3100/                     ← Serial Analyzer (ASTM)
│   │   │   └── 📄 ha3100.go               ├─ Connect/Disconnect (Serial)
│   │   │                                  ├─ Listen() goroutine
│   │   │                                  ├─ ParseData() ASTM parsing
│   │   │                                  └─ SendCommand()
│   │   │
│   │   ├── 📁 scitek/                     ← Network Analyzer (HL7)
│   │   │   └── 📄 scitek.go               ├─ Connect/Disconnect (TCP)
│   │   │                                  ├─ Listen() goroutine
│   │   │                                  ├─ ParseData() HL7 parsing
│   │   │                                  └─ MLLP frame handling
│   │   │
│   │   └── 📁 mindray/                    ← Hematology Analyzer (TCP)
│   │       └── 📄 mindray.go              ├─ Connect/Disconnect (TCP)
│   │                                      ├─ Listen() goroutine
│   │                                      ├─ ParseData() CBC parsing
│   │                                      └─ Proprietary protocol
│   │
│   ├── 📁 parser/                         ← Protocol parsing
│   │   └── 📄 parser.go                   ├─ HL7Parser
│   │                                      ├─ ASTMParser
│   │                                      ├─ Parse() methods
│   │                                      └─ Message structs
│   │
│   ├── 📁 repository/                     ← Data access layer
│   │   ├── 📁 owners/
│   │   │   └── 📄 repository.go           ├─ Owner CRUD
│   │   │
│   │   ├── 📁 pets/
│   │   │   └── 📄 repository.go           ├─ Pet CRUD
│   │   │
│   │   └── 📁 results/
│   │       └── 📄 repository.go           ├─ Result queries & storage
│   │
│   └── 📁 config/                         ← Configuration management
│       └── 📄 config.go                   ├─ ConfigManager
│                                          ├─ Load/Save YAML
│                                          ├─ Machine config CRUD
│                                          └─ Validation
│
├── 📁 frontend/                           ← React + Vite application
│   ├── 📄 index.html                      ← HTML entry point
│   ├── 📄 package.json                    ← NPM dependencies
│   ├── 📄 vite.config.ts                  ← Vite build config
│   ├── 📄 tailwind.config.js              ← Tailwind styling
│   ├── 📄 postcss.config.js               ← PostCSS config
│   ├── 📄 tsconfig.json                   ← TypeScript config
│   ├── 📄 tsconfig.node.json              ← Node TypeScript config
│   │
│   ├── 📁 src/
│   │   ├── 📄 index.css                   ← Global styles
│   │   ├── 📄 main.tsx                    ← React entry point
│   │   ├── 📄 App.tsx                     ← Root component
│   │   │
│   │   ├── 📁 features/                   ← Feature-based modules
│   │   │   │
│   │   │   ├── 📁 dashboard/
│   │   │   │   ├── 📄 Dashboard.tsx       ├─ Machine overview
│   │   │   │   └── 📄 index.ts            └─ Feature export
│   │   │   │
│   │   │   ├── 📁 machines/
│   │   │   │   ├── 📄 MachineList.tsx
│   │   │   │   ├── 📄 MachineDetail.tsx
│   │   │   │   └── 📄 index.ts
│   │   │   │
│   │   │   ├── 📁 settings/
│   │   │   │   ├── 📄 MachineSettings.tsx ├─ Configure connections
│   │   │   │   ├── 📄 DatabaseSettings.tsx
│   │   │   │   └── 📄 index.ts
│   │   │   │
│   │   │   ├── 📁 pets/
│   │   │   │   ├── 📄 PetList.tsx
│   │   │   │   ├── 📄 PetForm.tsx
│   │   │   │   └── 📄 index.ts
│   │   │   │
│   │   │   └── 📁 results/
│   │   │       ├── 📄 ResultsView.tsx
│   │   │       └── 📄 index.ts
│   │   │
│   │   ├── 📁 hooks/                      ← Custom React hooks
│   │   │   ├── 📄 useBackend.ts           ├─ Wails integration
│   │   │   │                              ├─ useMachineControl()
│   │   │   │                              ├─ useMachineConfig()
│   │   │   │                              └─ Calls Go methods
│   │   │   │
│   │   │   ├── 📄 useResults.ts           ├─ (Future) Result streaming
│   │   │   └── 📄 useAuth.ts              ├─ (Future) Authentication
│   │   │
│   │   ├── 📁 context/                    ← Global state management
│   │   │   ├── 📄 MachineContext.tsx      ├─ Machine status state
│   │   │   │                              ├─ Machine provider
│   │   │   │                              └─ useMachineContext() hook
│   │   │   │
│   │   │   ├── 📄 ResultsContext.tsx      ├─ (Future) Results state
│   │   │   └── 📄 UserContext.tsx         ├─ (Future) User session
│   │   │
│   │   ├── 📁 components/                 ← Shared UI components
│   │   │   ├── 📁 common/
│   │   │   │   ├── 📄 Header.tsx
│   │   │   │   ├── 📄 Sidebar.tsx
│   │   │   │   ├── 📄 LoadingSpinner.tsx
│   │   │   │   └── 📄 ErrorBoundary.tsx
│   │   │   │
│   │   │   ├── 📄 MachineCard.tsx         ├─ Machine status card
│   │   │   └── 📄 StatusBadge.tsx         ├─ Status indicator
│   │   │
│   │   ├── 📁 types/                      ← TypeScript types
│   │   │   └── 📄 index.ts                ├─ Machine, LabResult
│   │   │                                  ├─ Pet, Owner
│   │   │                                  └─ API response types
│   │   │
│   │   └── 📁 utils/                      ← Utility functions
│   │       ├── 📄 formatters.ts           ├─ Date, number formatting
│   │       └── 📄 validators.ts           ├─ Input validation
│   │
│   ├── 📁 public/                         ← Static assets
│   │   ├── 📄 favicon.ico
│   │   └── 📁 images/
│   │
│   ├── 📁 dist/                           ← Build output (generated)
│   └── 📁 wailsjs/                        ← Auto-generated Wails bindings (generated)
│       └── 📁 go/
│           └── 📁 app/
│               └── 📄 App.ts              ← Type-safe Go bindings
│
├── 📁 docs/                               ← Documentation
│   ├── 📄 ARCHITECTURE.md                 ← Design patterns & layers
│   ├── 📄 API.md                          ← (Future) API reference
│   └── 📄 SETUP.md                        ← (Future) Dev environment
│
└── 📁 scripts/                            ← Build & deployment
    ├── 📄 build.sh                        ← Cross-platform build
    └── 📁 docker/
        └── 📄 Dockerfile                  ← Container image
```

---

## File Statistics

```
Backend (Go):
├── Interfaces & Types: ~150 lines (drivers.go)
├── Implementations: ~600 lines (3 drivers × 200 lines each)
├── Business Logic: ~150 lines (core/business_logic.go)
├── Parser: ~200 lines (parser/parser.go)
├── Configuration: ~250 lines (config/config.go)
├── App Layer: ~400 lines (app/app.go)
└── Total Go: ~1,750 lines

Frontend (React + TypeScript):
├── Components: ~400 lines
├── Hooks: ~150 lines
├── Context: ~80 lines
├── Types: ~60 lines
└── Total Frontend: ~690 lines

Configuration:
├── config.yaml: ~35 lines
├── wails.json: ~15 lines
├── package.json: ~25 lines
├── vite.config.ts: ~15 lines
└── tailwind.config.js: ~20 lines

Documentation:
├── README.md: ~400 lines
└── ARCHITECTURE.md: ~300 lines
```

---

## Key Files by Responsibility

### Hardware Communication
- `internal/drivers/drivers.go` - Interface definition
- `internal/drivers/ha3100/ha3100.go` - Serial implementation
- `internal/drivers/scitek/scitek.go` - HL7/TCP implementation
- `internal/drivers/mindray/mindray.go` - Hematology implementation

### Data Processing
- `internal/parser/parser.go` - HL7 and ASTM parsing
- `internal/core/business_logic.go` - Validation logic
- `internal/config/config.go` - Settings management

### Frontend-Backend Communication
- `internal/app/app.go` - Wails bindings (exposed to frontend)
- `frontend/src/hooks/useBackend.ts` - Frontend integration
- `frontend/src/context/MachineContext.tsx` - State management

### User Interface
- `frontend/src/features/dashboard/Dashboard.tsx` - Main view
- `frontend/src/features/settings/MachineSettings.tsx` - Configuration
- `frontend/src/features/machines/` - Machine management

---

## Development Workflow

### Directory Navigation

```bash
# Backend development
cd veta-lis
go mod download
go run main.go

# Frontend development  
cd veta-lis/frontend
npm install
npm run dev

# Full app (with hot reload)
cd veta-lis
wails dev

# Production build
wails build -m production
```

### Package Imports (Backend)

```go
// ✅ Correct (relative to module root)
import "veta-lis/internal/drivers"
import "veta-lis/internal/core"
import "veta-lis/internal/config"

// ❌ Wrong (external imports)
// import "github.com/veta-lis/internal/drivers" // Private!
```

### Module Boundaries

```
✅ Can import:
- internal/app → internal/core ✓
- internal/app → internal/drivers ✓
- internal/drivers → internal/parser ✓

❌ Cannot import:
- internal/core → internal/drivers ✗
- internal/repository → internal/app ✗
- internal/drivers → internal/app ✗
```

---

## Next Steps for Implementation

### Phase 1: Core Infrastructure
- [ ] Set up PostgreSQL database & schema
- [ ] Implement Repository layer (CRUD operations)
- [ ] Add authentication & authorization
- [ ] Build database connection pooling

### Phase 2: Driver Enhancement
- [ ] Implement serial port communication (go-serial)
- [ ] TCP connection retry logic with exponential backoff
- [ ] Complete ASTM/HL7 parsing with edge cases
- [ ] Add result validation with reference ranges

### Phase 3: Frontend Features
- [ ] Machine connection/status UI
- [ ] Results viewing & filtering
- [ ] Pet database management
- [ ] Test result export (PDF/CSV)

### Phase 4: Production Hardening
- [ ] Comprehensive error handling & logging
- [ ] Performance optimization (connection pooling, caching)
- [ ] Security audit (SQL injection, CSRF protection)
- [ ] Docker deployment configuration
- [ ] Cross-platform build (Windows, macOS, Linux)

---

**Version:** 1.0.0  
**Architecture:** Clean Architecture + Adapter Pattern  
**Status:** Ready for Development ✅
