package utils

import "testing"

// normalizeURL should map every URL that points to the same page onto one
// canonical string, so that two URLs can be compared with ==.
//
// Expected rules:
//   - scheme is dropped (http and https are treated as the same page)
//   - host is lowercased; path case is preserved (paths are case-sensitive)
//   - default ports (:80 for http, :443 for https) are dropped; other ports are kept
//   - trailing slashes on the path are dropped
//   - fragments (#...) are dropped (they don't change the page)
//   - query strings are kept (they can change the page)
//   - surrounding whitespace is ignored
//   - input without a scheme or host is an error
func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		expected    string
		expectError bool
	}{
		// --- scheme and trailing slash ---
		{
			name:     "remove scheme",
			inputURL: "https://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "remove https scheme and trailing slash",
			inputURL: "https://www.boot.dev/blog/path/",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "remove http scheme",
			inputURL: "http://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "remove http scheme and trailing slash",
			inputURL: "http://www.boot.dev/blog/path/",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "remove multiple trailing slashes",
			inputURL: "https://www.boot.dev/blog/path//",
			expected: "www.boot.dev/blog/path",
		},

		// --- root / host only ---
		{
			name:     "host only",
			inputURL: "https://www.boot.dev",
			expected: "www.boot.dev",
		},
		{
			name:     "host with root slash",
			inputURL: "https://www.boot.dev/",
			expected: "www.boot.dev",
		},

		// --- case sensitivity ---
		{
			name:     "lowercase uppercase scheme",
			inputURL: "HTTPS://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "lowercase mixed-case host",
			inputURL: "https://WWW.Boot.DEV/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "preserve path case",
			inputURL: "https://www.boot.dev/Blog/Path",
			expected: "www.boot.dev/Blog/Path",
		},

		// --- ports ---
		{
			name:     "drop default https port",
			inputURL: "https://www.boot.dev:443/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "drop default http port",
			inputURL: "http://www.boot.dev:80/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "keep non-default port",
			inputURL: "http://localhost:8080/blog/path/",
			expected: "localhost:8080/blog/path",
		},

		// --- fragments ---
		{
			name:     "drop fragment",
			inputURL: "https://www.boot.dev/blog/path#section-2",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "drop fragment and trailing slash",
			inputURL: "https://www.boot.dev/blog/path/#section-2",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "drop empty fragment",
			inputURL: "https://www.boot.dev/blog/path#",
			expected: "www.boot.dev/blog/path",
		},

		// --- query strings ---
		{
			name:     "keep query string",
			inputURL: "https://www.boot.dev/blog/path?page=2",
			expected: "www.boot.dev/blog/path?page=2",
		},
		{
			name:     "keep query string, drop trailing slash before it",
			inputURL: "https://www.boot.dev/blog/path/?page=2",
			expected: "www.boot.dev/blog/path?page=2",
		},
		{
			name:     "keep query string, drop fragment",
			inputURL: "http://www.boot.dev/blog/path?page=2#top",
			expected: "www.boot.dev/blog/path?page=2",
		},

		// --- misc valid inputs ---
		{
			name:     "trim surrounding whitespace",
			inputURL: "  https://www.boot.dev/blog/path/  ",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "different subdomain is a different host",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "deeply nested path",
			inputURL: "https://www.boot.dev/a/b/c/d/e/",
			expected: "www.boot.dev/a/b/c/d/e",
		},
		{
			name:     "ip address host",
			inputURL: "http://127.0.0.1:3000/blog/path/",
			expected: "127.0.0.1:3000/blog/path",
		},

		// --- invalid inputs ---
		{
			name:        "empty string",
			inputURL:    "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			inputURL:    "   ",
			expectError: true,
		},
		{
			name:        "missing scheme",
			inputURL:    "www.boot.dev/blog/path",
			expectError: true,
		},
		{
			name:        "scheme without host",
			inputURL:    "https://",
			expectError: true,
		},
		{
			name:        "malformed url",
			inputURL:    "://www.boot.dev",
			expectError: true,
		},
		{
			name:        "invalid percent encoding",
			inputURL:    "https://www.boot.dev/%zz",
			expectError: true,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := NormalizeURL(tc.inputURL)
			if tc.expectError {
				if err == nil {
					t.Errorf("Test %v - '%s' FAIL: expected error for input %q, got result: %q", i, tc.name, tc.inputURL, actual)
				}
				return
			}
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

// Every variant of the same page must normalize to the same string.
func TestNormalizeURLEquivalence(t *testing.T) {
	variants := []string{
		"https://www.boot.dev/blog/path/",
		"https://www.boot.dev/blog/path",
		"http://www.boot.dev/blog/path/",
		"http://www.boot.dev/blog/path",
		"HTTPS://WWW.BOOT.DEV/blog/path",
		"https://www.boot.dev:443/blog/path",
		"http://www.boot.dev:80/blog/path/",
		"https://www.boot.dev/blog/path#intro",
	}

	const want = "www.boot.dev/blog/path"
	for _, v := range variants {
		got, err := NormalizeURL(v)
		if err != nil {
			t.Errorf("normalizeURL(%q) unexpected error: %v", v, err)
			continue
		}
		if got != want {
			t.Errorf("normalizeURL(%q) = %q, want %q", v, got, want)
		}
	}
}
