package main

import (
	"fmt"
	"maps"
	"os"
	"slices"

	"github.com/its-me-sv/rajaraja-web-crawler/internal/crawler"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL := args[0]
	crawler, err := crawler.CreateNewConfig(baseURL, 16)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	crawler.Crawl()

	pages := slices.Sorted(maps.Keys(crawler.GetPages()))
	fmt.Println(pages)
}
