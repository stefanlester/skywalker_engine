package main

import (
	"log"
	"myapp/handlers"
	"os"

	"github.com/stefanlester/skywalker"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	//init skywalker
	skywalker := &skywalker.Skywalker{}
	err = skywalker.New(path)
	if err != nil {
		log.Fatal(err)
	}

	skywalker.AppName = "myapp"

	//init handlers
	myhandlers := &handlers.Handlers{
		App: skywalker,
	}

	app := &application{
		App:      skywalker,
		Handlers: myhandlers,
	}

	app.App.Routes = app.routes()

	return app
}
