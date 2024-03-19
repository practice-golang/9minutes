package db

import (
	"9minutes/consts"
	"9minutes/model"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type SqlServer struct{ dsn string }

func (d *SqlServer) connect() (*sql.DB, error) {
	db, err := sql.Open("sqlserver", d.dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (d *SqlServer) CreateDB() error {
	sql := `
	USE master
	IF NOT EXISTS(
		SELECT name
		FROM sys.databases
		WHERE name=N'` + Info.DatabaseName + `'
	) CREATE DATABASE [` + Info.DatabaseName + `]
	`
	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateTable - Not use
// func (d *SqlServer) CreateTable() error {
// 	tableName := "books"

// 	sql := `
// 	USE "` + Info.DatabaseName + `"
// 	IF OBJECT_ID(N'` + tableName + `', N'U') IS NULL
// 	CREATE TABLE ` + tableName + ` (
// 		IDX INT NOT NULL IDENTITY PRIMARY KEY,
// 		TITLE VARCHAR(256) NULL,
// 		AUTHOR VARCHAR(256) NULL,
// 	)`

// 	_, err := Con.Exec(sql)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (d *SqlServer) Exec(sql string, colValues []interface{}, options string) (int64, int64, error) {
	var err error
	var count int64 = 0
	var idx int64 = 0

	sql += ` SELECT ID = CONVERT(bigint, ISNULL(SCOPE_IDENTITY(), -1)) `

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
func (d *SqlServer) CreateBoardTable() error {
	sql := `USE "` + Info.DatabaseName + `"`
	sql += `
	IF OBJECT_ID(N'` + Info.BoardTable + `', N'U') IS NULL
	CREATE TABLE "` + Info.BoardTable + `" (
		IDX           BIGINT NOT NULL IDENTITY PRIMARY KEY,
		BOARD_NAME    VARCHAR(128) NULL DEFAULT NULL,
		BOARD_CODE    VARCHAR(64) NULL DEFAULT NULL,
		BOARD_TYPE    VARCHAR(64) NULL DEFAULT NULL,
		BOARD_TABLE   VARCHAR(64) NULL DEFAULT NULL,
		COMMENT_TABLE VARCHAR(64) NULL DEFAULT NULL,
		GRANT_READ    VARCHAR(16) NULL DEFAULT NULL,
		GRANT_WRITE   VARCHAR(16) NULL DEFAULT NULL,
		GRANT_COMMENT VARCHAR(16) NULL DEFAULT NULL,
		GRANT_UPLOAD  VARCHAR(16) NULL DEFAULT NULL,
		FIELDS        TEXT NULL DEFAULT NULL,
	)`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUploadTable - Create upload table
func (d *SqlServer) CreateUploadTable() error {
	sql := `USE "` + Info.DatabaseName + `"`
	sql += `
	IF OBJECT_ID(N'` + Info.UploadTable + `', N'U') IS NULL
	CREATE TABLE "` + Info.UploadTable + `" (
		IDX             BIGINT       NOT NULL IDENTITY PRIMARY KEY,
		FILE_NAME       VARCHAR(512) NULL DEFAULT NULL,
		STORAGE_NAME    VARCHAR(512) NULL DEFAULT NULL
	)`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserTable - Create user table
func (d *SqlServer) CreateUserTable() error {
	// sql := `
	// USE master

	// IF NOT EXISTS(
	// 	SELECT name
	// 	FROM sys.databases
	// 	WHERE name=N'` + Info.DatabaseName + `'
	// )
	// CREATE DATABASE "` + Info.DatabaseName + `"`

	// _, err := Con.Exec(sql)
	// if err != nil {
	// 	return err
	// }

	sql := `USE "` + Info.DatabaseName + `"`
	sql += `
	IF OBJECT_ID(N'` + Info.UserTable + `', N'U') IS NULL
	CREATE TABLE "` + Info.UserTable + `" (
		IDX      INT       NOT NULL IDENTITY PRIMARY KEY,
		USERID VARCHAR(128) UNIQUE NULL DEFAULT NULL,
		PASSWORD VARCHAR(128) NULL DEFAULT NULL,
		EMAIL    VARCHAR(128) UNIQUE NULL DEFAULT NULL,
		GRADE    VARCHAR(24)  NULL DEFAULT NULL,
		APPROVAL VARCHAR(2)   NULL DEFAULT NULL,
		REGDATE VARCHAR(14)  NULL DEFAULT NULL,
	)`

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
	USE "` + Info.DatabaseName + `"
	IF NOT EXISTS (SELECT TOP 1 * FROM "` + Info.UserTable + `" WHERE USERID = 'admin')
	INSERT INTO "` + Info.UserTable + `" (
		USERID, "PASSWORD", EMAIL, "GRADE", APPROVAL, REGDATE
	) VALUES (
		'admin', '` + string(adminPassword) + `', 'admin@please.modify', 'admin', 'Y', '` + now + `'
	)`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	userfieldTable := "user_fields"
	sql = `USE "` + Info.DatabaseName + `"`
	sql += `
	IF OBJECT_ID(N'` + userfieldTable + `', N'U') IS NULL
	CREATE TABLE "` + userfieldTable + `" (
		IDX          INT NOT NULL IDENTITY PRIMARY KEY,
		DISPLAY_NAME VARCHAR(128) NULL DEFAULT NULL,
		COLUMN_CODE  VARCHAR(128) NULL DEFAULT NULL,
		COLUMN_TYPE  VARCHAR(128) NULL DEFAULT NULL,
		COLUMN_NAME  VARCHAR(128) UNIQUE NULL DEFAULT NULL,
		SORT_ORDER   INT NULL DEFAULT NULL,
	)`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	sql = `
	USE "` + Info.DatabaseName + `"
	INSERT INTO "` + Info.DatabaseName + `".dbo.` + userfieldTable + `
		(DISPLAY_NAME, COLUMN_CODE, COLUMN_TYPE, COLUMN_NAME, SORT_ORDER)
	VALUES
		('Idx', 'idx', 'integer', 'IDX', 1),
		('UserId', 'userid', 'text', 'USERID', 2),
		('Password', 'password', 'text', 'PASSWORD', 3),
		('Email', 'email', 'text', 'EMAIL', 4),
		('Grade', 'grade', 'text', 'GRADE', 5),
		('Approval', 'approval', 'text', 'APPROVAL', 6),
		('RegDate', 'regdate', 'text', 'REGDATE', 7)`

	_, err = Con.Exec(sql)
	if err != nil {
		if !strings.Contains(err.Error(), "mssql: UNIQUE KEY") {
			return err
		}
	}

	return nil
}

// CreateUserVerificationTable - Create user verification table
func (d *SqlServer) CreateUserVerificationTable() error {
	sql := `USE "` + Info.DatabaseName + `"`
	sql += `
	IF OBJECT_ID(N'` + Info.UserTable + `', N'U') IS NULL
	CREATE TABLE "` + Info.UserTable + `" (
		IDX      INT          NOT NULL IDENTITY PRIMARY KEY,
		USER_IDX INT          NULL DEFAULT NULL,
		TOKEN    VARCHAR(128) NULL DEFAULT NULL,
		REGDATE VARCHAR(14)  NULL DEFAULT NULL,
	)`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// AddTableColumn - Add table column
func (d *SqlServer) AddTableColumn(tableName string, column model.UserColumn) error {
	sql := `
	USE "` + Info.DatabaseName + `"
	ALTER TABLE "` + Info.UserTable + `"
		ADD`

	sql += "\n" + column.ColumnName.String + ` `

	switch column.ColumnType.String {
	case "text":
		sql += `VARCHAR(128)`
	case "long_text":
		sql += `VARCHAR(MAX)`
	case "number-integer":
		sql += `INTEGER`
	case "number-real":
		sql += `REAL`
	}

	sql += `;`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// EditTableColumn - Edit table column name
func (d *SqlServer) EditTableColumn(tableName string, columnOld model.UserColumn, columnNew model.UserColumn) error {
	sql := `
	USE "` + Info.DatabaseName + `"
	EXEC sp_rename 'dbo.` + tableName + `.` + columnOld.ColumnName.String + `', '` + columnNew.ColumnName.String + `', 'COLUMN';`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTableColumn - Delete table column
func (d *SqlServer) DeleteTableColumn(tableName string, column model.UserColumn) error {
	sql := `
	USE "` + Info.DatabaseName + `"
	ALTER TABLE "` + Info.UserTable + `"`

	sql += `DROP COLUMN ` + column.ColumnName.String + `;`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateBoard - Create list board table
func (d *SqlServer) CreateBoard(tableInfo model.Board, recreate bool) error {
	sql := ``
	if recreate {
		sql = `
		USE "` + Info.DatabaseName + `"
		IF OBJECT_ID('` + tableInfo.BoardTable.String + `','U') IS NOT NULL
		DROP TABLE "` + tableInfo.BoardTable.String + `"
		-- GO`

		_, err := Con.Exec(sql)
		if err != nil {
			return err
		}
	}

	sql = `
	USE "` + Info.DatabaseName + `"
	IF OBJECT_ID(N'` + tableInfo.BoardTable.String + `', N'U') IS NULL
	CREATE TABLE "` + tableInfo.BoardTable.String + `" (
		IDX           BIGINT       NOT NULL IDENTITY PRIMARY KEY,
		TITLE         VARCHAR(256) NULL DEFAULT NULL,
		TITLE_IMAGE   VARCHAR(256) NULL DEFAULT NULL,
		CONTENT       TEXT         NULL DEFAULT NULL,
		AUTHOR_IDX    BIGINT       NULL DEFAULT NULL,
		AUTHOR_NAME   VARCHAR(256) NULL DEFAULT NULL,
		AUTHOR_IP     VARCHAR(50)  NULL DEFAULT NULL,
		AUTHOR_IP_CUT VARCHAR(50)  NULL DEFAULT NULL,
		EDIT_PASSWORD VARCHAR(256) NULL DEFAULT NULL,
		FILES         TEXT         NULL DEFAULT NULL,
		IMAGES        TEXT         NULL DEFAULT NULL,
		VIEWS         VARCHAR(11)  NULL DEFAULT NULL,
		REGDATE       VARCHAR(14)  NULL DEFAULT NULL,
	)
	--GO`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateComment - Create comment table
func (d *SqlServer) CreateComment(tableInfo model.Board, recreate bool) error {
	sql := ``
	if recreate {
		sql = `
		USE "` + Info.DatabaseName + `"
		IF OBJECT_ID('` + tableInfo.CommentTable.String + `','U') IS NOT NULL
		DROP TABLE "` + tableInfo.CommentTable.String + `"
		-- GO`

		_, err := Con.Exec(sql)
		if err != nil {
			return err
		}
	}

	sql = `
	USE "` + Info.DatabaseName + `"
	IF OBJECT_ID(N'` + tableInfo.CommentTable.String + `', N'U') IS NULL
	CREATE TABLE "` + tableInfo.CommentTable.String + `" (
		IDX           BIGINT       NOT NULL IDENTITY PRIMARY KEY,
		TOPIC_IDX     BIGINT       NOT NULL,
		TITLE         VARCHAR(256) NULL DEFAULT NULL,
		TITLE_IMAGE   VARCHAR(256) NULL DEFAULT NULL,
		CONTENT       TEXT         NULL DEFAULT NULL,
		AUTHOR_IDX    BIGINT       NULL DEFAULT NULL,
		AUTHOR_NAME   VARCHAR(256) NULL DEFAULT NULL,
		AUTHOR_IP     VARCHAR(50)  NULL DEFAULT NULL,
		AUTHOR_IP_CUT VARCHAR(50)  NULL DEFAULT NULL,
		EDIT_PASSWORD VARCHAR(256) NULL DEFAULT NULL,
		FILES         TEXT         NULL DEFAULT NULL,
		IMAGES        TEXT         NULL DEFAULT NULL,
		VIEWS         VARCHAR(11)  NULL DEFAULT NULL,
		REGDATE       VARCHAR(14)  NULL DEFAULT NULL,
	)
	--GO`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// RenameBoard - Rename board table name
func (d *SqlServer) RenameBoard(tableInfoOld model.Board, tableInfoNew model.Board) error {
	sql := `
	USE "` + Info.DatabaseName + `"
	EXEC sp_rename "` + tableInfoOld.BoardTable.String + `", "` + tableInfoNew.BoardTable.String + `"`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// RenameComment - Rename comment table name
func (d *SqlServer) RenameComment(tableInfoOld model.Board, tableInfoNew model.Board) error {
	sql := `
	USE "` + Info.DatabaseName + `"
	EXEC sp_rename "` + tableInfoOld.CommentTable.String + `", "` + tableInfoNew.CommentTable.String + `"`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteBoard - Delete a board table
func (d *SqlServer) DeleteBoard(tableInfo model.Board) error {
	tableName := tableInfo.BoardTable.String

	sql := `USE "` + Info.DatabaseName + `"
	IF OBJECT_ID('` + tableName + `','U') IS NOT NULL
	DROP TABLE "` + tableName + `"
	-- GO`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteComment - Delete a comment table
func (d *SqlServer) DeleteComment(tableInfo model.Board) error {
	tableName := Info.DatabaseName + ".dbo." + tableInfo.BoardTable.String

	sql := `
	IF OBJECT_ID(N'` + tableName + `', N'U') IS NOT NULL
	DROP TABLE ` + tableName + `;`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// GetPagingQuery - Get paging query
func (d *SqlServer) GetPagingQuery(offset int, listCount int) string {
	sql := `
	OFFSET ` + strconv.Itoa(offset) + ` ROWS
	FETCH NEXT ` + strconv.Itoa(listCount) + ` ROWS ONLY`

	return sql
}
