package extractors

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
