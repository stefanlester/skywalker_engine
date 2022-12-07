module myapp

go 1.18

replace github.com/stefanlester/skywalker => ../skywalker

require (
	github.com/CloudyKit/jet/v6 v6.1.0
	github.com/go-chi/chi/v5 v5.0.7
	github.com/stefanlester/skywalker v0.0.0-20220813131536-52f59af98d97
	github.com/upper/db/v4 v4.6.0
)

require (
	github.com/CloudyKit/fastprinter v0.0.0-20200109182630-33d98a066a53 // indirect
	github.com/alexedwards/scs/v2 v2.5.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.13.0 // indirect
	github.com/jackc/pgx/v4 v4.17.2 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/segmentio/fasthash v1.0.3 // indirect
	golang.org/x/crypto v0.4.0 // indirect
	golang.org/x/text v0.5.0 // indirect
)
