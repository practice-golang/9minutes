package config

type DBpath struct {
	Type     string
	Server   string // mysql, postgres, sqlserver
	Port     int    // mysql, postgres, sqlserver
	User     string // mysql, postgres, sqlserver
	Password string // mysql, postgres, sqlserver
	Database string // mysql, postgres, sqlserver
	Schema   string // postgres
	Filename string // sqlite
}

// DB - sqlite
var (
	DbInfo = DBpath{
		Type:     "sqlite", // "sqlite" "mysql" "postgres" "sqlserver"
		Filename: "./9minutes.db",
	}
)

// DB - mysql
// var (
// 	DbInfo = DBpath{
// 		Type:     "mysql", // "sqlite" "mysql" "postgres" "sqlserver"
// 		Server:   "127.0.0.1",
// 		Port:     13306,
// 		User:     "root",
// 		Password: "",
// 		Database: "floating-shelf",
// 	}
// )

// DB - ms-sqlserver
// var (
// 	DbInfo = DBpath{
// 		Type:     "sqlserver", // "sqlite" "mysql" "postgres" "sqlserver"
// 		Server:   "127.0.0.1",
// 		Port:     1433,
// 		User:     "sa",
// 		Password: "mssql",
// 		Database: "floating-shelf",
// 	}
// )

// DB - postgresql
// var (
// 	DbInfo = DBpath{
// 		Type:     "postgres", // "sqlite" "mysql" "postgres" "sqlserver"
// 		Server:   "127.0.0.1",
// 		Port:     5432,
// 		User:     "root",
// 		Password: "pgsql",
// 		Database: "postgres",
// 		Schema:   "floating-shelf",
// 	}
// )
