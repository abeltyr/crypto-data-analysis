package main

// import (
// 	"fmt"
// 	"os"
// 	"trader/utils"
// )

// func main() {
// 	// // --- Configuration ---
// 	// symbol := "BTCUSDT"
// 	// interval := "1h" // 1 hour interval
// 	// // Define the date range (inclusive)
// 	// // Example: Fetch data from 2017-01-01 to 2017-12-31
// 	// startDateStr := "2025-04-01"
// 	// endDateStr := "2025-04-23"
// 	// outputDir := "./data/btc"
// 	// outputFileName := "202504.json"

// 	// // Call the function to fetch and save data
// 	// err := utils.FetchAndSaveKlines(symbol, interval, startDateStr, endDateStr, outputDir, outputFileName)
// 	// if err != nil {
// 	// 	fmt.Printf("Error: %v\n", err)
// 	// 	os.Exit(1) // Exit with a non-zero code to indicate failure
// 	// }

// 	// --- Merge JSON Files ---
// 	// inputDir := "./data/btc"
// 	// outputFilePath := "./data/btc/merged_eth_data.json"

// 	// // Read the input directory
// 	// files, err := os.ReadDir(inputDir)
// 	// if err != nil {
// 	// 	fmt.Printf("Error reading input directory '%s': %v\n", inputDir, err)
// 	// 	os.Exit(1)
// 	// }

// 	// var inputFilePaths []string
// 	// for _, file := range files {
// 	// 	if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
// 	// 		// Construct the full path for each JSON file
// 	// 		fullPath := filepath.Join(inputDir, file.Name())
// 	// 		inputFilePaths = append(inputFilePaths, fullPath)
// 	// 	}
// 	// }

// 	// // Check if any JSON files were found
// 	// if len(inputFilePaths) == 0 {
// 	// 	fmt.Printf("No JSON files found in directory '%s'\n", inputDir)
// 	// 	os.Exit(0) // Exit gracefully if no files to merge
// 	// }

// 	// // Call the merge function
// 	// err = utils.MergeJSONFiles(inputFilePaths, outputFilePath)
// 	// if err != nil {
// 	// 	fmt.Printf("Error merging JSON files: %v\n", err)
// 	// 	os.Exit(1)
// 	// }

// 	// fmt.Println("Successfully executed merge process.")

// 	// --- Simulate DCA ---
// 	fmt.Println("\n--- Starting DCA Simulation ---")
// 	jsonInputPath := "./data/btc/merged_eth_data.json" // Use the merged data
// 	// "01", "08" "22" "31"
// 	simStartDate := "2023-04-01"
// 	simEndDate := "2023-04-08"
// 	investmentAmount := 50.0 // Invest $100 each interval
// 	dailyReportOutputPath := "./reports/weekly/2023/04/dca_report_week1.txt"

// 	err := utils.SimulateDCA(jsonInputPath, simStartDate, simEndDate, utils.Daily, investmentAmount, dailyReportOutputPath)
// 	if err != nil {
// 		fmt.Printf("Error running DCA simulation: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// weeklyReportOutputPath := "./reports/Weekly/2025/100/dca_simulation_report_Weekly.txt"
// 	// err = utils.SimulateDCA(jsonInputPath, simStartDate, simEndDate, utils.Weekly, investmentAmount, weeklyReportOutputPath)
// 	// if err != nil {
// 	// 	fmt.Printf("Error running DCA simulation: %v\n", err)
// 	// 	os.Exit(1)
// 	// }

// 	// monthlyReportOutputPath := "./reports/Weekly/2025/100/dca_simulation_report_Monthly.txt"
// 	// err = utils.SimulateDCA(jsonInputPath, simStartDate, simEndDate, utils.Monthly, investmentAmount, monthlyReportOutputPath)
// 	// if err != nil {
// 	// 	fmt.Printf("Error running DCA simulation: %v\n", err)
// 	// 	os.Exit(1)
// 	// }

// 	fmt.Println("\n--- All Processes Completed Successfully ---")
// }
