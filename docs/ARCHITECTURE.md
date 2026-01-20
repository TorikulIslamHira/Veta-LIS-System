# Veta LIS - Architecture Deep Dive

## Design Principles

Veta LIS follows **Clean Architecture** with emphasis on:

1. **Separation of Concerns** - Each layer has a single responsibility
2. **Dependency Inversion** - High-level modules don't depend on low-level modules
3. **Testability** - Business logic decoupled from I/O and hardware
4. **Scalability** - Easy to add new machines, databases, or features
5. **Maintainability** - Clear code organization and responsibility

---

## Layered Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Frontend (React)                       │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Features: Dashboard, Settings, Machines, Pets      │    │
│  │  State: MachineContext, ResultsContext              │    │
│  │  Hooks: useBackend, useResults, useMachineControl   │    │
│  └─────────────────────────────────────────────────────┘    │
└──────────────────┬──────────────────────────────────────────┘
                   │ Wails Runtime (TypeScript ↔ Go)
┌──────────────────▼──────────────────────────────────────────┐
│                   Application Layer                         │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  App struct: ConnectMachine(), GetMachineStatus()   │    │
│  │  Event handlers, lifecycle management               │    │
│  │  Wails bindings (exposed to frontend)               │    │
│  └─────────────────────────────────────────────────────┘    │
└──────────────────┬──────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────────┐
│                  Business Logic Layer                       │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Core Services                                       │   │
│  │  - ResultValidator: Status determination            │   │
│  │  - PetService: Pet CRUD validation                  │   │
│  │  - OwnerService: Owner data validation              │   │
│  │                                                      │   │
│  │  No DB, no hardware - Pure logic                    │   │
│  └──────────────────────────────────────────────────────┘   │
└──────────────────┬──────────────────────────────────────────┘
                   │
        ┌──────────┴──────────┐
        │                     │
┌───────▼────────┐    ┌──────▼────────────┐
│  Data Layer    │    │  Hardware Layer   │
├────────────────┤    ├───────────────────┤
│ Repository     │    │ Drivers           │
│ - owners       │    │ - HA3100          │
│ - pets         │    │ - Scitek          │
│ - results      │    │ - Mindray         │
│                │    │                   │
│ PostgreSQL     │    │ Serial/TCP        │
│ + pgx          │    │ Communication     │
└────────────────┘    └───────────────────┘
```

---

## Clean Architecture Boundaries

### ✅ What Goes in Each Layer

| Layer | Responsibility | Examples |
|-------|---|---|
| **Frontend** | UI, user interactions, form validation | React components, Tailwind styling |
| **Application** | Wails bindings, orchestration | App.ConnectMachine(), lifecycle |
| **Business Logic** | Domain rules, validation | ResultValidator, PetService |
| **Repository** | Database CRUD, SQL queries | OwnerRepository.GetByID() |
| **Drivers** | Hardware communication, parsing | HA3100.ParseData(), Connect() |
| **Config** | Settings management | YAML loading, validation |
| **Parser** | Protocol parsing | HL7Parser, ASTMParser |

### ❌ What NOT to Do

```go
// ❌ DON'T: Database calls in driver
func (d *HA3100Driver) Listen() {
    if err := d.db.SaveResult(result) { ... } // ← WRONG
}

// ✅ DO: Return data via channel
func (d *HA3100Driver) Listen(ctx context.Context, dataChan chan<- *LabResult) {
    dataChan <- result  // ← Let App layer handle storage
}
```

```go
// ❌ DON'T: Business logic in repository
func (r *ResultRepository) GetResults() {
    // Validate status, calculate ranges, etc. // ← WRONG
}

// ✅ DO: Repository only queries
func (r *ResultRepository) GetResults() []LabResult {
    return query results only
}
```

---

## Adapter Pattern (Drivers)

### Problem It Solves

Without the adapter pattern, adding a new machine would require:
- Modifying existing code (Open/Closed Principle violation)
- Knowing all machine types upfront
- Duplication of connection logic

### How It Works

```
┌─────────────────────┐
│  MachineDriver      │ ← Interface contract
│  (interface)        │   All machines implement this
└──────────┬──────────┘
           │
     ┌─────┴─────┬─────────┬──────────┐
     │           │         │          │
