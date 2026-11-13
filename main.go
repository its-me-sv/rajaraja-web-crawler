package main

import (
	"fmt"
	"maps"
	"os"
	"slices"

	"github.com/its-me-sv/rajaraja-web-crawler/internal/crawler"
	"github.com/its-me-sv/rajaraja-web-crawler/internal/utils"
)

func main() {
	// parse arguments
	baseURL, maxConcurrency, maxPages, err := utils.ParseArgs(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// create crawler instance
	crawler, err := crawler.CreateNewConfig(baseURL, maxConcurrency, maxPages)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	crawler.Crawl()

	pages := slices.Sorted(maps.Keys(crawler.GetPages()))
	for _, page := range pages {
		fmt.Println(page)
	}
}
