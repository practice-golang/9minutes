package db

import (
	"database/sql"
	"errors"

	"github.com/practice-golang/9minutes/models"
)

const (
	_ = iota
	// SQLITE - sqlite
	SQLITE
	// SQLSERVER - mssql
	SQLSERVER
	// MYSQL - mysql
	MYSQL
	// POSTGRES - pgsql
	POSTGRES
)

type DBI interface {
	initDB() (*sql.DB, error)
	CreateDB() error
	CreateBoardManagerTable(recreate bool) error
	CreateUserFieldTable(recreate bool) error
	CreateUserTable(recreate bool) error
	CreateBasicBoard(tableInfo models.Board, recreate bool) error
	CreateCustomBoard(tableInfo models.Board, fields []models.Field, recreate bool) error
	EditBasicBoard(tableInfoOld models.Board, tableInfoNew models.Board) error
	EditCustomBoard(tableInfoOld models.Board, tableInfoNew models.Board) error
	DeleteBoard(tableName string) error
	CreateComment(tableInfo models.Board, recreate bool) error
	AddUserTableFields(fields []models.UserColumn) error
	EditUserTableFields(fieldsInfoOld []models.UserColumn, fieldsInfoNew []models.UserColumn, notUse []string) error
	DeleteUserTableFields(fieldsInfoRemove []models.UserColumn) error
	SelectColumnNames(table string) (sql.Result, error)
}

var (
	Dbi                       DBI    // DB Object Interface
	Dsn                       string // Data Source Name
	DatabaseName              = "9minutes"
	BoardManagerTable         = "BOARD_TABLES"
	BoardManagerTableNoQuotes = "BOARD_TABLES" // makeshift - postgres
	UserFieldTable            = "USER_FIELDS"
	UserFieldTableNoQuotes    = "USER_FIELDS" // makeshift - postgres
	UserTable                 = "USERS"
	UserTableNoQuotes         = "USERS" // makeshift - postgres
	Dbo                       *sql.DB
	DBType                    int
	UpdateScope               []string     // UPDATE ... WHERE IDX=?
	IgnoreScope               []string     // Ignore if nil or null
	listCount                 uint     = 3 // Default list count
	OrderScope                string       // Default order column name
)

// InitDB - Prepare DB
func InitDB(driver int) (DBI, error) {
	var err error
	var dbi DBI

	dbi, err = dbFactory(driver)

	return dbi, err
}

func dbFactory(driver int) (DBI, error) {
	switch driver {
	case SQLITE:
		dbi := &Sqlite{Dsn: Dsn}
		dbi.initDB()
		return dbi, nil
	case MYSQL:
		dbi := &Mysql{Dsn: Dsn}
		dbi.initDB()
		return dbi, nil
	case SQLSERVER:
		dbi := &Sqlserver{Dsn: Dsn}
		dbi.initDB()
		return dbi, nil
	case POSTGRES:
		dbi := &Postgres{Dsn: Dsn}
		dbi.initDB()
		return dbi, nil
	default:
		return nil, errors.New("nothing to support DB")
	}
}
