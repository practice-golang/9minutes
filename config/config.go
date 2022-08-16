package config

import (
	"9minutes/auth"
	"9minutes/db"
	"9minutes/email"
	"9minutes/model"
)

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
	DatabaseType: model.SQLITE,
	DatabaseName: "9m",
	FilePath:     "./9minutes.db",
}

var DatabaseInfoMySQL = db.DBInfo{
	DatabaseType:  model.MYSQL,
	Protocol:      "tcp",
	Addr:          "localhost",
	Port:          "13306",
	DatabaseName:  "9m",
	GrantID:       "root",
	GrantPassword: "",
}

var DatabaseInfoPgPublic = db.DBInfo{
	DatabaseType:  model.POSTGRES,
	Protocol:      "tcp",
	Addr:          "localhost",
	Port:          "5432",
	DatabaseName:  "postgres",
	SchemaName:    "public",
	GrantID:       "root",
	GrantPassword: "pgsql",
}

var DatabaseInfoPgSchema = db.DBInfo{
	DatabaseType:  model.POSTGRES,
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
	DatabaseType:  model.POSTGRES,
	Protocol:      "tcp",
	Addr:          "localhost",
	Port:          "5432",
	DatabaseName:  "9mdb",
	SchemaName:    "9m",
	GrantID:       "root",
	GrantPassword: "pgsql",
}

var DatabaseInfoSqlServer = db.DBInfo{
	DatabaseType:  model.SQLSERVER,
	Protocol:      "tcp",
	Addr:          "localhost",
	Port:          "1433",
	DatabaseName:  "9mdb",
	SchemaName:    "dbo",
	GrantID:       "sa",
	GrantPassword: "SQLServer1433",
}

// Oracle Local XE >= 12c
var DatabaseInfoOracle = db.DBInfo{
	DatabaseType:  model.ORACLE,
	Protocol:      "tcp",
	Addr:          "localhost",
	Port:          "1521",
	DatabaseName:  "XE",         // physical&on-premise: database name, cloud: service name
	SchemaName:    "",           // not use
	GrantID:       "myaccount",  // physical&on-premise: username only, cloud: username and database name
	GrantPassword: "mypassword", // password
	FilePath:      "",           // wallet file path for cloud using
}

// Oracle Local XE. >= 12c ID = system for create DB
// Require only for database(=user) creation. Comment out this after database creation is done
var DatabaseInfoOracleSystem = db.DBInfo{
	DatabaseType:  model.ORACLE,
	Protocol:      "tcp",
	Addr:          "localhost",
	Port:          "1521",
	DatabaseName:  "XE",
	SchemaName:    "",
	GrantID:       "system",
	GrantPassword: "oracle",
	FilePath:      "",
}

// Oracle Cloud 19c
var DatabaseInfoOracleCloud = db.DBInfo{
	DatabaseType:  model.ORACLE,
	Protocol:      "tcp",
	Addr:          "adb.ap-seoul-1.oraclecloud.com",
	Port:          "1522",
	DatabaseName:  "a12345abcde1_mydbname_low.adb.oraclecloud.com",
	SchemaName:    "",
	GrantID:       "myaccount",
	GrantPassword: "MyPassword!522",
	FilePath:      "./wallet_myaccount",
}

// Oracle Cloud 19c. ID = ADMIN for create DB
// Require only for database(=user) creation. Comment out this after database creation is done
var DatabaseInfoOracleCloudAdmin = db.DBInfo{
	DatabaseType:  model.ORACLE,
	Protocol:      "tcp",
	Addr:          "adb.ap-seoul-1.oraclecloud.com",
	Port:          "1522",
	DatabaseName:  "a12345abcde1_mydbname_low.adb.oraclecloud.com",
	SchemaName:    "",
	GrantID:       "admin",
	GrantPassword: "MyPassword!522",
	FilePath:      "./wallet_admin",
}

var StoreInfoMemory = auth.SessionStoreInfo{
	StoreType: auth.MEMSTORE,
	Address:   "",
	Port:      "",
}

var StoreInfoETCD = auth.SessionStoreInfo{
	StoreType: auth.ETCD,
	Address:   "localhost",
	Port:      "2379",
}

var StoreInfoRedis = auth.SessionStoreInfo{
	StoreType: auth.REDIS,
	Address:   "localhost",
	Port:      "6379",
}

var EmailServerDirect = email.Config{
	UseEmail:   false,
	Domain:     "http://localhost:8080",
	SendDirect: true,
	Service:    email.Service{KeyDKIM: ""},
	SenderInfo: email.From{Name: "Administrator", Email: "admin@domain.ext"},
}

var EmailServerSMTP = email.Config{
	UseEmail:   false,
	Domain:     "http://localhost:8080",
	SendDirect: false,
	Service:    email.Service{Host: "smtp.gmail.com", Port: "587", ID: "", Password: ""},
	SenderInfo: email.From{Name: "Administrator", Email: "admin@domain.ext"},
}
