package main

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type SecurityConfig struct {
	URL          string
	CSPHeader    bool
	XFrameHeader bool
	Creds        map[string]string
}

// ResolveURL constructs an absolute URL from a base and a relative path.
func resolveURL(baseURL, path string) string {
	// Ensure the base URL is correctly parsed
	parsedBase, err := url.Parse(baseURL)
	if err != nil || parsedBase.Scheme == "" || parsedBase.Host == "" {
		return path // Return the path as-is if the base URL is invalid
	}

	// Ensure the base URL has a trailing slash if needed
	if !strings.HasSuffix(parsedBase.Path, "/") {
		parsedBase.Path += "/"
	}

	// Handle paths that are protocol-relative (e.g., "//example.com")
	if strings.HasPrefix(path, "//") {
		return parsedBase.Scheme + ":" + path
	}

	// Handle absolute URLs (e.g., "http://..." or "https://...")
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}

	// Handle relative paths starting with "/" (e.g., "/s/.../vendor.js")
	if strings.HasPrefix(path, "/") {
		return parsedBase.Scheme + "://" + parsedBase.Host + path
	}

	// Handle relative URLs (e.g., "../path/to/file.js")
	ref, err := url.Parse(path)
	if err != nil {
		return path // Return the path as-is if parsing fails
	}
	return parsedBase.ResolveReference(ref).String()
}

func addHTTPSIfNeeded(url string) string {
	httpsPrefix := "https://"

	if strings.HasPrefix(url, httpsPrefix) {
		return url
	} else {
		return httpsPrefix + url
	}
}

func cleanInput(domains []string) []string {
	domainSet := make(map[string]struct{})
	var cleanedDomains []string

	for _, domain := range domains {
		domain = extractDomain(domain)

		// Add to set if not already present
		if _, exists := domainSet[domain]; !exists {
			domainSet[domain] = struct{}{}
			cleanedDomains = append(cleanedDomains, domain)
		}
	}

	return cleanedDomains
}

func extractDomain(url string) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimSuffix(url, "/")

	return url
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := d.Seconds() - float64(hours*3600+minutes*60)

	if hours > 0 {
		return fmt.Sprintf("%dh%dm%.2fs", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm%.2fs", minutes, seconds)
	}
	return fmt.Sprintf("%.2fs", seconds)
}
