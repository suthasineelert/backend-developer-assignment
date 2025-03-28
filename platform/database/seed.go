package database

import (
	"backend-developer-assignment/pkg/utils"
	"fmt"
	"log"
)

// SeedUserPINs initializes PINs for all users in the database. default pin is 123456
func SeedUserPINs() error {
	// Get database connection
	db, err := MysqlConnection()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	pin := "123456"
	// Hash the PIN
	hashedPIN, err := utils.HashPIN(pin)
	if err != nil {
		return fmt.Errorf("failed to create hash PIN: %w", err)
	}

	// Update all user PIN
	_, err = db.Exec("UPDATE users SET pin = ?", hashedPIN)
	if err != nil {
		return fmt.Errorf("failed to update PIN for all user %w", err)
	}

	log.Println("Successfully updated all users with hashed PINs")
	return nil
}
