package handlers

import (
	"encoding/json"
	"file-sharing/metadata"
	"file-sharing/models"
	"file-sharing/utils"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	Templates = template.Must(template.ParseGlob("templates/*.html"))
)

// IndexHandler serves the main page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if err := Templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UploadHandler handles file uploads
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form with a 10MB max memory
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get the files from the form
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	uploadedFiles := []models.FileInfo{}

	// Process each file
	for _, fileHeader := range files {
		// Open the uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to open uploaded file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Sanitize the filename
		filename := utils.SanitizeFilename(fileHeader.Filename)

		// Generate a unique ID for the file
		fileID := strconv.FormatInt(time.Now().UnixNano(), 10)

		// Create the file path
		filePath := filepath.Join(utils.UploadDir, fileID)

		// Create the destination file
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Failed to create destination file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		// Store the original filename in metadata
		metadata.MetadataMutex.Lock()
		metadata.FileMetadata[fileID] = filename
		metadata.MetadataMutex.Unlock()

		// Save metadata to disk
		if err := metadata.SaveMetadata(); err != nil {
			log.Printf("Warning: Failed to save metadata: %v", err)
			// Continue anyway
		}

		// Add file info to the list
		fileInfo := models.FileInfo{
			ID:         fileID,
			Name:       filename,
			Size:       fileHeader.Size,
			UploadedAt: time.Now(),
		}
		uploadedFiles = append(uploadedFiles, fileInfo)
	}

	// Return the uploaded files info as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(uploadedFiles); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// ListFilesHandler returns a list of all files
func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := utils.GetFilesList()
	if err != nil {
		http.Error(w, "Failed to get files list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DownloadHandler serves a file for download
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the file ID from the URL
	fileID := strings.TrimPrefix(r.URL.Path, "/download/")
	if fileID == "" {
		http.Error(w, "File ID is required", http.StatusBadRequest)
		return
	}

	// Get the file info
	fileInfo, err := utils.GetFileInfo(fileID)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Open the file
	filePath := filepath.Join(utils.UploadDir, fileID)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set the headers for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.Name))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size, 10))

	// Copy the file to the response
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
	}
}

// DeleteHandler deletes a single file
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the file ID from the request
	fileID := r.FormValue("id")
	if fileID == "" {
		http.Error(w, "File ID is required", http.StatusBadRequest)
		return
	}

	// Delete the file
	if err := utils.DeleteFile(fileID); err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File deleted successfully"))
}

// ViewHandler serves a file for viewing in the browser
func ViewHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the file ID from the URL
	fileID := strings.TrimPrefix(r.URL.Path, "/view/")
	if fileID == "" {
		http.Error(w, "File ID is required", http.StatusBadRequest)
		return
	}

	// Get the file info
	fileInfo, err := utils.GetFileInfo(fileID)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Open the file
	filePath := filepath.Join(utils.UploadDir, fileID)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Determine content type based on file extension
	contentType := getContentType(fileInfo.Name)

	// Set the headers for inline viewing
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", fileInfo.Name))
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size, 10))

	// Copy the file to the response
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
	}
}

// getContentType determines the content type based on file extension
func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".xml":
		return "application/xml"
	case ".mp4":
		return "video/mp4"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	default:
		return "application/octet-stream"
	}
}

// DeleteMultipleHandler deletes multiple files
func DeleteMultipleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form with a 10MB max memory
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get the file IDs from the request
	fileIDs := r.MultipartForm.Value["ids"]
	if len(fileIDs) == 0 {
		http.Error(w, "File IDs are required", http.StatusBadRequest)
		return
	}

	// Delete each file
	for _, fileID := range fileIDs {
		if err := utils.DeleteFile(fileID); err != nil {
			log.Printf("Failed to delete file %s: %v", fileID, err)
			// Continue with other files even if one fails
		}
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Files deleted successfully"))
}
