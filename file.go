package main

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// ListFiles function: Traverse all files in the specified directory and return a slice of the file path
func ListFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// Get Content function: Read the content of a specified file and return it
func GetFileContent(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get File Size
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	// Read File Content
	fileContent := make([]byte, fileSize)
	_, err = file.Read(fileContent)
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}

// fileHandler function: Handle file requests
func fileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "MeowMusicEmbeddedServer")
	// Obtain the path of the request
	filePath := r.URL.Path

	// Check if the request path starts with "/url/"
	if strings.HasPrefix(filePath, "/url/") {
		// Extract the URL after "/url/"
		urlPath := filePath[len("/url/"):]
		// Decode the URL path in case it's URL encoded
		decodedURL, err := url.QueryUnescape(urlPath)
		if err != nil {
			NotFoundHandler(w, r)
			return
		}
		// Determine the protocol based on the URL path
		var protocol string
		if strings.HasPrefix(decodedURL, "http/") {
			protocol = "http://"
		} else if strings.HasPrefix(decodedURL, "https/") {
			protocol = "https://"
		} else {
			NotFoundHandler(w, r)
			return
		}
		// Remove the protocol part from the decoded URL
		decodedURL = strings.TrimPrefix(decodedURL, "http/")
		decodedURL = strings.TrimPrefix(decodedURL, "https/")
		// Correctly concatenate the protocol with the decoded URL
		decodedURL = protocol + decodedURL
		// Create a new HTTP request to the decoded URL, without copying headers
		req, err := http.NewRequest("GET", decodedURL, nil)
		if err != nil {
			NotFoundHandler(w, r)
			return
		}
		// Send the request and get the response
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			NotFoundHandler(w, r)
			return
		}
		defer resp.Body.Close()
		// Read the response body into a byte slice
		fileContent, err := io.ReadAll(resp.Body)
		if err != nil {
			NotFoundHandler(w, r)
			return
		}
		// Set appropriate Content-Type based on file extension
		ext := filepath.Ext(decodedURL)
		switch ext {
		case ".mp3":
			w.Header().Set("Content-Type", "audio/mpeg")
		case ".wav":
			w.Header().Set("Content-Type", "audio/wav")
		case ".flac":
			w.Header().Set("Content-Type", "audio/flac")
		case ".aac":
			w.Header().Set("Content-Type", "audio/aac")
		case ".ogg":
			w.Header().Set("Content-Type", "audio/ogg")
		case ".m4a":
			w.Header().Set("Content-Type", "audio/mp4")
		case ".amr":
			w.Header().Set("Content-Type", "audio/amr")
		case ".jpg", ".jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".gif":
			w.Header().Set("Content-Type", "image/gif")
		case ".bmp":
			w.Header().Set("Content-Type", "image/bmp")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		case ".webp":
			w.Header().Set("Content-Type", "image/webp")
		case ".txt":
			w.Header().Set("Content-Type", "text/plain")
		case ".lrc":
			w.Header().Set("Content-Type", "text/plain")
		case ".mrc":
			w.Header().Set("Content-Type", "text/plain")
		case ".json":
			w.Header().Set("Content-Type", "application/json")
		default:
			w.Header().Set("Content-Type", "application/octet-stream")
		}
		// Write file content to response
		w.Write(fileContent)
		return
	}

	// Construct the complete file path
	fullFilePath := filepath.Join("./files", filePath)

	// Get file content
	fileContent, err := GetFileContent(fullFilePath)
	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	// Set appropriate Content-Type based on file extension
	ext := filepath.Ext(filePath)
	switch ext {
	case ".mp3":
		w.Header().Set("Content-Type", "audio/mpeg")
	case ".wav":
		w.Header().Set("Content-Type", "audio/wav")
	case ".flac":
		w.Header().Set("Content-Type", "audio/flac")
	case ".aac":
		w.Header().Set("Content-Type", "audio/aac")
	case ".ogg":
		w.Header().Set("Content-Type", "audio/ogg")
	case ".m4a":
		w.Header().Set("Content-Type", "audio/mp4")
	case ".amr":
		w.Header().Set("Content-Type", "audio/amr")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".bmp":
		w.Header().Set("Content-Type", "image/bmp")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case ".webp":
		w.Header().Set("Content-Type", "image/webp")
	case ".txt":
		w.Header().Set("Content-Type", "text/plain")
	case ".lrc":
		w.Header().Set("Content-Type", "text/plain")
	case ".mrc":
		w.Header().Set("Content-Type", "text/plain")
	case ".json":
		w.Header().Set("Content-Type", "application/json")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	// Write file content to response
	w.Write(fileContent)
}
