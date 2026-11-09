package utils

import "testing"

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
