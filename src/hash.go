package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
)

func calculateHash(busInfos []BusInfo) string {
	data, _ := json.Marshal(busInfos)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func loadLastHash() (string, error) {
	data, err := os.ReadFile("last_hash.txt")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func saveLastHash(hash string) error {
	return os.WriteFile("last_hash.txt", []byte(hash), 0644)
}

