package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"trader/utils"
)

func MergeKlines() {

	// --- Merge JSON Files ---
	inputDir := "./reports/weekly"
	outputFilePath := "./reports/weekly-merged/data.json"

	// Read the input directory
	files, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Printf("Error reading input directory '%s': %v\n", inputDir, err)
		os.Exit(1)
	}

	var inputFilePaths []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			// Construct the full path for each JSON file
			fullPath := filepath.Join(inputDir, file.Name())
			inputFilePaths = append(inputFilePaths, fullPath)
		}
	}

	// Check if any JSON files were found
	if len(inputFilePaths) == 0 {
		fmt.Printf("No JSON files found in directory '%s'\n", inputDir)
		os.Exit(0) // Exit gracefully if no files to merge
	}

	log.Println("inputFilePaths", len(inputFilePaths))
	// Call the merge function
	err = utils.MergeJSONFiles(inputFilePaths, outputFilePath)
	if err != nil {
		fmt.Printf("Error merging JSON files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully executed merge process.")
}
