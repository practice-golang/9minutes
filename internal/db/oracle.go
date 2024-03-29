package db

import (
	"9minutes/consts"
	"9minutes/model"
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Oracle struct {
	dsn     string
	Version int64
}

func (d *Oracle) connect() (*sql.DB, error) {
	db, err := sql.Open("oracle", d.dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(-1)
	db.SetConnMaxLifetime(-1)
	db.SetMaxIdleConns(-1)
	db.SetMaxOpenConns(-1)

	return db, nil
}

func (d *Oracle) CreateDB() error {
	return nil
}

func (d *Oracle) Exec(que string, colValues []interface{}, options string) (int64, int64, error) {
	var err error

	var count int64 = 0
	var idx *int64

	que += ` RETURNING ` + options + ` INTO :1`

	result, err := Con.Exec(que, sql.Out{Dest: &idx})
	if err != nil {
		log.Println("Oracle Exec:", err.Error())
		return 0, -1, err
	}

	count, _ = result.RowsAffected()

	return count, *idx, nil
}

// CreateBoardTable - Create board manager table
func (d *Oracle) CreateBoardTable() error {
	var err error

	var tableCNT int64 = 0

	sql := `
	SELECT
		COUNT(table_name) AS CNT
	FROM all_tables
	WHERE owner = '` + strings.ToUpper(Info.GrantID) + `'
		AND table_name = '` + strings.ToUpper(Info.BoardTable) + `'`

	rows, err := Con.Query(sql)
	if err != nil {
		log.Println("Oracle CreateBoardTable:", err.Error())
	} else {
		for rows.Next() {
			_ = rows.Scan(&tableCNT)
		}
	}

	if tableCNT > 0 {
		return nil
	}

	boardTable := strings.ToUpper(`"` + Info.GrantID + `"."` + Info.BoardTable + `"`)

	sql = `
	CREATE TABLE ` + boardTable + ` (
		IDX           NUMBER(11) GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1) NOT NULL,
		BOARD_NAME    VARCHAR2(128),
		BOARD_CODE    VARCHAR2(64),
		BOARD_TYPE    VARCHAR2(64),
		BOARD_TABLE   VARCHAR2(64),
		COMMENT_TABLE VARCHAR2(64),
		GRANT_READ    VARCHAR2(16),
		GRANT_WRITE   VARCHAR2(16),
		GRANT_COMMENT VARCHAR2(16),
		GRANT_UPLOAD  VARCHAR2(16),
		FIELDS        NCLOB,

		CONSTRAINT boards_pk PRIMARY KEY ("IDX")
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

	var tableCNT int64

	sql := `
	SELECT
		COUNT(table_name) AS CNT
	FROM all_tables
	WHERE owner = '` + strings.ToUpper(Info.GrantID) + `'
		AND table_name = '` + strings.ToUpper(Info.UploadTable) + `'`

	rows, _ := Con.Query(sql)
	for rows.Next() {
		_ = rows.Scan(&tableCNT)
	}

	if tableCNT > 0 {
		return nil
	}

	uploadTable := strings.ToUpper(`"` + Info.GrantID + `"."` + Info.UploadTable + `"`)

	sql = `
	CREATE TABLE ` + uploadTable + ` (
		IDX          NUMBER(11) GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1) NOT NULL,
		TOPIC_IDX    NUMBER(11),
		COMMENT_IDX  NUMBER(11),
		FILE_NAME    VARCHAR2(512),
		STORAGE_NAME VARCHAR2(512),
		REGDATE      VARCHAR2(14),

		CONSTRAINT uploads_pk PRIMARY KEY (IDX)
	)`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateMemberTable - Create board member table
func (d *Oracle) CreateMemberTable() error {
	var err error

	var tableCNT int64

	sql := `
	SELECT
		COUNT(table_name) AS CNT
	FROM all_tables
	WHERE owner = '` + strings.ToUpper(Info.GrantID) + `'
		AND table_name = '` + strings.ToUpper(Info.MemberTable) + `'`

	rows, _ := Con.Query(sql)
	for rows.Next() {
		_ = rows.Scan(&tableCNT)
	}

	if tableCNT > 0 {
		return nil
	}

	memberTable := strings.ToUpper(`"` + Info.GrantID + `"."` + Info.MemberTable + `"`)

	sql = `
	CREATE TABLE ` + memberTable + ` (
		IDX          NUMBER(11) GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1) NOT NULL,
		BOARD_IDX    NUMBER(11),
		USER_IDX     NUMBER(11),
		GRADE        VARCHAR2(24),
		REGDATE      VARCHAR2(14),

		CONSTRAINT member_pk PRIMARY KEY (IDX)
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

	var tableCNT int64

	sql := `
	SELECT
		COUNT(table_name) AS CNT
	FROM all_tables
	WHERE owner = '` + strings.ToUpper(Info.GrantID) + `'
		AND table_name = '` + strings.ToUpper(Info.UserTable) + `'`

	rows, _ := Con.Query(sql)
	for rows.Next() {
		_ = rows.Scan(&tableCNT)
	}

	if tableCNT > 0 {
		return nil
	}

	userTable := strings.ToUpper(`"` + Info.GrantID + `"."` + Info.UserTable + `"`)

	sql = `
	CREATE TABLE ` + userTable + ` (
		IDX      NUMBER(11) GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1) NOT NULL,
		USERID   VARCHAR2(128),
		PASSWORD VARCHAR2(128),
		EMAIL    VARCHAR2(128),
		GRADE    VARCHAR2(24),
		APPROVAL VARCHAR2(2),
		REGDATE  VARCHAR2(14),

		CONSTRAINT "` + strings.ToLower(Info.UserTable) + `_idx" PRIMARY KEY (IDX),
		CONSTRAINT "` + strings.ToLower(Info.UserTable) + `_userconstraint" UNIQUE (USERID, EMAIL)
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
		"USERID", "PASSWORD", "EMAIL",
		"GRADE", "APPROVAL",
		"REGDATE"
	)
	SELECT
		'admin', '` + string(adminPassword) + `', 'admin@please.modify',
		'admin', 'Y',
		'` + now + `'
	FROM DUAL
	WHERE NOT EXISTS (
		SELECT
			*
		FROM ` + userTable + `
		WHERE "USERID" = 'admin'
	)`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	userfieldTable := strings.ToUpper(`"` + Info.GrantID + `"."` + strings.ToUpper(`user_fields"`))

	sql = `
	CREATE TABLE ` + userfieldTable + ` (
		IDX          NUMBER(11) GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1) NOT NULL,
		DISPLAY_NAME VARCHAR2(128),
		COLUMN_CODE  VARCHAR2(128),
		COLUMN_TYPE  VARCHAR2(128),
		COLUMN_NAME  VARCHAR2(128),
		SORT_ORDER   NUMBER(5),

		CONSTRAINT user_fields_idx PRIMARY KEY (IDX),
		CONSTRAINT user_fields_columnconstraint UNIQUE (COLUMN_NAME)
	)`

	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	sql = `INSERT INTO ` + userfieldTable + ` ("DISPLAY_NAME", "COLUMN_CODE", "COLUMN_TYPE", "COLUMN_NAME", "SORT_ORDER")
	VALUES ('Idx', 'idx', 'integer', 'IDX', 1)`
	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}
	sql = `INSERT INTO ` + userfieldTable + ` ("DISPLAY_NAME", "COLUMN_CODE", "COLUMN_TYPE", "COLUMN_NAME", "SORT_ORDER")
	VALUES ('UserId', 'userid', 'text', 'USERID', 2)`
	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}
	sql = `INSERT INTO ` + userfieldTable + ` ("DISPLAY_NAME", "COLUMN_CODE", "COLUMN_TYPE", "COLUMN_NAME", "SORT_ORDER")
	VALUES ('Password', 'password', 'text', 'PASSWORD', 3)`
	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}
	sql = `INSERT INTO ` + userfieldTable + ` ("DISPLAY_NAME", "COLUMN_CODE", "COLUMN_TYPE", "COLUMN_NAME", "SORT_ORDER")
	VALUES ('Email', 'email', 'text', 'EMAIL', 4)`
	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}
	sql = `INSERT INTO ` + userfieldTable + ` ("DISPLAY_NAME", "COLUMN_CODE", "COLUMN_TYPE", "COLUMN_NAME", "SORT_ORDER")
	VALUES ('Grade', 'grade', 'text', 'GRADE', 5)`
	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}
	sql = `INSERT INTO ` + userfieldTable + ` ("DISPLAY_NAME", "COLUMN_CODE", "COLUMN_TYPE", "COLUMN_NAME", "SORT_ORDER")
	VALUES ('Approval', 'approval', 'text', 'APPROVAL', 6)`
	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}
	sql = `INSERT INTO ` + userfieldTable + ` ("DISPLAY_NAME", "COLUMN_CODE", "COLUMN_TYPE", "COLUMN_NAME", "SORT_ORDER")
	VALUES ('RegDate', 'regdate', 'text', 'REGDATE', 7)`
	_, err = Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserVerificationTable - Create user verification table
