package utils

import (
	"github.com/google/uuid"
)

func GenerateNIMBID(prefix string) string {
	uid := generateUID(12)
	return prefix + "_" + uid
}
func generateUID(length int) string {
	uid := uuid.New()
	// Remove hyphens for a more compact ID
	compact := ""
	for _, c := range uid.String() {
		if c != '-' {
			compact += string(c)
		}
	}
	if length > len(compact) {
		length = len(compact)
	}
	return compact[:length]
}
