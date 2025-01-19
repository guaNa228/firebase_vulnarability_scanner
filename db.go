package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	connStr := "user=postgres password=f1r2o3l4o5v dbname=vulnarability sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v\n", err)
	}
}

func insertScan(start, end time.Time) (int64, error) {
	var scanID int64
	err := db.QueryRow(`INSERT INTO scans (start, "end") VALUES ($1, $2) RETURNING id`, start, end).Scan(&scanID)
	if err != nil {
		return 0, fmt.Errorf("error inserting scan: %v", err)
	}
	return scanID, nil
}

func bulkInsertResultsAndCreds(scanID int64, results []SecurityConfig) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}

	resultStmt, err := tx.Prepare(`INSERT INTO results (scan, url, csp, xframe) VALUES ($1, $2, $3, $4) RETURNING id`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error preparing result statement: %v", err)
	}
	defer resultStmt.Close()

	credStmt, err := tx.Prepare(`INSERT INTO cred (res, key, value) VALUES ($1, $2, $3)`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error preparing cred statement: %v", err)
	}
	defer credStmt.Close()

	for _, result := range results {
		var resultID int64
		err := resultStmt.QueryRow(scanID, result.URL, result.CSPHeader, result.XFrameHeader).Scan(&resultID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error inserting result: %v", err)
		}

		for key, value := range result.Creds {
			_, err := credStmt.Exec(resultID, key, value)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error inserting credential: %v", err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func updateScanEndTime(scanID int64) error {
	_, err := db.Exec(`UPDATE scans SET "end" = $1 WHERE id = $2`, time.Now(), scanID)
	if err != nil {
		return fmt.Errorf("error updating scan end time: %v", err)
	}
	return nil
}
