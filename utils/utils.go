package utils

import (
	"file-sharing/metadata"
	"file-sharing/models"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Global variables
const (
	UploadDir = "./uploads"
)

// SanitizeFilename removes potentially dangerous characters from filenames
func SanitizeFilename(filename string) string {
	// Remove path components
	filename = filepath.Base(filename)

	// Replace problematic characters
	filename = strings.ReplaceAll(filename, "..", "")

	return filename
}

// GetFilesList returns a list of all files in the upload directory
func GetFilesList() ([]models.FileInfo, error) {
	files := []models.FileInfo{}

	// Read the upload directory
	entries, err := os.ReadDir(UploadDir)
	if err != nil {
		return nil, err
	}

	// Process each file
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Get file info
		fileID := entry.Name()
		fileInfo, err := GetFileInfo(fileID)
		if err != nil {
			log.Printf("Failed to get info for file %s: %v", fileID, err)
			continue
		}

		files = append(files, fileInfo)
	}

	return files, nil
}

// GetFileInfo returns information about a specific file
func GetFileInfo(fileID string) (models.FileInfo, error) {
	// Check for path traversal attempts
	if strings.Contains(fileID, "..") || strings.Contains(fileID, "/") || strings.Contains(fileID, "\\") {
		return models.FileInfo{}, fmt.Errorf("invalid file ID")
	}

	// Get file path
	filePath := filepath.Join(UploadDir, fileID)

	// Get file info
	info, err := os.Stat(filePath)
	if err != nil {
		return models.FileInfo{}, err
	}

	// Get the original filename from metadata
	metadata.MetadataMutex.Lock()
	name, exists := metadata.FileMetadata[fileID]
	metadata.MetadataMutex.Unlock()

	// If the filename doesn't exist in metadata, use the file ID as fallback
	if !exists {
		name = fileID
	}

	return models.FileInfo{
		ID:         fileID,
		Name:       name,
		Size:       info.Size(),
		UploadedAt: info.ModTime(),
	}, nil
}

// DeleteFile deletes a file from the upload directory
func DeleteFile(fileID string) error {
	// Check for path traversal attempts
	if strings.Contains(fileID, "..") || strings.Contains(fileID, "/") || strings.Contains(fileID, "\\") {
		return fmt.Errorf("invalid file ID")
	}

	// Get file path
	filePath := filepath.Join(UploadDir, fileID)

	// Delete the file
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	// Remove the file from metadata
	metadata.MetadataMutex.Lock()
	delete(metadata.FileMetadata, fileID)
	metadata.MetadataMutex.Unlock()

	// Save updated metadata
	if err := metadata.SaveMetadata(); err != nil {
		log.Printf("Warning: Failed to save metadata after deletion: %v", err)
		// Continue anyway since the file is already deleted
	}

	return nil
}
