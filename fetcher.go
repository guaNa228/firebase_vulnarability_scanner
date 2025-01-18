package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func checkCSPHeader(resp *http.Response) bool {
	cspHeader := resp.Header.Get("Content-Security-Policy")
	return cspHeader != ""
}

func checkXFrameOptionsHeader(resp *http.Response) bool {
	xFrameOptionsHeader := resp.Header.Get("X-Frame-Options")
	return xFrameOptionsHeader != ""
}

func fetchHTML(url string) (string, map[string]bool, string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(addHTTPSIfNeeded(url))
	if err != nil {
		return "", map[string]bool{}, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", map[string]bool{}, "", fmt.Errorf("failed to fetch URL: %s", url)
	}

	// Check headers
	cspHeaderPresent := checkCSPHeader(resp)
	xFrameOptionsHeaderPresent := checkXFrameOptionsHeader(resp)

	// Pass header check results to the channel
	results := map[string]bool{
		"CSP Header Present":             cspHeaderPresent,
		"X-Frame-Options Header Present": xFrameOptionsHeaderPresent,
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", map[string]bool{}, "", err
	}

	html, err := doc.Html()
	if err != nil {
		return "", map[string]bool{}, "", err
	}

	// Get the final URL after any redirects
	finalURL := resp.Request.URL.String()

	return html, results, finalURL, nil
}

// FetchScriptContent fetches the content of a script from a given URL.
func fetchScriptContent(url string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch script: %s", url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
