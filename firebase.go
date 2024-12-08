package main

import (
	"regexp"
	"strings"
)

// containsFirebaseConfig checks if the given script contains Firebase configuration keys.
func containsFirebaseConfig(script string) bool {
	// Key indicators for Firebase configuration
	keywords := []string{
		"apiKey",
		"authDomain",
		"projectId",
	}

	// Check if all keywords exist in the script
	for _, keyword := range keywords {
		if !strings.Contains(script, keyword) {
			return false
		}
	}

	return true
}

// extractFirebaseKeys extracts Firebase keys from a script using regular expressions.
func extractFirebaseKeys(script string) map[string]string {
	// Regular expressions for Firebase keys
	keyPatternsJSON := map[string]*regexp.Regexp{
		"apiKey":      regexp.MustCompile(`\\"apiKey\\"\s*:\s*\\"(.*?)\\"`),
		"authDomain":  regexp.MustCompile(`\\"authDomain\\"\s*:\s*\\"(.*?)\\"`),
		"databaseURL": regexp.MustCompile(`\\"databaseURL\\"\s*:\s*\\"(.*?)\\"`),
		"projectId":   regexp.MustCompile(`\\"projectId\\"\s*:\s*\\"(.*?)\\"`),
	}

	keyPatternsJSON2 := map[string]*regexp.Regexp{
		"apiKey":      regexp.MustCompile(`\"apiKey\"\s*:\s*\"(.*?)\"`),
		"authDomain":  regexp.MustCompile(`\"authDomain\"\s*:\s*\"(.*?)\"`),
		"databaseURL": regexp.MustCompile(`\"databaseURL\"\s*:\s*\"(.*?)\"`),
		"projectId":   regexp.MustCompile(`\"projectId\"\s*:\s*\"(.*?)\"`),
	}

	// Regular expressions for Firebase keys in JavaScript object literal format
	keyPatternsJS := map[string]*regexp.Regexp{
		"apiKey":      regexp.MustCompile(`apiKey\s*:\s*\"(.*?)\"`),
		"authDomain":  regexp.MustCompile(`authDomain\s*:\s*\"(.*?)\"`),
		"databaseURL": regexp.MustCompile(`databaseURL\s*:\s*\"(.*?)\"`),
		"projectId":   regexp.MustCompile(`projectId\s*:\s*\"(.*?)\"`),
	}

	results := make(map[string]string)

	// Extract matches for each key
	for key, pattern := range keyPatternsJSON {
		if match := pattern.FindStringSubmatch(script); len(match) > 1 {
			results[key] = match[1]
		}
	}

	for key, pattern := range keyPatternsJSON2 {
		if match := pattern.FindStringSubmatch(script); len(match) > 1 {
			results[key] = match[1]
		}
	}

	for key, pattern := range keyPatternsJS {
		if match := pattern.FindStringSubmatch(script); len(match) > 1 {
			results[key] = match[1]
		}
	}

	return results
}
