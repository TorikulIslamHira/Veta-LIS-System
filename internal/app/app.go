package app

import (
	"context"
	"fmt"
	"strconv"

	"veta-lis/internal/core"
	"veta-lis/internal/service"
)

// App is the main Wails application struct that bridges backend logic to React frontend.
// It holds the service layer and context for the application lifecycle.
type App struct {
	ctx        context.Context
	labService *service.LabService
}

// NewApp creates a new Wails application instance with service layer injection.
// The service parameter provides all business logic and database operations.
//
// Parameters:
//   - labService: The service layer instance that handles all business logic
//
// Returns:
//   - Configured App instance ready to be bound to Wails
func NewApp(labService *service.LabService) *App {
	return &App{
		labService: labService,
	}
}

// Startup is called when the Wails application starts.
// It receives and stores the application context for use in service calls.
//
// This is a Wails lifecycle method - do not call directly.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// Pass context to service layer for Wails event emission
	if a.labService != nil {
		a.labService.SetContext(ctx)
	}

	println("🚀 Veta LIS Application Started")
}

// Shutdown is called when the Wails application is closing.
// Perform cleanup operations here.
//
// This is a Wails lifecycle method - do not call directly.
func (a *App) Shutdown(ctx context.Context) {
	println("🛑 Veta LIS Application Shutting Down")
	// Additional cleanup can be added here
}

// ============================================================================
// FRONTEND EXPOSED METHODS - React can call these via Wails bindings
// ============================================================================

// RegisterPatient registers a new patient (pet) and owner in the system.
// This method is exposed to the React frontend via Wails binding.
//
// Parameters:
//   - name: Pet's name (e.g., "Max", "Whiskers")
//   - species: Pet's species (e.g., "Dog", "Cat", "Bird")
//   - age: Pet's age as string (will be converted to int)
//   - gender: Pet's gender (e.g., "Male", "Female", "Unknown")
//   - ownerName: Owner's full name
//   - ownerPhone: Owner's phone number (used as unique identifier)
//
// Returns:
//   - Pet ID if registration successful
//   - Error if validation or database operation fails
//
// Example from React:
//
//	const petId = await RegisterPatient("Max", "Dog", "3", "Male", "John Doe", "555-1234");
func (a *App) RegisterPatient(
	name string,
	species string,
	age string,
	gender string,
	ownerName string,
	ownerPhone string,
) (int, error) {
	// Validate inputs
	if name == "" {
		return 0, fmt.Errorf("pet name cannot be empty")
	}
	if species == "" {
		return 0, fmt.Errorf("species cannot be empty")
	}
	if ownerPhone == "" {
		return 0, fmt.Errorf("owner phone cannot be empty")
	}

	// Convert age string to integer
	ageInt, err := strconv.Atoi(age)
	if err != nil {
		return 0, fmt.Errorf("invalid age format: must be a number")
	}

	if ageInt < 0 {
		return 0, fmt.Errorf("age cannot be negative")
	}

	// Step 1: Create or find owner
	owner := &core.Owner{
		Name:    ownerName,
		Phone:   ownerPhone,
		Address: "", // Can be added later or made optional
	}

	// Register owner via service layer
	ownerID, err := a.labService.RegisterOwner(a.ctx, owner)
	if err != nil {
		return 0, fmt.Errorf("failed to register owner: %w", err)
	}

	// Step 2: Create pet record
	pet := &core.Pet{
		Name:    name,
		Species: species,
		Breed:   "Not specified", // Can be enhanced to accept breed parameter
		Age:     ageInt,
		Gender:  gender,
		OwnerID: ownerID,
	}

	// Register pet via service layer
	petID, err := a.labService.RegisterPatient(a.ctx, pet)
	if err != nil {
		return 0, fmt.Errorf("failed to register patient: %w", err)
	}

	println(fmt.Sprintf("✅ Patient registered: %s (ID: %d), Owner: %s", name, petID, ownerName))
	return petID, nil
}

