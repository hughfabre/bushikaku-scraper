package main

type Config struct {
	SearchURL     string `json:"search_url"`
	WebhookURL    string `json:"webhook_url"`
	IntervalHours int    `json:"interval_hours"`
	Language      string `json:"language"`
}

type BusInfo struct {
	Name            string
	Departure       string
	Arrival         string
	DepartureTime   string
	ArrivalTime     string
	Duration        string
	Company         string
	ReservationSite string
	Price           string
	SeatStatus      string
	Options         []string
	URL             string
}

type DiscordEmbed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Color       int          `json:"color"`
	Fields      []EmbedField `json:"fields"`
	URL         string       `json:"url,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type DiscordWebhook struct {
	Content string         `json:"content,omitempty"`
	Embeds  []DiscordEmbed `json:"embeds"`
}

type Messages struct {
	ManualRunStart     string
	ManualRunComplete  string
	ScheduledModeStart string
	ManualRunHint      string
	FetchingBusInfo    string
	NoBusInfoFound     string
	BusInfoFetched     string
	DiscordSent        string
	NoChanges          string
	PageFetching       string
	PageFetched        string
	PageNotFound       string
	Departure          string
	Arrival            string
	Duration           string
	Company            string
	ReservationSite    string
	SeatStatus         string
	Options            string
	Price              string
	BusReservationInfo string
	Items              string
}