func (d *Oracle) CreateUserVerificationTable() error {
	var err error

	var tableCNT int64

	sql := `
	SELECT
		COUNT(table_name) AS CNT
	FROM all_tables
	WHERE owner = '` + strings.ToUpper(Info.GrantID) + `'
		AND table_name = '` + strings.ToUpper(Info.UserTable+"_verification") + `'`

	rows, _ := Con.Query(sql)
	for rows.Next() {
		_ = rows.Scan(&tableCNT)
	}

	if tableCNT > 0 {
		return nil
	}

	verificationTable := strings.ToUpper(`"` + Info.GrantID + `"."` + Info.UserTable + `_verification"`)

	sql = `
	CREATE TABLE ` + verificationTable + ` (
		IDX      NUMBER(11) GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1) NOT NULL,
		USER_IDX NUMBER(11),
		TOKEN    VARCHAR2(128),
		REGDATE  VARCHAR2(14),

		UNIQUE(IDX)
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
		sql += ` VARCHAR2(256)`
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
		IDX           NUMBER(11) GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1) NOT NULL,
		TITLE         VARCHAR2(256),
		TITLE_IMAGE   VARCHAR2(256),
		CONTENT       NCLOB,
		AUTHOR_IDX    NUMBER(11),
		AUTHOR_NAME   VARCHAR2(256),
		AUTHOR_IP     VARCHAR2(50),
		AUTHOR_IP_CUT VARCHAR2(50),
		EDIT_PASSWORD VARCHAR2(256),
		FILES         NCLOB,
		IMAGES        NCLOB,
		VIEWS         VARCHAR2(11),
		REGDATE       VARCHAR2(14),

		UNIQUE(IDX)
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
		IDX           NUMBER(11) GENERATED ALWAYS AS IDENTITY (START WITH 1 INCREMENT BY 1) NOT NULL,
		TOPIC_IDX     NUMBER(11),
		CONTENT       NCLOB,
		AUTHOR_IDX    NUMBER(11),
		AUTHOR_NAME   VARCHAR2(256),
		AUTHOR_IP     VARCHAR2(50),
		AUTHOR_IP_CUT VARCHAR2(50),
		EDIT_PASSWORD VARCHAR2(256),
		FILES         NCLOB,
		IMAGES        NCLOB,
		REGDATE       VARCHAR2(14),
		
		UNIQUE(IDX)
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

	boardName := strings.ToUpper(`"` + Info.GrantID + `"."` + tableInfo.BoardTable.String + `"`)

	sql := `DROP TABLE ` + boardName
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

	commentName := strings.ToUpper(`"` + Info.GrantID + `"."` + tableInfo.CommentTable.String + `"`)

	sql := `DROP TABLE ` + commentName
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
