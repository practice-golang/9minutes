package db

import (
	"9minutes/consts"
	"9minutes/model"
	"database/sql"
	"log"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type SQLite struct{ dsn string }

func (d *SQLite) connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite", d.dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CreateDB - enough to use connect so, not use
func (d *SQLite) CreateDB() error { return nil }

// CreateTable - Not use
// func (d *SQLite) CreateTable() error {
// 	sql := `
// 	CREATE TABLE IF NOT EXISTS "` + Info.TableName + `" (
// 		"IDX"			INTEGER,
// 		"TITLE"			TEXT,
// 		"AUTHOR"		TEXT,
// 		PRIMARY KEY("IDX" AUTOINCREMENT)
// 	);`

// 	_, err := Con.Exec(sql)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (d *SQLite) Exec(sql string, colValues []interface{}, options string) (int64, int64, error) {
	var err error
	var count int64 = 0
	var idx int64 = 0

	result, err := Con.Exec(sql, colValues...)
	if err != nil {
		return count, idx, err
	}

	count, _ = result.RowsAffected()
	idx, _ = result.LastInsertId()

	return count, idx, nil
}

// CreateBoardTable - Create board manager table
func (d *SQLite) CreateBoardTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS "` + Info.BoardTable + `" (
		"IDX"           INTEGER,
		"BOARD_NAME"    TEXT,
		"BOARD_CODE"    TEXT,
		"BOARD_TYPE"    TEXT,
		"BOARD_TABLE"   TEXT UNIQUE,
		"COMMENT_TABLE" TEXT UNIQUE,
		"GRANT_READ"    TEXT,
		"GRANT_WRITE"   TEXT,
		"GRANT_COMMENT" TEXT,
		"GRANT_UPLOAD"	TEXT,
		"FIELDS"        TEXT,

		PRIMARY KEY("IDX" AUTOINCREMENT)
	);`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserTable - Create user table
func (d *SQLite) CreateUserTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS "` + Info.UserTable + `" (
		"IDX"      INTEGER,
		"USERNAME" TEXT UNIQUE,
		"PASSWORD" TEXT,
		"EMAIL"    TEXT UNIQUE,
		"GRADE"    TEXT,
		"APPROVAL" TEXT,
		"REG_DTTM" TEXT,

		PRIMARY KEY("IDX" AUTOINCREMENT)
	);`

	// Add temp admin
	adminPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), consts.BcryptCost)
	if err != nil {
		return err
	}

	now := time.Now().Format("20060102150405")

	sql += `
	INSERT OR IGNORE INTO "` + Info.UserTable + `" (
		USERNAME, PASSWORD, EMAIL, GRADE, APPROVAL, REG_DTTM
	) VALUES (
		"admin", "` + string(adminPassword) + `", "admin@please.modify", "admin", "Y", "` + now + `"
	)`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	userfieldTable := "user_fields"
	sql = `
	CREATE TABLE IF NOT EXISTS "` + userfieldTable + `" (
		"IDX"          INTEGER,
		"DISPLAY_NAME" TEXT,
		"COLUMN_CODE"  TEXT,
		"COLUMN_TYPE"  TEXT,
		"COLUMN_NAME"  TEXT UNIQUE,
		"SORT_ORDER"   INTEGER,

		PRIMARY KEY("IDX" AUTOINCREMENT)
	);`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	sql = `
	INSERT OR IGNORE INTO ` + userfieldTable + `
		(DISPLAY_NAME, COLUMN_CODE, COLUMN_TYPE, COLUMN_NAME, SORT_ORDER)
	VALUES
		("Idx", "idx", "integer", "IDX", 1),
		("Username", "username", "text", "USERNAME", 2),
		("Password", "password", "text", "PASSWORD", 3),
		("Email", "email", "text", "EMAIL", 4),
		("Grade", "grade", "text", "GRADE", 5),
		("Approval", "approval", "text", "APPROVAL", 6),
		("Registered datetime", "regdate", "text", "REG_DTTM", 7);`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserVerificationTable - Create user validation table
