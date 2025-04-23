package services

import (
	"fmt"
	"os"
	"trader/utils"
)

func FetchKlines() {
	// --- Configuration ---
	symbol := "BTCUSDT"
	interval := "1h" // 1 hour interval
	// Define the date range (inclusive)
	// Example: Fetch data from 2017-01-01 to 2017-12-31
	startDateStr := "2025-04-01"
	endDateStr := "2025-04-23"
	outputDir := "./data/btc"
	outputFileName := "202504.json"

	// Call the function to fetch and save data
	err := utils.FetchAndSaveKlines(symbol, interval, startDateStr, endDateStr, outputDir, outputFileName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1) // Exit with a non-zero code to indicate failure
	}

}
