package db

import (
	"9minutes/consts"
	"9minutes/model"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Postgres struct{ dsn string }

func (d *Postgres) connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", d.dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CreateDB - Create DB and Schema
func (d *Postgres) CreateDB() error {
	var err error

	sql := `CREATE SCHEMA IF NOT EXISTS ` + Info.SchemaName + `;`
	_, err = Con.Exec(sql)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			msg := "Database does not exist\n" +
				"With Postgres, create database yourself like below SQL query:" +
				"\nCREATE DATABASE " + Info.DatabaseName + ";"
			return errors.New(msg)
		}
		return err
	}

	return nil
}

// CreateTable - Not use
// func (d *Postgres) CreateTable() error {
// 	sql := `
// 	CREATE TABLE IF NOT EXISTS ` + Info.SchemaName + `.` + Info.TableName + ` (
// 		"IDX" SERIAL PRIMARY KEY,
// 		"TITLE" VARCHAR(255) NULL DEFAULT NULL,
// 		"AUTHOR" VARCHAR(255) NULL DEFAULT NULL
// 	)`
// 	// "TITLE" VARCHAR(255) UNIQUE NULL DEFAULT NULL

// 	_, err := Con.Exec(sql)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (d *Postgres) Exec(sql string, colValues []interface{}, options string) (int64, int64, error) {
	var err error
	var count int64 = 0
	var idx int64 = 0

	sql += ` RETURNING "` + options + `";`
	err = Con.QueryRow(sql, colValues...).Scan(&idx)
	if err != nil {
		return count, idx, err
	}

	if idx > 0 {
		count = 1
	}

	return count, idx, nil
}

// CreateBoardTable - Create board manager table
func (d *Postgres) CreateBoardTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS ` + Info.SchemaName + `.` + Info.BoardTable + ` (
		"IDX"           SERIAL       PRIMARY KEY,
		"BOARD_NAME"    VARCHAR(128) NULL DEFAULT NULL,
		"BOARD_CODE"    VARCHAR(64)  NULL DEFAULT NULL,
		"BOARD_TYPE"    VARCHAR(64)  NULL DEFAULT NULL,
		"BOARD_TABLE"   VARCHAR(64)  NULL DEFAULT NULL,
		"COMMENT_TABLE" VARCHAR(64)  NULL DEFAULT NULL,
		"GRANT_READ"    VARCHAR(16)  NULL DEFAULT NULL,
		"GRANT_WRITE"   VARCHAR(16)  NULL DEFAULT NULL,
		"GRANT_COMMENT" VARCHAR(16)  NULL DEFAULT NULL,
		"GRANT_UPLOAD"  VARCHAR(16)  NULL DEFAULT NULL,
		"FIELDS"        TEXT         NULL DEFAULT NULL
	);`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUploadTable - Create upload table
func (d *Postgres) CreateUploadTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS ` + Info.SchemaName + `.` + Info.UploadTable + ` (
		"IDX"             SERIAL        PRIMARY KEY,
		"FILE_NAME"       VARCHAR(512)  NULL DEFAULT NULL,
		"STORAGE_NAME"    VARCHAR(512)  NULL DEFAULT NULL
	);`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserTable - Create user table
func (d *Postgres) CreateUserTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS ` + Info.SchemaName + `.` + Info.UserTable + ` (
		"IDX"      SERIAL       PRIMARY KEY,
		"USERID" VARCHAR(128) UNIQUE NULL DEFAULT NULL,
		"PASSWORD" VARCHAR(128) NULL DEFAULT NULL,
		"EMAIL"    VARCHAR(128) UNIQUE NULL DEFAULT NULL,
		"GRADE"    VARCHAR(24)  NULL DEFAULT NULL,
		"APPROVAL" VARCHAR(2)   NULL DEFAULT NULL,
		"REGDATE" VARCHAR(14)  NULL DEFAULT NULL,

		CONSTRAINT "USERS_UQ" UNIQUE ("USERID", "EMAIL")
	);`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	// Add temp admin
	adminPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), consts.BcryptCost)
	if err != nil {
		return err
	}

	now := time.Now().Format("20060102150405")

	sql = `
	INSERT INTO ` + Info.SchemaName + `.` + Info.UserTable + ` (
		"USERID", "PASSWORD", "EMAIL", "GRADE", "APPROVAL", "REGDATE"
	) VALUES (
		'admin', '` + string(adminPassword) + `', 'admin@please.modify', 'admin', 'Y', '` + now + `'
	)
	ON CONFLICT DO NOTHING;`
	// ON CONFLICT ON CONSTRAINT "USERS_UQ" DO NOTHING;`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	userfieldTable := "user_fields"
	sql = `
	CREATE TABLE IF NOT EXISTS ` + Info.SchemaName + `.` + userfieldTable + ` (
		"IDX"          SERIAL PRIMARY KEY,
		"DISPLAY_NAME" VARCHAR(128) NULL DEFAULT NULL,
		"COLUMN_CODE"  VARCHAR(128) NULL DEFAULT NULL,
		"COLUMN_TYPE"  VARCHAR(128) NULL DEFAULT NULL,
		"COLUMN_NAME"  VARCHAR(128) NULL DEFAULT NULL,
		"SORT_ORDER"   INTEGER NULL DEFAULT NULL,

		CONSTRAINT "COLUMN_NAME_UQ" UNIQUE ("COLUMN_NAME")
	);`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	sql = `
	INSERT INTO ` + Info.SchemaName + `.` + userfieldTable + `
		("DISPLAY_NAME", "COLUMN_CODE", "COLUMN_TYPE", "COLUMN_NAME", "SORT_ORDER")
	VALUES
		('Idx', 'idx', 'integer', 'IDX', 1),
		('UserId', 'userid', 'text', 'USERID', 2),
		('Password', 'password', 'text', 'PASSWORD', 3),
		('Email', 'email', 'text', 'EMAIL', 4),
		('Grade', 'grade', 'text', 'GRADE', 5),
		('Approval', 'approval', 'text', 'APPROVAL', 6),
		('RegDate', 'regdate', 'text', 'REGDATE', 7)
	ON CONFLICT ON CONSTRAINT "COLUMN_NAME_UQ" DO NOTHING;`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserVerificationTable - Create user verification table
