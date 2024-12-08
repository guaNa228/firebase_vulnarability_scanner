package scraper

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sync"

	"golang.org/x/net/html"
)

// Function to fetch the HTML content of a URL
func fetchHTML(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// Function to extract href attributes from <a> elements with specific id and class
func extractLinks(node *html.Node, idFilter, classFilter string) []string {
	var links []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			id := ""
			class := ""
			href := ""
			for _, a := range n.Attr {
				if a.Key == "id" {
					id = a.Val
				}
				if a.Key == "class" {
					class = a.Val
				}
				if a.Key == "href" {
					href = a.Val
				}
			}
			if id == idFilter && class == classFilter && href != "" {
				links = append(links, href)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)

	return links
}

// Function to resolve a short link to its actual URL
func resolveShortLink(shortURL string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow following redirects
			return nil
		},
	}

	resp, err := client.Head(shortURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Return the final URL after redirects
	return resp.Request.URL.String(), nil
}

// Function to scrape the second-level links
func scrapeSecondLevelLinks(url string) (string, error) {
	node, err := fetchHTML(url)
	if err != nil {
		return "", err
	}

	// Find the element with id="specialButton" and class="button-plus2"
	links := extractLinks(node, "specialButton", "button-plus2")
	if len(links) > 0 {
		return links[0], nil // Return the first match
	}
	return "", nil
}

// ScrapeLinks function to scrape links and follow to extract second-level links
func ScrapeLinks() {
	var wg sync.WaitGroup
	linksChan := make(chan string)
	var finalLinks []string

	for i := 1; i <= 20; i++ {
		wg.Add(1)
		url := fmt.Sprintf("https://www.aixploria.com/en/last-ai/page/%d/", i)
		go func(url string) {
			defer wg.Done()
			node, err := fetchHTML(url)
			if err != nil {
				fmt.Printf("Failed to fetch %s: %v\n", url, err)
				return
			}

			// Extract initial links
			firstLevelLinks := extractLinks(node, "specialButton", "dark-title")
			for _, link := range firstLevelLinks {
				// Fetch second-level link
				secondLevelLink, err := scrapeSecondLevelLinks(link)
				if err != nil {
					fmt.Printf("Failed to fetch second-level link for %s: %v\n", link, err)
					continue
				}

				if secondLevelLink != "" {
					// Resolve short link to actual URL
					resolvedLink, err := resolveShortLink(secondLevelLink)
					if err != nil {
						fmt.Printf("Failed to resolve short link %s: %v\n", secondLevelLink, err)
						continue
					}
					linksChan <- resolvedLink
				}
			}
		}(url)
	}

	go func() {
		wg.Wait()
		close(linksChan)
	}()

	for link := range linksChan {
		finalLinks = append(finalLinks, link)
	}

	// Write links to a .txt file
	file, err := os.Create("links.txt")
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, link := range finalLinks {
		_, err := writer.WriteString(link + "\n")
		if err != nil {
			fmt.Printf("Failed to write to file: %v\n", err)
			return
		}
	}
	writer.Flush()

	fmt.Println("Links have been written to links.txt")
}
