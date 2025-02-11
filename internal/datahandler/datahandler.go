package datahandler

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nielsjaspers/cli-sky/bluesky"
)

// WriteAuthResponseToFile takes a BlueskyAuthResponse struct and writes it to a JSON file
// in the "api_responses" directory, with the filename containing the user handle.
func WriteAuthResponseToFile(authResponse *bluesky.BlueskyAuthResponse) error {
	dirPath := "api_responses"

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}
		log.Println("Created directory:", dirPath)
	}

	fileName := fmt.Sprintf("auth_response_%s.json", authResponse.Handle)
	filePath := filepath.Join(dirPath, fileName)

	jsonData, err := json.MarshalIndent(authResponse, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	log.Printf("Successfully wrote auth response to %s\n", filePath)
	return nil
}

// ReadAuthResponseFromFile reads a BlueskyAuthResponse from a JSON file.
// If handle is provided, it tries to read the file with that handle in the filename.
// If handle is empty, it tries to find a file in the api_responses directory and reads the first one it finds.
func ReadAuthResponseFromFile(handle string) (*bluesky.BlueskyAuthResponse, error) {
	dirPath := "api_responses"
	var filePath string

	if handle != "" {
		// Construct the filename with the provided handle.
		fileName := fmt.Sprintf("auth_response_%s.json", handle)
		filePath = filepath.Join(dirPath, fileName)

		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found for handle %s: %w", handle, err)
		}
	} else {
		// If no handle is provided, search for a file in the directory.
		files, err := filepath.Glob(filepath.Join(dirPath, "auth_response_*.json"))
		if err != nil {
			return nil, fmt.Errorf("error listing files in directory: %w", err)
		}

		if len(files) == 0 {
			return nil, fmt.Errorf("no auth_response files found in directory %s", dirPath)
		}

		filePath = files[0] // Use the first file found.
		log.Printf("No handle provided. Reading from the first file found: %s\n", filePath)
	}

	// Read the file.
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal the JSON data into a BlueskyAuthResponse struct.
	var authResponse bluesky.BlueskyAuthResponse
	err = json.Unmarshal(jsonData, &authResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	log.Printf("Successfully read auth response from %s\n", filePath)
	return &authResponse, nil
}
