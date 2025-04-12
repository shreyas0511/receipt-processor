package utils

import "github.com/google/uuid"

// Generate a unique ID
func GenerateUniqueId() string {
	return uuid.New().String()
}
