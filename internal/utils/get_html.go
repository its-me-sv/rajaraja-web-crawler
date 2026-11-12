package utils

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
)

func GetHTML(rawURL string) (string, error) {
	httpClient := &http.Client{}

	// building request
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, rawURL, nil)
	if err != nil {
		return "", errors.New("failed to create request")
	}
	req.Header.Set("User-Agent", "RajarajaWebCrawler/1.0")

	// making request
	res, err := httpClient.Do(req)
	if err != nil {
		return "", errors.New("failed to make request")
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", errors.New("failed to get request")
	}
	if !strings.Contains(res.Header.Get("content-type"), "text/html") {
		return "", errors.New("invalid response content")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("failed to read body")
	}

	return string(body), nil
}
