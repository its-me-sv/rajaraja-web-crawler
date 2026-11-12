package utils

import "net/url"

func isSameDomain(baseURL, currentURL string) bool {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return false
	}

	parsedCurrentURL, err := url.Parse(currentURL)
	if err != nil {
		return false
	}

	return parsedBaseURL.Hostname() == parsedCurrentURL.Hostname()
}
