package http_index

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pluveto/flydav/internal/auth"
	"github.com/pluveto/flydav/internal/config"
	"github.com/pluveto/flydav/internal/logger"
	"github.com/pluveto/flydav/pkg/storage"
)

type HTTPIndexModule struct {
	Config  config.HTTPIndexConfig
	Storage storage.Storage
	Auth    *auth.AuthModule
}

func NewHTTPIndexModule(cfg config.HTTPIndexConfig, store storage.Storage, auth *auth.AuthModule) *HTTPIndexModule {
	return &HTTPIndexModule{
		Config:  cfg,
		Storage: store,
		Auth:    auth,
	}
}

func (his *HTTPIndexModule) RegisterRoutes(router *mux.Router) {
	logger.Info("Registering HTTP Index module routes on " + his.Config.Path)
	router.PathPrefix(his.Config.Path).HandlerFunc(his.handleHTTPIndex)
}

// Directory listing template.
const dirListTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Index of {{ .Path }}</title>
</head>
<body>
    <h1>Index of {{ .Path }}</h1>
    <ul>
        {{ range .Contents }}
        <li><a href="{{ .FullName }}">{{ .Name }}</a></li>
        {{ end }}
    </ul>
</body>
</html>
`

// TemplateData holds data for rendering the directory listing template.
type TemplateData struct {
	Path     string
	Contents []storage.Metadata
}

func (his *HTTPIndexModule) handleHTTPIndex(w http.ResponseWriter, r *http.Request) {
	requestPath, err := his.getRequestPath(r)
	logger.Info("HTTP Index request for path: " + requestPath)
	if err != nil {
		http.Error(w, "Invalid request path", http.StatusBadRequest)
		return
	}

	username := "anonymous"
	permission := config.ReadPermission
	ok, err := his.Auth.Authenticator.Authorize(username, requestPath, permission)
	if err != nil {
		http.Error(w, "Error checking permissions", http.StatusInternalServerError)
		logger.Error("Error checking permissions: ", err)
		return
	}
	if !ok {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Check if the path is a file or a directory.
	metadata, err := his.Storage.Stat(requestPath)
	if err != nil {
		http.Error(w, "Error accessing path", http.StatusInternalServerError)
		logger.Error("Error accessing path: ", err)
		return
	}

	// If the metadata indicates it's a file, handle file download.
	if !metadata.IsDir {
		his.handleFileDownload(w, r, requestPath, metadata)
		return
	}

	contents, err := his.Storage.List(requestPath)
	if err != nil {
		http.Error(w, "Error listing directory", http.StatusInternalServerError)
		return
	}

	// Check the Accept header to determine the response format.
	acceptHeader := r.Header.Get("Accept")
	if strings.Contains(acceptHeader, "text/html") {
		// Respond with HTML.
		tmpl, err := template.New("index").Parse(dirListTemplate)
		if err != nil {
			http.Error(w, "Error creating template", http.StatusInternalServerError)
			return
		}

		data := TemplateData{
			Path:     requestPath,
			Contents: contents,
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, data)
	} else {
		// Default to JSON response.
		jsonContents, err := json.Marshal(contents)
		if err != nil {
			http.Error(w, "Error encoding directory contents", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonContents)
	}
}

func (his *HTTPIndexModule) getRequestPath(r *http.Request) (string, error) {
	// Extract and sanitize the path from the URL.
	requestPath := r.URL.Path
	if requestPath == "" {
		requestPath = "/"
	}
	// trim url prefix
	// /data/xxx -> /xxx
	requestPath = strings.TrimPrefix(requestPath, his.Config.Path)
	// Clean the path to prevent directory traversal attacks.
	requestPath = path.Clean(requestPath)

	if requestPath == "." {
		requestPath = "/"
	}

	return requestPath, nil
}

func (his *HTTPIndexModule) handleFileDownload(w http.ResponseWriter, r *http.Request, filePath string, metadata storage.Metadata) {
	fileSize := metadata.Size

	// Set the response header for file download.
	w.Header().Set("Content-Disposition", "attachment; filename="+path.Base(filePath))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Accept-Ranges", "bytes")

	// Check if the request is for a range of bytes.
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		// Parse range header to determine the byte range.
		start, end, err := parseRange(rangeHeader, fileSize)
		if err != nil {
			http.Error(w, "Invalid Range Header", http.StatusBadRequest)
			return
		}

		// Adjust the status code to 206 Partial Content.
		w.WriteHeader(http.StatusPartialContent)

		// Set the Content-Range header indicating the range of bytes we are providing.
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))

		// Read and write the specified range of bytes to the response.
		data, err := his.Storage.Read(filePath, start, end-start+1)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		w.Write(data)
	} else {
		// No range header, send the whole file.
		w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
		w.WriteHeader(http.StatusOK)

		// Read and write the whole file to the response.
		data, err := his.Storage.ReadAll(filePath)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

// parseRange parses the Range header string and returns the start and end byte positions.
func parseRange(rangeHeader string, fileSize int64) (start, end int64, err error) {
	// Example Range header: "bytes=0-499"
	rangeHeader = strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(rangeHeader, "-")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid range")
	}

	start, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, errors.New("invalid range start")
	}

	if parts[1] == "" {
		// If no end is specified, use the entire remaining file size.
		end = fileSize - 1
	} else {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, errors.New("invalid range end")
		}
	}

	if start >= fileSize || end >= fileSize {
		return 0, 0, errors.New("invalid range: out of bounds")
	}

	return start, end, nil
}

func (his *HTTPIndexModule) Start() error {
	logger.Info("Starting HTTP Index module")
	return nil
}