┌────▼────┐ ┌────▼────┐ ┌─▼────┐ ┌──▼─────┐
│HA3100   │ │ Scitek  │ │Mind- │ │Future? │
│Adapter  │ │Adapter  │ │ray   │ │        │
│         │ │         │ │Adapt │ │        │
└────┬────┘ └────┬────┘ └──┬───┘ └────┬───┘
     │           │         │         │
     ↓           ↓         ↓         ↓
   SERIAL       TCP       TCP       ???
```

### Factory Pattern for Dynamic Creation

```go
factory := NewDriverFactory()
factory.Register("HA3100", ha3100.NewHA3100Driver)
factory.Register("Scitek", scitek.NewScitekDriver)

// Later, runtime configuration determines which driver to use
driver, _ := factory.Create(config.Type) // ← No hard-coded dependencies!
```

---

## Data Flow Examples

### Example 1: HA3100 Test Result Reception

```
1. Frontend: useMachineControl() → connectMachine("machine_1")
                    ↓
2. App.ConnectMachine() → driver.Connect()
                    ↓
3. HA3100Driver spawns Listen() goroutine
                    ↓
4. Machine sends: "GLUCOSE|85|mg/dL|60-110\n"
                    ↓
5. Listen() reads and calls ParseData()
                    ↓
6. LabResult sent to dataChan
                    ↓
7. App.GetRecentResults() retrieves from channel
                    ↓
8. Frontend receives via useMachineControl().getRecentResults()
                    ↓
9. Display in Dashboard/Results UI
```

### Example 2: Machine Configuration Save

```
1. Frontend: MachineSettings.tsx → useMachineConfig().saveMachineConfig()
                    ↓
2. App.SaveMachineConfig()
                    ↓
3. ConfigManager.UpdateMachine() → ConfigManager.Save()
                    ↓
4. Write to config.yaml
                    ↓
5. On next app startup, load new configuration
```

---

## Error Handling Strategy

### Typed Errors for Better Error Handling

```go
// ✅ GOOD: Typed errors allow pattern matching
if errors.Is(err, drivers.ErrNotConnected) {
    return fmt.Errorf("please connect machine first")
}

// ❌ AVOID: String matching
if err.Error() == "not connected" { ... }
```

### Error Propagation

```
Hardware Error
    ↓
Driver layer: ErrConnectionFailed
    ↓
App layer: Wraps with context
    ↓
Frontend: Receives user-friendly message
```

---

## Concurrency Model

### Goroutine Architecture

```go
main()
    ├── Wails Event Loop (UI)
    │
    └── For each connected machine:
        └── Listen() goroutine
            └── Continuously reads data
                └── Sends results via channel
```

### Thread Safety

- **Channels** used for communication (type-safe)
- **Mutex** for shared state (activeDrivers map)
- **Context** for graceful cancellation

```go
type App struct {
    activeDriversMutex sync.RWMutex
    activeDrivers      map[string]drivers.MachineDriver
    resultChannel      chan *drivers.LabResult
    ctx                context.Context
    cancel             context.CancelFunc
}
```

---

## Configuration Management Strategy

### Layered Configuration

```
config.yaml
    ├── Static: Machine definitions, DB host, server port
    │   (committed to git, environment-specific)
    │
    └── Runtime overrides
        ├── Environment variables (secrets)
        ├── SQLite cache (user preferences)
        └── Command-line flags
```

### Secret Management

```yaml
# config.yaml (committed)
database:
  host: localhost
  password: "${DB_PASSWORD}"  # ← Loaded from env at runtime

# .env (NOT committed)
export DB_PASSWORD="secure_password_123"
```

---

## Testing Strategy

### Unit Testing Each Layer

```go
// Test core business logic
func TestResultValidator_ValidateResult(t *testing.T) {
    rv := core.NewResultValidator()
    result := &drivers.LabResult{
        ResultValue: "85",
        TestCode:    "GLU",
    }
    rv.ValidateResult(result)
    assert.Equal(t, "NORMAL", result.Status)
}

