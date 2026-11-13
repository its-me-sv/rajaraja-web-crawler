package crawler

import (
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/its-me-sv/rajaraja-web-crawler/internal/extractor"
	"github.com/its-me-sv/rajaraja-web-crawler/internal/utils"
)

type config struct {
	pages              map[string]extractor.PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
	httpClient         *http.Client
}

func CreateNewConfig(baseURL string, maxConcurrency, maxPages int) (config, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return config{}, errors.New("not an valid url")
	}

	if maxConcurrency < 1 {
		return config{}, errors.New("maxConcurrency must at least be 1")
	}
	if maxPages < 1 {
		return config{}, errors.New("maxPages must at least be 1")
	}

	// configuring http client's no. of idle connections
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = maxConcurrency

	return config{
		pages:              make(map[string]extractor.PageData),
		baseURL:            parsedBaseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
		httpClient: &http.Client{
			Timeout:   7 * time.Second,
			Transport: transport,
		},
	}, nil
}

func (cfg *config) GetPages() map[string]extractor.PageData {
	return cfg.pages
}

func (cfg *config) setPageDataForURL(normalizedURL string, pageData extractor.PageData) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.pages[normalizedURL] = pageData
}

func (cfg *config) unsetPageDataForURL(normalizedURL string) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	delete(cfg.pages, normalizedURL)
}

func (cfg *config) markVisitedWhenMaxPageNotReached(normalizedURL string) (limitReached, claimed bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if len(cfg.pages) >= cfg.maxPages {
		return true, false
	}

	if _, exists := cfg.pages[normalizedURL]; exists {
		return false, false
	}

	cfg.pages[normalizedURL] = extractor.PageData{}
	return false, true
}

func (cfg *config) fetchHTMLContent(rawURL string) (string, error) {
	return utils.GetHTML(cfg.httpClient, rawURL)
}
