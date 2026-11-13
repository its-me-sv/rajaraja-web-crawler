package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/its-me-sv/rajaraja-web-crawler/internal/crawler"
	"github.com/its-me-sv/rajaraja-web-crawler/internal/utils"
)

func main() {
	// parse arguments
	baseURL, maxConcurrency, maxPages, err := utils.ParseArgs(os.Args)
	if errors.Is(err, utils.ErrHelp) {
		fmt.Println(utils.Usage)
		os.Exit(0)
	}
	if errors.Is(err, utils.ErrUsage) {
		fmt.Println(utils.Usage)
		os.Exit(1)
	}
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
