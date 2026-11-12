package extractor

import (
	"errors"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Returns all un-normalised list of URLs found
// within the given HTML and an error if one
// occurs
func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	if baseURL == nil {
		return nil, errors.New("invalid base url")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, errors.New("invalid html body")
	}

	urls := []string{}
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		href = strings.TrimSpace(href)
		if href == "" {
			return
		}

		parsedURL, err := url.Parse(href)
		if err != nil {
			return
		}

		resolved := baseURL.ResolveReference(parsedURL)
		if resolved.Scheme != "http" && resolved.Scheme != "https" {
			return
		}

		urls = append(urls, resolved.String())
	})

	return urls, nil
}
