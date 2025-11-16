package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		logError("Failed to load config file: %v", err)
		os.Exit(1)
	}

	initMessages(config.Language)

	if len(os.Args) > 1 && os.Args[1] == "run" {
		logInfo(messages.ManualRunStart)
		if err := fetchAndSend(config); err != nil {
			logError("Error: %v", err)
			os.Exit(1)
		}
		logInfo(messages.ManualRunComplete)
		return
	}

	logInfo(fmt.Sprintf(messages.ScheduledModeStart, config.IntervalHours))
	logInfo(messages.ManualRunHint)

	if err := fetchAndSend(config); err != nil {
		logWarn("Initial run error: %v", err)
	}

	ticker := time.NewTicker(time.Duration(config.IntervalHours) * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		if err := fetchAndSend(config); err != nil {
			logWarn("Scheduled run error: %v", err)
		}
	}
}

func fetchAndSend(config *Config) error {
	logInfo(messages.FetchingBusInfo)

	busInfos, err := fetchBusInfo(config.SearchURL)
	if err != nil {
		return fmt.Errorf("failed to fetch bus information: %v", err)
	}

	if len(busInfos) == 0 {
		logInfo(messages.NoBusInfoFound)
		return nil
	}

	logInfo(fmt.Sprintf(messages.BusInfoFetched, len(busInfos)))

	hash := calculateHash(busInfos)
	lastHash, err := loadLastHash()
	if err == nil && hash == lastHash {
		logInfo(messages.NoChanges)
		return nil
	}

	if err := saveLastHash(hash); err != nil {
		logWarn("Failed to save hash: %v", err)
	}

	if err := sendToDiscord(config.WebhookURL, busInfos); err != nil {
		return fmt.Errorf("failed to send to Discord: %v", err)
	}

	logInfo(messages.DiscordSent)
	return nil
}
