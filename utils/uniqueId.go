package utils

import "github.com/google/uuid"

func GenerateUniqueId() string {
	return uuid.New().String()
}
