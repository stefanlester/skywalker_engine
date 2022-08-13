module myapp

go 1.18

replace github.com/stefanlester/skywalker => ../skywalker

require (
	github.com/go-chi/chi/v5 v5.0.7
	github.com/stefanlester/skywalker v0.0.0-20220813131536-52f59af98d97
)

require (
	github.com/CloudyKit/fastprinter v0.0.0-20200109182630-33d98a066a53 // indirect
	github.com/CloudyKit/jet/v6 v6.1.0 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
)
