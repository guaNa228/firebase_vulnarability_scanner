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
	if strings.HasPrefix(relativePath, "//") {
		// Prepend the scheme from the base URL
		parsedBase, err := url.Parse(baseURL)
		if err != nil || parsedBase.Scheme == "" {
			return "https:" + relativePath // Default to HTTPS
		}
		return parsedBase.Scheme + ":" + relativePath
	}

	if strings.HasPrefix(relativePath, "http://") || strings.HasPrefix(relativePath, "https://") {
		return relativePath
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return relativePath // Fallback to the relative path
	}
	ref, err := url.Parse(relativePath)
	if err != nil {
		return relativePath // Fallback to the relative path
	}
	return base.ResolveReference(ref).String()
}

// GetBaseURL extracts the base URL from a given URL.
func getBaseURL(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
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
