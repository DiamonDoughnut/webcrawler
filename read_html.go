package main

import (
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
