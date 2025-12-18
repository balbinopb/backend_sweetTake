package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func GenerateRandomToken(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return strings.ToUpper(hex.EncodeToString(bytes))
}
