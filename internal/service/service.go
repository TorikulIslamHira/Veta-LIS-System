package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"veta-lis/internal/core"
	"veta-lis/internal/drivers"
	"veta-lis/internal/repository"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// LabService provides business logic for processing lab results from hardware drivers
// and persisting them to the database. It acts as the bridge between the driver layer
// and the repository layer, following Clean Architecture principles.
type LabService struct {
	repo *repository.PetRepository
	ctx  context.Context // Wails context for event emission
}

// NewLabService creates a new instance of LabService with the provided repository.
// The repository dependency is injected to maintain loose coupling and testability.
func NewLabService(repo *repository.PetRepository) *LabService {
	return &LabService{
		repo: repo,
	}
}

// StartProcessing begins processing lab results from the provided channel.
// This method starts a background goroutine that continuously listens for incoming
// lab results from hardware drivers and persists them to the database.
//
// Parameters:
//   - resultsChan: A read-only channel that receives lab results from machine drivers
//
// Architecture:
//   - Non-blocking: Runs in a goroutine, allowing caller to continue execution
//   - Context-aware: Can be stopped by canceling the provided context
//   - Error handling: Logs errors but continues processing (fault-tolerant)
//
// Data Flow:
//  1. Machine Driver → Channel → This Service
//  2. Convert driver.LabResult to core.LabResult
//  3. Validate and transform data (PatientID → PetID)
//  4. Persist to database via repository
func (s *LabService) StartProcessing(ctx context.Context, resultsChan <-chan *drivers.LabResult) {
	log.Println("🚀 LabService: Starting result processing goroutine...")

	// Start a background goroutine to process results asynchronously
	go func() {
		for {
			select {
			case <-ctx.Done():
				// Context canceled, stop processing
				log.Println("⏹️  LabService: Stopping result processing (context canceled)")
				return

			case result, ok := <-resultsChan:
				if !ok {
					// Channel closed, exit gracefully
					log.Println("⏹️  LabService: Results channel closed, stopping processing")
					return
				}

				// Process the received result
				if err := s.processResult(ctx, result); err != nil {
					log.Printf("❌ LabService: Failed to process result: %v\n", err)
					// Continue processing despite error (fault tolerance)
				}
			}
		}
	}()
}

// processResult handles the business logic for a single lab result.
// It converts the driver's result format to the database format and persists it.
func (s *LabService) processResult(ctx context.Context, driverResult *drivers.LabResult) error {
	// Log incoming result
	log.Printf("📥 Received result: Test=%s, Value=%s %s, Patient=%s",
		driverResult.TestName,
		driverResult.ResultValue,
		driverResult.ResultUnit,
		driverResult.PetID,
	)

	// Business Logic: Convert PatientID (string) to PetID (int)
	// The hardware driver uses string IDs, but our database uses integer primary keys
	petID, err := s.convertPatientIDToPetID(driverResult.PetID)
	if err != nil {
		// If conversion fails, log the error and mark as unknown patient
		log.Printf("⚠️  Warning: Could not convert PatientID '%s' to PetID: %v", driverResult.PetID, err)
		log.Printf("   Result will be logged but may need manual review")
		// Could implement fallback logic here, e.g., create an "Unknown Patient" entry
		// For now, we'll skip this result
		return fmt.Errorf("invalid patient ID: %w", err)
	}

	// Convert drivers.LabResult to core.LabResult (database format)
	coreResult := s.convertToCoreLabResult(driverResult, petID)

	// Validate the result before saving
	if err := s.validateResult(coreResult); err != nil {
		return fmt.Errorf("result validation failed: %w", err)
	}

	// Persist to database via repository
	if s.repo == nil {
		log.Printf("⚠️  Warning: Repository not initialized, result not saved")
		return fmt.Errorf("repository not initialized")
	}

	if err := s.repo.SaveResult(coreResult); err != nil {
		return fmt.Errorf("database save failed: %w", err)
	}

	log.Printf("✅ Successfully saved result: Test=%s, PetID=%d, Status=%s",
		coreResult.TestName,
		coreResult.PetID,
		coreResult.Status,
	)

	// Emit Wails event to notify frontend of new result
	if s.ctx != nil {
		if err := s.emitNewResultEvent(coreResult); err != nil {
			log.Printf("⚠️  Warning: Failed to emit new_result event: %v", err)
		}
	}

	return nil
}

// convertPatientIDToPetID converts the string PatientID from the driver
// to an integer PetID for the database.
//
// Business Rules:
//   - If PatientID is empty or "0", return error (unknown patient)
//   - If PatientID is not a valid integer, return error
//   - Valid PetID must be > 0
func (s *LabService) convertPatientIDToPetID(patientID string) (int, error) {
	if patientID == "" || patientID == "0" {
		return 0, fmt.Errorf("patient ID is empty or zero")
	}

	petID, err := strconv.Atoi(patientID)
	if err != nil {
		return 0, fmt.Errorf("patient ID '%s' is not a valid number: %w", patientID, err)
	}

	if petID <= 0 {
		return 0, fmt.Errorf("patient ID must be positive, got %d", petID)
	}

	return petID, nil
}

// convertToCoreLabResult transforms a driver.LabResult into a core.LabResult
// suitable for database persistence.
func (s *LabService) convertToCoreLabResult(driverResult *drivers.LabResult, petID int) *core.LabResult {
	return &core.LabResult{
		PetID:       petID,
		MachineName: driverResult.MachineType,
		TestName:    driverResult.TestName,
		ResultValue: driverResult.ResultValue,
		Unit:        driverResult.ResultUnit,
		Status:      s.normalizeStatus(driverResult.Status),
		CreatedAt:   s.parseTimestamp(driverResult.Timestamp),
	}
}

