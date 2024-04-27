package handlers

import (
	"myapp/data"
	"net/http"
	"net/url"

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
