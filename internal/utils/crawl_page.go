package utils

import (
	"fmt"

	"github.com/its-me-sv/rajaraja-web-crawler/internal/extractors"
)

func CrawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	// checking if current url is part of base url
	if !isSameDomain(rawBaseURL, rawCurrentURL) {
		return
	}

	// normalising current url
	normalisedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	// incrementing visit count for normalised url
	// if already visited, count is increased and exited
	pages[normalisedCurrentURL]++
	if pages[normalisedCurrentURL] > 1 {
		return
	}

	fmt.Println("starting crawl of:", rawCurrentURL)

	// fetching html content of page
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		fmt.Println("failed to get HTML, page SKIPPED")
		return
	}

	// extracting page data
	pageData := extractors.ExtractPageData(html, rawCurrentURL)

	// crawling outgoing links
	for _, outgoingLink := range pageData.OutgoingLinks {
		CrawlPage(rawBaseURL, outgoingLink, pages)
	}
}
