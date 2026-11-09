package utils

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Returns the text content of the first h1 tag
// and fallbacks to first h2 tag. Returns
// empty string when neither is found.
func getHeadingFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	return strings.TrimSpace(doc.Find("h1, h2").First().Text())
}

// Returns the text content of the first p1 tag
// in main and fallbacks to first p1 tag. Returns
// empty string when neither is found.
func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	main := doc.Find("main")
	p := ""
	if main.Length() > 0 {
		p = main.Find("p").First().Text()
	} else {
		p = doc.Find("p").First().Text()
	}

	return strings.TrimSpace(p)
}
