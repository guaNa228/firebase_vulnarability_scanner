package main

import (
	"fmt"
	"sync"
)

// func main() {
// 	scraper.ScrapeLinks()
// }

// Main function to orchestrate the scraping process
func main() {
	// Read links from links.txt
	links, err := readLinksFromFile("links.txt")
	if err != nil {
		fmt.Printf("Error reading links: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	resultsChan := make(chan string, len(links))
	sem := make(chan struct{}, 10) // Semaphore to limit concurrent connections to 10

	for _, link := range links {
		wg.Add(1)
		sem <- struct{}{} // Acquire a semaphore slot
		go analyzeLink(link, resultsChan, &wg, sem)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Write results to a .txt file
	err = writeResultsToFile("results.txt", resultsChan)
	if err != nil {
		fmt.Printf("Error writing results: %v\n", err)
	}

	fmt.Println("Results have been written to results.txt")
}