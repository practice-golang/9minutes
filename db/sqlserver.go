package db

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver"
	"github.com/practice-golang/9minutes/models"
)

type Sqlserver struct{ Dsn string }

// initDB - Prepare DB
func (d *Sqlserver) initDB() (*sql.DB, error) {
	var err error

	Dbo, err = sql.Open("sqlserver", d.Dsn)
	if err != nil {
		return nil, err
	}

	return Dbo, nil
}

func (d *Sqlserver) CreateDB() error {
	return nil
}

// CreateBoardManagerTable - Create board manager table
func (d *Sqlserver) CreateBoardManagerTable(recreate bool) error {
	sql := `
	USE master
	-- GO

	IF NOT EXISTS(
		SELECT name
		FROM sys.databases
		WHERE name=N'#DATABASE'
	)
	CREATE DATABASE "#DATABASE"
	-- GO
	`
	sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)

	log.Println("Sqlserver/CreateBoardManagerTable: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	if recreate {
		sql = `USE "#DATABASE"`
		sql += `
		IF OBJECT_ID('#TABLE_NAME','U') IS NOT NULL
		DROP TABLE "#TABLE_NAME"
		-- GO
		`

		sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
		sql = strings.ReplaceAll(sql, "#TABLE_NAME", BoardManagerTable)

		_, err := Dbo.Exec(sql)
		if err != nil {
			return err
		}
	}

	sql = `USE "#DATABASE"`
	sql += `
	IF OBJECT_ID(N'#TABLE_NAME', N'U') IS NULL
	CREATE TABLE "#TABLE_NAME" (
		IDX INT NOT NULL IDENTITY PRIMARY KEY,
		NAME VARCHAR(128) NULL DEFAULT NULL,
		CODE VARCHAR(64) NULL DEFAULT NULL,
		TYPE VARCHAR(64) NULL DEFAULT NULL,
		"TABLE" VARCHAR(64) NULL DEFAULT NULL,
		GRANT_READ VARCHAR(16) NULL DEFAULT NULL,
		GRANT_WRITE VARCHAR(16) NULL DEFAULT NULL,
		GRANT_COMMENT VARCHAR(16) NULL DEFAULT NULL,
		FILE_UPLOAD VARCHAR(2) NULL DEFAULT NULL,
		FIELDS TEXT NULL DEFAULT NULL,
	)
	--GO`

	sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", BoardManagerTableName)

	log.Println("Sqlserver/CreateBoardManagerTable: ", sql)

	_, err = Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserTable - Create user table
func (d *Sqlserver) CreateUserTable(recreate bool) error {
	sql := `
	USE master
	-- GO

	IF NOT EXISTS(
		SELECT name
		FROM sys.databases
		WHERE name=N'#DATABASE'
	)
	CREATE DATABASE "#DATABASE"
	-- GO
	`
	sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
	log.Println("Sqlserver/CreateUserFieldTable: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	if recreate {
		sql = `USE "#DATABASE"`
		sql += `
		IF OBJECT_ID('#TABLE_NAME','U') IS NOT NULL
		DROP TABLE "#TABLE_NAME"
		-- GO
		`

		sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
		sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTableName)
		log.Println("Sqlserver/CreateUserFieldTable: ", sql)

		_, err := Dbo.Exec(sql)
		if err != nil {
			return err
		}
	}

	sql = `USE "#DATABASE"`
	sql += `
	IF OBJECT_ID(N'#TABLE_NAME', N'U') IS NULL
	CREATE TABLE "#TABLE_NAME" (
		IDX INT NOT NULL IDENTITY PRIMARY KEY,
		USERNAME VARCHAR(128) NULL DEFAULT NULL,
		PASSWORD VARCHAR(128) NULL DEFAULT NULL,
		EMAIL VARCHAR(128) NULL DEFAULT NULL,
		ADMIN VARCHAR(2) NULL DEFAULT NULL,
		APPROVAL VARCHAR(2) NULL DEFAULT NULL,
		REG_DTTM BIGINT NULL DEFAULT NULL,
	)
	--GO`

	sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTableName)
	log.Println("Sqlserver/CreateUserTable: ", sql)

	_, err = Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserFieldTable - Create user manager table
