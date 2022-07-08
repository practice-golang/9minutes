package config

import "9minutes/db"

var (
	AdminUserCountPerPage  int = 10
	AdminBoardCountPerPage int = 10
	// ContentsCountPerPage   int = 25
	ContentsCountPerPage int = 5
	CommentCountPerPage  int = 3
)

var (
	StaticPath = "../static"
	UploadPath = "../upload"
	HtmlPath   = "../html"
)

var (
	UserGrades = newCollection("admin", "manager", "regular_user", "pending_user", "guest", "banned_user")
)

var DatabaseInfoSQLite = db.DBInfo{
	DatabaseType: db.SQLITE,
	DatabaseName: "9m",
	FilePath:     "./9minutes.db",
}

var DatabaseInfoMySQL = db.DBInfo{
	DatabaseType:  db.MYSQL,
	Protocol:      "tcp",
	Addr:          "localhost",
	Port:          "13306",
	DatabaseName:  "9m",
	GrantID:       "root",
	GrantPassword: "",
}

var DatabaseInfoPgPublic = db.DBInfo{
	DatabaseType:  db.POSTGRES,
	Protocol:      "tcp",
	Addr:          "localhost",
	Port:          "5432",
	DatabaseName:  "postgres",
	SchemaName:    "public",
	GrantID:       "root",
	GrantPassword: "pgsql",
}

var DatabaseInfoPgSchema = db.DBInfo{
	DatabaseType:  db.POSTGRES,
	Protocol:      "tcp",
	Addr:          "localhost",
	Port:          "5432",
	DatabaseName:  "postgres",
	SchemaName:    "9m",
	GrantID:       "root",
	GrantPassword: "pgsql",
}

// For not using database name 'postgres', you should create database yourself
var DatabaseInfoPgOtherDatabase = db.DBInfo{
	DatabaseType:  db.POSTGRES,
	Protocol:      "tcp",
	Addr:          "localhost",
	Port:          "5432",
	DatabaseName:  "9mdb",
	SchemaName:    "9m",
	GrantID:       "root",
	GrantPassword: "pgsql",
}

var DatabaseInfoSqlServer = db.DBInfo{
	DatabaseType:  db.SQLSERVER,
	Protocol:      "tcp",
	Addr:          "localhost",
	Port:          "1433",
	DatabaseName:  "9mdb",
	SchemaName:    "dbo",
	GrantID:       "sa",
	GrantPassword: "SQLServer1433",
}
