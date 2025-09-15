package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(inputUrl string) (string, error) {
	inputUrl = strings.TrimSpace(inputUrl)
	if inputUrl == "" {
		return "", fmt.Errorf("url cannot be blank")
	}

	parsedURL, err := url.Parse(inputUrl)
	if err != nil {
		return "", err
	}

	host := strings.ToLower(parsedURL.Host)
	if host == "" {
		host = "blog.boot.dev"
	}

	host = strings.TrimSuffix(host, ":443")
	result := ""
	if strings.HasPrefix(parsedURL.Path, "blog.boot.dev") {
		result = strings.TrimSuffix(parsedURL.Path, "/")
	} else {
		if !strings.HasPrefix(parsedURL.Path, "/") && !strings.HasSuffix(host, "/") {
			host += "/"
		}
		result = host + strings.TrimSuffix(parsedURL.Path, "/")
	}

	if parsedURL.RawQuery != "" {
		result += "?" + parsedURL.RawQuery
	}

	result = strings.TrimSuffix(result, "/")

	return result, nil
}