func (d *SQLite) CreateUserVerificationTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS "` + Info.UserTable + `_VERIFICATION` + `" (
		"IDX"      INTEGER,
		"USER_IDX" INTEGER UNIQUE,
		"REG_DTTM" TEXT,

		PRIMARY KEY("IDX" AUTOINCREMENT)
	);`

	log.Println(sql)

	return nil
}

// AddTableColumn - Add table column
func (d *SQLite) AddTableColumn(tableName string, column model.UserColumn) error {
	sql := `ALTER TABLE "` + tableName + `" ADD COLUMN ` + column.ColumnName.String + ` `

	switch column.ColumnType.String {
	case "text":
		sql += ` TEXT`
	case "long_text":
		sql += ` TEXT`
	case "number":
		sql += ` INTEGER`
	case "real":
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
func (d *SQLite) EditTableColumn(tableName string, columnOld model.UserColumn, columnNew model.UserColumn) error {
	sql := `
	ALTER TABLE "` + tableName + `"
	RENAME COLUMN "` + columnOld.ColumnName.String + `" TO "` + columnNew.ColumnName.String + `"; `

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTableColumn - Delete table column
func (d *SQLite) DeleteTableColumn(tableName string, column model.UserColumn) error {
	sql := `ALTER TABLE "` + tableName + `" DROP COLUMN "` + column.ColumnName.String + `";`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateBoard - Create list board table
func (d *SQLite) CreateBoard(tableInfo model.Board, recreate bool) error {
	tableName := tableInfo.BoardTable.String

	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS "` + tableName + `";`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS "` + tableName + `" (
		"IDX"         INTEGER,
		"TITLE"       TEXT,
		"TITLE_IMAGE" TEXT,
		"CONTENT"     TEXT,
		"AUTHOR_IDX"  INTEGER,
		"FILES"       TEXT,
		"VIEWS"       TEXT,
		"REG_DTTM"    TEXT,
		PRIMARY KEY("IDX" AUTOINCREMENT)
	);`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateComment - Create comment table
func (d *SQLite) CreateComment(tableInfo model.Board, recreate bool) error {
	tableName := tableInfo.CommentTable.String

	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS "` + tableName + `";`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS "` + tableName + `" (
		"IDX"         INTEGER,
		"BOARD_IDX"   INTEGER,
		"CONTENT"     TEXT,
		"AUTHOR_IDX"  INTEGER,
		"FILES"       TEXT,
		"REG_DTTM"    TEXT,
		PRIMARY KEY("IDX" AUTOINCREMENT)
	);`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// RenameBoard - Rename board table name
func (d *SQLite) RenameBoard(tableInfoOLD model.Board, tableInfoNEW model.Board) error {
	sql := `
	ALTER TABLE
		"` + tableInfoOLD.BoardTable.String + `"
	RENAME TO
		"` + tableInfoNEW.BoardTable.String + `";`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// RenameComment - Rename comment table name
func (d *SQLite) RenameComment(tableInfoOLD model.Board, tableInfoNEW model.Board) error {
	sql := `
	ALTER TABLE
		"` + tableInfoOLD.CommentTable.String + `"
	RENAME TO
		"` + tableInfoNEW.CommentTable.String + `";`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteBoard - Delete a board table
func (d *SQLite) DeleteBoard(tableInfo model.Board) error {
	sql := `DROP TABLE IF EXISTS "` + tableInfo.BoardTable.String + `";`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteComment - Delete a comment table
func (d *SQLite) DeleteComment(tableInfo model.Board) error {
	sql := `DROP TABLE IF EXISTS "` + tableInfo.CommentTable.String + `";`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// GetPagingQuery - Get paging query
func (d *SQLite) GetPagingQuery(offset int, listCount int) string {
	sql := `
	LIMIT ` + strconv.Itoa(listCount) + `
	OFFSET ` + strconv.Itoa(offset)

	return sql
}
