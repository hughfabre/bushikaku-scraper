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

func sendToDiscord(webhookURL string, busInfos []BusInfo) error {
	if len(busInfos) <= maxEmbeds {
		return sendBatch(webhookURL, busInfos, 1, len(busInfos))
	}

	for i := 0; i < len(busInfos); i += maxEmbeds {
		end := i + maxEmbeds
		if end > len(busInfos) {
			end = len(busInfos)
		}
		if err := sendBatch(webhookURL, busInfos[i:end], i+1, len(busInfos)); err != nil {
			return err
		}
		time.Sleep(batchDelay)
	}
	return nil
}

func sendBatch(webhookURL string, busInfos []BusInfo, startIdx, total int) error {
	embeds := make([]DiscordEmbed, 0, len(busInfos))

	for _, bus := range busInfos {
		embed := buildEmbed(bus)
		embeds = append(embeds, embed)
	}

	content := fmt.Sprintf("%s (%s)", messages.BusReservationInfo, fmt.Sprintf(messages.Items, total, startIdx, startIdx+len(busInfos)-1))
	webhook := DiscordWebhook{
		Content: content,
		Embeds:  embeds,
	}

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

func buildEmbed(bus BusInfo) DiscordEmbed {
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

