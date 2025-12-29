package utils

import (
	"github.com/google/uuid"
)

func GenerateNIMBID(prefix string) string {
	uid := generateUID(12)
	return prefix + "_" + uid
}
func generateUID(length int) string {
	return uuid.New().String()[:length]
}
