package app

import (
	"context"
	"sync"

	"veta-lis/internal/config"
	"veta-lis/internal/core"
	"veta-lis/internal/drivers"
	"veta-lis/internal/drivers/ha3100"
	"veta-lis/internal/drivers/scitek"
)

// App is the main Wails application struct that will be bound to the frontend.
type App struct {
	ctx                context.Context
	cancel             context.CancelFunc
	configManager      *config.ConfigManager
	resultValidator    *core.ResultValidator
	petService         *core.PetService
	ownerService       *core.OwnerService
	driverFactory      *drivers.DriverFactory
	activeDrivers      map[string]drivers.MachineDriver
	activeDriversMutex sync.RWMutex
	resultChannel      chan *drivers.LabResult
}

// NewApp creates a new Wails application instance.
func NewApp(configPath string) *App {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize driver factory
	factory := drivers.NewDriverFactory()
	factory.Register("HA3100", ha3100.NewHA3100Driver)
	factory.Register("Scitek", scitek.NewScitekDriver)

	app := &App{
		ctx:            ctx,
		cancel:         cancel,
		configManager:  config.NewConfigManager(configPath),
		resultValidator: core.NewResultValidator(),
		petService:     core.NewPetService(),
		ownerService:   core.NewOwnerService(),
		driverFactory:  factory,
		activeDrivers:  make(map[string]drivers.MachineDriver),
		resultChannel:  make(chan *drivers.LabResult, 100),
	}

	return app
}

// Startup is called at application startup.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// Load configuration
	if err := a.configManager.Load(); err != nil {
		// Log error but don't crash
		println("Error loading config:", err)
	}

	// Initialize reference ranges
	a.initializeReferenceRanges()
}

// Shutdown is called at application shutdown.
func (a *App) Shutdown(ctx context.Context) {
	a.cancel()
	// Disconnect all active drivers
	for _, driver := range a.activeDrivers {
		_ = driver.Disconnect(ctx)
	}
}

// ConnectMachine connects to a lab machine.
// Called from the frontend via Wails binding.
func (a *App) ConnectMachine(machineID string) (string, error) {
	// Get machine configuration
	machineConfig, err := a.configManager.GetMachineByID(machineID)
	if err != nil {
		return "", err
	}

	// Create driver instance
	driver, err := a.driverFactory.Create(machineConfig.Type)
	if err != nil {
		return "", err
	}

	// Set configuration (type assertion depends on implementation)
	// In a real implementation, we'd have a SetConfig method on all drivers
	if err := driver.Connect(a.ctx); err != nil {
		return "", err
	}

	// Store active driver
	a.activeDriversMutex.Lock()
	a.activeDrivers[machineID] = driver
	a.activeDriversMutex.Unlock()

	// Start listening in a goroutine
	go func() {
		_ = driver.Listen(a.ctx, a.resultChannel)
	}()

	info, err := driver.GetMachineInfo(a.ctx)
	if err != nil {
		return "", err
	}

	return info.Status, nil
}

// DisconnectMachine disconnects from a lab machine.
func (a *App) DisconnectMachine(machineID string) error {
	a.activeDriversMutex.Lock()
	driver, exists := a.activeDrivers[machineID]
	a.activeDriversMutex.Unlock()

	if !exists {
		return drivers.ErrNotConnected
	}

	if err := driver.Disconnect(a.ctx); err != nil {
		return err
	}

	a.activeDriversMutex.Lock()
	delete(a.activeDrivers, machineID)
	a.activeDriversMutex.Unlock()

	return nil
}

// GetMachineStatus returns the status of a connected machine.
func (a *App) GetMachineStatus(machineID string) (map[string]interface{}, error) {
	a.activeDriversMutex.RLock()
	driver, exists := a.activeDrivers[machineID]
	a.activeDriversMutex.RUnlock()

	if !exists {
		return nil, drivers.ErrNotConnected
	}

	info, err := driver.GetMachineInfo(a.ctx)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"type":             info.Type,
		"model":            info.Model,
		"connected":        driver.IsConnected(),
		"status":           info.Status,
		"connectionType":   info.ConnectionType,
		"port":             info.Port,
		"firmwareVersion":  info.FirmwareVersion,
		"lastTestTime":     info.LastTestTime,
		"totalTestsRun":    info.TotalTestsRun,
	}, nil
}

// GetAllMachines returns configuration of all machines.
func (a *App) GetAllMachines() []map[string]interface{} {
	machines := a.configManager.GetMachines()
	var result []map[string]interface{}

	for _, m := range machines {
		a.activeDriversMutex.RLock()
		_, connected := a.activeDrivers[m.ID]
		a.activeDriversMutex.RUnlock()

		result = append(result, map[string]interface{}{
			"id":             m.ID,
			"type":           m.Type,
			"connectionType": m.ConnectionType,
			"port":           m.Port,
			"enabled":        m.Enabled,
			"connected":      connected,
		})
	}

	return result
}

// SaveMachineConfig saves machine configuration.
func (a *App) SaveMachineConfig(id string, connectionType string, port string, baudRate int) error {
	machineConfig, err := a.configManager.GetMachineByID(id)
	if err != nil {
		return err
	}

	machineConfig.ConnectionType = connectionType
	machineConfig.Port = port
	machineConfig.BaudRate = baudRate

	if err := a.configManager.UpdateMachine(id, *machineConfig); err != nil {
		return err
	}

	return a.configManager.Save()
}

// GetRecentResults returns recent lab results (from result channel).
func (a *App) GetRecentResults() []*drivers.LabResult {
	var results []*drivers.LabResult

	// Non-blocking drain of result channel
	for {
		select {
		case result := <-a.resultChannel:
			results = append(results, result)
		default:
			return results
		}
	}
}

// initializeReferenceRanges sets up default reference ranges for common tests.
func (a *App) initializeReferenceRanges() {
	// Example reference ranges for common veterinary tests
	a.resultValidator.RegisterReferenceRange(core.ReferenceRange{
		TestCode:    "GLU",
		TestName:    "Glucose",
		Unit:        "mg/dL",
		MinValue:    70,
		MaxValue:    100,
		CriticalMin: 40,
		CriticalMax: 400,
	})

	a.resultValidator.RegisterReferenceRange(core.ReferenceRange{
		TestCode:    "BUN",
		TestName:    "Blood Urea Nitrogen",
		Unit:        "mg/dL",
		MinValue:    7,
		MaxValue:    20,
		CriticalMin: 1,
		CriticalMax: 100,
	})

	a.resultValidator.RegisterReferenceRange(core.ReferenceRange{
		TestCode:    "CRE",
		TestName:    "Creatinine",
		Unit:        "mg/dL",
		MinValue:    0.5,
		MaxValue:    1.5,
		CriticalMin: 0.1,
		CriticalMax: 10,
	})
}
