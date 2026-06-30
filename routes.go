package main

import (
	"log"
	"net/http"

	"github.com/stefanlester/skywalker"
	"github.com/stefanlester/skywalker/filesystems"

	"github.com/go-chi/chi/v5"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes

	// add routes here
	a.get("/", a.Handlers.Home)

	// list fs route
	a.get("/list-fs", a.Handlers.ListFS)

	// upload to fs routes
	a.get("/files/upload", a.Handlers.UploadToFS)
	a.post("/files/upload", a.Handlers.PostUploadToFS)

	// delete from fs route
	a.get("/delete-from-fs", a.Handlers.DeleteFromFS)

	// minio fs test route
	a.get("/test-minio", func(w http.ResponseWriter, r *http.Request) {
		f, ok := a.App.FileSystems["MINIO"].(filesystems.FS)
		if !ok {
			log.Println("unknown or unconfigured filesystem: MINIO")
			return
		}

		files, err := f.List("")
		if err != nil {
			log.Println(err)
			return
		}

		for _, file := range files {
			log.Println("File:" + file.Key)
		}
	})

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	// routes from celeritas
	a.App.Routes.Mount("/skywalker", skywalker.Routes())
	a.App.Routes.Mount("/api", a.ApiRoutes())

	return a.App.Routes
}
