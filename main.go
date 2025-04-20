package main

import (
	"file-sharing/handlers"
	"file-sharing/metadata"
	"file-sharing/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Global variables
const (
	port = 8080
)

func main() {
	// Set timezone to Asia/Kolkata
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Printf("Warning: Failed to load Asia/Kolkata timezone: %v", err)
	} else {
		time.Local = loc
		log.Printf("Timezone set to Asia/Kolkata")
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(utils.UploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Load file metadata
	if err := metadata.LoadMetadata(); err != nil {
		log.Printf("Warning: Failed to load metadata: %v", err)
		// Continue anyway with empty metadata
	}

	// Create templates directory if it doesn't exist
	if err := os.MkdirAll("templates", 0755); err != nil {
		log.Fatalf("Failed to create templates directory: %v", err)
	}

	// Create static directory if it doesn't exist
	if err := os.MkdirAll("static", 0755); err != nil {
		log.Fatalf("Failed to create static directory: %v", err)
	}

	// Create CSS directory if it doesn't exist
	if err := os.MkdirAll("static/css", 0755); err != nil {
		log.Fatalf("Failed to create CSS directory: %v", err)
	}

	// Create JS directory if it doesn't exist
	if err := os.MkdirAll("static/js", 0755); err != nil {
		log.Fatalf("Failed to create JS directory: %v", err)
	}

	// Set up routes
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/upload", handlers.UploadHandler)
	http.HandleFunc("/files", handlers.ListFilesHandler)
	http.HandleFunc("/download/", handlers.DownloadHandler)
	http.HandleFunc("/view/", handlers.ViewHandler)
	http.HandleFunc("/delete", handlers.DeleteHandler)
	http.HandleFunc("/delete-multiple", handlers.DeleteMultipleHandler)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the server
	log.Printf("Server starting on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
