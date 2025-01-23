package main

import (
	"database/sql"
	"fmt"
	"time"
)

func getScanResults(scanID int64) ([]map[string]interface{}, error) {
	rows, err := db.Query(`
        SELECT r.url, r.csp, r.xframe, c.key, c.value
        FROM results r
        LEFT JOIN cred c ON r.id = c.res
        WHERE r.scan = $1
    `, scanID)
	if err != nil {
		return nil, fmt.Errorf("error querying scan results: %v", err)
	}
	defer rows.Close()

	resultsMap := make(map[string]map[string]interface{})
	for rows.Next() {
		var url string
		var csp, xframe bool
		var key, value sql.NullString

		if err := rows.Scan(&url, &csp, &xframe, &key, &value); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		if _, exists := resultsMap[url]; !exists {
			resultsMap[url] = map[string]interface{}{
				"domain":      url,
				"csp":         csp,
				"xframe":      xframe,
				"credentials": make(map[string]string),
			}
		}

		if key.Valid && value.Valid {
			resultsMap[url]["credentials"].(map[string]string)[key.String] = value.String
		}
	}

	var results []map[string]interface{}
	for _, result := range resultsMap {
		results = append(results, result)
	}

	return results, nil
}

func getAllScans() ([]map[string]interface{}, error) {
	rows, err := db.Query(`
        SELECT s.id, s.start, s.end, COUNT(r.id) as domain_count
        FROM scans s
        LEFT JOIN results r ON s.id = r.scan
        GROUP BY s.id
        ORDER BY s.start DESC
    `)
	if err != nil {
		return nil, fmt.Errorf("error querying scans: %v", err)
	}
	defer rows.Close()

	var scans []map[string]interface{}
	for rows.Next() {
		var id int64
		var startTime, endTime time.Time
		var domainCount int

		if err := rows.Scan(&id, &startTime, &endTime, &domainCount); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		duration := endTime.Sub(startTime).String()
		scans = append(scans, map[string]interface{}{
			"id":           id,
			"start_time":   startTime,
			"duration":     duration,
			"domain_count": domainCount,
		})
	}

	return scans, nil
}
