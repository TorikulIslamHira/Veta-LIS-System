package core

import (
	"errors"
	"time"
)

// Owner represents a pet owner in the veterinary system
type Owner struct {
	ID      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Phone   string `json:"phone" db:"phone"`
	Address string `json:"address" db:"address"`
}

// Pet represents a pet in the veterinary system
type Pet struct {
	ID      int    `json:"id" db:"id"`
	OwnerID int    `json:"owner_id" db:"owner_id"`
	Name    string `json:"name" db:"name"`
	Species string `json:"species" db:"species"`
	Breed   string `json:"breed" db:"breed"`
	Age     int    `json:"age" db:"age"`
	Gender  string `json:"gender" db:"gender"`
}

// Validate checks if the Pet struct has valid required fields
func (p *Pet) Validate() error {
	if p.Name == "" {
		return errors.New("pet name cannot be empty")
	}
	if p.Species == "" {
		return errors.New("pet species cannot be empty")
	}
	return nil
}

// LabResult represents a laboratory test result for a pet
type LabResult struct {
	ID          int       `json:"id" db:"id"`
	PetID       int       `json:"pet_id" db:"pet_id"`
	MachineName string    `json:"machine_name" db:"machine_name"`
	TestName    string    `json:"test_name" db:"test_name"`
	ResultValue string    `json:"result_value" db:"result_value"`
	Unit        string    `json:"unit" db:"unit"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
