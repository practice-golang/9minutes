package db

import (
	"9minutes/consts"
	"9minutes/model"
	"database/sql"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	go_ora "github.com/sijms/go-ora/v2"
	"golang.org/x/crypto/bcrypt"
)

type Oracle struct {
	dsn     string
	Version int64
}

func (d *Oracle) createAccount() {
	var err error

	tableSpace := "USERS"
	port, _ := strconv.Atoi(InfoOracleAdmin.Port)
	dsn := go_ora.BuildUrl(InfoOracleAdmin.Addr, port, InfoOracleAdmin.DatabaseName, InfoOracleAdmin.GrantID, InfoOracleAdmin.GrantPassword, nil)
	if InfoOracleAdmin.FilePath != "" {
		tableSpace = "DATA"
		dsn += "?SSL=enable&SSL Verify=false&WALLET=" + url.QueryEscape(InfoOracleAdmin.FilePath)
	}

	conn, err := sql.Open("oracle", dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	versionSTR := ""
	sql := `SELECT version FROM V$INSTANCE`
	err = conn.QueryRow(sql).Scan(&versionSTR)
	if err != nil {
		panic(err)
	}

	d.Version, _ = strconv.ParseInt(strings.Split(versionSTR, ".")[0], 10, 64)

	if d.Version < 12 {
		panic("oracle version is lower than 12")
	}

	sql = `
	SELECT COUNT(USERID) AS COUNT
	FROM ALL_USERS
	WHERE USERID = '` + strings.ToUpper(Info.GrantID) + `'`

	var count int64
	_ = conn.QueryRow(sql).Scan(&count)
	if count > 0 {
		return
	}

	sql = `CREATE USER "` + strings.ToUpper(Info.GrantID) + `" IDENTIFIED BY "` + Info.GrantPassword + `"`
	if InfoOracleAdmin.FilePath != "" {
		sql += `
		DEFAULT TABLESPACE ` + tableSpace + `
		TEMPORARY TABLESPACE TEMP`
	}
	_, err = conn.Exec(sql)
	if err != nil {
		panic(err)
	}

	sql = `GRANT CONNECT, RESOURCE TO ` + Info.GrantID
	_, err = conn.Exec(sql)
	if err != nil {
		panic(err)
	}

	sql = `ALTER USER ` + Info.GrantID + ` DEFAULT TABLESPACE ` + tableSpace + ` QUOTA UNLIMITED ON ` + tableSpace
	_, err = conn.Exec(sql)
	if err != nil {
		panic(err)
	}
}

func (d *Oracle) connect() (*sql.DB, error) {
	db, err := sql.Open("oracle", d.dsn)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			time.Sleep(time.Minute * 1)
			err = db.Ping()
			if err != nil {
				log.Println("Ping:", err)
				db.Close()
				db, err = sql.Open("oracle", d.dsn)
				if err != nil {
					log.Println("Oracle reconnect:", err)
				}
				Con = db
			}
		}
	}()

	return db, nil
}

func (d *Oracle) CreateDB() error {
	err := Con.Ping()
	if err != nil {
		if strings.Contains(err.Error(), "ORA-01017") {
			d.createAccount()
			return nil
		}
	}

	return err
}

