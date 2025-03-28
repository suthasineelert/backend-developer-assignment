package main

import (
	"backend-developer-assignment/platform/database"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

func main() {
	fmt.Println("Starting database seeding...")

	// Seed user PINs
	if err := database.SeedUserPINs(); err != nil {
		log.Fatalf("Error seeding user PINs: %v", err)
		os.Exit(1)
	}

	fmt.Println("Database seeding completed successfully!")
}
