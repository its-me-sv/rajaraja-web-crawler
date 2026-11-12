package utils

import (
	"errors"
	"net/url"
	"strings"
)

// normalizeURL returns a canonical form of an absolute URL so that URLs
// pointing to the same page compare equal. It drops the scheme, default
// ports, trailing slashes and fragments, lowercases the host, and keeps
// the path case and query string.
func NormalizeURL(inputURL string) (string, error) {
	trimmed := strings.TrimSpace(inputURL)
	if trimmed == "" {
		return "", errors.New("empty URL")
	}

	u, err := url.Parse(trimmed)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" {
		return "", errors.New("URL is missing a scheme: " + trimmed)
	}
	if u.Host == "" {
		return "", errors.New("URL is missing a host: " + trimmed)
	}

	host := strings.ToLower(u.Host)
	if port := u.Port(); (u.Scheme == "http" && port == "80") || (u.Scheme == "https" && port == "443") {
		host = strings.TrimSuffix(host, ":"+port)
	}

	path := strings.TrimRight(u.Path, "/")
	if u.RawQuery != "" {
		return host + path + "?" + u.RawQuery, nil
	}
	return host + path, nil
}
