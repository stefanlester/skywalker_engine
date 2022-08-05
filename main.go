package main

import (
	"github.com/stefanlester/skywalker"

	"myapp/handlers"
)

type application struct {
	App      *skywalker.Skywalker
	Handlers *handlers.Handlers
}

func main() {
	c := initApplication()
	c.App.ListenAndServe()
}
