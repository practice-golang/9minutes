package db

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/lib/pq"
	"github.com/practice-golang/9minutes/models"
)

type Postgres struct{ Dsn string }

// initDB - Prepare DB
func (d *Postgres) initDB() (*sql.DB, error) {
	var err error

	Dbo, err = sql.Open("postgres", d.Dsn)
	if err != nil {
		return nil, err
	}

	return Dbo, nil
}

func (d *Postgres) CreateDB() error {
	sql := `CREATE DATABASE #DATABASE_NAME;`
	sql = strings.ReplaceAll(sql, "#DATABASE_NAME", DatabaseName)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserTable - Create user table
func (d *Postgres) CreateUserTable(recreate bool) error {
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
		"COLUMN_NAME"	TEXT UNIQUE,
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

// CreateBoardManagerTable - Create board manager table
func (d *Postgres) CreateBoardManagerTable(recreate bool) error {
	sql := `CREATE SCHEMA IF NOT EXISTS #SCHEMA_NAME;`

	if recreate {
		sql += `DROP TABLE IF EXISTS #TABLE_NAME;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS #TABLE_NAME (
		"IDX" SERIAL PRIMARY KEY,
		"NAME" VARCHAR(128) NULL DEFAULT NULL,
		"PRICE" NUMERIC(10,2) NULL DEFAULT NULL,
		"AUTHOR" VARCHAR(128) NULL DEFAULT NULL,
		"ISBN" VARCHAR(13) UNIQUE NULL DEFAULT NULL
	);`

	sql = strings.ReplaceAll(sql, "#SCHEMA_NAME", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", BoardManagerTable)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserFieldTable - Create user manager table
func (d *Postgres) CreateUserFieldTable(recreate bool) error {
	sql := `CREATE SCHEMA IF NOT EXISTS #SCHEMA_NAME;`

	if recreate {
		sql += `DROP TABLE IF EXISTS #TABLE_NAME;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS #TABLE_NAME (
		"IDX" SERIAL PRIMARY KEY,
		"NAME" VARCHAR(128) NULL DEFAULT NULL,
		"PRICE" NUMERIC(10,2) NULL DEFAULT NULL,
		"AUTHOR" VARCHAR(128) NULL DEFAULT NULL,
		"ISBN" VARCHAR(13) UNIQUE NULL DEFAULT NULL
	);`

	sql = strings.ReplaceAll(sql, "#SCHEMA_NAME", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserFieldTable)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateBasicBoard - Create board table
func (d *Postgres) CreateBasicBoard(tableInfo models.Board, recreate bool) error {
	return nil
}

// CreateCustomBoard - Create board table
func (d *Postgres) CreateCustomBoard(tableInfo models.Board, fields []models.Field, recreate bool) error {
	return nil
}

// EditBasicBoard - Create board table
func (d *Postgres) EditBasicBoard(tableInfoOld models.Board, tableInfoNew models.Board) error {
	return nil
}

// EditBasicBoard - Create board table
func (d *Postgres) EditCustomBoard(tableInfoOld models.Board, tableInfoNew models.Board) error {
	return nil
}

// DeleteBoard - Delete a board table
func (d *Postgres) DeleteBoard(tableName string) error {
	return nil
}

// CreateComment - Create comment table
func (d *Postgres) CreateComment(tableInfo models.Board, recreate bool) error {
	return nil
}

// EditUserTableFields - Edit user table schema
func (d *Postgres) EditUserTableFields(fieldsInfoOld []models.UserColumn, fieldsInfoNew []models.UserColumn) error {
	return nil
}

// DeleteUserTableFields - Delete user table field
func (d *Postgres) DeleteUserTableFields(fieldsInfoRemove []models.UserColumn) error {
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
