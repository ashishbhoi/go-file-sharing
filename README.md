# Go File Sharing Web Application

[![Go Version](https://img.shields.io/github/go-mod/go-version/ashishbhoi/go-file-sharing)](https://golang.org)
[![License](https://img.shields.io/github/license/ashishbhoi/go-file-sharing)](https://raw.githubusercontent.com/ashishbhoi/go-file-sharing/refs/heads/master/LICENSE)
[![Docker](https://img.shields.io/docker/v/ashishbhoi/go-file-sharing)](https://hub.docker.com/repository/docker/ashishbhoi/go-file-sharing/general)

A simple, lightweight, and modern web-based file sharing application built with Go. This application allows users to easily upload, view, download, and delete files through an intuitive web interface without requiring user authentication. Perfect for temporary file sharing in local networks or controlled environments.

![Go File Sharing App (Dark Mode)](https://raw.githubusercontent.com/ashishbhoi/go-file-sharing/refs/heads/master/screenshot/dark.png)
![Go File Sharing App (Light Mode)](https://raw.githubusercontent.com/ashishbhoi/go-file-sharing/refs/heads/master/screenshot/light.png)

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Installation](#installation)
  - [Docker Setup](#docker-setup)
  - [Manual Setup](#manual-setup)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Project Structure](#project-structure)
- [License](#license)
- [Contact](#contact)

## Features

- **File Upload**: Upload files via button or drag-and-drop interface
- **File Management**: 
  - List files with details (name, size, upload date)
  - Download files
  - View files directly in browser
  - Delete individual files
  - Delete multiple files at once
- **User Interface**:
  - Clean, responsive design
  - Dark mode support
  - Visual feedback for actions
- **Performance**:
  - Fast file uploads with progress indication
  - Efficient file handling
- **Security**:
  - Filename sanitization
  - Path traversal prevention

## Technologies Used

- **Backend**:
  - [Go](https://golang.org/) - Core programming language
  - Standard Go libraries for HTTP server, file handling, and JSON processing
- **Frontend**:
  - HTML5, CSS3, JavaScript (ES6+)
  - [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS framework
- **Containerization**:
  - [Docker](https://www.docker.com/) - For easy deployment and distribution
- **Timezone**:
  - Asia/Kolkata - Default timezone for file timestamps

## Installation

### Docker Setup

#### Building the Docker Image

To build the Docker image, run the following command from the project root:

```bash
docker build -t go-file-sharing .
```

#### Running the Docker Container

To run the Docker container, use the following command:

```bash
docker run -p 8080:8080 -v $(pwd)/uploads:/app/uploads go-file-sharing
```

This will:
- Map port 8080 from the container to port 8080 on your host
- Mount the `uploads` directory from your host to the container for persistent storage

For Windows PowerShell, use this command instead:

```powershell
docker run -p 8080:8080 -v ${PWD}/uploads:/app/uploads go-file-sharing
```

#### Docker Hub

You can also pull the image directly from Docker Hub:

```bash
docker pull ashishbhoi/go-file-sharing
docker run -p 8080:8080 -v $(pwd)/uploads:/app/uploads ashishbhoi/go-file-sharing
```

### Manual Setup

#### Prerequisites

- Go 1.24 or later

#### Installation Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/ashishbhoi/go-file-sharing.git
   cd go-file-sharing
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

3. Access the application at `http://localhost:8080`

## Usage

### Uploading Files

1. Click the "Select Files" button or drag and drop files onto the designated area
2. Wait for the upload to complete (progress bar will indicate status)
3. Uploaded files will appear in the files list

### Managing Files

- **Download**: Click on the file name or the "Download" button
- **View**: Click the "View" button to open the file in a new browser tab
- **Delete**: Click the "Delete" button next to a file
- **Delete Multiple**: Select multiple files using checkboxes and click "Delete Selected"

### Dark Mode

- Toggle between light and dark mode using the sun/moon icon in the top right corner

## API Endpoints

The application provides the following API endpoints:

- `GET /` - Serves the main page
- `POST /upload` - Handles file uploads
- `GET /files` - Returns a list of all files
- `GET /download/:id` - Downloads a file
- `GET /view/:id` - Views a file in the browser
- `POST /delete` - Deletes a single file
- `POST /delete-multiple` - Deletes multiple files

## Project Structure

- `main.go`: Entry point of the application
- `handlers/`: HTTP request handlers
- `metadata/`: File metadata management
- `models/`: Data models
- `utils/`: Utility functions
- `templates/`: HTML templates
- `static/`: Static assets (CSS, JavaScript)
- `uploads/`: Directory for uploaded files
- `Dockerfile`: Docker configuration

## License

This project is open source and available under the [MIT License](https://raw.githubusercontent.com/ashishbhoi/go-file-sharing/refs/heads/master/LICENSE).

## Contact


Project Link: [https://github.com/ashishbhoi/go-file-sharing](https://github.com/ashishbhoi/go-file-sharing)

---

Made with ❤️ using Go
