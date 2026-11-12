package crawler

import (
	"fmt"
	"net/url"

	"github.com/its-me-sv/rajaraja-web-crawler/internal/extractor"
	"github.com/its-me-sv/rajaraja-web-crawler/internal/utils"
)

func (cfg *config) Crawl() {
	cfg.wg.Add(1)
	go cfg.worker(cfg.baseURL.String())

	cfg.wg.Wait()
}

func (cfg *config) worker(rawCurrentURL string) {
	defer cfg.wg.Done()
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
	cfg.mu.Lock()
	if _, visited := cfg.pages[normalisedCurrentURL]; visited {
		fmt.Println("SKIPPED - already visited")
		cfg.mu.Unlock()
		return
	}
	cfg.pages[normalisedCurrentURL] = extractor.PageData{}
	cfg.mu.Unlock()

	// fetching html content of page
	cfg.concurrencyControl <- struct{}{}
	html, err := utils.GetHTML(rawCurrentURL)
	<-cfg.concurrencyControl
	if err != nil {
		fmt.Println("SKIPPED - failed to get page content")
		return
	}

	// extracting page data
	pageData := extractor.ExtractPageData(html, rawCurrentURL)
	cfg.mu.Lock()
	cfg.pages[normalisedCurrentURL] = pageData
	cfg.mu.Unlock()

	// crawling outgoing links
	for _, outgoingLink := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.worker(outgoingLink)
	}
}
