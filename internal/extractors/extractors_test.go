package extractors

import (
	"reflect"
	"testing"
)

func TestExtractPageData(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`

	actual := ExtractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://crawler-test.com",
		Heading:        "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://crawler-test.com/link1"},
		ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

// ExtractPageData should combine the individual extractors into one PageData.
//
// Expected rules:
//   - URL is the page URL exactly as passed in (the crawl key is the caller's job)
//   - Heading, FirstParagraph, OutgoingLinks and ImageURLs follow the rules of
//     their own extractor functions
//   - relative links and image srcs are resolved against the full page URL,
//     including its scheme, port and path (a trailing slash on the page URL matters)
//   - Heading and FirstParagraph never depend on the page URL, so they are
//     always extracted, even when the page URL is invalid
//   - an empty, unparseable or scheme-less page URL still sets URL, Heading and
//     FirstParagraph, but OutgoingLinks and ImageURLs are empty because
//     relative URLs can't be resolved against it
func TestExtractPageDataCases(t *testing.T) {
	tests := []struct {
		name      string
		pageURL   string
		inputBody string
		expected  PageData
	}{
		// --- combining the extractors ---
		{
			name:    "multiple links and images, mixed relative and absolute",
			pageURL: "https://crawler-test.com",
			inputBody: `<html><body>
				<h1>Title</h1>
				<p>Intro.</p>
				<a href="/one">One</a>
				<a href="https://other.com/two">Two</a>
				<img src="/a.png">
				<img src="https://cdn.other.com/b.png">
			</body></html>`,
			expected: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "Title",
				FirstParagraph: "Intro.",
				OutgoingLinks:  []string{"https://crawler-test.com/one", "https://other.com/two"},
				ImageURLs:      []string{"https://crawler-test.com/a.png", "https://cdn.other.com/b.png"},
			},
		},
		{
			name:    "heading falls back to h2 and paragraph prefers main",
			pageURL: "https://crawler-test.com",
			inputBody: `<html><body>
				<h2>Sub Title</h2>
				<p>Outside.</p>
				<main><p>Inside main.</p></main>
			</body></html>`,
			expected: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "Sub Title",
				FirstParagraph: "Inside main.",
			},
		},
		{
			name:    "image inside a link is reported as both link and image",
			pageURL: "https://crawler-test.com",
			inputBody: `<html><body>
				<a href="/gallery"><img src="/thumb.png"></a>
			</body></html>`,
			expected: PageData{
				URL:           "https://crawler-test.com",
				OutgoingLinks: []string{"https://crawler-test.com/gallery"},
				ImageURLs:     []string{"https://crawler-test.com/thumb.png"},
			},
		},
		{
			name:    "non-http links and images are skipped",
			pageURL: "https://crawler-test.com",
			inputBody: `<html><body>
				<a href="mailto:hi@crawler-test.com">Mail</a>
				<a href="/kept">Kept</a>
				<img src="data:image/png;base64,iVBORw0KGgo=">
				<img src="/kept.png">
			</body></html>`,
			expected: PageData{
				URL:           "https://crawler-test.com",
				OutgoingLinks: []string{"https://crawler-test.com/kept"},
				ImageURLs:     []string{"https://crawler-test.com/kept.png"},
			},
		},

		{
			name:    "realistic page with head, nav, article and footer",
			pageURL: "https://crawler-test.com/blog/go-tips",
			inputBody: `<!DOCTYPE html>
			<html lang="en">
			<head>
				<title>Go Tips</title>
				<link rel="stylesheet" href="/style.css">
				<meta property="og:image" content="/og.png">
			</head>
			<body>
				<header>
					<a href="/"><img src="/logo.svg" alt="Logo"></a>
					<nav><a href="/blog">Blog</a> <a href="/about">About</a></nav>
				</header>
				<p>Site-wide banner.</p>
				<main>
					<article>
						<h1>  Ten Go Tips  </h1>
						<p>Go is <strong>simple</strong> &amp; fast.</p>
						<p>Second paragraph.</p>
						<img src="images/gopher.png" alt="Gopher">
						<a href="go-errors">Next: errors</a>
						<a href="https://go.dev/doc">Go docs</a>
					</article>
				</main>
				<footer><a href="mailto:hi@crawler-test.com">Email</a></footer>
			</body>
			</html>`,
			expected: PageData{
				URL:            "https://crawler-test.com/blog/go-tips",
				Heading:        "Ten Go Tips",
				FirstParagraph: "Go is simple & fast.",
				OutgoingLinks: []string{
					"https://crawler-test.com/",
					"https://crawler-test.com/blog",
					"https://crawler-test.com/about",
					"https://crawler-test.com/blog/go-errors",
					"https://go.dev/doc",
				},
				ImageURLs: []string{
					"https://crawler-test.com/logo.svg",
					"https://crawler-test.com/blog/images/gopher.png",
				},
			},
		},
		{
			name:    "duplicate links and images are kept",
			pageURL: "https://crawler-test.com",
			inputBody: `<html><body>
				<a href="/same">A</a><a href="/same">B</a>
				<img src="/same.png"><img src="/same.png">
			</body></html>`,
			expected: PageData{
				URL:           "https://crawler-test.com",
				OutgoingLinks: []string{"https://crawler-test.com/same", "https://crawler-test.com/same"},
				ImageURLs:     []string{"https://crawler-test.com/same.png", "https://crawler-test.com/same.png"},
			},
		},
		{
			name:      "links only, no images",
			pageURL:   "https://crawler-test.com",
			inputBody: `<html><body><a href="/only-link">Link</a></body></html>`,
			expected: PageData{
				URL:           "https://crawler-test.com",
				OutgoingLinks: []string{"https://crawler-test.com/only-link"},
			},
		},
		{
			name:      "images only, no links",
			pageURL:   "https://crawler-test.com",
			inputBody: `<html><body><img src="/only-image.png"></body></html>`,
			expected: PageData{
				URL:       "https://crawler-test.com",
				ImageURLs: []string{"https://crawler-test.com/only-image.png"},
			},
		},

		// --- resolving against the page URL ---
		{
			name:    "page URL fragment is ignored when resolving",
			pageURL: "https://crawler-test.com/docs/page#section",
			inputBody: `<html><body>
				<a href="other">Other</a>
				<img src="pic.png">
			</body></html>`,
			expected: PageData{
				URL:           "https://crawler-test.com/docs/page#section",
				OutgoingLinks: []string{"https://crawler-test.com/docs/other"},
				ImageURLs:     []string{"https://crawler-test.com/docs/pic.png"},
			},
		},
		{
			name:    "page URL is not normalised before resolving",
			pageURL: "https://Crawler-Test.com/Blog/",
			inputBody: `<html><body>
				<a href="Post">Post</a>
				<img src="Pic.png">
			</body></html>`,
			expected: PageData{
				URL:           "https://Crawler-Test.com/Blog/",
				OutgoingLinks: []string{"https://Crawler-Test.com/Blog/Post"},
				ImageURLs:     []string{"https://Crawler-Test.com/Blog/Pic.png"},
			},
		},
		{
			name:    "relative URLs resolve against page path ending in slash",
			pageURL: "https://crawler-test.com/blog/",
			inputBody: `<html><body>
				<a href="post">Post</a>
				<img src="img/cover.png">
			</body></html>`,
			expected: PageData{
				URL:           "https://crawler-test.com/blog/",
				OutgoingLinks: []string{"https://crawler-test.com/blog/post"},
				ImageURLs:     []string{"https://crawler-test.com/blog/img/cover.png"},
			},
		},
		{
			name:    "relative URLs resolve against page path not ending in slash",
			pageURL: "https://crawler-test.com/blog/current",
			inputBody: `<html><body>
				<a href="next">Next</a>
				<img src="../logo.png">
			</body></html>`,
			expected: PageData{
				URL:           "https://crawler-test.com/blog/current",
				OutgoingLinks: []string{"https://crawler-test.com/blog/next"},
				ImageURLs:     []string{"https://crawler-test.com/logo.png"},
			},
		},
		{
			name:    "page scheme and port are used for relative URLs",
			pageURL: "http://localhost:8080/docs/",
			inputBody: `<html><body>
				<a href="intro">Intro</a>
				<img src="//cdn.localhost/pic.png">
			</body></html>`,
			expected: PageData{
				URL:           "http://localhost:8080/docs/",
				OutgoingLinks: []string{"http://localhost:8080/docs/intro"},
				ImageURLs:     []string{"http://cdn.localhost/pic.png"},
			},
		},
		{
			name:    "page URL with query string",
			pageURL: "https://crawler-test.com/search?q=go",
			inputBody: `<html><body>
				<a href="?q=go&amp;page=2">Next page</a>
			</body></html>`,
			expected: PageData{
				URL:           "https://crawler-test.com/search?q=go",
				OutgoingLinks: []string{"https://crawler-test.com/search?q=go&page=2"},
			},
		},
		{
			name:      "page URL is kept exactly as given",
			pageURL:   "https://Crawler-Test.com/Blog/Post/",
			inputBody: `<html><body><h1>Title</h1></body></html>`,
			expected: PageData{
				URL:     "https://Crawler-Test.com/Blog/Post/",
				Heading: "Title",
			},
		},

		// --- empty pages ---
		{
			name:      "page with no content",
			pageURL:   "https://crawler-test.com",
			inputBody: `<html><body><div></div></body></html>`,
			expected:  PageData{URL: "https://crawler-test.com"},
		},
		{
			name:      "empty HTML body",
			pageURL:   "https://crawler-test.com",
			inputBody: "",
			expected:  PageData{URL: "https://crawler-test.com"},
		},

		// --- invalid page URLs ---
		{
			name:      "empty page URL keeps text but has no links or images",
			pageURL:   "",
			inputBody: `<html><body><h1>Title</h1><p>Text.</p><a href="/x">X</a><img src="/x.png"></body></html>`,
			expected: PageData{
				URL:            "",
				Heading:        "Title",
				FirstParagraph: "Text.",
			},
		},
		{
			name:      "page URL without scheme keeps text but has no links or images",
			pageURL:   "crawler-test.com/blog",
			inputBody: `<html><body><h1>Title</h1><p>Text.</p><a href="/x">X</a><img src="/x.png"></body></html>`,
			expected: PageData{
				URL:            "crawler-test.com/blog",
				Heading:        "Title",
				FirstParagraph: "Text.",
			},
		},
		{
			name:      "malformed page URL keeps text but has no links or images",
			pageURL:   "://crawler-test.com",
			inputBody: `<html><body><h1>Title</h1><p>Text.</p><a href="/x">X</a><img src="/x.png"></body></html>`,
			expected: PageData{
				URL:            "://crawler-test.com",
				Heading:        "Title",
				FirstParagraph: "Text.",
			},
		},
		{
			name:      "malformed page URL drops even absolute links and images",
			pageURL:   "://crawler-test.com",
			inputBody: `<html><body><a href="https://other.com/page">X</a><img src="https://cdn.other.com/pic.png"></body></html>`,
			expected: PageData{
				URL: "://crawler-test.com",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := ExtractPageData(tc.inputBody, tc.pageURL)
			assertPageDataEqual(t, tc.expected, actual)
		})
	}
}

// assertPageDataEqual compares two PageData values, treating nil and empty
// slices as equal so tests don't depend on how "no results" is represented.
func assertPageDataEqual(t *testing.T, expected, actual PageData) {
	t.Helper()

	if actual.URL != expected.URL {
		t.Errorf("URL: expected %q, got %q", expected.URL, actual.URL)
	}
	if actual.Heading != expected.Heading {
		t.Errorf("Heading: expected %q, got %q", expected.Heading, actual.Heading)
	}
	if actual.FirstParagraph != expected.FirstParagraph {
		t.Errorf("FirstParagraph: expected %q, got %q", expected.FirstParagraph, actual.FirstParagraph)
	}
	if !sameStrings(actual.OutgoingLinks, expected.OutgoingLinks) {
		t.Errorf("OutgoingLinks: expected %q, got %q", expected.OutgoingLinks, actual.OutgoingLinks)
	}
	if !sameStrings(actual.ImageURLs, expected.ImageURLs) {
		t.Errorf("ImageURLs: expected %q, got %q", expected.ImageURLs, actual.ImageURLs)
	}
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	return len(a) == 0 || reflect.DeepEqual(a, b)
}
