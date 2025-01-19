package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	initDB()
	defer db.Close()

	startTime := time.Now()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the name of the txt file (without extension): ")
	filename, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading filename: %v\n", err)
		return
	}
	filename = strings.TrimSpace(filename) // Remove any trailing newline or whitespace

	// Construct the full path to the file
	filePath := "domains/" + filename + ".txt"

	// Read links from the file
	links, err := readLinksFromFile(filePath)
	if err != nil {
		fmt.Printf("Error reading links: %v\n", err)
		return
	}

	// Insert scan record
	scanID, err := insertScan(startTime, time.Time{})
	if err != nil {
		fmt.Printf("Error inserting scan: %v\n", err)
		return
	}

	// Split links into batches of 500
	batchSize := 500
	for i := 0; i < len(links); i += batchSize {
		end := i + batchSize
		if end > len(links) {
			end = len(links)
		}
		batch := links[i:end]

		var wg sync.WaitGroup
		resultsChan := make(chan SecurityConfig, len(batch))
		sem := make(chan struct{}, batchSize)

		for _, link := range batch {
			wg.Add(1)
			sem <- struct{}{}
			go analyzeLink(link, resultsChan, &wg, sem)
		}

		// Close resultsChan when all links are processed
		go func() {
			wg.Wait()
			close(resultsChan)
		}()

		var results []SecurityConfig
		for result := range resultsChan {
			results = append(results, result)
		}

		err = bulkInsertResultsAndCreds(scanID, results)
		if err != nil {
			fmt.Printf("Error performing bulk insert: %v\n", err)
		}

		err = updateScanEndTime(scanID)
		if err != nil {
			fmt.Printf("Error updating scan end time: %v\n", err)
		}
	}

	fmt.Println("Results have been written to the database")
	fmt.Println("Current time is:", time.Now())
}
