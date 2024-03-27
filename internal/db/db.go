package db

import (
	"9minutes/model"
	"database/sql"
	"errors"
	"net/url"
	"strconv"

	go_ora "github.com/sijms/go-ora/v2"
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
		UploadTable   string
		MemberTable   string
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
		CreateUploadTable() error
		CreateMemberTable() error
		CreateUserTable() error
		CreateUserVerificationTable() error

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
	Info.UploadTable = "uploads"
	Info.UserTable = "users"
	Info.MemberTable = "members"

	switch Info.DatabaseType {

	case model.SQLITE:
		dsn := Info.FilePath
		Obj = &SQLite{dsn: dsn}

		Con, err = Obj.connect()
		if err != nil {
			return err
		}

	case model.MYSQL:
		// dsn := Info.GrantID + ":" + Info.GrantPassword + "@" + Info.Protocol + "(" + Info.Addr + ":" + Info.Port + ")/" + Info.DatabaseName
		dsn := Info.GrantID + ":" + Info.GrantPassword + "@" + Info.Protocol + "(" + Info.Addr + ":" + Info.Port + ")/"
		Obj = &Mysql{dsn: dsn}

		Con, err = Obj.connect()
		if err != nil {
			return err
		}

	case model.POSTGRES:
		dsn := `host=` + Info.Addr + ` port=` + Info.Port + ` user=` + Info.GrantID + ` password=` + Info.GrantPassword + ` dbname=` + Info.DatabaseName + ` sslmode=disable`
		Obj = &Postgres{dsn: dsn}

		Con, err = Obj.connect()
		if err != nil {
			return err
		}

	case model.SQLSERVER:
		dsn := "sqlserver://" + Info.GrantID + ":" + Info.GrantPassword + "@" + Info.Addr + ":" + Info.Port + "?" + Info.DatabaseName + "&connction+timeout=30&encrypt=disable"
		Obj = &SqlServer{dsn: dsn}

		Con, err = Obj.connect()
		if err != nil {
			return err
		}

	case model.ORACLE:
		port, _ := strconv.Atoi(Info.Port)
		dsn := go_ora.BuildUrl(Info.Addr, port, Info.DatabaseName, Info.GrantID, Info.GrantPassword, nil)
		if Info.FilePath != "" {
			// dsn += "?SSL=enable&SSL Verify=false&WALLET=" + url.QueryEscape(Info.FilePath)
			// dsn += "?TRACE FILE=trace.log&prefetch_rows=10&SSL=enable&SSL Verify=false&WALLET=" + url.QueryEscape(Info.FilePath)
			dsn += "?prefetch_rows=10&SSL=enable&SSL Verify=false&WALLET=" + url.QueryEscape(Info.FilePath)
		}
		Obj = &Oracle{dsn: dsn}

		Con, err = Obj.connect()
		if err != nil {
			return err
		}

	default:
		return errors.New("database type not supported")
	}

	return nil
}
