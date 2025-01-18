package snyk_scraper

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// FetchHTML fetches the HTML content of a given URL.
func fetchHTML(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL: %s", url)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// ParseVulnerableVersion parses the vulnerable version from the given URL.
func parseVulnerableVersion(url string) (string, error) {
	doc, err := fetchHTML(url)
	if err != nil {
		return "", err
	}

	version := doc.Find("table tr:first-child td:first-child").Text()
	return strings.TrimSpace(version), nil
}

// ScrapeSnyk scrapes the Snyk vulnerability pages and extracts the library name and max vulnerable version.
func ScrapeSnyk() {
	var wg sync.WaitGroup
	resultsChan := make(chan map[string]string, 30)

	for i := 1; i <= 30; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			url := fmt.Sprintf("https://security.snyk.io/vuln/npm/%d", page)
			doc, err := fetchHTML(url)
			if err != nil {
				fmt.Printf("Failed to fetch page %d: %v\n", page, err)
				return
			}

			doc.Find("table tr").Each(func(index int, row *goquery.Selection) {
				if index%2 == 1 {
					link, exists := row.Find("td:nth-child(2) a").Attr("href")
					if exists {
						fullLink := "https://security.snyk.io" + link
						version, err := parseVulnerableVersion(fullLink)
						if err != nil {
							fmt.Printf("Failed to parse version for %s: %v\n", fullLink, err)
							return
						}

						libraryName := row.Find("td:nth-child(2) a").Text()
						resultsChan <- map[string]string{libraryName: version}
					}
				}
			})
		}(i)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	results := make(map[string]string)
	for result := range resultsChan {
		for lib, version := range result {
			if existingVersion, exists := results[lib]; !exists || version > existingVersion {
				results[lib] = version
			}
		}
	}

	// Print or store the results
	for lib, version := range results {
		fmt.Printf("Library: %s, Max Vulnerable Version: %s\n", lib, version)
	}
}
