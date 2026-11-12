package extractor

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetImagesFromHTMLRelative(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body><img src="/logo.png" alt="Logo"></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getImagesFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://crawler-test.com/logo.png"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

// getImagesFromHTML should return the src of every <img> tag in the page, in
// document order, with relative srcs resolved against baseURL.
//
// Expected rules:
//   - every <img src> in the document is found, wherever it is nested
//   - relative srcs are resolved to absolute URLs using baseURL
//   - absolute srcs are returned as-is (no normalisation)
//   - duplicates are kept
//   - only the src attribute counts: alt, width, srcset, data-src etc. are
//     neither required nor used
//   - <img> without src, or with an empty/whitespace-only src, is skipped
//   - non-http(s) srcs (data:, blob:, javascript:, ...) are skipped
//   - srcs that fail to parse are skipped, not fatal
//   - only <img> tags count (<source>, <link>, <a>, CSS backgrounds are ignored)
func TestGetImagesFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		inputBody string
		expected  []string
	}{
		// --- finding <img> tags ---
		{
			name:      "absolute src",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img src="https://crawler-test.com/logo.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/logo.png"},
		},
		{
			name:      "multiple images in document order",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img src="/one.png"><img src="/two.jpg"><img src="/three.gif"></body></html>`,
			expected:  []string{"https://crawler-test.com/one.png", "https://crawler-test.com/two.jpg", "https://crawler-test.com/three.gif"},
		},
		{
			name:    "images nested anywhere in the page",
			baseURL: "https://crawler-test.com",
			inputBody: `<html><body>
				<header><a href="/"><img src="/logo.png"></a></header>
				<main><article><figure><img src="/hero.jpg"><figcaption>Hero</figcaption></figure></article></main>
				<footer><div><span><img src="/badge.svg"></span></div></footer>
			</body></html>`,
			expected: []string{"https://crawler-test.com/logo.png", "https://crawler-test.com/hero.jpg", "https://crawler-test.com/badge.svg"},
		},
		{
			name:      "img inside picture element",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><picture><source srcset="/photo.webp" type="image/webp"><img src="/photo.jpg"></picture></body></html>`,
			expected:  []string{"https://crawler-test.com/photo.jpg"},
		},
		{
			name:      "uppercase tag and attribute names",
			baseURL:   "https://crawler-test.com",
			inputBody: `<HTML><BODY><IMG SRC="/upper.png"></BODY></HTML>`,
			expected:  []string{"https://crawler-test.com/upper.png"},
		},
		{
			name:      "self-closing, single-quoted and unquoted src",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img src="/a.png" /><img src='/b.png'><img src=/c.png></body></html>`,
			expected:  []string{"https://crawler-test.com/a.png", "https://crawler-test.com/b.png", "https://crawler-test.com/c.png"},
		},
		{
			name:      "duplicate images are kept",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img src="/same.png"><img src="/same.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/same.png", "https://crawler-test.com/same.png"},
		},
		{
			name:      "works on a fragment without html/body tags",
			baseURL:   "https://crawler-test.com",
			inputBody: `<img src="/fragment.png">`,
			expected:  []string{"https://crawler-test.com/fragment.png"},
		},
		{
			name:      "handles malformed HTML with unclosed tags",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><div><p><img src="/unclosed.png">`,
			expected:  []string{"https://crawler-test.com/unclosed.png"},
		},

		// --- relative URL resolution ---
		{
			name:      "root-relative src",
			baseURL:   "https://crawler-test.com/blog/post",
			inputBody: `<html><body><img src="/static/logo.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/static/logo.png"},
		},
		{
			name:      "relative src against base ending in slash",
			baseURL:   "https://crawler-test.com/blog/",
			inputBody: `<html><body><img src="cover.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/blog/cover.png"},
		},
		{
			name:      "relative src against base not ending in slash replaces last segment",
			baseURL:   "https://crawler-test.com/blog/post",
			inputBody: `<html><body><img src="cover.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/blog/cover.png"},
		},
		{
			name:      "dot-slash relative src",
			baseURL:   "https://crawler-test.com/blog/",
			inputBody: `<html><body><img src="./images/cover.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/blog/images/cover.png"},
		},
		{
			name:      "parent directory relative src",
			baseURL:   "https://crawler-test.com/blog/2024/post",
			inputBody: `<html><body><img src="../../assets/pic.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/assets/pic.png"},
		},
		{
			name:      "parent directory beyond root stays at root",
			baseURL:   "https://crawler-test.com/blog/",
			inputBody: `<html><body><img src="../../../top.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/top.png"},
		},
		{
			name:      "base query string is not carried over to relative src",
			baseURL:   "https://crawler-test.com/blog/post?ref=home",
			inputBody: `<html><body><img src="pic.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/blog/pic.png"},
		},
		{
			name:      "protocol-relative src uses base scheme",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img src="//cdn.crawler-test.com/pic.png"></body></html>`,
			expected:  []string{"https://cdn.crawler-test.com/pic.png"},
		},
		{
			name:      "protocol-relative src on http page uses http",
			baseURL:   "http://crawler-test.com",
			inputBody: `<html><body><img src="//cdn.crawler-test.com/pic.png"></body></html>`,
			expected:  []string{"http://cdn.crawler-test.com/pic.png"},
		},
		{
			name:      "relative src keeps base port",
			baseURL:   "http://localhost:8080/docs/",
			inputBody: `<html><body><img src="diagram.png"></body></html>`,
			expected:  []string{"http://localhost:8080/docs/diagram.png"},
		},
		{
			name:      "mix of relative and absolute srcs",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img src="/local.png"><img src="https://images.other.com/remote.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/local.png", "https://images.other.com/remote.png"},
		},

		// --- src values ---
		{
			name:      "absolute src on another domain is kept as-is",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img src="https://Images.Other.com/Pics/Cat.PNG"></body></html>`,
			expected:  []string{"https://Images.Other.com/Pics/Cat.PNG"},
		},
		{
			name:      "surrounding whitespace in src is trimmed",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img src="  /spaced.png  "></body></html>`,
			expected:  []string{"https://crawler-test.com/spaced.png"},
		},
		{
			name:      "query string and HTML entities in src",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img src="/resize?w=200&amp;h=100"></body></html>`,
			expected:  []string{"https://crawler-test.com/resize?w=200&h=100"},
		},
		{
			name:      "src with spaces is percent-encoded",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img src="/my photo.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/my%20photo.png"},
		},

		// --- missing or extra attributes ---
		{
			name:      "src is used when alt is missing",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img src="/no-alt.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/no-alt.png"},
		},
		{
			name:      "src is used alongside many other attributes",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img class="hero" alt="Hero" width="100" height="50" loading="lazy" src="/hero.png" title="Hero"></body></html>`,
			expected:  []string{"https://crawler-test.com/hero.png"},
		},
		{
			name:      "img without src is skipped",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img alt="No source"><img src="/kept.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/kept.png"},
		},
		{
			name:      "img with no attributes at all is skipped",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img><img src="/kept.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/kept.png"},
		},
		{
			name:      "empty, valueless and whitespace-only src are skipped",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img src=""><img src><img src="   "><img src="/kept.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/kept.png"},
		},
		{
			name:      "srcset-only img is skipped",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img srcset="/small.png 1x, /large.png 2x"><img src="/kept.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/kept.png"},
		},
		{
			name:      "lazy-load data-src without src is skipped",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img data-src="/lazy.png" class="lazyload"><img src="/kept.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/kept.png"},
		},
		{
			name:      "src is used and srcset ignored when both are present",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img src="/default.png" srcset="/small.png 1x, /large.png 2x"></body></html>`,
			expected:  []string{"https://crawler-test.com/default.png"},
		},

		// --- skipped images ---
		{
			name:    "non-http srcs are skipped",
			baseURL: "https://crawler-test.com",
			inputBody: `<html><body>
				<img src="data:image/png;base64,iVBORw0KGgo=">
				<img src="blob:https://crawler-test.com/1234">
				<img src="javascript:alert(1)">
				<img src="ftp://crawler-test.com/pic.png">
				<img src="/kept.png">
			</body></html>`,
			expected: []string{"https://crawler-test.com/kept.png"},
		},
		{
			name:      "unparseable src is skipped without failing",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><img src="http://[::1"><img src="/bad%zz.png"><img src="/kept.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/kept.png"},
		},
		{
			name:    "only img tags count",
			baseURL: "https://crawler-test.com",
			inputBody: `<html><head>
				<link rel="icon" href="/favicon.ico">
				<meta property="og:image" content="/og.png">
			</head><body>
				<a href="/photo.png">Photo link</a>
				<div style="background-image: url(/bg.png)"></div>
				<video poster="/poster.png"></video>
				<script src="/app.js"></script>
				<iframe src="/embed"></iframe>
				<img src="/kept.png">
			</body></html>`,
			expected: []string{"https://crawler-test.com/kept.png"},
		},
		{
			name:      "images inside comments are ignored",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><!-- <img src="/hidden.png"> --><img src="/visible.png"></body></html>`,
			expected:  []string{"https://crawler-test.com/visible.png"},
		},

		// --- nothing to find ---
		{
			name:      "no img tags returns empty list",
			baseURL:   "https://crawler-test.com",
			inputBody: `<html><body><p>No images here.</p></body></html>`,
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

			actual, err := getImagesFromHTML(tc.inputBody, baseURL)
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

func TestGetImagesFromHTMLNilBaseURL(t *testing.T) {
	_, err := getImagesFromHTML(`<html><body><img src="/logo.png"></body></html>`, nil)
	if err == nil {
		t.Error("expected error for nil base URL, got nil")
	}
}
