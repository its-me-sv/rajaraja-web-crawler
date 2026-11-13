package utils

import (
	"encoding/json"
	"maps"
	"os"
	"slices"

	"github.com/its-me-sv/rajaraja-web-crawler/internal/extractor"
)

func WriteJSONReport(pages map[string]extractor.PageData, filename string) error {
	// creating array of page data in sorted by keys
	reportData := make([]extractor.PageData, 0, len(pages))
	for _, key := range slices.Sorted(maps.Keys(pages)) {
		reportData = append(reportData, pages[key])
	}

	// marshalling and writing to given file
	fileData, err := json.MarshalIndent(reportData, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, fileData, 0644)

	return err
}