// GetRecentPets retrieves all registered pets with their owner information.
// This method is exposed to the React frontend for displaying patient lists.
//
// Returns:
//   - Array of Pet objects
//   - Error if database query fails
//
// Example from React:
//
//	const pets = await GetRecentPets();
//	pets.forEach(pet => console.log(pet.name, pet.species));
func (a *App) GetRecentPets() ([]*core.Pet, error) {
	// Get all pets from service layer (0 means get all)
	pets, err := a.labService.GetPetsByOwner(a.ctx, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pets: %w", err)
	}

	println(fmt.Sprintf("📋 Retrieved %d pets", len(pets)))
	return pets, nil
}

// GetPetsByOwner retrieves all pets belonging to a specific owner.
// Useful for filtering patient lists by owner.
//
// Parameters:
//   - ownerID: The owner's database ID
//
// Returns:
//   - Array of Pet objects belonging to that owner
//   - Error if query fails
func (a *App) GetPetsByOwner(ownerID int) ([]*core.Pet, error) {
	pets, err := a.labService.GetPetsByOwner(a.ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve owner's pets: %w", err)
	}

	return pets, nil
}

// GetPetResults retrieves all lab results for a specific pet.
// Used to display patient test history in the frontend.
//
// Parameters:
//   - petID: The pet's database ID
//
// Returns:
//   - Array of LabResult objects
//   - Error if query fails
//
// Example from React:
//
//	const results = await GetPetResults(petId);
func (a *App) GetPetResults(petID int) ([]*core.LabResult, error) {
	if petID <= 0 {
		return nil, fmt.Errorf("invalid pet ID")
	}

	results, err := a.labService.GetResultsByPet(a.ctx, petID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve test results: %w", err)
	}

	println(fmt.Sprintf("🔬 Retrieved %d test results for Pet ID: %d", len(results), petID))
	return results, nil
}

// GetDashboardStats returns system statistics for the dashboard display.
// Provides an overview of system activity and current state.
//
// Returns:
//   - Map containing various statistics
//   - Error if stats cannot be generated
//
// Example from React:
//
//	const stats = await GetDashboardStats();
//	console.log(`Total Tests: ${stats.total_tests}`);
func (a *App) GetDashboardStats() (map[string]any, error) {
	// TODO: Implement real statistics from database queries
	// For now, return dummy data as requested

	stats := map[string]any{
		"total_tests": 150,
		"pending":     5,
		"completed":   145,
		// Additional stats that might be useful
		"active_machines":  0,
		"registered_pets":  0,
		"critical_results": 0,
	}

	println("📊 Dashboard stats retrieved")
	return stats, nil
}

// ============================================================================
// ADDITIONAL HELPER METHODS
// ============================================================================

// RegisterPatientWithBreed is an extended version that includes breed information.
// Use this for more detailed patient registration.
func (a *App) RegisterPatientWithBreed(
	name string,
	species string,
	breed string,
	age int,
	gender string,
	ownerName string,
	ownerPhone string,
	ownerAddress string,
) (int, error) {
	// Create owner
	owner := &core.Owner{
		Name:    ownerName,
		Phone:   ownerPhone,
		Address: ownerAddress,
	}

	ownerID, err := a.labService.RegisterOwner(a.ctx, owner)
	if err != nil {
		return 0, fmt.Errorf("failed to register owner: %w", err)
	}

	// Create pet
	pet := &core.Pet{
		Name:    name,
		Species: species,
		Breed:   breed,
		Age:     age,
		Gender:  gender,
		OwnerID: ownerID,
	}

	petID, err := a.labService.RegisterPatient(a.ctx, pet)
	if err != nil {
		return 0, fmt.Errorf("failed to register patient: %w", err)
	}

	return petID, nil
}

// Greet returns a greeting message.
// This is the default Wails example method - can be removed in production.
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, welcome to Veta LIS!", name)
}
