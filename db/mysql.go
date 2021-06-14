package db

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/practice-golang/9minutes/models"
)

type Mysql struct{ Dsn string }

// initDB - Prepare DB
func (d *Mysql) initDB() (*sql.DB, error) {
	var err error

	Dbo, err = sql.Open("mysql", d.Dsn)
	if err != nil {
		return nil, err
	}

	return Dbo, nil
}

func (d *Mysql) CreateDB() error {
	return nil
}

// CreateBoardManagerTable - Create board manager table
func (d *Mysql) CreateBoardManagerTable(recreate bool) error {
	sql := `CREATE DATABASE IF NOT EXISTS #DATABASE_NAME;`
	sql = strings.ReplaceAll(sql, "#DATABASE_NAME", DatabaseName)
	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	sql = ""
	if recreate {
		sql += `DROP TABLE IF EXISTS #TABLE_NAME;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS #TABLE_NAME (
		IDX INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
		NAME VARCHAR(128) NULL DEFAULT NULL COLLATE 'utf8_general_ci',
		PRICE DOUBLE NULL DEFAULT NULL,
		AUTHOR VARCHAR(128) NULL DEFAULT NULL COLLATE 'utf8_general_ci',
		ISBN VARCHAR(13) NULL DEFAULT NULL COLLATE 'utf8_general_ci',
		PRIMARY KEY (IDX),
		UNIQUE INDEX ISBN (ISBN),
		INDEX IDX (IDX)
	)
	COLLATE='utf8_general_ci'
	ENGINE=InnoDB;`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", BoardManagerTable)

	_, err = Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserFieldTable - Create user manager table
func (d *Mysql) CreateUserFieldTable(recreate bool) error {
	sql := `CREATE DATABASE IF NOT EXISTS #DATABASE_NAME;`
	sql = strings.ReplaceAll(sql, "#DATABASE_NAME", DatabaseName)
	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	sql = ""
	if recreate {
		sql += `DROP TABLE IF EXISTS #TABLE_NAME;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS #TABLE_NAME (
		IDX INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
		NAME VARCHAR(128) NULL DEFAULT NULL COLLATE 'utf8_general_ci',
		PRICE DOUBLE NULL DEFAULT NULL,
		AUTHOR VARCHAR(128) NULL DEFAULT NULL COLLATE 'utf8_general_ci',
		ISBN VARCHAR(13) NULL DEFAULT NULL COLLATE 'utf8_general_ci',
		PRIMARY KEY (IDX),
		UNIQUE INDEX ISBN (ISBN),
		INDEX IDX (IDX)
	)
	COLLATE='utf8_general_ci'
	ENGINE=InnoDB;`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserFieldTable)

	_, err = Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserTable - Create user table
func (d *Mysql) CreateUserTable(recreate bool) error {
	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS "#TABLE_NAME";`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS "#TABLE_NAME" (
		"IDX"			INTEGER,
		"NAME"			TEXT,
		"CODE"			TEXT,
		"TYPE"			TEXT,
		"FIELD_NAME"	TEXT UNIQUE,
		"ORDER"			INTEGER,
		PRIMARY KEY("IDX" AUTOINCREMENT)
	);`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserFieldTable)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateBasicBoard - Create board table
func (d *Mysql) CreateBasicBoard(tableInfo models.Board, recreate bool) error {
	return nil
}

// CreateCustomBoard - Create board table
func (d *Mysql) CreateCustomBoard(tableInfo models.Board, fields []models.Field, recreate bool) error {
	return nil
}

// EditBasicBoard - Create board table
func (d *Mysql) EditBasicBoard(tableInfoOld models.Board, tableInfoNew models.Board) error {
	return nil
}

// EditCustomBoard - Create custom table
func (d *Mysql) EditCustomBoard(tableInfoOld models.Board, tableInfoNew models.Board) error {
	return nil
}

// DeleteBoard - Delete a board table
func (d *Mysql) DeleteBoard(tableName string) error {
	return nil
}

// CreateComment - Create comment table
func (d *Mysql) CreateComment(tableInfo models.Board, recreate bool) error {
	return nil
}

// EditUserTableFields - Edit user table schema
func (d *Mysql) EditUserTableFields(fieldsInfoOld []models.UserColumn, fieldsInfoNew []models.UserColumn) error {
	return nil
}

// DeleteUserTableFields - Delete user table field
func (d *Mysql) DeleteUserTableFields(fieldsInfoRemove []models.UserColumn) error {
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

	log.Println("Sqlite/DeleteUserTableFields: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
