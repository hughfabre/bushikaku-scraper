package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		logError("Failed to load config file: %v", err)
		os.Exit(1)
	}

	initMessages(config.Language)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		logInfo(messages.Exiting)
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

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

	urlBusInfos := fetchBusInfosFromURLs(config)
	if len(urlBusInfos) == 0 {
		return handleNoBusInfoFound(config.WebhookURL)
	}

	allBusInfos := flattenBusInfos(urlBusInfos)
	logInfo(fmt.Sprintf(messages.BusInfoFetched, len(allBusInfos)))

	if !hasChanges(allBusInfos, config.WebhookURL) {
		return nil
	}

	if err := saveLastHash(calculateHash(allBusInfos)); err != nil {
		logWarn("Failed to save hash: %v", err)
	}

	return sendBusInfosToDiscord(config, urlBusInfos, allBusInfos)
}

func fetchBusInfosFromURLs(config *Config) []URLBusInfo {
	urls := parseSearchURLs(config.SearchURL, config.Language)
	urlBusInfos := fetchFromURLs(urls, false)

	if len(urlBusInfos) == 0 && config.Language == "en" {
		logInfo(messages.TryingJapanese)
		urlsJA := parseSearchURLs(config.SearchURL, "ja")
		urlBusInfos = fetchFromURLs(urlsJA, true)
	}

	return urlBusInfos
}

func fetchFromURLs(urls []string, fromJA bool) []URLBusInfo {
	var urlBusInfos []URLBusInfo
	for _, url := range urls {
		busInfos, err := fetchBusInfo(url)
		if err != nil {
			logWarn("Failed to fetch from %s: %v", url, err)
			continue
		}
		if len(busInfos) > 0 {
			urlBusInfos = append(urlBusInfos, URLBusInfo{
				URL:      url,
				BusInfos: busInfos,
				FromJA:   fromJA,
			})
		}
	}
	return urlBusInfos
}

func flattenBusInfos(urlBusInfos []URLBusInfo) []BusInfo {
	var allBusInfos []BusInfo
	for _, ubi := range urlBusInfos {
		allBusInfos = append(allBusInfos, ubi.BusInfos...)
	}
	return allBusInfos
}

func hasChanges(allBusInfos []BusInfo, webhookURL string) bool {
	hash := calculateHash(allBusInfos)
	lastHash, err := loadLastHash()
	if err == nil && hash == lastHash {
		logInfo(messages.NoChanges)
		if err := sendNoChangesNotification(webhookURL); err != nil {
			logWarn("Failed to send no changes notification: %v", err)
		}
		return false
	}
	return true
}

func sendBusInfosToDiscord(config *Config, urlBusInfos []URLBusInfo, allBusInfos []BusInfo) error {
	originalURLs := parseSearchURLs(config.SearchURL, config.Language)
	if len(originalURLs) >= 2 {
		return sendMultipleURLs(config.WebhookURL, urlBusInfos)
	}
	return sendSingleURL(config.WebhookURL, allBusInfos, urlBusInfos[0].FromJA)
}

func sendMultipleURLs(webhookURL string, urlBusInfos []URLBusInfo) error {
	for _, ubi := range urlBusInfos {
		if err := sendToDiscordByURL(webhookURL, ubi); err != nil {
			logWarn("Failed to send to Discord for URL %s: %v", ubi.URL, err)
		}
	}
	logInfo(messages.DiscordSent)
	return nil
}

func sendSingleURL(webhookURL string, busInfos []BusInfo, fromJA bool) error {
	if err := sendToDiscord(webhookURL, busInfos, fromJA); err != nil {
		return fmt.Errorf("failed to send to Discord: %v", err)
	}
	logInfo(messages.DiscordSent)
	return nil
}

func handleNoBusInfoFound(webhookURL string) error {
	logInfo(messages.NoBusInfoFound)
	if err := sendNoBusInfoNotification(webhookURL); err != nil {
		logWarn("Failed to send no bus info notification: %v", err)
	}
	return nil
}

func parseSearchURLs(searchURL string, language string) []string {
	urls := strings.Split(searchURL, ",")
	var result []string
	for _, url := range urls {
		trimmed := strings.TrimSpace(url)
		if trimmed != "" {
			trimmed = addEnSubdomain(trimmed, language)
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		result = append(result, addEnSubdomain(searchURL, language))
	}
	return result
}

func addEnSubdomain(url string, language string) string {
	if language != "en" || strings.Contains(url, "://en.") {
		return url
	}

	replacements := map[string]string{
		"https://www.": "https://en.",
		"http://www.":  "http://en.",
		"https://":     "https://en.",
		"http://":     "http://en.",
	}

	for old, new := range replacements {
		if strings.HasPrefix(url, old) {
			return strings.Replace(url, old, new, 1)
		}
	}

	return url
}
