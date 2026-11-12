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
}

func CreateNewConfig(baseURL string, maxConcurrency int) (config, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return config{}, errors.New("not an valid url")
	}

	return config{
		pages:              make(map[string]extractor.PageData),
		baseURL:            parsedBaseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
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

func (cfg *config) markVisited(normalizedURL string) bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, visited := cfg.pages[normalizedURL]; visited {
		return false
	}

	cfg.pages[normalizedURL] = extractor.PageData{}
	return true
}
