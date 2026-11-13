package extractor

import (
	"net/url"
)

type PageData struct {
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"image_urls"`
}

func ExtractPageData(html, pageURL string) PageData {
	heading := getHeadingFromHTML(html)
	firstPara := getFirstParagraphFromHTML(html)

	pageData := PageData{
		URL:            pageURL,
		Heading:        heading,
		FirstParagraph: firstPara,
		OutgoingLinks:  nil,
		ImageURLs:      nil,
	}

	parsedUrl, err := url.Parse(pageURL)
	if err != nil {
		return pageData
	}

	if outgoingLinks, err := getURLsFromHTML(html, parsedUrl); err == nil {
		pageData.OutgoingLinks = outgoingLinks
	}
	if imageURLs, err := getImagesFromHTML(html, parsedUrl); err == nil {
		pageData.ImageURLs = imageURLs
	}

	return pageData
}