func (d *Sqlserver) CreateUserFieldTable(recreate bool) error {
	sql := `
	USE master
	-- GO

	IF NOT EXISTS(
		SELECT name
		FROM sys.databases
		WHERE name=N'#DATABASE'
	)
	CREATE DATABASE "#DATABASE"
	-- GO
	`
	sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
	log.Println("Sqlserver/CreateUserFieldTable: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	if recreate {
		sql = `USE "#DATABASE"`
		sql += `
		IF OBJECT_ID('#TABLE_NAME','U') IS NOT NULL
		DROP TABLE "#TABLE_NAME"
		-- GO
		`

		sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
		sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserFieldTableName)
		log.Println("Sqlserver/CreateUserFieldTable: ", sql)

		_, err := Dbo.Exec(sql)
		if err != nil {
			return err
		}
	}

	sql = `USE "#DATABASE"`
	sql += `
	IF OBJECT_ID(N'#TABLE_NAME', N'U') IS NULL
	CREATE TABLE "#TABLE_NAME" (
		IDX INT NOT NULL IDENTITY PRIMARY KEY,
		TITLE VARCHAR(256) NULL DEFAULT NULL,
		CONTENT TEXT NULL DEFAULT NULL,
		IS_MEMBER VARCHAR(2) NULL DEFAULT NULL,
		WRITER_IDX VARCHAR(11) NULL DEFAULT NULL,
		WRITER_NAME VARCHAR(64) NULL DEFAULT NULL,
		WRITER_PASSWORD VARCHAR(128) NULL DEFAULT NULL,
		FILES TEXT NULL DEFAULT NULL,
		REG_DTTM BIGINT NULL DEFAULT NULL,
	)
	--GO`

	sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserFieldTableName)
	log.Println("Sqlserver/CreateUserFieldTable: ", sql)

	_, err = Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateBasicBoard - Create board table
func (d *Sqlserver) CreateBasicBoard(tableInfo models.Board, recreate bool) error {
	sql := ``
	if recreate {
		sql = `USE "#DATABASE"`
		sql += `
		IF OBJECT_ID('#TABLE_NAME','U') IS NOT NULL
		DROP TABLE "#TABLE_NAME"
		-- GO
		`

		sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
		sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTableName)
		log.Println("Sqlserver/CreateBasicBoard: ", sql)

		_, err := Dbo.Exec(sql)
		if err != nil {
			return err
		}
	}

	sql = `USE "#DATABASE"`
	sql += `
	IF OBJECT_ID(N'#TABLE_NAME', N'U') IS NULL
	CREATE TABLE "#TABLE_NAME" (
		IDX INT NOT NULL IDENTITY PRIMARY KEY,
		TITLE VARCHAR(256) NULL DEFAULT NULL,
		CONTENT TEXT NULL DEFAULT NULL,
		IS_MEMBER VARCHAR(2) NULL DEFAULT NULL,
		WRITER_IDX VARCHAR(11) NULL DEFAULT NULL,
		WRITER_NAME VARCHAR(64) NULL DEFAULT NULL,
		WRITER_PASSWORD VARCHAR(128) NULL DEFAULT NULL,
		FILES TEXT NULL DEFAULT NULL,
		REG_DTTM BIGINT NULL DEFAULT NULL,
	)
	--GO`

	sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", tableInfo.Table.String)
	log.Println("Sqlserver/CreateBasicBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateCustomBoard - Create board table
func (d *Sqlserver) CreateCustomBoard(tableInfo models.Board, fields []models.Field, recreate bool) error {
	return nil
}

// EditBasicBoard - Create board table
func (d *Sqlserver) EditBasicBoard(tableInfoOld models.Board, tableInfoNew models.Board) error {
	return nil
}

// EditBasicBoard - Create board table
func (d *Sqlserver) EditCustomBoard(tableInfoOld models.Board, tableInfoNew models.Board) error {
	return nil
}

// DeleteBoard - Delete a board table
func (d *Sqlserver) DeleteBoard(tableName string) error {
	return nil
}

// CreateComment - Create comment table
func (d *Sqlserver) CreateComment(tableInfo models.Board, recreate bool) error {
	return nil
}

// EditUserTableFields - Edit user table schema
func (d *Sqlserver) EditUserTableFields(fieldsInfoOld []models.UserColumn, fieldsInfoNew []models.UserColumn, notUse []string) error {
	return nil
}

// DeleteUserTableFields - Delete user table field
func (d *Sqlserver) DeleteUserTableFields(fieldsInfoRemove []models.UserColumn) error {
	remove := fieldsInfoRemove
	sql := ""

	if len(remove) > 0 {
		sqlRemove := `ALTER TABLE "#TABLE_NAME" `
		for _, r := range remove {
			sqlRemove += ` DROP COLUMN ` + r.ColumnName.String + `, `
		}
		if strings.Contains(sqlRemove, "DROP COLUMN") {
			sqlRemove = sqlRemove[:len(sqlRemove)-2]
		}
		sql += sqlRemove + `; `
	}

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTable)

	log.Println("Sqlserver/DeleteUserTableFields: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// AddUserTableFields - Add user column
func (d *Sqlserver) AddUserTableFields(fields []models.UserColumn) error {
	sql := ""
	for _, a := range fields {
		sql += `ALTER TABLE "#TABLE_NAME" ADD COLUMN ` + a.ColumnName.String + ` `
		switch a.Type.String {
		case "text":
			sql += ` TEXT`
		case "number":
			sql += ` INTEGER`
		case "real":
			sql += ` REAL`
		}

		sql += `; `
	}

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTable)

	log.Println("Sqlite/EditUserTableFields: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// SelectColumnNames - Get column names of table
func (d *Sqlserver) SelectColumnNames(table string) (sql.Result, error) {
	result, err := Dbo.Exec("PRAGMA TABLE_INFO(" + table + ")")
	if err != nil {
		return nil, err
	}

	return result, nil
}
