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
	"github.com/stefanlester/skywalker/filesystems/miniofilesystem"
)

// Handlers is the type for handlers, and gives access to Skywalker and its models
// Handlers struct contains application-wide dependencies that handlers might need.
type Handlers struct {
	App    *skywalker.Skywalker
	Models data.Models
}

// Home is the handler to render the home page
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.render(w, r, "home", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

// List File System Handler
// ListFS handles the HTTP request to list files in a specified filesystem. It supports different filesystem types and directories.
func (h *Handlers) ListFS(w http.ResponseWriter, r *http.Request) {
	var fs filesystems.FS
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
		switch fsType {
		case "MINIO":
			f := h.App.FileSystems["MINIO"].(miniofilesystem.Minio) // Type assertion to Minio filesystem.
			fs = &f
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

	// Handle file upload based on the specified filesystem type.
	switch uploadType {
	case "MINIO":
		fs := h.App.FileSystems["MINIO"].(miniofilesystem.Minio) // Type assertion to Minio filesystem.
		err = fs.Put(fileName, "")                               // Attempt to upload the file.
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // Respond with an error if upload fails.
			return
		}
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
	var fs filesystems.FS
	fsType := r.URL.Query().Get("fs_type") // Retrieve the filesystem type from the URL query.
	item := r.URL.Query().Get("file")      // Retrieve the file name to be deleted from the URL query.

	switch fsType {
	case "MINIO":
		f := h.App.FileSystems["MINIO"].(miniofilesystem.Minio) // Type assertion to Minio filesystem.
		fs = &f
	}

	// Attempt to delete the specified file and handle the response.
	deleted := fs.Delete([]string{item})
	if deleted {
		// If deletion is successful, store a success message and redirect.
		h.App.Session.Put(r.Context(), "flash", fmt.Sprintf("%s was deleted", item))
		http.Redirect(w, r, "/list-fs?fs-type="+fsType, http.StatusSeeOther)
	}
}
