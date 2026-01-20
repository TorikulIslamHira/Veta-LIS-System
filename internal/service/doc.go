// Package service provides the business logic layer for the Veta LIS system.
//
// Example Usage:
//
//	// 1. Initialize database connection
//	dbPool, err := pgxpool.New(ctx, databaseURL)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer dbPool.Close()
//
//	// 2. Create repository
//	repo := repository.NewPetRepository(dbPool)
//
//	// 3. Create service with repository dependency
//	labService := service.NewLabService(repo)
//
//	// 4. Create a channel for lab results
//	resultsChan := make(chan *drivers.LabResult, 100)
//
//	// 5. Start processing results in background
//	ctx := context.Background()
//	labService.StartProcessing(ctx, resultsChan)
//
//	// 6. Initialize and start machine driver
//	driver := scitek.NewScitekDriver("COM3", 9600)
//	if err := driver.Connect(ctx); err != nil {
//	    log.Fatal(err)
//	}
//
//	// 7. Driver sends results to the channel
//	go driver.Listen(ctx, resultsChan)
//
//	// Results flow: Machine → Driver → Channel → Service → Database
//
// Architecture:
//   - Service Layer sits between Drivers and Repository
//   - Handles data transformation and business logic
//   - Provides clean API for frontend integration
//   - Supports multiple concurrent drivers
package service

// Example integration in Wails app.go:
//
// type App struct {
//     ctx        context.Context
//     labService *service.LabService
//     drivers    map[string]drivers.MachineDriver
// }
//
// func NewApp(labService *service.LabService) *App {
//     return &App{
//         labService: labService,
//         drivers:    make(map[string]drivers.MachineDriver),
//     }
// }
//
// // Exposed to React frontend
// func (a *App) RegisterPet(name, species, breed string, age int, ownerID int) (int, error) {
//     pet := &core.Pet{
//         Name:    name,
//         Species: species,
//         Breed:   breed,
//         Age:     age,
//         OwnerID: ownerID,
//     }
//     return a.labService.RegisterPatient(a.ctx, pet)
// }
//
// func (a *App) GetPetResults(petID int) ([]*core.LabResult, error) {
//     return a.labService.GetResultsByPet(a.ctx, petID)
// }
