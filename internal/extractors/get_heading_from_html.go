package extractors

import (
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
