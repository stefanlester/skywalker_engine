package main

import (
	"github.com/stefanlester/skywalker"

	"myapp/data"
	"myapp/handlers"
)

type application struct {
	App      *skywalker.Skywalker
	Handlers *handlers.Handlers
	Models   data.Models
}

func main() {
	c := initApplication()
	c.App.ListenAndServe()
}
