package repository

import (
	"context"
	"fmt"

	"veta-lis/internal/core"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PetRepository handles all database operations for pets, owners, and lab results
type PetRepository struct {
	db *pgxpool.Pool
}

// NewPetRepository creates a new instance of PetRepository
func NewPetRepository(db *pgxpool.Pool) *PetRepository {
	return &PetRepository{
		db: db,
	}
}

// PetWithOwner represents a pet with associated owner information
type PetWithOwner struct {
	core.Pet
	OwnerName    string `json:"owner_name"`
	OwnerPhone   string `json:"owner_phone"`
	OwnerAddress string `json:"owner_address"`
}

// CreatePet creates a new pet record in the database
// If the pet has an OwnerID, it uses that; otherwise it would need owner creation logic
func (r *PetRepository) CreatePet(pet *core.Pet) (int, error) {
	ctx := context.Background()

	// Validate pet data
	if err := pet.Validate(); err != nil {
		return 0, fmt.Errorf("validation failed: %w", err)
	}

	// Insert pet record with parameterized query to prevent SQL injection
	query := `
		INSERT INTO pets (owner_id, name, species, breed, age, gender)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var petID int
	err := r.db.QueryRow(
		ctx,
		query,
		pet.OwnerID,
		pet.Name,
		pet.Species,
		pet.Breed,
		pet.Age,
		pet.Gender,
	).Scan(&petID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert pet: %w", err)
	}

	return petID, nil
}

// CreateOwner creates a new owner record in the database
func (r *PetRepository) CreateOwner(owner *core.Owner) (int, error) {
	ctx := context.Background()

	// Insert owner record with parameterized query
	query := `
		INSERT INTO owners (name, phone, address)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var ownerID int
	err := r.db.QueryRow(
		ctx,
		query,
		owner.Name,
		owner.Phone,
		owner.Address,
	).Scan(&ownerID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert owner: %w", err)
	}

	return ownerID, nil
}

// GetOwnerByPhone retrieves an owner by phone number (useful for linking existing owners)
func (r *PetRepository) GetOwnerByPhone(phone string) (*core.Owner, error) {
	ctx := context.Background()

	query := `
		SELECT id, name, phone, address
		FROM owners
		WHERE phone = $1
		LIMIT 1
	`

	var owner core.Owner
	err := r.db.QueryRow(ctx, query, phone).Scan(
		&owner.ID,
		&owner.Name,
		&owner.Phone,
		&owner.Address,
	)

	if err != nil {
		return nil, fmt.Errorf("owner not found: %w", err)
	}

	return &owner, nil
}

// GetPets retrieves all pets with their owner information using JOIN
func (r *PetRepository) GetPets() ([]*core.Pet, error) {
	ctx := context.Background()

	// Use JOIN to get pet and owner data together
	query := `
		SELECT 
			p.id, p.owner_id, p.name, p.species, p.breed, p.age, p.gender,
			o.name as owner_name, o.phone as owner_phone, o.address as owner_address
		FROM pets p
		INNER JOIN owners o ON p.owner_id = o.id
		ORDER BY p.id DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pets: %w", err)
	}
	defer rows.Close()

	var pets []*core.Pet
	for rows.Next() {
		var pet core.Pet
		var ownerName, ownerPhone, ownerAddress string

		err := rows.Scan(
			&pet.ID,
			&pet.OwnerID,
			&pet.Name,
			&pet.Species,
			&pet.Breed,
			&pet.Age,
			&pet.Gender,
			&ownerName,
			&ownerPhone,
			&ownerAddress,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan pet row: %w", err)
		}

		pets = append(pets, &pet)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating pet rows: %w", err)
	}

	return pets, nil
}

// GetPetsWithOwnerInfo retrieves all pets with their owner information
// This version returns PetWithOwner structs that include owner details
func (r *PetRepository) GetPetsWithOwnerInfo() ([]*PetWithOwner, error) {
	ctx := context.Background()

	query := `
		SELECT 
			p.id, p.owner_id, p.name, p.species, p.breed, p.age, p.gender,
			o.name as owner_name, o.phone as owner_phone, o.address as owner_address
		FROM pets p
		INNER JOIN owners o ON p.owner_id = o.id
		ORDER BY p.id DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pets: %w", err)
	}
	defer rows.Close()

	var pets []*PetWithOwner
	for rows.Next() {
		var pet PetWithOwner

		err := rows.Scan(
			&pet.Pet.ID,
			&pet.Pet.OwnerID,
			&pet.Pet.Name,
			&pet.Pet.Species,
			&pet.Pet.Breed,
			&pet.Pet.Age,
			&pet.Pet.Gender,
			&pet.OwnerName,
			&pet.OwnerPhone,
			&pet.OwnerAddress,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan pet row: %w", err)
		}

		pets = append(pets, &pet)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating pet rows: %w", err)
	}

	return pets, nil
}

// SaveResult saves a lab result to the database
func (r *PetRepository) SaveResult(result *core.LabResult) error {
	ctx := context.Background()

	// Insert lab result with parameterized query to prevent SQL injection
	query := `
		INSERT INTO lab_results (pet_id, machine_name, test_name, result_value, unit, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var resultID int
	err := r.db.QueryRow(
		ctx,
		query,
		result.PetID,
		result.MachineName,
		result.TestName,
		result.ResultValue,
		result.Unit,
		result.Status,
		result.CreatedAt,
	).Scan(&resultID)

	if err != nil {
		return fmt.Errorf("failed to insert lab result: %w", err)
	}

	result.ID = resultID
	return nil
}

// GetResultsByPetID retrieves all lab results for a specific pet
func (r *PetRepository) GetResultsByPetID(petID int) ([]*core.LabResult, error) {
	ctx := context.Background()

	query := `
		SELECT id, pet_id, machine_name, test_name, result_value, unit, status, created_at
		FROM lab_results
		WHERE pet_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, petID)
	if err != nil {
		return nil, fmt.Errorf("failed to query lab results: %w", err)
	}
	defer rows.Close()

	var results []*core.LabResult
	for rows.Next() {
		var result core.LabResult

		err := rows.Scan(
			&result.ID,
			&result.PetID,
			&result.MachineName,
			&result.TestName,
			&result.ResultValue,
			&result.Unit,
			&result.Status,
			&result.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan lab result row: %w", err)
		}

		results = append(results, &result)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating lab result rows: %w", err)
	}

	return results, nil
}
