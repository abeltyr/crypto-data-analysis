package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Binance API base URL
const apiBaseURL = "https://api.binance.com/api/v3/klines"

// fetchAndSaveKlines fetches K-line data from Binance API for a given symbol and date range,
// then saves it to a JSON file.
func FetchAndSaveKlines(symbol, interval, startDateStr, endDateStr, outputDir, outputFileName string) error {
	// Parse start and end dates
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return fmt.Errorf("error parsing start date '%s': %v", startDateStr, err)
	}
	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return fmt.Errorf("error parsing end date '%s': %v", endDateStr, err)
	}

	// Ensure output directory exists
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating output directory '%s': %v", outputDir, err)
	}

	// Ensure end date is after start date
	if endDate.Before(startDate) {
		return fmt.Errorf("end date must be after start date")
	}

	// Use a map to store klines, keyed by date string "YYYY-MM-DD"
	// The value is a slice of *Kline pointers, size 24, indexed by hour (0-23)
	allKlinesByDate := make(map[string][]*Kline)

	// Loop through each day in the range
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		startTimeMs := d.UnixMilli()
		// Binance API endTime is exclusive, so we need the start of the next day
		endTimeMs := d.AddDate(0, 0, 1).UnixMilli()
		currentDateStr := d.Format("2006-01-02") // Store the current date string for the map key

		fmt.Printf("Fetching data for %s...\n", currentDateStr)

		// Initialize the slice for the current date with 24 nil pointers
		allKlinesByDate[currentDateStr] = make([]*Kline, 24)

		// Construct API URL for the day
		apiURL := fmt.Sprintf("%s?symbol=%s&interval=%s&startTime=%d&endTime=%d&limit=1000",
			apiBaseURL, symbol, interval, startTimeMs, endTimeMs)

		// Fetch data from Binance API
		resp, err := http.Get(apiURL)
		if err != nil {
			fmt.Printf("Error fetching data for %s: %v\n", currentDateStr, err)
			continue // Skip to the next day on error
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Printf("API error for %s: Status %d, Body: %s\n", currentDateStr, resp.StatusCode, string(bodyBytes))
			continue // Skip to the next day on API error
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body for %s: %v\n", currentDateStr, err)
			continue
		}

		// Parse JSON response (array of arrays)
		var rawKlines [][]interface{}
		err = json.Unmarshal(body, &rawKlines)
		if err != nil {
			// Handle potential API error message format (JSON object instead of array)
			var apiError map[string]interface{}
			if json.Unmarshal(body, &apiError) == nil {
				fmt.Printf("API returned error object for %s: %v\n", currentDateStr, apiError)
			} else {
				fmt.Printf("Error parsing JSON for %s: %v\n", currentDateStr, err)
			}
			continue
		}
		// Convert raw data to Kline struct
		for _, rawKline := range rawKlines {
			if len(rawKline) < 7 {
				fmt.Printf("Skipping incomplete kline data: %v\n", rawKline)
				continue
			}

			// Safely convert interface{} to expected types
			openTime, ok1 := rawKline[0].(float64)
			open, ok2 := rawKline[1].(string)
			high, ok3 := rawKline[2].(string)
			low, ok4 := rawKline[3].(string)
			closePrice, ok5 := rawKline[4].(string)
			closeTime, ok6 := rawKline[6].(float64)

			if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
				fmt.Printf("Skipping kline due to type assertion error: %v\n", rawKline)
				continue
			}

			// Calculate date and hour from open time
			ts := time.UnixMilli(int64(openTime))
			klineDateStr := ts.Format("2006-01-02")
			klineHour := ts.Hour() // Get hour (0-23)

			kline := &Kline{ // Create a pointer to Kline
				OpenTime:  int64(openTime),
				Open:      open,
				High:      high,
				Low:       low,
				Close:     closePrice,
				CloseTime: int64(closeTime),
				Date:      klineDateStr,
				Hour:      klineHour,
			}
			// Place the kline pointer at the index corresponding to its hour
			allKlinesByDate[currentDateStr][klineHour] = kline

		}

		// Optional: Add a small delay to avoid hitting API rate limits
		time.Sleep(200 * time.Millisecond)
	}

	// Prepare JSON output from the map
	jsonData, err := json.MarshalIndent(allKlinesByDate, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON data: %v", err)
	}

	// Write JSON data to file
	outputPath := filepath.Join(outputDir, outputFileName)
	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON to file '%s': %v", outputPath, err)
	}

	fmt.Printf("Successfully fetched data and saved to %s\n", outputPath)
	return nil
}
