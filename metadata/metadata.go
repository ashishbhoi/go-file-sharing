package metadata

import (
	"encoding/json"
	"os"
	"sync"
)

// Global variables
const (
	MetadataFile = "./uploads/metadata.json"
)

var (
	FileMetadata  = make(map[string]string) // Maps file IDs to original filenames
	MetadataMutex = &sync.Mutex{}           // Mutex to protect access to FileMetadata
)

// LoadMetadata loads the file metadata from disk
func LoadMetadata() error {
	MetadataMutex.Lock()
	defer MetadataMutex.Unlock()

	// Check if metadata file exists
	if _, err := os.Stat(MetadataFile); os.IsNotExist(err) {
		// If it doesn't exist, initialize with empty metadata
		FileMetadata = make(map[string]string)
		return nil
	}

	// Read the metadata file
	data, err := os.ReadFile(MetadataFile)
	if err != nil {
		return err
	}

	// Parse the JSON data
	return json.Unmarshal(data, &FileMetadata)
}

// SaveMetadata saves the file metadata to disk
func SaveMetadata() error {
	MetadataMutex.Lock()
	defer MetadataMutex.Unlock()

	// Marshal the metadata to JSON
	data, err := json.MarshalIndent(FileMetadata, "", "  ")
	if err != nil {
		return err
	}

	// Write the metadata to file
	return os.WriteFile(MetadataFile, data, 0644)
}
