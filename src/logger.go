package main

import "fmt"

func logInfo(format string, args ...interface{}) {
	fmt.Printf("[INFO] %s\n", fmt.Sprintf(format, args...))
}

func logWarn(format string, args ...interface{}) {
	fmt.Printf("[WARN] %s\n", fmt.Sprintf(format, args...))
}

func logError(format string, args ...interface{}) {
	fmt.Printf("[ERROR] %s\n", fmt.Sprintf(format, args...))
}

