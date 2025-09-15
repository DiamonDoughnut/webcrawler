package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	args := os.Args[1:]
	if len(args) < 3 {
		fmt.Println("usage: crawler <url> <maxPages> <concurrency>")
		os.Exit(1)
	}

	baseUrl := args[0]
	concurrency, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("invalid concurrency: %v\n", err)
		os.Exit(1)
	}
	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Printf("invalid maxPages: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %s\n", baseUrl)
	pages := make(map[string]PageData)
	urlStr, err := normalizeURL(baseUrl)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	urlStr = "https://" + urlStr
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	waitGroup := sync.WaitGroup{}
	config := config{
		pages:              pages,
		baseURL:            parsedUrl,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, concurrency),
		wg:                 &waitGroup,
		maxPages:           maxPages,
	}
	config.wg.Add(1)
	go config.crawlPage(baseUrl)
	config.wg.Wait()
	for key, value := range pages {
		fmt.Printf("%s: %d\n", key, value.Visits)
	}
}

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "BootCrawler/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("response status code: %d", resp.StatusCode)
	}
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return "", fmt.Errorf("site content not html text: %v", resp.Header.Get("Content-Type"))
	}
	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	html := string(htmlBytes)
	return html, nil
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()
	if !strings.Contains(rawCurrentURL, cfg.baseURL.String()) {
		return
	}
	currentUrl, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Println(err.Error())
		return
	}
	fetch_url := "https://" + currentUrl
	unvisited := cfg.addPageVisit(currentUrl, PageData{})
	if !unvisited {
		return
	}

	html, err := getHTML(fetch_url)
	if err != nil {
		log.Println(err.Error())
		return
	}
	pageData := extractPageData(html, rawCurrentURL)
	cfg.updatePageData(currentUrl, pageData)

	parsedURL, err := url.Parse(fetch_url)
	if err != nil {
		log.Println(err.Error())
		return
	}
	links, err := getURLsFromHTML(html, parsedURL)
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Printf("Page Crawled: %s\n", rawCurrentURL)
	for _, link := range links {
		cfg.wg.Add(1)
		go cfg.crawlPage(link)
	}

}

func (cfg *config) addPageVisit(normalizedURL string, pageData PageData) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if len(cfg.pages) >= cfg.maxPages {
		return false
	}

	if existing, ok := cfg.pages[normalizedURL]; ok {
		existing.Visits++
		cfg.pages[normalizedURL] = existing
		return false
	}

	pageData.Visits = 1
	cfg.pages[normalizedURL] = pageData
	return true
}

func (cfg *config) updatePageData(normalizedURL string, pageData PageData) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if existing, ok := cfg.pages[normalizedURL]; ok {
		pageData.Visits = existing.Visits
		cfg.pages[normalizedURL] = pageData
	}
}
