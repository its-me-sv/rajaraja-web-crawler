package main

import (
	"fmt"
	"os"

	"github.com/its-me-sv/rajaraja-web-crawler/internal/utils"
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
	pages := map[string]int{}
	utils.CrawlPage(baseURL, baseURL, pages)

	for page, count := range pages {
		fmt.Println(page, ":", count)
	}
}
