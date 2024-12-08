package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
)

// ResolveURL constructs an absolute URL from a base and a relative path.
func resolveURL(baseURL, relativePath string) string {
	if strings.HasPrefix(relativePath, "http://") || strings.HasPrefix(relativePath, "https://") {
		return relativePath
	}
	return strings.TrimSuffix(baseURL, "/") + "/" + strings.TrimPrefix(relativePath, "/")
}

// GetBaseURL extracts the base URL from a given URL.
func getBaseURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL // Fallback to the original URL
	}
	return fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
}

// Read links from a file.
func readLinksFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var links []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		links = append(links, scanner.Text())
	}
	return links, scanner.Err()
}

// Write results to a file.
func writeResultsToFile(filePath string, resultsChan <-chan string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for result := range resultsChan {
		_, err := writer.WriteString(result)
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
