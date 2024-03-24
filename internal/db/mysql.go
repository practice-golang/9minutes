package db

import (
	"9minutes/consts"
	"9minutes/model"
	"database/sql"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Mysql struct{ dsn string }

func (d *Mysql) connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", d.dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (d *Mysql) CreateDB() error {
	sql := `CREATE DATABASE IF NOT EXISTS ` + Info.DatabaseName + `;`
	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func (d *Mysql) Exec(sql string, colValues []interface{}, options string) (int64, int64, error) {
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
func (d *Mysql) CreateBoardTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS ` + Info.DatabaseName + `.` + Info.BoardTable + ` (
		IDX                       INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
		BOARD_NAME                VARCHAR(128) NULL DEFAULT NULL,
		BOARD_CODE                VARCHAR(64) NULL DEFAULT NULL,
		BOARD_TYPE                VARCHAR(64) NULL DEFAULT NULL,
		` + "`BOARD_TABLE`" + `   VARCHAR(64) NULL DEFAULT NULL,
		` + "`COMMENT_TABLE`" + ` VARCHAR(64) NULL DEFAULT NULL,
		GRANT_READ                VARCHAR(16) NULL DEFAULT NULL,
		GRANT_WRITE               VARCHAR(16) NULL DEFAULT NULL,
		GRANT_COMMENT             VARCHAR(16) NULL DEFAULT NULL,
		GRANT_UPLOAD              VARCHAR(16) NULL DEFAULT NULL,
		FIELDS                    TEXT NULL DEFAULT NULL,

		PRIMARY KEY (IDX),
		INDEX   IDX (IDX)
	)
	COLLATE='utf8_general_ci'
	ENGINE=InnoDB;`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUploadTable - Create upload table
func (d *Mysql) CreateUploadTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS ` + Info.DatabaseName + `.` + Info.UploadTable + ` (
		IDX           INT(11)      UNSIGNED NOT NULL AUTO_INCREMENT,
		TOPIC_IDX     INT(11)      UNSIGNED NOT NULL,
		COMMENT_IDX   INT(11)      UNSIGNED NOT NULL,
		FILE_NAME     VARCHAR(512) NULL DEFAULT NULL,
		STORAGE_NAME  VARCHAR(512) NULL DEFAULT NULL,

		PRIMARY KEY (IDX),
		INDEX   IDX (IDX)
	)
	COLLATE='utf8_general_ci'
	ENGINE=InnoDB;`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserTable - Create user table
func (d *Mysql) CreateUserTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS ` + Info.DatabaseName + `.` + Info.UserTable + ` (
		IDX      INT(11)      UNSIGNED NOT NULL AUTO_INCREMENT,
		USERID VARCHAR(128) NULL DEFAULT NULL,
		PASSWORD VARCHAR(128) NULL DEFAULT NULL,
		EMAIL    VARCHAR(128) NULL DEFAULT NULL,
		GRADE    VARCHAR(24)  NULL DEFAULT NULL,
		APPROVAL VARCHAR(2)   NULL DEFAULT NULL,
		REGDATE VARCHAR(14)  NULL DEFAULT NULL,

		PRIMARY  KEY(IDX),
		UNIQUE   INDEX USERID (USERID),
		UNIQUE   INDEX EMAIL (EMAIL),
		INDEX    IDX (IDX)
	)
	COLLATE='utf8_general_ci'
	ENGINE=InnoDB;`

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
	INSERT IGNORE INTO ` + Info.DatabaseName + `.` + Info.UserTable + ` (
		USERID, ` + "`PASSWORD`" + `, EMAIL, ` + "`GRADE`" + `, APPROVAL, REGDATE
	) VALUES (
		"admin", "` + string(adminPassword) + `", "admin@please.modify", "admin", "Y", "` + now + `"
	);`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	userfieldTable := "user_fields"
	sql = `
	CREATE TABLE IF NOT EXISTS ` + Info.DatabaseName + `.` + userfieldTable + ` (
		IDX          INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
		DISPLAY_NAME VARCHAR(128) NULL DEFAULT NULL,
		COLUMN_CODE  VARCHAR(128) NULL DEFAULT NULL,
		COLUMN_TYPE  VARCHAR(128) NULL DEFAULT NULL,
		COLUMN_NAME  VARCHAR(128) NULL DEFAULT NULL,
		SORT_ORDER   INT(5) UNSIGNED NULL DEFAULT NULL,

		PRIMARY      KEY(IDX),
		UNIQUE       INDEX COLUMN_NAME (COLUMN_NAME),
		INDEX        IDX (IDX)
	)
	COLLATE='utf8_general_ci'
	ENGINE=InnoDB;`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	sql = `
	INSERT IGNORE INTO ` + Info.DatabaseName + `.` + userfieldTable + `
		(DISPLAY_NAME, COLUMN_CODE, COLUMN_TYPE, COLUMN_NAME, SORT_ORDER)
	VALUES
		("Idx", "idx", "integer", "IDX", 1),
		("UserID", "userid", "text", "USERID", 2),
		("Password", "password", "text", "PASSWORD", 3),
		("Email", "email", "text", "EMAIL", 4),
		("Grade", "grade", "text", "GRADE", 5),
		("Approval", "approval", "text", "APPROVAL", 6),
		("RegDate", "regdate", "text", "REGDATE", 7);`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserVerificationTable - Create user verification table
func (d *Mysql) CreateUserVerificationTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS ` + Info.DatabaseName + `.` + Info.UserTable + `_verification` + ` (
		IDX      INT(11)      UNSIGNED NOT NULL AUTO_INCREMENT,
		USER_IDX INT(11)      NULL DEFAULT NULL,
		TOKEN    VARCHAR(128) NULL DEFAULT NULL,
		REGDATE VARCHAR(14)  NULL DEFAULT NULL,

		PRIMARY  KEY(IDX),
		INDEX    IDX (IDX)
	)
	COLLATE='utf8_general_ci'
	ENGINE=InnoDB;`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// AddTableColumn - Add table column
func (d *Mysql) AddTableColumn(tableName string, column model.UserColumn) error {
	sql := "ALTER TABLE " + Info.DatabaseName + `.` + tableName + " ADD COLUMN `" + column.ColumnName.String + "`"

	switch column.ColumnType.String {
	case "text":
		sql += ` VARCHAR(256)`
	case "long_text":
		sql += ` TEXT`
	case "number-integer":
		sql += ` INT(16)`
	case "number-real":
		sql += ` DECIMAL(20,20)`
	}

	sql += " NULL"
	sql += `;`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// EditTableColumn - Edit table column name
func (d *Mysql) EditTableColumn(tableName string, columnOld model.UserColumn, columnNew model.UserColumn) error {
	sql := `
	ALTER TABLE ` + Info.DatabaseName + `.` + tableName + `
	RENAME COLUMN ` + "`" + columnOld.ColumnName.String + "`" + ` TO ` + "`" + columnNew.ColumnName.String + "`" + `; `

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTableColumn - Delete table column
func (d *Mysql) DeleteTableColumn(tableName string, column model.UserColumn) error {
	sql := "ALTER TABLE " + Info.DatabaseName + `.` + tableName + ""
	sql += ` DROP COLUMN ` + column.ColumnName.String + `;`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateBoard - Create list board table
func (d *Mysql) CreateBoard(tableInfo model.Board, recreate bool) error {
	tableName := Info.DatabaseName + "." + tableInfo.BoardTable.String

	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS ` + tableName + `;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS ` + tableName + ` (
		IDX           INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
		TITLE         VARCHAR(256)     NULL DEFAULT NULL,
		TITLE_IMAGE   VARCHAR(256)     NULL DEFAULT NULL,
		CONTENT       TEXT             NULL DEFAULT NULL,
		AUTHOR_IDX    INT(11)          NULL DEFAULT NULL,
		AUTHOR_NAME   VARCHAR(256)     NULL DEFAULT NULL,
		AUTHOR_IP     VARCHAR(50)      NULL DEFAULT NULL,
		AUTHOR_IP_CUT VARCHAR(50)      NULL DEFAULT NULL,
		EDIT_PASSWORD VARCHAR(256)     NULL DEFAULT NULL,
		FILES         TEXT             NULL DEFAULT NULL,
		IMAGES        TEXT             NULL DEFAULT NULL,
		VIEWS         VARCHAR(11)      NULL DEFAULT NULL,
		REGDATE       VARCHAR(14)      NULL DEFAULT NULL,

		PRIMARY KEY(IDX),
		INDEX IDX (IDX)
	);`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateComment - Create comment table
func (d *Mysql) CreateComment(tableInfo model.Board, recreate bool) error {
	tableName := Info.DatabaseName + "." + tableInfo.CommentTable.String

	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS ` + tableName + `;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS ` + tableName + ` (
		IDX           INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
		TOPIC_IDX     INT(11) UNSIGNED NOT NULL,
		CONTENT       TEXT             NULL DEFAULT NULL,
		AUTHOR_IDX    INT(11)          NULL DEFAULT NULL,
		AUTHOR_NAME   VARCHAR(256)     NULL DEFAULT NULL,
		AUTHOR_IP     VARCHAR(50)      NULL DEFAULT NULL,
		AUTHOR_IP_CUT VARCHAR(50)      NULL DEFAULT NULL,
		EDIT_PASSWORD VARCHAR(256)     NULL DEFAULT NULL,
		FILES         TEXT             NULL DEFAULT NULL,
		IMAGES        TEXT             NULL DEFAULT NULL,
		REGDATE       VARCHAR(14)      NULL DEFAULT NULL,

		PRIMARY KEY(IDX),
		INDEX IDX (IDX)
	);`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteBoard - Delete a board table
func (d *Mysql) DeleteBoard(tableInfo model.Board) error {
	tableName := Info.DatabaseName + "." + tableInfo.BoardTable.String

	sql := `DROP TABLE IF EXISTS ` + tableName + `;`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// RenameBoard - Rename board table name
func (d *Mysql) RenameBoard(tableInfoOLD model.Board, tableInfoNEW model.Board) error {
	sql := `
	ALTER TABLE
		` + Info.DatabaseName + "." + tableInfoOLD.BoardTable.String + `
	RENAME TO
		` + Info.DatabaseName + "." + tableInfoNEW.BoardTable.String + `;`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// RenameComment - Rename comment table name
func (d *Mysql) RenameComment(tableInfoOLD model.Board, tableInfoNEW model.Board) error {
	sql := `
	ALTER TABLE
		` + Info.DatabaseName + "." + tableInfoOLD.CommentTable.String + `
	RENAME TO
		` + Info.DatabaseName + "." + tableInfoNEW.CommentTable.String + `;`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteComment - Delete a comment table
func (d *Mysql) DeleteComment(tableInfo model.Board) error {
	tableName := Info.DatabaseName + "." + tableInfo.CommentTable.String

	sql := `DROP TABLE IF EXISTS ` + tableName + `;`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// GetPagingQuery - Get paging query
func (d *Mysql) GetPagingQuery(offset int, listCount int) string {
	sql := `
	LIMIT ` + strconv.Itoa(listCount) + `
	OFFSET ` + strconv.Itoa(offset)

	return sql
}
