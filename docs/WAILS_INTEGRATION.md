# Wails App Integration Guide

## Updated Architecture

The `internal/app/app.go` has been updated to integrate with the service layer and expose backend functionality to the React frontend.

## Changes Made

### 1. App Struct Updated
```go
type App struct {
    // ... existing fields ...
    labService *service.LabService  // NEW: Service layer for DB operations
}
```

### 2. Constructor Updated
```go
func NewApp(configPath string, labService *service.LabService) *App
```
**Important:** You must now pass a `*service.LabService` instance when creating the app.

### 3. Auto-start Processing
The `Startup()` method now automatically starts the lab result processing goroutine:
```go
func (a *App) Startup(ctx context.Context) {
    // ... existing code ...
    a.labService.StartProcessing(ctx, a.resultChannel)
}
```

## Methods Exposed to React Frontend

### Patient Management

#### `RegisterPatient(name, species, age, ownerPhone) (int, error)`
Simple patient registration with minimal fields.

**Example from React:**
```typescript
const petId = await RegisterPatient("Max", "Dog", "3", "555-1234");
```

#### `RegisterPatientFull(...) (int, error)`
Complete patient registration with all details.

**Example from React:**
```typescript
const petId = await RegisterPatientFull(
    "Max",           // petName
    "Dog",           // species
    "Golden Retriever", // breed
    3,               // age
    "Male",          // gender
    "John Doe",      // ownerName
    "555-1234",      // ownerPhone
    "123 Main St"    // ownerAddress
);
```

#### `GetRecentPets() ([]map[string]interface{}, error)`
Retrieves all registered pets.

**Returns:**
```json
[
  {
    "id": 1,
    "name": "Max",
    "species": "Dog",
    "breed": "Golden Retriever",
    "age": 3,
    "gender": "Male",
    "owner_id": 1
  }
]
```

#### `GetPetResults(petID int) ([]map[string]interface{}, error)`
Gets all lab results for a specific pet.

**Returns:**
```json
[
  {
    "id": 1,
    "pet_id": 1,
    "machine_name": "Scitek",
    "test_name": "Glucose",
    "result_value": "105",
    "unit": "mg/dL",
    "status": "NORMAL",
    "created_at": "2026-01-20 15:30:00"
  }
]
```

### Dashboard

#### `GetDashboardStats() (map[string]interface{}, error)`
Returns system statistics for dashboard display.

**Returns:**
```json
{
  "total_tests": 120,
  "pending": 5,
  "completed_today": 25,
  "active_machines": 2,
  "registered_pets": 45,
  "critical_results": 2
}
```

## Main.go Integration Example

Update your `main.go` to provide the service layer:

```go
package main

import (
    "context"
    "embed"
    "log"

    "github.com/wailsapp/wails/v2"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
    
    "veta-lis/internal/app"
    "veta-lis/internal/repository"
    "veta-lis/internal/service"
    
    "github.com/jackc/pgx/v5/pgxpool"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    // 1. Initialize database connection
    dbPool, err := pgxpool.New(context.Background(), "postgresql://user:pass@localhost/veta_lis")
    if err != nil {
        log.Fatal("Database connection failed:", err)
    }
    defer dbPool.Close()

    // 2. Create repository layer
    repo := repository.NewPetRepository(dbPool)

    // 3. Create service layer
    labService := service.NewLabService(repo)

    // 4. Create Wails app with service injection
    wailsApp := app.NewApp("config.yaml", labService)

    // 5. Run Wails application
    err = wails.Run(&options.App{
        Title:  "Veta LIS",
        Width:  1024,
        Height: 768,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        Bind: []interface{}{
            wailsApp,
        },
        OnStartup: wailsApp.Startup,
        OnShutdown: wailsApp.Shutdown,
    })

    if err != nil {
        log.Fatal(err)
    }
}
```

## React Frontend Usage

After running `wails dev`, Wails generates TypeScript bindings in `frontend/wailsjs/go/app/App.js`.

### Example React Component

```typescript
import { RegisterPatient, GetRecentPets, GetDashboardStats } from '../wailsjs/go/app/App';

function Dashboard() {
  const [stats, setStats] = useState(null);
  const [pets, setPets] = useState([]);

  useEffect(() => {
    // Load dashboard stats
    GetDashboardStats().then(setStats);
    
    // Load recent pets
    GetRecentPets().then(setPets);
  }, []);

  const handleRegister = async () => {
    try {
      const petId = await RegisterPatient("Max", "Dog", "3", "555-1234");
      alert(`Pet registered with ID: ${petId}`);
      // Refresh pet list
      GetRecentPets().then(setPets);
    } catch (error) {
      alert(`Error: ${error}`);
    }
  };

  return (
    <div>
      <h1>Dashboard</h1>
      <div>Total Tests: {stats?.total_tests}</div>
      <div>Pending: {stats?.pending}</div>
      
      <h2>Recent Pets</h2>
      <ul>
        {pets.map(pet => (
          <li key={pet.id}>{pet.name} ({pet.species})</li>
        ))}
      </ul>
      
      <button onClick={handleRegister}>Register New Pet</button>
    </div>
  );
}
```

## Data Flow

```
React Frontend
    ↓ (calls RegisterPatient)
Wails Binding
    ↓
App.RegisterPatient() [app.go]
    ↓
LabService.RegisterPatient() [service/service.go]
    ↓
PetRepository.CreatePet() [repository/pet_repo.go]
    ↓
PostgreSQL Database
```

## Auto-Processing Results

Lab results from machines are automatically processed:

```
Machine Driver → resultChannel → LabService.StartProcessing() → Database
```

No manual intervention needed - results are automatically saved when machines send data!
