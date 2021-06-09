package db

import (
	"database/sql"
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

// CreateTable - Create table
func (d *Postgres) CreateTable(recreate bool) error {
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
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", TableName)

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
