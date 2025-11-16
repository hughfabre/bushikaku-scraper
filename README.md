A scraper for [this website](https://www.bushikaku.net/).
It will send you the scraped info via discord webhook.

you can configure files [here](/config.json).

Use this command to build:
<br />
`go build -o busyoyaku.exe ./src`

Explanation of [config.json](/config.json)

- `search_url`: Search URL (example: https://www.bushikaku.net/search/tokyo_osaka/20260101/)
- `webhook_url`: Discord webhook URL
- `interval_hours`: Check interval in hours. Default is 6.
- `language`: Language setting. "en" or "ja. "Default is "ja".
