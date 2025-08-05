# File: backend/go.mod
module github.com/thisisjackii/quicklynks/backend

go 1.21

require (
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/labstack/echo/v4 v4.11.3
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/spf13/viper v1.18.1
	github.com/stretchr/testify v1.8.4
	golang.org/x/crypto v0.16.0
	golang.org/x/exp v0.0.0-20231214170342-aacd6d4b4611
	rs.zerolog.com/log v0.0.0-20221120224412-6d3843468595
)

# ... (indirect dependencies will be added by `go mod tidy`)
