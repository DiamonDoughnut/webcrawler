package main

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		expected    string
		expectError bool
	}{
		{
			name:     "remove https scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme and trailing /",
			inputURL: "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove http scheme",
			inputURL: "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme and trailing /",
			inputURL: "http://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "ensure path integrity",
			inputURL: "blog.boot.dev/path/dir/target",
			expected: "blog.boot.dev/path/dir/target",
		},
		{
			name:     "ensure extension integrity",
			inputURL: "blog.boot.dev/path/target.ext",
			expected: "blog.boot.dev/path/target.ext",
		},
		{
			name:     "remove scheme and ensure integrity",
			inputURL: "https://blog.boot.dev/path/dir/target.ext",
			expected: "blog.boot.dev/path/dir/target.ext",
		},
		{
			name:     "add domain to absolute path",
			inputURL: "/path/dir/target.ext",
			expected: "blog.boot.dev/path/dir/target.ext",
		},
		{
			name:     "add domain to relative path",
			inputURL: "path/dir/target.ext",
			expected: "blog.boot.dev/path/dir/target.ext",
		},
		{
			name:     "root path with trailing slash",
			inputURL: "https://blog.boot.dev/",
			expected: "blog.boot.dev",
		},
		{
			name:     "root path without trailing slash",
			inputURL: "https://blog.boot.dev",
			expected: "blog.boot.dev",
		},
		{
			name:     "preserve query parameters",
			inputURL: "https://blog.boot.dev/path?param=value",
			expected: "blog.boot.dev/path?param=value",
		},
		{
			name:     "trim fragments",
			inputURL: "https://blog.boot.dev/path#section",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove standard https port",
			inputURL: "https://blog.boot.dev:443/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "preserve custom port",
			inputURL: "https://blog.boot.dev:8080/path",
			expected: "blog.boot.dev:8080/path",
		},
		{
			name:     "lowercase domain",
			inputURL: "https://BLOG.BOOT.DEV/Path",
			expected: "blog.boot.dev/Path",
		},
	}

	errorTests := []struct {
		name     string
		inputURL string
	}{
		{
			name:     "empty string",
			inputURL: "",
		},
		{
			name:     "whitespace only",
			inputURL: "   ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("FAIL: unexpected error: %v", err)
				return
			}
			if actual != tc.expected {
				t.Errorf("FAIL: expected URL: %v, actual: %v", tc.expected, actual)
			}
		})
	}

	for _, tc := range errorTests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := normalizeURL(tc.inputURL)
			if err != nil {
				return
			}
			t.Errorf("FAIL: expected error but got none")
		})
	}
}
