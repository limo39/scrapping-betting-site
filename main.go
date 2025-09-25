package main

import (
	"log"
	"os"

	"betting-odds-scraper/internal/api"
	"betting-odds-scraper/internal/config"
	"betting-odds-scraper/internal/scraper"
	"betting-odds-scraper/internal/scheduler"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize scraper manager
	scraperManager := scraper.NewManager(cfg)

	// Initialize scheduler for periodic scraping
	scheduler := scheduler.New(scraperManager)
	scheduler.Start()

	// Initialize and start API server
	server := api.NewServer(scraperManager)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := server.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}