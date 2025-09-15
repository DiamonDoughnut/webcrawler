package main

import (
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
