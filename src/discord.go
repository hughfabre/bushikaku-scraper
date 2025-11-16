package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const maxEmbeds = 10
const batchDelay = 500 * time.Millisecond

func sendToDiscord(webhookURL string, busInfos []BusInfo, fromJA bool) error {
	return sendBatches(webhookURL, busInfos, "", fromJA)
}

func sendToDiscordByURL(webhookURL string, ubi URLBusInfo) error {
	urlLink := formatURLLink(ubi.URL)
	return sendBatches(webhookURL, ubi.BusInfos, urlLink, ubi.FromJA)
}

func sendBatches(webhookURL string, busInfos []BusInfo, urlLink string, fromJA bool) error {
	total := len(busInfos)
	for i := 0; i < total; i += maxEmbeds {
		end := i + maxEmbeds
		if end > total {
			end = total
		}
		if err := sendBatch(webhookURL, busInfos[i:end], i+1, total, urlLink, fromJA); err != nil {
			return err
		}
		if i+maxEmbeds < total {
			time.Sleep(batchDelay)
		}
	}
	return nil
}

func formatURLLink(url string) string {
	displayText := extractDomain(url)
	if displayText == "" {
		displayText = url
		if len(displayText) > 50 {
			displayText = displayText[:47] + "..."
		}
	}
	return fmt.Sprintf("[%s](%s)", displayText, url)
}

func extractDomain(url string) string {
	prefixes := []string{"https://", "http://"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(url, prefix) {
			url = url[len(prefix):]
			break
		}
	}

	subdomains := []string{"en.", "www."}
	for _, subdomain := range subdomains {
		if strings.HasPrefix(url, subdomain) {
			url = url[len(subdomain):]
			break
		}
	}

	if idx := strings.Index(url, "/"); idx > 0 {
		return url[:idx]
	}
	return url
}

func sendBatch(webhookURL string, busInfos []BusInfo, startIdx, total int, urlLink string, fromJA bool) error {
	embeds := make([]DiscordEmbed, 0, len(busInfos))
	for _, bus := range busInfos {
		embeds = append(embeds, buildEmbed(bus))
	}

	itemsText := fmt.Sprintf(messages.Items, total, startIdx, startIdx+len(busInfos)-1)
	if fromJA {
		itemsText += messages.FromJapanese
	}

	var content string
	if urlLink != "" {
		content = fmt.Sprintf("%s %s (%s)", messages.BusReservationInfo, urlLink, itemsText)
	} else {
		content = fmt.Sprintf("%s (%s)", messages.BusReservationInfo, itemsText)
	}

	webhook := DiscordWebhook{
		Content: content,
		Embeds:  embeds,
	}

	return sendWebhook(webhookURL, webhook)
}

func buildEmbed(bus BusInfo) DiscordEmbed {
	fields := buildEmbedFields(bus)
	embed := DiscordEmbed{
		Title:       bus.Name,
		Description: fmt.Sprintf("**%s** → **%s**", bus.Departure, bus.Arrival),
		Color:       0x3498db,
		Fields:      fields,
	}
	if bus.URL != "" {
		embed.URL = bus.URL
	}
	return embed
}

func buildEmbedFields(bus BusInfo) []EmbedField {
	fields := []EmbedField{
		{Name: messages.Departure, Value: fmt.Sprintf("%s %s", bus.DepartureTime, bus.Departure), Inline: true},
		{Name: messages.Arrival, Value: fmt.Sprintf("%s %s", bus.ArrivalTime, bus.Arrival), Inline: true},
	}

	if bus.Duration != "" {
		fields = append(fields, EmbedField{Name: messages.Duration, Value: bus.Duration, Inline: true})
	}
	if bus.Company != "" {
		fields = append(fields, EmbedField{Name: messages.Company, Value: bus.Company, Inline: true})
	}
	if bus.ReservationSite != "" {
		fields = append(fields, EmbedField{Name: messages.ReservationSite, Value: bus.ReservationSite, Inline: true})
	}
	if bus.SeatStatus != "" {
		fields = append(fields, EmbedField{Name: messages.SeatStatus, Value: bus.SeatStatus, Inline: true})
	}
	if len(bus.Options) > 0 {
		fields = append(fields, EmbedField{Name: messages.Options, Value: strings.Join(bus.Options, ", "), Inline: false})
	}
	if bus.Price != "" {
		fields = append(fields, EmbedField{Name: messages.Price, Value: bus.Price, Inline: false})
	}

	return fields
}

func sendNoChangesNotification(webhookURL string) error {
	webhook := DiscordWebhook{
		Content: messages.NoChanges,
	}
	return sendWebhook(webhookURL, webhook)
}

func sendNoBusInfoNotification(webhookURL string) error {
	webhook := DiscordWebhook{
		Content: messages.NoBusInfoFound,
	}
	return sendWebhook(webhookURL, webhook)
}

func sendWebhook(webhookURL string, webhook DiscordWebhook) error {
	jsonData, err := json.Marshal(webhook)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord API error: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

