package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func shouldFilterScript(scriptURL string) bool {
	// Define patterns to filter out
	patterns := []string{
		`\.min\.`,
		`\d+\.\d+\.\d+`,
		`ver=`,
		`googletagmanager`,
	}

	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, scriptURL)
		if err != nil {
			fmt.Printf("Error matching pattern %s: %v\n", pattern, err)
			continue
		}
		if matched {
			return true
		}
	}
	return false
}

func analyzeLink(url string, resultsChan chan<- SecurityConfig, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	defer func() { <-sem }() // Release semaphore slot

	htmlContent, headers, baseURL, err := fetchHTML(url)

	res := SecurityConfig{
		URL:          baseURL,
		CSPHeader:    headers["CSP Header Present"],
		XFrameHeader: headers["X-Frame-Options Header Present"],
		Creds:        map[string]string{},
	}

	if err != nil {
		fmt.Printf("Error fetching HTML for %s: %v\n", url, err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Printf("Error parsing HTML for %s: %v\n", url, err)
		return
	}

	var scriptWg sync.WaitGroup
	scriptSem := make(chan struct{}, 10) // Limit to 10 concurrent script fetches

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		scriptWg.Add(1)
		go func(s *goquery.Selection) {
			defer scriptWg.Done()
			scriptSem <- struct{}{}        // Acquire semaphore slot
			defer func() { <-scriptSem }() // Release semaphore slot

			if src, exists := s.Attr("src"); exists {
				scriptURL := resolveURL(baseURL, src)

				if shouldFilterScript(scriptURL) {
					fmt.Printf("Filtering out script %s\n", scriptURL)
					return
				}

				scriptContent, err := fetchScriptContent(scriptURL)
				if err != nil {
					fmt.Printf("Error fetching script content for %s: %v\n", scriptURL, err)
					return
				}
				fmt.Printf("Scanning script %s\n", scriptURL)
				config := findSensitiveData(scriptContent)
				for key, value := range config {
					res.Creds[key] = value
				}
			} else {
				scriptContent := s.Text()
				config := findSensitiveData(scriptContent)
				for key, value := range config {
					res.Creds[key] = value
				}
			}
		}(s)
	})

	scriptWg.Wait()

	resultsChan <- res
}
