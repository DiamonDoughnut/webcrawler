package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getH1FromHTML(html string) string {
	doc := getGoDocFromHTML(html)
	header := doc.Find("h1")
	return header.First().Text()
}

func getFirstParagraphFromHTML(html string) string {
	doc := getGoDocFromHTML(html)
	main := doc.Find("main")
	firstP := main.First().Find("p")
	result := firstP.First().Text()
	if result == "" {
		fallback := doc.Find("p")
		result = fallback.First().Text()
	}
	return result
}

func getGoDocFromHTML(html string) *goquery.Document {
	reader := strings.NewReader(html)
	godoc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil
	}
	return godoc
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc := getGoDocFromHTML(htmlBody)
	var urls []string

	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		absoluteURL, err := baseURL.Parse(href)
		if err != nil {
			return
		}
		urls = append(urls, absoluteURL.String())
	})

	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc := getGoDocFromHTML(htmlBody)
	var urls []string

	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}
		absoluteURL, err := baseURL.Parse(src)
		if err != nil {
			return
		}
		urls = append(urls, absoluteURL.String())
	})

	return urls, nil
}

type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func extractPageData(html, pageURL string) PageData {
	var data PageData
	urlString, err := normalizeURL(pageURL)
	if err != nil {
		return PageData{}
	}
	data.URL = urlString
	header := getH1FromHTML(html)
	data.H1 = header
	firstP := getFirstParagraphFromHTML(html)
	data.FirstParagraph = firstP
	urlString = "https://" + urlString
	urlObj, err := url.Parse(urlString)
	if err != nil {
		return PageData{}
	}
	hrefs, err := getURLsFromHTML(html, urlObj)
	if err != nil {
		return PageData{}
	}
	data.OutgoingLinks = hrefs
	imgs, err := getImagesFromHTML(html, urlObj)
	if err != nil {
		return PageData{}
	}
	data.ImageURLs = imgs
	return data
}
