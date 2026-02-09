package main

import (
	"context"
	"embed"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"veta-lis/internal/app"
	"veta-lis/internal/repository"
	"veta-lis/internal/service"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Database connection string
	// TODO: Move to config file or environment variable
	databaseURL := "postgresql://postgres:password@localhost:5432/veta_lis?sslmode=disable"

	// Initialize database connection pool
	dbPool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Printf("⚠️  Warning: Could not connect to database: %v", err)
		log.Println("📝 Running in demo mode without database")
		// You can create a mock service here or continue without DB
		// For now, we'll create service with nil repo (will cause errors if used)
		// In production, you'd want to handle this more gracefully
	}
	if dbPool != nil {
		defer dbPool.Close()
		log.Println("✅ Database connection established")
	}

	// Create repository layer (handles database operations)
	var repo *repository.PetRepository
	if dbPool != nil {
		repo = repository.NewPetRepository(dbPool)
	}

	// Create service layer (handles business logic)
	labService := service.NewLabService(repo)

	// Create Wails application with service injection
	application := app.NewApp(labService)

	// Run Wails application
	err = wails.Run(&options.App{
		Title:      "Veta LIS - Veterinary Laboratory Information System",
		Width:      1400,
		Height:     900,
		MinWidth:   1000,
		MinHeight:  600,
		OnStartup:  application.Startup,
		OnShutdown: application.Shutdown,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Bind: []interface{}{
			application,
		},
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
