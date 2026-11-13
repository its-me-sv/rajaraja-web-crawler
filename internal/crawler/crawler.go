package crawler

import (
	"errors"
	"net/url"
	"sync"

	"github.com/its-me-sv/rajaraja-web-crawler/internal/extractor"
)

type config struct {
	pages              map[string]extractor.PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
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

	return config{
		pages:              make(map[string]extractor.PageData),
		baseURL:            parsedBaseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}, nil
}

func (cfg *config) isPagesLimitReached() bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.pages) >= cfg.maxPages
}

func (cfg *config) GetPages() map[string]extractor.PageData {
	return cfg.pages
}

func (cfg *config) setPageDataForURL(normalizedURL string, pageData extractor.PageData) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.pages[normalizedURL] = pageData
}

func (cfg *config) markVisited(normalizedURL string) bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, visited := cfg.pages[normalizedURL]; visited {
		return false
	}

	cfg.pages[normalizedURL] = extractor.PageData{}
	return true
}