func (d *Postgres) CreateUserVerificationTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS ` + Info.SchemaName + `.` + Info.UserTable + ` (
		"IDX"      SERIAL       PRIMARY KEY,
		"USER_IDX" INTEGER      NULL DEFAULT NULL,
		"TOKEN"    VARCHAR(128) NULL DEFAULT NULL,
		"REGDATE" VARCHAR(14)  NULL DEFAULT NULL
	);`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// AddTableColumn - Add table column
func (d *Postgres) AddTableColumn(tableName string, column model.UserColumn) error {
	sql := `ALTER TABLE ` + tableName + ` ADD COLUMN "` + column.ColumnName.String + `" `

	switch column.ColumnType.String {
	case "text":
		sql += ` VARCHAR(128)`
	case "long_text":
		sql += ` VARCHAR(65535)`
	case "number-integer":
		sql += ` INTEGER`
	case "number-real":
		sql += ` REAL`
	}

	sql += `; `

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// EditTableColumn - Edit table column name
func (d *Postgres) EditTableColumn(tableName string, columnOld model.UserColumn, columnNew model.UserColumn) error {
	sql := `
	ALTER TABLE ` + tableName + `
	RENAME COLUMN "` + columnOld.ColumnName.String + `" TO "` + columnNew.ColumnName.String + `"; `

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTableColumn - Delete table column
func (d *Postgres) DeleteTableColumn(tableName string, column model.UserColumn) error {
	sql := `ALTER TABLE ` + Info.UserTable + ` `
	sql += ` DROP COLUMN "` + column.ColumnName.String + `";`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateBoard - Create list board table
func (d *Postgres) CreateBoard(tableInfo model.Board, recreate bool) error {
	tableName := Info.SchemaName + "." + tableInfo.BoardTable.String

	sql := ``
	if recreate {
		sql += `DROP TABLE IF EXISTS ` + tableName + `;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS ` + tableName + ` (
		"IDX"          BIGSERIAL    PRIMARY KEY,
		"TITLE"        VARCHAR(256) NULL DEFAULT NULL,
		"TITLE_IMAGE"  VARCHAR(256) NULL DEFAULT NULL,
		"CONTENT"      TEXT         NULL DEFAULT NULL,
		"AUTHOR_IDX"   BIGINT       NULL DEFAULT NULL,
		"AUTHOR_NAME"  VARCHAR(256) NULL DEFAULT NULL,
		"FILES"        TEXT         NULL DEFAULT NULL,
		"IMAGES"       TEXT         NULL DEFAULT NULL,
		"VIEWS"        VARCHAR(11)  NULL DEFAULT NULL,
		"REGDATE"      VARCHAR(14)  NULL DEFAULT NULL
	);`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateComment - Create comment table
func (d *Postgres) CreateComment(tableInfo model.Board, recreate bool) error {
	tableName := Info.SchemaName + "." + tableInfo.CommentTable.String

	sql := ``
	if recreate {
		sql += `DROP TABLE IF EXISTS ` + tableName + `;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS ` + tableName + ` (
		"IDX"         BIGSERIAL    PRIMARY KEY,
		"POSTING_IDX" BIGINT,
		"CONTENT"     TEXT         NULL DEFAULT NULL,
		"AUTHOR_IDX"  BIGINT       NULL DEFAULT NULL,
		"AUTHOR_NAME" VARCHAR(256) NULL DEFAULT NULL,
		"FILES"       TEXT         NULL DEFAULT NULL,
		"IMAGES"      TEXT         NULL DEFAULT NULL,
		"REGDATE"     VARCHAR(14)  NULL DEFAULT NULL
	);`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// RenameBoard - Rename board table name
func (d *Postgres) RenameBoard(tableInfoOld model.Board, tableInfoNew model.Board) error {
	sql := `
	ALTER TABLE
		"` + Info.DatabaseName + `"."` + tableInfoOld.BoardTable.String + `"
	RENAME TO
		"` + tableInfoNew.BoardTable.String + `";`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// RenameComment - Rename comment table name
func (d *Postgres) RenameComment(tableInfoOld model.Board, tableInfoNew model.Board) error {
	sql := `
	ALTER TABLE
		"` + Info.DatabaseName + `"."` + tableInfoOld.CommentTable.String + `"
	RENAME TO
		"` + tableInfoNew.CommentTable.String + `";`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteBoard - Delete a board table
func (d *Postgres) DeleteBoard(tableInfo model.Board) error {
	tableName := Info.SchemaName + "." + tableInfo.BoardTable.String

	sql := `DROP TABLE IF EXISTS ` + tableName + `;`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteComment - Delete comment table
func (d *Postgres) DeleteComment(tableInfo model.Board) error {
	tableName := Info.SchemaName + "." + tableInfo.CommentTable.String

	sql := `DROP TABLE IF EXISTS ` + tableName + `;`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// GetPagingQuery - Get paging query
func (d *Postgres) GetPagingQuery(offset int, listCount int) string {
	sql := `
	OFFSET ` + strconv.Itoa(offset) + `
	LIMIT ` + strconv.Itoa(listCount)

	return sql
}
