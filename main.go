package main

import (
	"fmt"
	"os"

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

	// writing report to file
	filename := "report.json"
	if err := utils.WriteJSONReport(crawler.GetPages(), filename); err != nil {
		fmt.Println("failed to save report to file:", filename)
		os.Exit(1)
	}
	fmt.Println("report saved to file:", filename)
}
