package extractor

import "testing"

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
