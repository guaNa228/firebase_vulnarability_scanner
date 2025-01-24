package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()

	// Use CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.POST("/scan", func(c *gin.Context) {
		var request struct {
			Links []string `json:"domains"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Clean the input domains
		cleanedLinks := cleanInput(request.Links)

		go performScan(cleanedLinks)

		c.JSON(http.StatusOK, gin.H{"status": "scan started"})
	})

	r.GET("/scan/:scan_id", func(c *gin.Context) {
		scanID, err := strconv.ParseInt(c.Param("scan_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scan_id"})
			return
		}

		results, err := getScanResults(scanID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	})

	r.GET("/scans", func(c *gin.Context) {
		scans, err := getAllScans()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, scans)
	})

	r.Run(":8080")
}

func performScan(links []string) {
	startTime := time.Now()

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

		// Update scan end time after each batch
		err = updateScanEndTime(scanID)
		if err != nil {
			fmt.Printf("Error updating scan end time: %v\n", err)
		}
	}

	fmt.Println("Results have been written to the database")
	fmt.Println("Current time is:", time.Now())
}
