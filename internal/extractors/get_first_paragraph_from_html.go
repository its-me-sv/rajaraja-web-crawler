package extractors

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
