package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// AnalyzeLink analyzes a URL for Firebase configuration.
func analyzeLink(url string, resultsChan chan<- string, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	defer func() { <-sem }() // Release semaphore slot

	htmlContent, err := fetchHTML(url)
	if err != nil {
		fmt.Printf("Error fetching HTML for %s: %v\n", url, err)
		resultsChan <- fmt.Sprintf("%s: No config\n", url)
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Printf("Error parsing HTML for %s: %v\n", url, err)
		resultsChan <- fmt.Sprintf("%s: No config\n", url)
		return
	}

	baseURL := getBaseURL(url)
	foundConfig := false

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			scriptURL := resolveURL(baseURL, src)
			scriptContent, err := fetchScriptContent(scriptURL)
			if err != nil {
				fmt.Printf("Error fetching script content for %s: %v\n", scriptURL, err)
				return
			}
			fmt.Printf("Scanning script %s\n", scriptURL)
			if containsFirebaseConfig(scriptContent) {
				config := extractFirebaseKeys(scriptContent)
				resultsChan <- fmt.Sprintf("%s: %v\n", url, config)
				foundConfig = true
			}
		} else {
			scriptContent := s.Text()
			if containsFirebaseConfig(scriptContent) {
				config := extractFirebaseKeys(scriptContent)
				resultsChan <- fmt.Sprintf("%s: %v\n", url, config)
				foundConfig = true
			}
		}
	})

	if !foundConfig {
		resultsChan <- fmt.Sprintf("%s: No config\n", url)
	}
}
