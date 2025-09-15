# Web Crawler

A concurrent web crawler built in Go that extracts structured data from websites and exports results to CSV format.

## Features

- **Concurrent crawling** with configurable concurrency limits
- **Domain-restricted crawling** - only crawls pages within the specified base domain
- **Content extraction** - extracts H1 titles, first paragraphs, outgoing links, and image URLs
- **HTML content filtering** - only processes `text/html` content types
- **CSV export** - generates spreadsheet-compatible reports
- **Configurable limits** - set maximum pages to crawl
- **Duplicate prevention** - tracks visited pages to avoid re-crawling

## Installation

```bash
git clone <repository-url>
cd webcrawler
go mod tidy
```

## Usage

```bash
go run . <base_url> <concurrency> <max_pages>
```

### Parameters

- `base_url` - The starting URL to crawl (e.g., `https://example.com`)
- `concurrency` - Number of concurrent crawlers (e.g., `5`)
- `max_pages` - Maximum number of pages to crawl (e.g., `50`)

### Example

```bash
go run . https://blog.boot.dev 5 25
```

This will:
- Start crawling from `https://blog.boot.dev`
- Use 5 concurrent crawlers
- Stop after crawling 25 pages
- Generate a `report.csv` file with the results

## Output

The crawler generates a CSV file (`report.csv`) with the following columns:

- **page_url** - Normalized URL of the crawled page
- **h1** - First H1 heading found on the page
- **first_paragraph** - First paragraph text (prioritizes content in `<main>` tags)
- **outgoing_link_urls** - Semicolon-separated list of outgoing links
- **image_urls** - Semicolon-separated list of image URLs

The "report.csv" file in this repo is an example of the given output

## Technical Details

- **Language**: Go
- **Dependencies**: `github.com/PuerkitoBio/goquery` for HTML parsing
- **Concurrency**: Uses goroutines with channel-based concurrency control
- **URL Normalization**: Removes schemes and trailing slashes for consistent tracking
- **Thread Safety**: Mutex-protected shared data structures

## Limitations

- Only crawls pages within the same domain as the base URL
- Only processes HTML content (skips PDFs, images, etc.)
- Respects the configured maximum page limit
- Does not follow robots.txt or implement crawl delays