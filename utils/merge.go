package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os" // Added for file operations
)

// MergeJSONObjects merges multiple JSON objects (represented as map[string]interface{}) into one.
// It handles nested objects recursively. If keys conflict:
// - If both values are maps, they are recursively merged.
// - Otherwise, the value from the later object overwrites the earlier one.
func MergeJSONObjects(objects ...map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for _, obj := range objects {
		if obj == nil {
			continue // Skip nil maps
		}
		for key, value := range obj {
			existingValue, keyExists := result[key]

			if keyExists {
				existingMap, existingIsMap := existingValue.(map[string]interface{})
				newMap, newIsMap := value.(map[string]interface{})

				if existingIsMap && newIsMap {
					// Recursively merge if both are maps
					mergedMap, err := MergeJSONObjects(existingMap, newMap)
					if err != nil {
						return nil, fmt.Errorf("error merging nested object for key '%s': %v", key, err)
					}
					result[key] = mergedMap
				} else {
					// Overwrite if types don't match or aren't both maps
					result[key] = value
				}
			} else {
				// Key doesn't exist, just add it
				result[key] = value
			}
		}
	}

	return result, nil
}

// Helper function to merge JSON data from byte slices
func MergeJSONData(jsonData ...[]byte) ([]byte, error) {
	var objects []map[string]interface{}

	for i, data := range jsonData {
		var obj map[string]interface{}
		if len(data) == 0 || string(data) == "null" { // Handle empty or null JSON
			continue
		}
		err := json.Unmarshal(data, &obj)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON data at index %d: %v", i, err)
		}
		objects = append(objects, obj)
	}

	mergedObject, err := MergeJSONObjects(objects...)
	if err != nil {
		return nil, fmt.Errorf("error merging JSON objects: %v", err)
	}

	mergedJSON, err := json.MarshalIndent(mergedObject, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling merged JSON: %v", err)
	}

	return mergedJSON, nil
}

// MergeJSONFiles reads JSON data from multiple files, merges them, and writes the result to an output file.
func MergeJSONFiles(inputFilePaths []string, outputFilePath string) error {
	var jsonDataSlices [][]byte

	for _, inputPath := range inputFilePaths {
		data, err := os.ReadFile(inputPath)
		if err != nil {
			// If the file doesn't exist, we might want to skip it or return an error.
			// Here, we return an error to be explicit.
			if os.IsNotExist(err) {
				return fmt.Errorf("input file not found: %s", inputPath)
			}
			return fmt.Errorf("error reading input file '%s': %v", inputPath, err)
		}
		log.Println("data", data)
		jsonDataSlices = append(jsonDataSlices, data)
	}

	// Merge the data from all files
	mergedData, err := MergeJSONData(jsonDataSlices...)
	if err != nil {
		return fmt.Errorf("error merging JSON data from files: %v", err)
	}

	// Write the merged data to the output file
	// Ensure the output directory exists (optional, depending on requirements)
	// outputDir := filepath.Dir(outputFilePath)
	// if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
	// 	 return fmt.Errorf("error creating output directory for '%s': %v", outputFilePath, err)
	// }

	err = os.WriteFile(outputFilePath, mergedData, 0644)
	if err != nil {
		return fmt.Errorf("error writing merged JSON to output file '%s': %v", outputFilePath, err)
	}

	fmt.Printf("Successfully merged JSON files into %s\n", outputFilePath)
	return nil
}
