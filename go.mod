module github.com/practice-golang/9minutes

go 1.16

replace github.com/practice-golang/ngjson => ./ngjson

require (
	github.com/denisenkom/go-mssqldb v0.10.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/doug-martin/goqu/v9 v9.12.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/go-cmp v0.5.5
	github.com/json-iterator/go v1.1.11
	github.com/labstack/echo/v4 v4.2.2
	github.com/lib/pq v1.10.1
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/thoas/go-funk v0.8.0
	github.com/tidwall/gjson v1.7.5
	gopkg.in/guregu/null.v4 v4.0.0
	gopkg.in/ini.v1 v1.62.0
	modernc.org/sqlite v1.10.6
)
