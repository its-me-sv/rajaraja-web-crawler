package crawler

import (
	"fmt"
	"net/url"

	"github.com/its-me-sv/rajaraja-web-crawler/internal/extractor"
	"github.com/its-me-sv/rajaraja-web-crawler/internal/utils"
)

func (cfg *config) Crawl(rawCurrentURL string) {
	fmt.Println("starting crawl of:", rawCurrentURL)

	parsedCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Println("SKIPPED - invalid page URL")
		return
	}

	// checking if current url is part of base url
	if cfg.baseURL.Hostname() != parsedCurrentURL.Hostname() {
		fmt.Println("SKIPPED - not part of base URL")
		return
	}

	// normalising current url
	normalisedCurrentURL, err := utils.NormalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	// do not visit if already visited
	if _, visited := cfg.pages[normalisedCurrentURL]; visited {
		fmt.Println("SKIPPED - already visited")
		return
	}

	// fetching html content of page
	html, err := utils.GetHTML(rawCurrentURL)
	if err != nil {
		fmt.Println("SKIPPED - failed to get page content")
		return
	}

	// extracting page data
	pageData := extractor.ExtractPageData(html, rawCurrentURL)
	cfg.pages[normalisedCurrentURL] = pageData

	// crawling outgoing links
	for _, outgoingLink := range pageData.OutgoingLinks {
		cfg.Crawl(outgoingLink)
	}
}