func (d *Oracle) Exec(sql string, colValues []interface{}, options string) (int64, int64, error) {
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
func (d *Oracle) CreateBoardTable() error {
	var err error
	var count int64

	boardTable := strings.ToUpper(`"` + Info.GrantID + `"."` + Info.BoardTable + `"`)

	sql := `
	SELECT COUNT(TABLE_NAME) AS COUNT
	FROM user_tables
	WHERE TABLE_NAME = '` + strings.ToUpper(Info.BoardTable) + `'`

	err = Con.QueryRow(sql).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	sql = `
	CREATE TABLE ` + boardTable + ` (
		IDX              NUMBER(11) GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1) NOT NULL,
		BOARD_NAME       VARCHAR(128),
		BOARD_CODE       VARCHAR(64),
		BOARD_TYPE       VARCHAR(64),
		BOARD_TABLE      VARCHAR(64),
		COMMENT_TABLE    VARCHAR(64),
		GRANT_READ       VARCHAR(16),
		GRANT_WRITE      VARCHAR(16),
		GRANT_COMMENT    VARCHAR(16),
		GRANT_UPLOAD     VARCHAR(16),
		FIELDS           NCLOB,

		UNIQUE("IDX")
	)`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUploadTable - Create upload table
func (d *Oracle) CreateUploadTable() error {
	var err error
	var count int64

	uploadTable := strings.ToUpper(`"` + Info.GrantID + `"."` + Info.UploadTable + `"`)

	sql := `
	SELECT COUNT(TABLE_NAME) AS COUNT
	FROM user_tables
	WHERE TABLE_NAME = '` + strings.ToUpper(Info.UploadTable) + `'`

	err = Con.QueryRow(sql).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	sql = `
	CREATE TABLE ` + uploadTable + ` (
		IDX             NUMBER(11) GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1) NOT NULL,
		FILE_NAME       VARCHAR(512),
		STORAGE_NAME    VARCHAR(512),
		BOARD_IDX       NUMBER(11),
		POST_IDX        NUMBER(11),

		UNIQUE("IDX")
	)`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserTable - Create user table
func (d *Oracle) CreateUserTable() error {
	var err error
	var count int64

	userTable := strings.ToUpper(`"` + Info.GrantID + `"."` + Info.UserTable + `"`)

	sql := `
	SELECT COUNT(TABLE_NAME) AS COUNT
	FROM user_tables
	WHERE TABLE_NAME = '` + strings.ToUpper(Info.UserTable) + `'`

	err = Con.QueryRow(sql).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	sql = `
	CREATE TABLE ` + userTable + ` (
		IDX         NUMBER(11) GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1) NOT NULL,
		USERID    VARCHAR(128),
		PASSWORD    VARCHAR(128),
		EMAIL       VARCHAR(128),
		GRADE       VARCHAR(24),
		APPROVAL    VARCHAR(2),
		REGDATE    VARCHAR(14),

		CONSTRAINT "` + Info.UserTable + `_idx" PRIMARY KEY ("IDX"),
		CONSTRAINT "` + Info.UserTable + `_userconstraint" UNIQUE ("USERID", "EMAIL")
	)`

	_, err = Con.Exec(sql)
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
	INSERT INTO ` + userTable + ` (
		"USERID", "PASSWORD", "EMAIL", "GRADE", "APPROVAL", "REGDATE"
	)
	SELECT 'admin', '` + string(adminPassword) + `', 'admin@please.modify', 'admin', 'Y', '` + now + `'
	FROM DUAL
	WHERE NOT EXISTS (SELECT * FROM ` + userTable + ` WHERE "USERID" = 'admin')`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	userfieldTable := strings.ToUpper(`"` + Info.GrantID + `"."` + strings.ToUpper(`user_fields"`))

	sql = `
	CREATE TABLE ` + userfieldTable + ` (
		IDX          NUMBER(11) GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1) NOT NULL,
		DISPLAY_NAME VARCHAR(128),
		COLUMN_CODE  VARCHAR(128),
		COLUMN_TYPE  VARCHAR(128),
		COLUMN_NAME  VARCHAR(128),
		SORT_ORDER   NUMBER(5),

		CONSTRAINT "user_fields_idx" PRIMARY KEY ("IDX"),
		CONSTRAINT "user_fields_columnconstraint" UNIQUE ("COLUMN_NAME")
	)`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	sql = `
	INSERT INTO ` + userfieldTable + ` ("DISPLAY_NAME", "COLUMN_CODE", "COLUMN_TYPE", "COLUMN_NAME", "SORT_ORDER")
	SELECT 'Idx' AS "DISPLAY_NAME", 'idx' AS "COLUMN_CODE", 'integer' AS "COLUMN_TYPE", 'IDX' AS "COLUMN_NAME", 1 AS "SORT_ORDER" FROM DUAL UNION ALL
	SELECT 'UserId' AS "DISPLAY_NAME", 'userid' AS "COLUMN_CODE", 'text' AS "COLUMN_TYPE", 'USERID' AS "COLUMN_NAME", 2 AS "SORT_ORDER" FROM DUAL UNION ALL
	SELECT 'Password' AS "DISPLAY_NAME", 'password' AS "COLUMN_CODE", 'text' AS "COLUMN_TYPE", 'PASSWORD' AS "COLUMN_NAME", 3 AS "SORT_ORDER" FROM DUAL UNION ALL
	SELECT 'Email' AS "DISPLAY_NAME", 'email' AS "COLUMN_CODE", 'text' AS "COLUMN_TYPE", 'EMAIL' AS "COLUMN_NAME", 4 AS "SORT_ORDER" FROM DUAL UNION ALL
	SELECT 'Grade' AS "DISPLAY_NAME", 'grade' AS "COLUMN_CODE", 'text' AS "COLUMN_TYPE", 'GRADE' AS "COLUMN_NAME", 5 AS "SORT_ORDER" FROM DUAL UNION ALL
	SELECT 'Approval' AS "DISPLAY_NAME", 'approval' AS "COLUMN_CODE", 'text' AS "COLUMN_TYPE", 'APPROVAL' AS "COLUMN_NAME", 6 AS "SORT_ORDER" FROM DUAL UNION ALL
	SELECT 'RegDate' AS "DISPLAY_NAME", 'regdate' AS "COLUMN_CODE", 'text' AS "COLUMN_TYPE", 'REGDATE' AS "COLUMN_NAME", 7 AS "SORT_ORDER" FROM DUAL`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserVerificationTable - Create user verification table
func (d *Oracle) CreateUserVerificationTable() error {
	var err error
	var count int64

	verificationTable := strings.ToUpper(`"` + Info.GrantID + `"."` + Info.UserTable + `_verification"`)

	sql := `
	SELECT COUNT(TABLE_NAME) AS COUNT
	FROM user_tables
	WHERE TABLE_NAME = '` + strings.ToUpper(Info.UserTable+`_verification`) + `'`

	err = Con.QueryRow(sql).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	sql = `
	CREATE TABLE ` + verificationTable + ` (
		IDX         NUMBER(11) GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1) NOT NULL,
		USER_IDX    NUMBER(11),
		TOKEN       VARCHAR(128),
		REGDATE    VARCHAR(14),

		UNIQUE("IDX")
	)`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// AddTableColumn - Add table column
func (d *Oracle) AddTableColumn(tableName string, column model.UserColumn) error {
	targetTable := strings.ToUpper(`"` + Info.GrantID + `"."` + tableName + `"`)

	sql := `ALTER TABLE ` + targetTable + ` ADD "` + column.ColumnName.String + `"`

	switch column.ColumnType.String {
	case "text":
		sql += ` VARCHAR(256)`
	case "long_text":
		sql += ` NCLOB`
	case "number-integer":
		sql += ` NUMBER(16)`
	case "number-real":
		sql += ` NUMBER(20,20)`
	}

	sql += " NULL"

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// EditTableColumn - Edit table column name
func (d *Oracle) EditTableColumn(tableName string, columnOld model.UserColumn, columnNew model.UserColumn) error {
	targetTable := strings.ToUpper(`"` + Info.GrantID + `"."` + tableName + `"`)

	sql := `
	ALTER TABLE ` + targetTable + `
	RENAME COLUMN "` + columnOld.ColumnName.String + `" TO "` + columnNew.ColumnName.String + `"`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTableColumn - Delete table column
func (d *Oracle) DeleteTableColumn(tableName string, column model.UserColumn) error {
	targetTable := strings.ToUpper(`"` + Info.GrantID + `"."` + tableName + `"`)

	sql := `
	ALTER TABLE ` + targetTable + `
	DROP COLUMN "` + column.ColumnName.String + `"`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateBoard - Create list board table
func (d *Oracle) CreateBoard(tableInfo model.Board, recreate bool) error {
	var err error

	boardName := strings.ToUpper(`"` + Info.GrantID + `"."` + tableInfo.BoardTable.String + `"`)

	sql := ``
	if recreate {
		sql = `DROP TABLE ` + boardName
		_, err = Con.Exec(sql)
		if err != nil {
			return err
		}
	}

	sql = `
	CREATE TABLE ` + boardName + ` (
		IDX         NUMBER(11) GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1) NOT NULL,
		TITLE       VARCHAR(256),
		TITLE_IMAGE VARCHAR(256),
		CONTENT     NCLOB,
		AUTHOR_IDX  NUMBER(11),
		FILES       NCLOB,
		VIEWS       VARCHAR(11),
		REGDATE    VARCHAR(14),

		UNIQUE("IDX")
	)`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateComment - Create comment table
func (d *Oracle) CreateComment(tableInfo model.Board, recreate bool) error {
	commentName := strings.ToUpper(`"` + Info.GrantID + `"."` + tableInfo.CommentTable.String + `"`)

	sql := ``
	if recreate {
		sql = `DROP TABLE ` + commentName
		_, err := Con.Exec(sql)
		if err != nil {
			return err
		}
	}

	sql = `
	CREATE TABLE ` + commentName + ` (
		IDX         NUMBER(11) GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1) NOT NULL,
		BOARD_IDX   NUMBER(11),
		CONTENT     NCLOB,
		AUTHOR_IDX  NUMBER(11),
		FILES       NCLOB,
		REGDATE    VARCHAR(14),
		
		UNIQUE("IDX")
	)`

	_, err := Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteBoard - Delete a board table
func (d *Oracle) DeleteBoard(tableInfo model.Board) error {
	var err error
	var count int64

	boardName := strings.ToUpper(`"` + Info.GrantID + `"."` + tableInfo.BoardTable.String + `"`)

	sql := `
	SELECT COUNT(TABLE_NAME) AS COUNT
	FROM user_tables
	WHERE TABLE_NAME = '` + strings.ToUpper(tableInfo.BoardTable.String) + `'`

	err = Con.QueryRow(sql).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	sql = `DROP TABLE ` + boardName
	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// RenameBoard - Rename board table name
func (d *Oracle) RenameBoard(tableInfoOLD model.Board, tableInfoNEW model.Board) error {
	var err error

	tableName := strings.ToUpper(`"` + tableInfoOLD.BoardTable.String + `"`)
	tableNameRename := strings.ToUpper(`"` + tableInfoNEW.BoardTable.String + `"`)

	sql := `ALTER TABLE ` + tableName + ` RENAME TO ` + tableNameRename
	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// RenameComment - Rename comment table name
func (d *Oracle) RenameComment(tableInfoOLD model.Board, tableInfoNEW model.Board) error {
	var err error

	tableName := strings.ToUpper(`"` + tableInfoOLD.CommentTable.String + `"`)
	tableNameRename := strings.ToUpper(`"` + tableInfoNEW.CommentTable.String + `"`)

	sql := `ALTER TABLE "` + tableName + `" RENAME TO "` + tableNameRename + `"`
	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteComment - Delete a comment table
func (d *Oracle) DeleteComment(tableInfo model.Board) error {
	var err error
	var count int64

	commentName := strings.ToUpper(`"` + Info.GrantID + `"."` + tableInfo.CommentTable.String + `"`)

	sql := `
	SELECT COUNT(TABLE_NAME) AS COUNT
	FROM user_tables
	WHERE TABLE_NAME = '` + strings.ToUpper(tableInfo.CommentTable.String) + `'`

	err = Con.QueryRow(sql).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	sql = `DROP TABLE ` + commentName
	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// GetPagingQuery - Get paging query
func (d *Oracle) GetPagingQuery(offset int, listCount int) string {
	sql := `
	OFFSET ` + strconv.Itoa(offset) + ` ROWS
	FETCH NEXT ` + strconv.Itoa(listCount) + ` ROWS ONLY`

	return sql
}