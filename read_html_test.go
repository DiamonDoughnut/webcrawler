package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetH1FromHTMLBasic(t *testing.T) {
	testCases := []struct {
		name      string
		inputBody string
		expected  string
	}{
		{
			name:      "basic h1",
			inputBody: "<html><body><h1>Test Title</h1></body></html>",
			expected:  "Test Title",
		},
		{
			name:      "no h1 tag",
			inputBody: "<html><body><h2>Not H1</h2></body></html>",
			expected:  "",
		},
		{
			name:      "multiple h1 tags - returns first",
			inputBody: "<html><body><h1>First</h1><h1>Second</h1></body></html>",
			expected:  "First",
		},
		{
			name:      "empty h1 tag",
			inputBody: "<html><body><h1></h1></body></html>",
			expected:  "",
		},
		{
			name:      "h1 with whitespace",
			inputBody: "<html><body><h1>  Spaced Title  </h1></body></html>",
			expected:  "  Spaced Title  ",
		},
		{
			name:      "h1 with nested elements",
			inputBody: "<html><body><h1>Title <span>with</span> nested</h1></body></html>",
			expected:  "Title with nested",
		},
		{
			name:      "malformed html",
			inputBody: "<h1>No closing tag",
			expected:  "No closing tag",
		},
		{
			name:      "empty string",
			inputBody: "",
			expected:  "",
		},
		{
			name:      "h1 with attributes",
			inputBody: `<html><body><h1 class="title" id="main">Attributed Title</h1></body></html>`,
			expected:  "Attributed Title",
		},
	}
	for _, tc := range testCases {
		actual := getH1FromHTML(tc.inputBody)
		if actual != tc.expected {
			t.Errorf("expected %q, got %q", tc.expected, actual)
		}
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	testCases := []struct {
		name      string
		inputBody string
		expected  string
	}{
		{
			name: "main paragraph priority",
			inputBody: `<html><body>
				<p>Outside paragraph.</p>
				<main>
					<p>Main paragraph.</p>
				</main>
			</body></html>`,
			expected: "Main paragraph.",
		},
		{
			name: "no main tag - fallback to any p",
			inputBody: `<html><body>
				<p>First paragraph.</p>
				<p>Second paragraph.</p>
			</body></html>`,
			expected: "First paragraph.",
		},
		{
			name: "main with no paragraphs - fallback to any p",
			inputBody: `<html><body>
				<p>Outside paragraph.</p>
				<main>
					<div>No paragraphs here</div>
				</main>
			</body></html>`,
			expected: "Outside paragraph.",
		},
		{
			name:      "no paragraphs at all",
			inputBody: "<html><body><div>No paragraphs</div></body></html>",
			expected:  "",
		},
		{
			name:      "empty paragraph in main",
			inputBody: `<html><body><main><p></p></main></body></html>`,
			expected:  "",
		},
		{
			name: "multiple main tags - first main wins",
			inputBody: `<html><body>
				<main><p>First main paragraph</p></main>
				<main><p>Second main paragraph</p></main>
			</body></html>`,
			expected: "First main paragraph",
		},
		{
			name: "nested paragraphs in main",
			inputBody: `<html><body>
				<main>
					<div><p>Nested paragraph</p></div>
				</main>
			</body></html>`,
			expected: "Nested paragraph",
		},
		{
			name:      "empty string",
			inputBody: "",
			expected:  "",
		},
		{
			name: "paragraph with nested elements",
			inputBody: `<html><body>
				<main>
					<p>Text with <strong>bold</strong> and <em>italic</em></p>
				</main>
			</body></html>`,
			expected: "Text with bold and italic",
		},
		{
			name: "paragraph with whitespace",
			inputBody: `<html><body>
				<main>
					<p>  Spaced paragraph  </p>
				</main>
			</body></html>`,
			expected: "  Spaced paragraph  ",
		},
	}

	for _, tc := range testCases {
		actual := getFirstParagraphFromHTML(tc.inputBody)
		if actual != tc.expected {
			t.Errorf("expected %q, got %q", tc.expected, actual)
		}
	}
}

func TestGetImagesFromHTMLRelative(t *testing.T) {
	testCases := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:      "relative img src",
			inputURL:  "https://blog.boot.dev",
			inputBody: `<html><body><img src="/logo.png" alt="Logo"></body></html>`,
			expected:  []string{"https://blog.boot.dev/logo.png"},
		},
		{
			name:      "multiple relative img src",
			inputURL:  "https://blog.boot.dev",
			inputBody: `<html><body><img src="/logo.png"><img src="/banner.jpg"></body></html>`,
			expected:  []string{"https://blog.boot.dev/logo.png", "https://blog.boot.dev/banner.jpg"},
		},
		{
			name:      "no img tags",
			inputURL:  "https://blog.boot.dev",
			inputBody: `<html><body><p>No images</p></body></html>`,
			expected:  []string{},
		},
		{
			name:      "img without src",
			inputURL:  "https://blog.boot.dev",
			inputBody: `<html><body><img alt="No src"></body></html>`,
			expected:  []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}

			actual, err := getImagesFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(tc.expected) == 0 && len(actual) == 0 {
				return
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTMLAbsolute(t *testing.T) {
	testCases := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:      "absolute anchor href",
			inputURL:  "https://blog.boot.dev",
			inputBody: `<html><body><a href="https://blog.boot.dev"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://blog.boot.dev"},
		},
		{
			name:      "relative anchor href",
			inputURL:  "https://blog.boot.dev",
			inputBody: `<html><body><a href="/path">Link</a></body></html>`,
			expected:  []string{"https://blog.boot.dev/path"},
		},
		{
			name:      "multiple anchors",
			inputURL:  "https://blog.boot.dev",
			inputBody: `<html><body><a href="/about">About</a><a href="/contact">Contact</a></body></html>`,
			expected:  []string{"https://blog.boot.dev/about", "https://blog.boot.dev/contact"},
		},
		{
			name:      "no anchor tags",
			inputURL:  "https://blog.boot.dev",
			inputBody: `<html><body><p>No links</p></body></html>`,
			expected:  []string{},
		},
		{
			name:      "anchor without href",
			inputURL:  "https://blog.boot.dev",
			inputBody: `<html><body><a>No href</a></body></html>`,
			expected:  []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}

			actual, err := getURLsFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tc.expected) == 0 && len(actual) == 0 {
				return
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestExtractPageData(t *testing.T) {
	testCases := []struct {
		name     string
		html     string
		pageURL  string
		expected PageData
	}{
		{
			name: "complete page data",
			html: `<html><body>
				<h1>Test Title</h1>
				<main><p>Main paragraph content</p></main>
				<a href="/about">About</a>
				<img src="/logo.png" alt="Logo">
			</body></html>`,
			pageURL: "https://blog.boot.dev/path",
			expected: PageData{
				URL:            "blog.boot.dev/path",
				H1:             "Test Title",
				FirstParagraph: "Main paragraph content",
				OutgoingLinks:  []string{"https://blog.boot.dev/about"},
				ImageURLs:      []string{"https://blog.boot.dev/logo.png"},
			},
		},
		{
			name:    "missing elements",
			html:    `<html><body><div>No structured content</div></body></html>`,
			pageURL: "https://blog.boot.dev",
			expected: PageData{
				URL:            "blog.boot.dev",
				H1:             "",
				FirstParagraph: "",
				OutgoingLinks:  nil,
				ImageURLs:      nil,
			},
		},
		{
			name: "multiple elements - first wins",
			html: `<html><body>
				<h1>First Title</h1>
				<h1>Second Title</h1>
				<main><p>First paragraph</p><p>Second paragraph</p></main>
				<a href="/link1">Link1</a>
				<a href="/link2">Link2</a>
			</body></html>`,
			pageURL: "https://blog.boot.dev",
			expected: PageData{
				URL:            "blog.boot.dev",
				H1:             "First Title",
				FirstParagraph: "First paragraph",
				OutgoingLinks:  []string{"https://blog.boot.dev/link1", "https://blog.boot.dev/link2"},
				ImageURLs:      nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := extractPageData(tc.html, tc.pageURL)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %+v, got %+v", tc.expected, actual)
			}
		})
	}
}