// normalizeStatus ensures status values are consistent
// ("NORMAL", "ABNORMAL", "CRITICAL", "PENDING")
func (s *LabService) normalizeStatus(status string) string {
	switch status {
	case "NORMAL", "ABNORMAL", "CRITICAL", "PENDING":
		return status
	case "N":
		return "NORMAL"
	case "H", "L": // High or Low
		return "ABNORMAL"
	case "C":
		return "CRITICAL"
	default:
		return "PENDING" // Default fallback
	}
}

// parseTimestamp converts the ISO 8601 timestamp from the driver to time.Time
func (s *LabService) parseTimestamp(timestamp string) time.Time {
	// Try parsing ISO 8601 format
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		// If parsing fails, use current time
		log.Printf("⚠️  Warning: Could not parse timestamp '%s', using current time", timestamp)
		return time.Now()
	}
	return t
}

// validateResult performs business validation on a lab result before saving
func (s *LabService) validateResult(result *core.LabResult) error {
	if result.PetID <= 0 {
		return fmt.Errorf("invalid pet ID: %d", result.PetID)
	}

	if result.TestName == "" {
		return fmt.Errorf("test name cannot be empty")
	}

	if result.ResultValue == "" {
		return fmt.Errorf("result value cannot be empty")
	}

	return nil
}

// RegisterPatient is a wrapper method for registering a new patient (pet).
// It delegates to the repository layer but provides a cleaner API for the frontend.
//
// This method can be called from the Wails backend to expose patient registration
// functionality to the React frontend.
func (s *LabService) RegisterPatient(ctx context.Context, pet *core.Pet) (int, error) {
	log.Printf("📝 Registering new patient: Name=%s, Species=%s, Owner=%d",
		pet.Name,
		pet.Species,
		pet.OwnerID,
	)
	// Check repository is initialized
	if s.repo == nil {
		return 0, fmt.Errorf("repository not initialized")
	}
	// Delegate to repository
	petID, err := s.repo.CreatePet(pet)
	if err != nil {
		return 0, fmt.Errorf("failed to register patient: %w", err)
	}

	log.Printf("✅ Patient registered successfully: ID=%d, Name=%s", petID, pet.Name)
	return petID, nil
}

// RegisterOwner is a wrapper method for registering a new pet owner.
// It delegates to the repository layer.
func (s *LabService) RegisterOwner(ctx context.Context, owner *core.Owner) (int, error) {
	log.Printf("📝 Registering new owner: Name=%s, Phone=%s",
		owner.Name,
		owner.Phone,
	)

	// Check repository is initialized
	if s.repo == nil {
		return 0, fmt.Errorf("repository not initialized")
	}

	// Delegate to repository
	ownerID, err := s.repo.CreateOwner(owner)
	if err != nil {
		return 0, fmt.Errorf("failed to register owner: %w", err)
	}

	log.Printf("✅ Owner registered successfully: ID=%d, Name=%s", ownerID, owner.Name)
	return ownerID, nil
}

// GetPetsByOwner retrieves all pets belonging to a specific owner
func (s *LabService) GetPetsByOwner(ctx context.Context, ownerID int) ([]*core.Pet, error) {
	// Check repository is initialized
	if s.repo == nil {
		return nil, fmt.Errorf("repository not initialized")
	}

	// This would require a new repository method
	// For now, we can filter from GetPets() result
	allPets, err := s.repo.GetPets()
	if err != nil {
		return nil, fmt.Errorf("failed to get pets: %w", err)
	}

	var ownerPets []*core.Pet
	for _, pet := range allPets {
		if pet.OwnerID == ownerID {
			ownerPets = append(ownerPets, pet)
		}
	}

	return ownerPets, nil
}

// GetResultsByPet retrieves all lab results for a specific pet
func (s *LabService) GetResultsByPet(ctx context.Context, petID int) ([]*core.LabResult, error) {
	log.Printf("🔍 Fetching results for PetID=%d", petID)

	// Handle nil repository
	if s.repo == nil {
		return nil, fmt.Errorf("repository not initialized")
	}

	results, err := s.repo.GetResultsByPetID(petID)
	if err != nil {
		return nil, fmt.Errorf("failed to get results: %w", err)
	}

	log.Printf("✅ Found %d results for PetID=%d", len(results), petID)
	return results, nil
}

// SetContext sets the Wails runtime context for event emission.
// This must be called from the Wails App.Startup method to enable real-time events.
func (s *LabService) SetContext(ctx context.Context) {
	s.ctx = ctx
	log.Println("🔗 LabService: Wails context set for event emission")
}

// emitNewResultEvent emits a Wails event to notify the React frontend of a new result.
// This enables real-time updates in the UI without polling.
func (s *LabService) emitNewResultEvent(result *core.LabResult) error {
	if s.ctx == nil {
		return fmt.Errorf("context not initialized")
	}

	// Emit event to React frontend via Wails runtime
	runtime.EventsEmit(s.ctx, "new_result", result)

	log.Printf("📡 Event 'new_result' emitted for PetID=%d", result.PetID)
	return nil
}

// StopProcessing provides a way to gracefully stop the service
// This is useful for cleanup during application shutdown
func (s *LabService) StopProcessing(ctx context.Context) error {
	log.Println("🛑 LabService: Initiating graceful shutdown...")
	// Context cancellation will trigger the goroutine to exit
	// Additional cleanup logic can be added here if needed
	return nil
}
