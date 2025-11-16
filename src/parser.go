package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseBusCard(s *goquery.Selection) BusInfo {
	var busInfo BusInfo

	busInfo.Name = strings.TrimSpace(s.Find("a.SearchCardDirect_bus-title-link__G6_zR").Text())
	busInfo.URL = extractURL(s)
	busInfo.Company = extractCompany(s)
	busInfo.ReservationSite = extractReservationSite(s)
	busInfo.Options = extractOptions(s)
	extractPlatforms(s, &busInfo)
	busInfo.Duration = strings.TrimSpace(s.Find("span.SearchCardDirect_time-box__pTbvH").Text())
	extractPriceAndSeat(s, &busInfo)

	return busInfo
}

func extractURL(s *goquery.Selection) string {
	href, exists := s.Find("a.SearchCardDirect_bus-title-link__G6_zR").Attr("href")
	if !exists {
		return ""
	}
	if strings.HasPrefix(href, "http") {
		return href
	}
	return "https://www.bushikaku.net" + href
}

func extractCompany(s *goquery.Selection) string {
	text := strings.TrimSpace(s.Find("li.SearchCardDirect_company-list-item__BWgqU:contains('バス会社：')").Text())
	return strings.Replace(text, "バス会社：", "", 1)
}

func extractReservationSite(s *goquery.Selection) string {
	text := strings.TrimSpace(s.Find("li.SearchCardDirect_company-list-item__BWgqU:contains('予約サイト：')").Text())
	return strings.Replace(text, "予約サイト：", "", 1)
}

func extractOptions(s *goquery.Selection) []string {
	var options []string
	selector := "li.OptionListItem_night__Is5zV, li.OptionListItem_seat4default__H6th9, li.OptionListItem_reserved_seat__K6LBu, li.OptionListItem_female_safety__PQzmL, li.OptionListItem_student_discount__5RT1t"
	s.Find(selector).Each(func(i int, opt *goquery.Selection) {
		text := strings.TrimSpace(opt.Text())
		if text != "" {
			options = append(options, text)
		}
	})
	return options
}

func extractPlatforms(s *goquery.Selection, busInfo *BusInfo) {
	var departureTimes, arrivalTimes []string
	var departurePlaces, arrivalPlaces []string

	s.Find("div.SearchBusPlatformsVertical_platform-item__tyVQB").Each(func(i int, platform *goquery.Selection) {
		platformType, _ := platform.Attr("data-platform-type")
		timeStr := strings.TrimSpace(platform.Find("div.SearchBusPlatformsVertical_time__0nFfW").Text())
		place := strings.TrimSpace(platform.Find("a.SearchBusPlatformsVertical_platform__qd_wD").Text())
		pref := strings.TrimSpace(platform.Find("div.SearchBusPlatformsVertical_pref__tA_f7").Text())

		if platformType == "geton" {
			departureTimes = append(departureTimes, timeStr)
			departurePlaces = append(departurePlaces, place+" "+pref)
		} else if platformType == "getout" {
			arrivalTimes = append(arrivalTimes, timeStr)
			arrivalPlaces = append(arrivalPlaces, place+" "+pref)
		}
	})

	if len(departurePlaces) > 0 {
		busInfo.Departure = departurePlaces[0]
		busInfo.DepartureTime = departureTimes[0]
	}
	if len(arrivalPlaces) > 0 {
		busInfo.Arrival = arrivalPlaces[len(arrivalPlaces)-1]
		busInfo.ArrivalTime = arrivalTimes[len(arrivalTimes)-1]
	}
}

func extractPriceAndSeat(s *goquery.Selection, busInfo *BusInfo) {
	var prices, seatStatuses []string

	s.Find("tr.SearchCardStructure_structure-table-trlink__o7UzM").Each(func(i int, row *goquery.Selection) {
		plan := strings.TrimSpace(row.Find("span:first-child").Text())
		price := strings.TrimSpace(row.Find("span.SearchCardStructure_structure-table-planamount-text__NXUJI").Text())
		seatStatus := strings.TrimSpace(row.Find("td.SearchCardStructure_structure-table-leftseat-td__i_vZQ a").Text())

		if plan != "" && price != "" {
			prices = append(prices, fmt.Sprintf("%s: %s", plan, price))
		}
		if seatStatus != "" {
			seatStatuses = append(seatStatuses, seatStatus)
		}
	})

	busInfo.Price = strings.Join(prices, "\n")
	if len(seatStatuses) > 0 {
		busInfo.SeatStatus = seatStatuses[0]
	}
}

