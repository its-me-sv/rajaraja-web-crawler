package extractors

import (
	"errors"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Returns all un-normalised list of image URLs found
// within the HTML and an error if one occurs
func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	if baseURL == nil {
		return nil, errors.New("invalid base url")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, errors.New("invalid html body")
	}

	urls := []string{}
	doc.Find("img[src]").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		src = strings.TrimSpace(src)
		if src == "" {
			return
		}

		parsedSrc, err := url.Parse(src)
		if err != nil {
			return
		}

		resolvedSrc := baseURL.ResolveReference(parsedSrc)
		if resolvedSrc.Scheme != "http" && resolvedSrc.Scheme != "https" {
			return
		}

		urls = append(urls, resolvedSrc.String())
	})

	return urls, nil
}
