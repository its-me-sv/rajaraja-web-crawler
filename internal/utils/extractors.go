package utils

import (
	"errors"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Returns the text content of the first h1 tag
// or fallbacks to first h2 tag. Returns
// empty string when neither is found.
func getHeadingFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	if h1 := strings.TrimSpace(doc.Find("h1").First().Text()); len(h1) > 0 {
		return h1
	}

	return strings.TrimSpace(doc.Find("h2").First().Text())
}

// Returns the text content of the first p1 tag
// in main or fallbacks to first p1 tag. Returns
// empty string when neither is found.
func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	main := doc.Find("main")
	p := ""
	if main.Length() > 0 {
		p = strings.TrimSpace(main.Find("p").First().Text())
	}
	if len(p) < 1 {
		p = strings.TrimSpace(doc.Find("p").First().Text())
	}

	return p
}

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
