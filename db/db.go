package db

import (
	"9minutes/model"
	"database/sql"
	"errors"
	"fmt"
)

const (
	_         = iota
	SQLITE    // SQLite
	MYSQL     // MySQL
	POSTGRES  // PostgreSQL
	SQLSERVER // MS SQL Server
)

type (
	DBInfo struct {
		DatabaseType  int
		Protocol      string
		Addr          string
		Port          string
		DatabaseName  string
		SchemaName    string
		TableName     string
		BoardTable    string
		UserTable     string
		GrantID       string
		GrantPassword string
		FilePath      string
	}

	DBObject interface {
		connect() (*sql.DB, error)
		CreateDB() error
		// CreateTable() error // Not use
		// Exec - Almost Same as sql.Exec()
		// Because of PostgreSQL and MS SQL Server, INSERT query and RETURN id way is not enough to use sql.Exec()
		// Return affected rows, last insert id, error
		// Not return sql.Result
		Exec(string, []interface{}, string) (int64, int64, error)
		CreateBoardTable() error
		CreateUserTable() error

		AddTableColumn(tableName string, column model.UserColumn) error
		EditTableColumn(tableName string, columnOld model.UserColumn, columnNew model.UserColumn) error
		DeleteTableColumn(tableName string, column model.UserColumn) error

		CreateBoard(tableInfo model.Board, recreate bool) error
		CreateComment(tableInfo model.Board, recreate bool) error

		RenameBoard(tableInfoOLD model.Board, tableInfoNEW model.Board) error
		RenameComment(tableInfoOLD model.Board, tableInfoNEW model.Board) error

		DeleteBoard(tableInfo model.Board) error
		DeleteComment(tableInfo model.Board) error

		GetPagingQuery(offset, count int) string
	}
)

var (
	Info DBInfo   // DB connection info
	Obj  DBObject // Duck interface
	Con  *sql.DB  // DB connection
)

func SetupDB() error {
	var err error

	Info.BoardTable = "boards"
	Info.UserTable = "users"

	switch Info.DatabaseType {

	case SQLITE:
		dsn := Info.FilePath
		Obj = &SQLite{dsn: dsn}

		Con, err = Obj.connect()
		if err != nil {
			return err
		}

	case MYSQL:
		// dsn := Info.GrantID + ":" + Info.GrantPassword + "@" + Info.Protocol + "(" + Info.Addr + ":" + Info.Port + ")/" + Info.DatabaseName
		dsn := Info.GrantID + ":" + Info.GrantPassword + "@" + Info.Protocol + "(" + Info.Addr + ":" + Info.Port + ")/"
		Obj = &Mysql{dsn: dsn}

		Con, err = Obj.connect()
		if err != nil {
			return err
		}

	case POSTGRES:
		dsn := `host=` + Info.Addr + ` port=` + Info.Port + ` user=` + Info.GrantID + ` password=` + Info.GrantPassword + ` dbname=` + Info.DatabaseName + ` sslmode=disable`
		Obj = &Postgres{dsn: dsn}

		Con, err = Obj.connect()
		if err != nil {
			return err
		}

	case SQLSERVER:
		dsn := "sqlserver://" + Info.GrantID + ":" + Info.GrantPassword + "@" + Info.Addr + ":" + Info.Port + "?" + Info.DatabaseName + "&connction+timeout=30"
		Obj = &SqlServer{dsn: dsn}

		Con, err = Obj.connect()
		if err != nil {
			return err
		}

	default:
		return errors.New("database type not supported " + fmt.Sprint(Info.DatabaseType))
	}

	return nil
}
