package utils

import "golang.org/x/crypto/bcrypt"

// HashPIN securely hashes the PIN
func HashPIN(pin string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	return string(hashed), err
}

// VerifyPIN checks if the provided PIN matches the stored hash
func VerifyPIN(storedHash, inputPIN string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(inputPIN))
	return err == nil
}
