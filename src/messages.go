package main

var messages Messages

func initMessages(lang string) {
	if lang == "en" {
		messages = Messages{
			ManualRunStart:     "Starting manual run...",
			ManualRunComplete:  "Completed!",
			ScheduledModeStart: "Starting scheduled mode (every %d hours)",
			ManualRunHint:      "To run manually, execute 'go run . run'",
			FetchingBusInfo:    "Fetching bus information...",
			NoBusInfoFound:     "No bus information found",
			BusInfoFetched:     "Fetched %d bus information items",
			DiscordSent:        "Sent to Discord",
			NoChanges:          "No changes from previous run",
			PageFetching:       "Fetching page %d: %s",
			PageFetched:        "Page %d: Fetched %d bus information items",
			PageNotFound:       "Page %d fetch failed (page may not exist): %v",
			Departure:          "Departure",
			Arrival:            "Arrival",
			Duration:           "Duration",
			Company:            "Company",
			ReservationSite:    "Reservation Site",
			SeatStatus:         "Seat Status",
			Options:            "Options",
			Price:              "Price",
			BusReservationInfo: "Bus Reservation Information",
			Items:              "items %d-%d of %d",
			Exiting:            "Exiting...",
			TryingJapanese:     "No results from English version, trying Japanese version...",
			FromJapanese:        " (from Japanese website)",
		}
	} else {
		messages = Messages{
			ManualRunStart:     "手動実行を開始します...",
			ManualRunComplete:  "完了しました！",
			ScheduledModeStart: "定期実行モードを開始します（%d時間ごと）",
			ManualRunHint:      "手動実行する場合は 'go run . run' を実行してください",
			FetchingBusInfo:    "バス情報を取得中...",
			NoBusInfoFound:     "バス情報が見つかりませんでした",
			BusInfoFetched:     "%d件のバス情報を取得しました",
			DiscordSent:        "Discordに送信しました",
			NoChanges:          "前回から変更はありません",
			PageFetching:       "ページ %d を取得中: %s",
			PageFetched:        "ページ %d: %d件のバス情報を取得",
			PageNotFound:       "ページ %d の取得に失敗（ページが存在しない可能性）: %v",
			Departure:          "出発",
			Arrival:            "到着",
			Duration:           "乗車時間",
			Company:            "バス会社",
			ReservationSite:    "予約サイト",
			SeatStatus:         "残席",
			Options:            "オプション",
			Price:              "料金",
			BusReservationInfo: "バス予約情報",
			Items:              "%d件中 %d-%d件目",
			Exiting:            "終了中...",
			TryingJapanese:     "英語版で結果が見つからなかったため、日本語版で再検索します...",
			FromJapanese:        " (日本語サイトから)",
		}
	}
}

