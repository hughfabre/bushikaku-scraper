package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const maxPages = 20
const requestDelay = 1 * time.Second

func fetchBusInfo(baseURL string) ([]BusInfo, error) {
	var allBusInfos []BusInfo

	for page := 1; page <= maxPages; page++ {
		url := buildPageURL(baseURL, page)
		logInfo(fmt.Sprintf(messages.PageFetching, page, url))

		busInfos, err := fetchPage(url)
		if err != nil {
			if page == 1 {
				return nil, err
			}
			logWarn(fmt.Sprintf(messages.PageNotFound, page, err))
			break
		}

		if len(busInfos) == 0 {
			break
		}

		allBusInfos = append(allBusInfos, busInfos...)
		logInfo(fmt.Sprintf(messages.PageFetched, page, len(busInfos)))

		time.Sleep(requestDelay)
	}

	return allBusInfos, nil
}

func buildPageURL(baseURL string, page int) string {
	if page == 1 {
		return baseURL
	}
	if strings.HasSuffix(baseURL, "/") {
		return fmt.Sprintf("%spage-%d/", baseURL, page)
	}
	return fmt.Sprintf("%s/page-%d/", baseURL, page)
}

func fetchPage(url string) ([]BusInfo, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("page not found (404)")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var busInfos []BusInfo
	doc.Find("li.SearchCardDirect_search-card__PPng1").Each(func(i int, s *goquery.Selection) {
		busInfo := parseBusCard(s)
		if busInfo.Name != "" {
			busInfos = append(busInfos, busInfo)
		}
	})

	return busInfos, nil
}

