package main

import (
	"fmt"
	"time"
)

func logInfo(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[INFO] [%s] %s\n", timestamp, fmt.Sprintf(format, args...))
}

func logWarn(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[WARN] [%s] %s\n", timestamp, fmt.Sprintf(format, args...))
}

func logError(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[ERROR] [%s] %s\n", timestamp, fmt.Sprintf(format, args...))
}

