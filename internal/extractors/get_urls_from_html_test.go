package extractors

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetURLsFromHTMLAbsolute(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body><a href="https://crawler-test.com"><span>Boot.dev</span></a></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://crawler-test.com"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

// getURLsFromHTML should return every crawlable link from the <a> tags in the
// page, in document order, with relative links resolved against baseURL.
//
// Expected rules:
//   - every <a href> in the document is found, wherever it is nested
//   - relative hrefs are resolved to absolute URLs using baseURL
//   - absolute hrefs are returned as-is (no normalisation)
//   - duplicates are kept (deduplication is the crawler's job)
//   - <a> without href, or with an empty/whitespace-only href, is skipped
//   - non-http(s) links (mailto:, tel:, javascript:, ...) are skipped
//   - hrefs that fail to parse are skipped, not fatal
//   - only <a> tags count (<link>, <area>, <img src> etc. are ignored)
func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		inputBody string
		expected  []string
	}{
		// --- finding <a> tags ---
		{
			name:      "multiple a tags in document order",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><a href="https://crawler-test.com/one">1</a><a href="https://crawler-test.com/two">2</a><a href="https://crawler-test.com/three">3</a></body></html>`,
			expected:  []string{"https://crawler-test.com/one", "https://crawler-test.com/two", "https://crawler-test.com/three"},
		},
		{
			name:    "a tags nested anywhere in the page",
			baseURL: "https://crawler-test.com",
			inputBody: `<html><body>
				<header><nav><ul><li><a href="/nav">Nav</a></li></ul></nav></header>
				<main><article><section><div><p>Read <a href="/deep">this</a>.</p></div></section></article></main>
				<footer><a href="/footer">Footer</a></footer>
			</body></html>`,
			expected: []string{"https://crawler-test.com/nav", "https://crawler-test.com/deep", "https://crawler-test.com/footer"},
		},
		{
			name:      "uppercase tag and attribute names",
			baseURL:   "https://crawler-test.com",
			inputBody: `<HTML><BODY><A HREF="/upper">Upper</A></BODY></HTML>`,
			expected:  []string{"https://crawler-test.com/upper"},
		},
		{
			name:      "duplicate links are kept",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><a href="/same">A</a><a href="/same">B</a></body></html>`,
			expected:  []string{"https://crawler-test.com/same", "https://crawler-test.com/same"},
		},
		{
			name:      "works on a fragment without html/body tags",
			baseURL:   "https://crawler-test.com",
			inputBody: `<a href="/fragment">Fragment</a>`,
			expected:  []string{"https://crawler-test.com/fragment"},
		},
		{
			name:      "handles malformed HTML with unclosed tags",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><div><a href="/unclosed">Unclosed`,
			expected:  []string{"https://crawler-test.com/unclosed"},
		},

		// --- relative URL resolution ---
		{
			name:      "root-relative path",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><a href="/path/one">Link</a></body></html>`,
			expected:  []string{"https://crawler-test.com/path/one"},
		},
		{
			name:      "relative path against base ending in slash",
			baseURL:   "https://crawler-test.com/blog/",
			inputBody: `<html><body><a href="post">Post</a></body></html>`,
			expected:  []string{"https://crawler-test.com/blog/post"},
		},
		{
			name:      "relative path against base not ending in slash replaces last segment",
			baseURL:   "https://crawler-test.com/blog/current",
			inputBody: `<html><body><a href="post">Post</a></body></html>`,
			expected:  []string{"https://crawler-test.com/blog/post"},
		},
		{
			name:      "dot-slash relative path",
			baseURL:   "https://crawler-test.com/blog/",
			inputBody: `<html><body><a href="./post">Post</a></body></html>`,
			expected:  []string{"https://crawler-test.com/blog/post"},
		},
		{
			name:      "parent directory relative path",
			baseURL:   "https://crawler-test.com/blog/2024/post",
			inputBody: `<html><body><a href="../archive">Archive</a></body></html>`,
			expected:  []string{"https://crawler-test.com/blog/archive"},
		},
		{
			name:      "parent directory beyond root stays at root",
			baseURL:   "https://crawler-test.com/blog/",
			inputBody: `<html><body><a href="../../../top">Top</a></body></html>`,
			expected:  []string{"https://crawler-test.com/top"},
		},
		{
			name:      "base query string is not carried over to relative link",
			baseURL:   "https://crawler-test.com/blog/post?ref=home",
			inputBody: `<html><body><a href="other">Other</a></body></html>`,
			expected:  []string{"https://crawler-test.com/blog/other"},
		},
		{
			name:      "query-only href keeps base path",
			baseURL:   "https://crawler-test.com/search",
			inputBody: `<html><body><a href="?page=2">Next</a></body></html>`,
			expected:  []string{"https://crawler-test.com/search?page=2"},
		},
		{
			name:      "fragment-only href keeps base path",
			baseURL:   "https://crawler-test.com/blog/post",
			inputBody: `<html><body><a href="#comments">Comments</a></body></html>`,
			expected:  []string{"https://crawler-test.com/blog/post#comments"},
		},
		{
			name:      "protocol-relative href uses base scheme",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><a href="//cdn.crawler-test.com/file">CDN</a></body></html>`,
			expected:  []string{"https://cdn.crawler-test.com/file"},
		},
		{
			name:      "relative link keeps base port",
			baseURL:   "http://localhost:8080/docs/",
			inputBody: `<html><body><a href="intro">Intro</a></body></html>`,
			expected:  []string{"http://localhost:8080/docs/intro"},
		},
		{
			name:      "mix of relative and absolute links",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><a href="/local">Local</a><a href="https://other.com/page">External</a></body></html>`,
			expected:  []string{"https://crawler-test.com/local", "https://other.com/page"},
		},

		// --- href values ---
		{
			name:      "absolute link to another domain is kept as-is",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><a href="https://Other.com/Some/Path/">External</a></body></html>`,
			expected:  []string{"https://Other.com/Some/Path/"},
		},
		{
			name:      "http link on https page keeps its own scheme",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><a href="http://crawler-test.com/insecure">Insecure</a></body></html>`,
			expected:  []string{"http://crawler-test.com/insecure"},
		},
		{
			name:      "surrounding whitespace in href is trimmed",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><a href="  /spaced  ">Spaced</a></body></html>`,
			expected:  []string{"https://crawler-test.com/spaced"},
		},
		{
			name:      "HTML entities in href are decoded",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><a href="/search?a=1&amp;b=2">Search</a></body></html>`,
			expected:  []string{"https://crawler-test.com/search?a=1&b=2"},
		},
		{
			name:      "query string and fragment in href are kept",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><a href="/page?x=1#part">Page</a></body></html>`,
			expected:  []string{"https://crawler-test.com/page?x=1#part"},
		},

		// --- skipped links ---
		{
			name:      "a without href is skipped",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><a name="anchor">No href</a><a href="/kept">Kept</a></body></html>`,
			expected:  []string{"https://crawler-test.com/kept"},
		},
		{
			name:      "empty and whitespace-only href are skipped",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><a href="">Empty</a><a href="   ">Blank</a><a href="/kept">Kept</a></body></html>`,
			expected:  []string{"https://crawler-test.com/kept"},
		},
		{
			name:    "non-http schemes are skipped",
			baseURL: "https://crawler-test.com",
			inputBody: `<html><body>
				<a href="mailto:hi@crawler-test.com">Mail</a>
				<a href="tel:+123456789">Call</a>
				<a href="javascript:void(0)">JS</a>
				<a href="ftp://crawler-test.com/file">FTP</a>
				<a href="data:text/plain,hello">Data</a>
				<a href="/kept">Kept</a>
			</body></html>`,
			expected: []string{"https://crawler-test.com/kept"},
		},
		{
			name:      "unparseable href is skipped without failing",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><a href="http://[::1">Bad</a><a href="/bad%zz">Bad escape</a><a href="/kept">Kept</a></body></html>`,
			expected:  []string{"https://crawler-test.com/kept"},
		},
		{
			name:    "only a tags count",
			baseURL: "https://crawler-test.com",
			inputBody: `<html><head><link rel="stylesheet" href="/style.css"><base href="/ignored/"></head><body>
				<img src="/image.png">
				<map><area href="/area"></map>
				<script src="/app.js"></script>
				<a href="/kept">Kept</a>
			</body></html>`,
			expected: []string{"https://crawler-test.com/kept"},
		},
		{
			name:      "a tags inside comments are ignored",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><!-- <a href="/hidden">Hidden</a> --><a href="/visible">Visible</a></body></html>`,
			expected:  []string{"https://crawler-test.com/visible"},
		},

		// --- nothing to find ---
		{
			name:      "no a tags returns empty list",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><p>No links here.</p></body></html>`,
			expected:  []string{},
		},
		{
			name:      "empty input returns empty list",
			baseURL:   "https://crawler-test.com",
			inputBody: "",
			expected:  []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.baseURL)
			if err != nil {
				t.Fatalf("couldn't parse base URL %q: %v", tc.baseURL, err)
			}

			actual, err := getURLsFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// nil and empty slices are treated as equal
			if len(actual) != len(tc.expected) || (len(actual) > 0 && !reflect.DeepEqual(actual, tc.expected)) {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTMLNilBaseURL(t *testing.T) {
	_, err := getURLsFromHTML(`<html><body><a href="/path">Link</a></body></html>`, nil)
	if err == nil {
		t.Error("expected error for nil base URL, got nil")
	}
}
