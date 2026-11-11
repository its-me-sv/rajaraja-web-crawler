package utils

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetHeadingFromHTMLBasic(t *testing.T) {
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	actual := getHeadingFromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetHeadingFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputBody string
		expected  string
	}{
		{
			name:      "falls back to h2 when there is no h1",
			inputBody: "<html><body><h2>Sub Title</h2></body></html>",
			expected:  "Sub Title",
		},
		{
			name:      "prefers h1 even when h2 comes first",
			inputBody: "<html><body><h2>Sub Title</h2><h1>Main Title</h1></body></html>",
			expected:  "Main Title",
		},
		{
			name:      "returns the first h1 when there are several",
			inputBody: "<html><body><h1>First</h1><h1>Second</h1></body></html>",
			expected:  "First",
		},
		{
			name:      "returns the first h2 when falling back and there are several",
			inputBody: "<html><body><h2>First Sub</h2><h2>Second Sub</h2></body></html>",
			expected:  "First Sub",
		},
		{
			name:      "falls back to h2 when h1 is empty",
			inputBody: "<html><body><h1></h1><h2>Sub Title</h2></body></html>",
			expected:  "Sub Title",
		},
		{
			name:      "includes text of nested inline tags",
			inputBody: "<html><body><h1>Hello <span>big</span> <em>World</em></h1></body></html>",
			expected:  "Hello big World",
		},
		{
			name:      "finds h1 nested deep in the document",
			inputBody: "<html><body><div><header><section><h1>Deep Title</h1></section></header></div></body></html>",
			expected:  "Deep Title",
		},
		{
			name:      "trims surrounding whitespace",
			inputBody: "<html><body><h1>\n\t  Spaced Title  \n</h1></body></html>",
			expected:  "Spaced Title",
		},
		{
			name:      "decodes HTML entities",
			inputBody: "<html><body><h1>Tom &amp; Jerry &lt;3</h1></body></html>",
			expected:  "Tom & Jerry <3",
		},
		{
			name:      "handles malformed HTML with unclosed tag",
			inputBody: "<html><body><h1>Unclosed Title",
			expected:  "Unclosed Title",
		},
		{
			name:      "works on a fragment without html/body tags",
			inputBody: "<h1>Fragment Title</h1>",
			expected:  "Fragment Title",
		},
		{
			name:      "ignores lower-level headings",
			inputBody: "<html><body><h3>Three</h3><h4>Four</h4></body></html>",
			expected:  "",
		},
		{
			name:      "ignores headings inside comments",
			inputBody: "<html><body><!-- <h1>Hidden</h1> --><h2>Visible</h2></body></html>",
			expected:  "Visible",
		},
		{
			name:      "returns empty string when there are no headings",
			inputBody: "<html><body><p>Just a paragraph.</p></body></html>",
			expected:  "",
		},
		{
			name:      "returns empty string for empty input",
			inputBody: "",
			expected:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getHeadingFromHTML(tc.inputBody)
			if actual != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	inputBody := `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "Main paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputBody string
		expected  string
	}{
		{
			name:      "falls back to first p when there is no main",
			inputBody: "<html><body><p>First paragraph.</p><p>Second paragraph.</p></body></html>",
			expected:  "First paragraph.",
		},
		{
			name: "returns the first p inside main when there are several",
			inputBody: `<html><body><main>
				<p>Main first.</p>
				<p>Main second.</p>
			</main></body></html>`,
			expected: "Main first.",
		},
		{
			name: "falls back to first p outside when main has no p",
			inputBody: `<html><body>
				<p>Outside paragraph.</p>
				<main><div>No paragraphs here.</div></main>
			</body></html>`,
			expected: "Outside paragraph.",
		},
		{
			name: "falls back to first p outside when main's p is empty",
			inputBody: `<html><body>
				<p>Outside paragraph.</p>
				<main><p></p></main>
			</body></html>`,
			expected: "Outside paragraph.",
		},
		{
			name: "prefers main even when main comes after other paragraphs",
			inputBody: `<html><body>
				<header><p>Header text.</p></header>
				<p>Intro text.</p>
				<main><p>Main text.</p></main>
				<footer><p>Footer text.</p></footer>
			</body></html>`,
			expected: "Main text.",
		},
		{
			name: "finds p nested deep inside main",
			inputBody: `<html><body>
				<p>Outside paragraph.</p>
				<main><article><section><div><p>Nested main paragraph.</p></div></section></article></main>
			</body></html>`,
			expected: "Nested main paragraph.",
		},
		{
			name:      "includes text of nested inline tags",
			inputBody: "<html><body><main><p>Read <a href=\"/x\">this</a> and <strong>that</strong>.</p></main></body></html>",
			expected:  "Read this and that.",
		},
		{
			name:      "trims surrounding whitespace",
			inputBody: "<html><body><main><p>\n\t  Spaced paragraph.  \n</p></main></body></html>",
			expected:  "Spaced paragraph.",
		},
		{
			name:      "decodes HTML entities",
			inputBody: "<html><body><p>Fish &amp; chips &gt; salad</p></body></html>",
			expected:  "Fish & chips > salad",
		},
		{
			name:      "handles malformed HTML with unclosed tags",
			inputBody: "<html><body><main><p>Unclosed paragraph",
			expected:  "Unclosed paragraph",
		},
		{
			name:      "works on a fragment without html/body tags",
			inputBody: "<p>Fragment paragraph.</p>",
			expected:  "Fragment paragraph.",
		},
		{
			name:      "does not match tags that merely start with p",
			inputBody: "<html><body><pre>code block</pre><param name=\"x\"></body></html>",
			expected:  "",
		},
		{
			name:      "ignores paragraphs inside comments",
			inputBody: "<html><body><!-- <p>Hidden.</p> --><p>Visible.</p></body></html>",
			expected:  "Visible.",
		},
		{
			name:      "returns empty string when there are no paragraphs",
			inputBody: "<html><body><h1>Only a heading</h1><main><div>No p</div></main></body></html>",
			expected:  "",
		},
		{
			name:      "returns empty string for empty input",
			inputBody: "",
			expected:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tc.inputBody)
			if actual != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

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
