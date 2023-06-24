package handlers

import (
	"myapp/data"
	"net/http"

	"github.com/stefanlester/skywalker"
	"github.com/stefanlester/skywalker/filesystems"
)

// Handlers is the type for handlers, and gives access to Celeritas and models
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

func (h *Handlers) ListFS(w http.ResponseWriter, r *http.Request) {
	var fs filesystems.FS
	var list []filesystems.Listing

	fsType := ""
	if r.URL.Query().Get("fs-type") != "" {
		fsType = r.URL.Query().Get("fs-type")
	}

}