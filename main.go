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
	fmt.Println("Current time:", time.Now())

	// Clear the results.txt file
	file, err := os.OpenFile("results.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error clearing results file: %v\n", err)
		return
	}
	file.Close()

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

	// Split links into batches of 1000
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

		// Open the results file in append mode
		file, err := os.OpenFile("results.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening results file: %v\n", err)
			return
		}
		defer file.Close()

		writer := bufio.NewWriter(file)

		for result := range resultsChan {
			line := fmt.Sprintf("%s => CSPHeader: %v, XFrameHeader: %v\n", result.URL, result.CSPHeader, result.XFrameHeader)
			if result.URL != "" {
				if _, err := writer.WriteString(line); err != nil {
					fmt.Printf("Error writing to results file: %v\n", err)
					return
				}
			}
			for key, value := range result.Creds {
				credLine := fmt.Sprintf("%s : %s\n", key, value)
				if _, err := writer.WriteString(credLine); err != nil {
					fmt.Printf("Error writing to results file: %v\n", err)
					return
				}
			}
		}
		writer.Flush()
	}

	fmt.Println("Results have been written to results.txt")
	fmt.Println("Current time is:", time.Now())
}
