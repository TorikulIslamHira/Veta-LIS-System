package core

import (
	"fmt"
	"strconv"
	"strings"

	"veta-lis/internal/drivers"
)

// ResultValidator provides business logic for validating and interpreting lab results.
type ResultValidator struct {
	referenceRanges map[string]ReferenceRange
}

// ReferenceRange defines the normal range for a specific test.
type ReferenceRange struct {
	TestCode  string
	TestName  string
	Unit      string
	MinValue  float64
	MaxValue  float64
	CriticalMin float64 // Below this is critical
	CriticalMax float64 // Above this is critical
}

// NewResultValidator creates a new result validator.
func NewResultValidator() *ResultValidator {
	return &ResultValidator{
		referenceRanges: make(map[string]ReferenceRange),
	}
}

// RegisterReferenceRange registers a reference range for a test.
func (rv *ResultValidator) RegisterReferenceRange(rr ReferenceRange) {
	rv.referenceRanges[rr.TestCode] = rr
}

// ValidateResult validates a lab result and sets its status.
func (rv *ResultValidator) ValidateResult(result *drivers.LabResult) error {
	if result == nil {
		return fmt.Errorf("result is nil")
	}

	// Try to parse the result value as float
	value, err := strconv.ParseFloat(strings.TrimSpace(result.ResultValue), 64)
	if err != nil {
		// If parsing fails, mark as pending (could be a text result)
		result.Status = "PENDING"
		return nil
	}

	// Look up reference range
	refRange, exists := rv.referenceRanges[result.TestCode]
	if !exists {
		// No reference range found, mark as pending
		result.Status = "PENDING"
		return nil
	}

	// Determine status based on value
	if value < refRange.CriticalMin || value > refRange.CriticalMax {
		result.Status = "CRITICAL"
	} else if value < refRange.MinValue || value > refRange.MaxValue {
		result.Status = "ABNORMAL"
	} else {
		result.Status = "NORMAL"
	}

	return nil
}

// PetService handles pet-related business logic.
type PetService struct {
	// In production, this would have repository dependencies
}

// Pet represents a pet in the system.
type Pet struct {
	ID       string
	Name     string
	Species  string
	Breed    string
	Age      int
	Weight   float64
	OwnerID  string
}

// NewPetService creates a new pet service.
func NewPetService() *PetService {
	return &PetService{}
}

// ValidatePetData validates pet information before saving.
func (ps *PetService) ValidatePetData(pet *Pet) error {
	if pet.Name == "" {
		return fmt.Errorf("pet name is required")
	}

	if pet.Species == "" {
		return fmt.Errorf("pet species is required")
	}

	if pet.OwnerID == "" {
		return fmt.Errorf("owner ID is required")
	}

	if pet.Age < 0 {
		return fmt.Errorf("pet age cannot be negative")
	}

	if pet.Weight <= 0 {
		return fmt.Errorf("pet weight must be positive")
	}

	return nil
}

// Owner represents a pet owner in the system.
type Owner struct {
	ID       string
	Name     string
	Email    string
	Phone    string
	Address  string
}

// OwnerService handles owner-related business logic.
type OwnerService struct{}

// NewOwnerService creates a new owner service.
func NewOwnerService() *OwnerService {
	return &OwnerService{}
}

// ValidateOwnerData validates owner information before saving.
func (os *OwnerService) ValidateOwnerData(owner *Owner) error {
	if owner.Name == "" {
		return fmt.Errorf("owner name is required")
	}

	if owner.Email == "" && owner.Phone == "" {
		return fmt.Errorf("email or phone is required")
	}

	return nil
}
