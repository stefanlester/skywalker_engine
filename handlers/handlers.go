package handlers

import (
	"fmt"
	"io"
	"myapp/data"
	"net/http"
	"net/url"
	"os"

	"github.com/CloudyKit/jet/v6"
	"github.com/stefanlester/skywalker"
	"github.com/stefanlester/skywalker/filesystems"
)

// Handlers is the type for handlers, and gives access to Skywalker and its models
// Handlers struct contains application-wide dependencies that handlers might need.
type Handlers struct {
	App    *skywalker.Skywalker
	Models data.Models
}

// Home renders the dashboard landing page. It reports the framework version and
// which storage backends are actually configured (built from .env) so the page
// reflects real runtime state instead of static marketing copy.
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	configured := make(map[string]bool, len(h.App.FileSystems))
	for k := range h.App.FileSystems {
		configured[k] = true
	}

	type backend struct {
		Name   string
		Label  string
		Blurb  string
		Active bool
	}
	backends := []backend{
		{"MINIO", "MinIO", "S3-compatible object storage", configured["MINIO"]},
		{"S3", "Amazon S3", "AWS object storage", configured["S3"]},
		{"SFTP", "SFTP", "File transfer over SSH", configured["SFTP"]},
		{"WEBDAV", "WebDAV", "Remote files over HTTP", configured["WEBDAV"]},
	}

	vars := make(jet.VarMap)
	vars.Set("version", h.App.Version)
	vars.Set("backends", backends)
	vars.Set("activeCount", len(configured))

	if err := h.render(w, r, "home", vars, nil); err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

// List File System Handler
// ListFS handles the HTTP request to list files in a specified filesystem. It supports different filesystem types and directories.
func (h *Handlers) ListFS(w http.ResponseWriter, r *http.Request) {
	var list []filesystems.Listing

	// Retrieve filesystem type from URL query parameters, defaulting to an empty string if not provided.
	fsType := ""
	if r.URL.Query().Get("fs-type") != "" {
		fsType = r.URL.Query().Get("fs-type")
	}

	// Retrieve current path from URL query parameters, default to the root directory if not provided.
	curPath := "/"
	if r.URL.Query().Get("curPath") != "" {
		curPath = r.URL.Query().Get("curPath")
		curPath, _ = url.QueryUnescape(curPath) // Decode the path to handle any URL-encoded characters.
	}

	// Handle filesystem listing based on the type.
	if fsType != "" {
		// Look up the configured backend by name and use it through the FS interface,
		// so any backend (MINIO, S3, SFTP, WEBDAV) works without a concrete switch.
		fs, ok := h.App.FileSystems[fsType].(filesystems.FS)
		if !ok {
			h.App.ErrorLog.Printf("unknown or unconfigured filesystem: %s", fsType)
			return
		}

		// List the contents of the directory.
		l, err := fs.List(curPath)
		if err != nil {
			h.App.ErrorLog.Println(err) // Log error and return early on failure.
			return
		}

		list = l
	}

	// Prepare variables for rendering the template.
	vars := make(jet.VarMap)
	vars.Set("list", list)
	vars.Set("fs_type", fsType)
	vars.Set("curPath", curPath)

	// Render the list-fs template with the file listing.
	err := h.render(w, r, "list-fs", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println(err) // Log rendering errors.
	}
}

// UploadToFS renders a simple form for uploading files to a filesystem.
func (h *Handlers) UploadToFS(w http.ResponseWriter, r *http.Request) {
	fsType := r.URL.Query().Get("type")

	// Prepare variables for rendering the upload form.
	vars := make(jet.VarMap)
	vars.Set("fs_type", fsType)

	// Render the upload template.
	err := h.render(w, r, "upload", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println(err) // Log rendering errors.
	}
}

// PostUploadToFS handles the HTTP POST request to upload files to the specified filesystem.
func (h *Handlers) PostUploadToFS(w http.ResponseWriter, r *http.Request) {
	// Retrieve the file from the form input using the getFileToUpload function.
	fileName, err := getFileToUpload(r, "formFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // Respond with an error if file retrieval fails.
		return
	}

	// Get the upload type from form data.
	uploadType := r.Form.Get("upload-type")

	// Look up the configured backend by name and use it through the FS interface.
	fs, ok := h.App.FileSystems[uploadType].(filesystems.FS)
	if !ok {
		h.App.ErrorLog.Printf("unknown or unconfigured filesystem: %s", uploadType)
		http.Error(w, "unknown or unconfigured filesystem", http.StatusBadRequest)
		return
	}

	err = fs.Put(fileName, "") // Attempt to upload the file.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // Respond with an error if upload fails.
		return
	}

	// Set a success message in the session and redirect the user to the upload page.
	h.App.Session.Put(r.Context(), "flash", "File uploaded!")
	http.Redirect(w, r, "/files/upload?type="+uploadType, http.StatusSeeOther)
}

// getFileToUpload processes the file upload from the form, saves it locally, and returns the local file path.
func getFileToUpload(r *http.Request, fieldName string) (string, error) {
	// Parse the multipart form with a max memory of 10 MB.
	_ = r.ParseMultipartForm(10 << 20)

	// Retrieve the file and its header from the form field.
	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", err // Return an error if the file is not found in the form.
	}
	defer file.Close() // Ensure the file is closed after processing.

	// Create a destination file in the ./tmp directory to store the uploaded file.
	dst, err := os.Create(fmt.Sprintf("./tmp/%s", header.Filename))
	if err != nil {
		return "", err // Return an error if the file creation fails.
	}
	defer dst.Close() // Ensure the destination file is closed after processing.

	// Copy the uploaded file to the destination file.
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err // Return an error if the copy operation fails.
	}

	// Return the path to the stored file.
	return fmt.Sprintf("./tmp/%s", header.Filename), nil
}

// DeleteFromFS handles the HTTP request to delete a file from a specified filesystem.
func (h *Handlers) DeleteFromFS(w http.ResponseWriter, r *http.Request) {
	fsType := r.URL.Query().Get("fs_type") // Retrieve the filesystem type from the URL query.
	item := r.URL.Query().Get("file")      // Retrieve the file name to be deleted from the URL query.

	// Look up the configured backend by name and use it through the FS interface.
	fs, ok := h.App.FileSystems[fsType].(filesystems.FS)
	if !ok {
		h.App.ErrorLog.Printf("unknown or unconfigured filesystem: %s", fsType)
		return
	}

	// Attempt to delete the specified file and handle the response.
	deleted := fs.Delete([]string{item})
	if deleted {
		// If deletion is successful, store a success message and redirect.
		h.App.Session.Put(r.Context(), "flash", fmt.Sprintf("%s was deleted", item))
		http.Redirect(w, r, "/list-fs?fs-type="+fsType, http.StatusSeeOther)
	}
}