// Test driver interface (mock implementation)
type MockDriver struct {
    connectCalled bool
}
func (m *MockDriver) Connect(ctx context.Context) error {
    m.connectCalled = true
    return nil
}

// Test App layer
func TestApp_ConnectMachine(t *testing.T) {
    app := NewApp("test-config.yaml")
    status, err := app.ConnectMachine("machine_1")
    // Assert...
}
```

### Integration Testing

```bash
# Start PostgreSQL in Docker
docker run -d -e POSTGRES_PASSWORD=test postgres:15

# Run integration tests
go test -tags=integration ./...
```

---

## Dependency Injection

### Constructor Injection

```go
// ✅ Dependencies injected via constructor
func NewApp(configPath string) *App {
    configManager := config.NewConfigManager(configPath)
    factory := drivers.NewDriverFactory()
    
    return &App{
        configManager: configManager,
        driverFactory: factory,
        // ...
    }
}
```

### Benefits

- Easy to test (inject mocks)
- Clear dependencies visible in constructor signature
- No global state
- Flexible swapping of implementations

---

## Monitoring and Debugging

### Logging Strategy

```go
// Each layer can log at appropriate level
driver.Listen(): log.Debug("Received data from HA3100")
app.ConnectMachine(): log.Info("Connected to machine_1")
parser.Parse(): log.Error("Failed to parse ASTM data")
```

### Raw Data Storage

```go
type LabResult struct {
    RawData []byte  // ← Store original bytes for debugging
}

// In debugging, can re-parse failed results
```

---

## Performance Considerations

### Channel Buffering

```go
// 100-result buffer prevents blocking if App slow to process
resultChannel: make(chan *drivers.LabResult, 100)
```

### Connection Pooling

```go
// For PostgreSQL via pgx
pool, _ := pgxpool.New(ctx, connString)
defer pool.Close()
```

### Goroutine Limits

- One Listen() goroutine per machine
- Max machines typically 5-10 (configurable)
- Use context for graceful shutdown

---

## Security Architecture

### Input Validation

```go
// Validate configuration
func (cm *ConfigManager) Validate() error {
    if cm.config.Database.Host == "" {
        return fmt.Errorf("database host not configured")
    }
    // ...
}

// Validate user input (frontend)
export const validatePetForm = (pet: Pet): Error[] => {
    if (!pet.name) return ["Pet name required"];
    if (pet.age < 0) return ["Age cannot be negative"];
    // ...
}
```

### SQL Injection Prevention

```go
// ✅ SAFE: Parameterized query with pgx
rows, err := pool.Query(ctx, 
    "SELECT * FROM results WHERE pet_id = $1", 
    petID)

// ❌ UNSAFE: String concatenation
rows, err := pool.Query(ctx, 
    "SELECT * FROM results WHERE pet_id = " + petID)
```

### HTTPS for Network Communication

```go
// For Scitek over network (implement TLS)
config := &tls.Config{
    Certificates: []tls.Certificate{cert},
}
listener, _ := tls.Listen("tcp", ":5000", config)
```

---

## Future Extensibility

### Adding New Features

1. **New Lab Machine**: Create driver + register in factory
2. **New Database**: Implement Repository interface
3. **New Protocol**: Add parser in `internal/parser/`
4. **New Frontend Screen**: Create feature folder with components

### Zero Impact on Existing Code

All changes isolated to new packages/directories. Existing code untouched.

---

## Build & Release

### Development

```bash
wails dev  # Hot reload, debug mode
```

### Production Build

```bash
wails build -m production -m windows -arch amd64
# Outputs: veta-lis.exe
```

### Cross-Platform

```bash
# macOS
wails build -m production -platform darwin/universal

# Linux
wails build -m production -platform linux/amd64

# Windows
wails build -m production -platform windows/amd64
```

---

**Version:** 1.0.0  
**Last Updated:** January 2026
