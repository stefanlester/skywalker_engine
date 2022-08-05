module myapp

go 1.18

replace github.com/stefanlester/skywalker => ../skywalker

require github.com/stefanlester/skywalker v0.0.0-20220731153529-aa0f58820857

require (
	github.com/go-chi/chi/v5 v5.0.7 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
)
