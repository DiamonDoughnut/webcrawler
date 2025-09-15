package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	baseUrl := args[0]
	fmt.Printf("starting crawl of: %s\n", baseUrl)
	data, err := getHTML(baseUrl)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(data)
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
