package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
		IDX BIGINT NOT NULL IDENTITY PRIMARY KEY,
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
		IDX BIGINT NOT NULL IDENTITY PRIMARY KEY,
		NAME VARCHAR(128) NULL DEFAULT NULL,
		CODE VARCHAR(64) NULL DEFAULT NULL,
		TYPE VARCHAR(64) NULL DEFAULT NULL,
		COLUMN_NAME VARCHAR(64) NULL DEFAULT NULL,
		"ORDER" INTEGER NULL DEFAULT NULL,
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
		IDX BIGINT NOT NULL IDENTITY PRIMARY KEY,
		USERNAME VARCHAR(128) UNIQUE NULL DEFAULT NULL,
		PASSWORD VARCHAR(128) NULL DEFAULT NULL,
		EMAIL VARCHAR(128) UNIQUE NULL DEFAULT NULL,
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

	// Add temp admin
	sql = `
	USE "#DATABASE"
	IF NOT EXISTS (SELECT TOP 1 * FROM #TABLE_NAME WHERE USERNAME = 'admin')
	INSERT INTO #TABLE_NAME (USERNAME, "PASSWORD", EMAIL, "ADMIN", APPROVAL)
		VALUES ('admin', 'admin', 'admin@please.modify', 'Y', 'Y')
	--GO`

	sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTableName)

	_, err = Dbo.Exec(sql)
	if err != nil {
		log.Println(sql)
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
		IDX BIGINT NOT NULL IDENTITY PRIMARY KEY,
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
	sql := `USE "#DATABASE"`
	if recreate {
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
		IDX BIGINT NOT NULL IDENTITY PRIMARY KEY,
		IS_MEMBER VARCHAR(2) NULL DEFAULT NULL,
		WRITER_IDX VARCHAR(11) NULL DEFAULT NULL,
		WRITER_NAME VARCHAR(64) NULL DEFAULT NULL,
		WRITER_PASSWORD VARCHAR(128) NULL DEFAULT NULL,
		FILES TEXT NULL DEFAULT NULL,
		REG_DTTM BIGINT NULL DEFAULT NULL,`

	if len(fields) > 0 {
		commentCount := 0

		for _, f := range fields {
			// log.Println(f.Name.String, f.Type.String, f.Order.Int64)
			if f.Type.String == "comment" {
				commentCount++
			}

			colType := ""
			switch f.Type.String {
			// cusom-tablelist
			case "text":
				colType = "TEXT"
			case "number":
				colType = "INT"
			case "real", "double":
				colType = "DECIMAL(20,20)"

			// cusom-board
			case "title", "author", "input":
				colType = "VARCHAR(512)"
			case "editor":
				colType = "TEXT"
			case "comment":
				colType = "VARCHAR(4)"
				_ = d.CreateComment(tableInfo, false)

			default:
				colType = "VARCHAR(128)"
			}

			sql += fmt.Sprintf(`%s		`+`"%s"`+`		%s,`, "\n", f.ColumnName.String, colType)
		}

		if commentCount > 1 {
			return errors.New("available only 1 comment")
		}
	}

	sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", tableInfo.Table.String)

	sql += `
	)
	--GO`

	log.Println("Sqlserver/CreateCustomBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// EditBasicBoard - Create board table
func (d *Sqlserver) EditBasicBoard(tableInfoOld models.Board, tableInfoNew models.Board) error {
	sql := `
	USE "#DATABASE"
	EXEC sp_rename "#TABLE_NAME_OLD", "#TABLE_NAME_NEW"`

	sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME_OLD", tableInfoOld.Table.String)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME_NEW", tableInfoNew.Table.String)

	log.Println("Sqlserver/EditBasicBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// EditBasicBoard - Create board table
func (d *Sqlserver) EditCustomBoard(tableInfoOld models.Board, tableInfoNew models.Board) error {
	// var err error
	// log.Println(tableInfoOld.Table.String, tableInfoNew.Table.String)
	// log.Println(tableInfoOld.Fields, tableInfoNew.Fields)

	var newFieldITF []map[string]interface{}
	_ = json.Unmarshal([]byte(tableInfoNew.Fields.(string)), &newFieldITF)

	commentCount := 0
	for _, f := range newFieldITF {
		if f["type"] == "comment" {
			commentCount++
		}
	}

	if commentCount > 1 {
		return errors.New("available only 1 comment")
	}

	var oldFieldITF []map[string]interface{}

	for _, f := range tableInfoOld.Fields.([]interface{}) {
		oldF := f.(map[string]interface{})
		oldFieldITF = append(oldFieldITF, oldF)
	}

	add, remove, modify := diffCustomBoardFields(oldFieldITF, newFieldITF)
	log.Println("Add: ", add)
	log.Println("Remove: ", remove)
	log.Println("Modify: ", modify)

	sql := ""
	if tableInfoOld.Table.String != tableInfoNew.Table.String {
		sql = `
		USE "#DATABASE"
		EXEC sp_rename "#TABLE_NAME_OLD", "#TABLE_NAME_NEW"`

		sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
		sql = strings.ReplaceAll(sql, "#TABLE_NAME_OLD", tableInfoOld.Table.String)
		sql = strings.ReplaceAll(sql, "#TABLE_NAME_NEW", tableInfoNew.Table.String)
	}

	sqlAdd := ``
	if len(add) > 0 {
		sqlAdd = `
		USE "#DATABASE"
		ALTER TABLE "#TABLE_NAME"
			ADD`
		for _, c := range add {
			sqlAdd += "\n" + c["column"].(string) + ` `
			switch c["type"].(string) {
			// cusom-tablelist
			case "text":
				sqlAdd += `TEXT`
			case "number":
				sqlAdd += `INTEGER`
			case "real":
				sqlAdd += `DECIMAL(20,20)`

			// cusom-board
			case "title", "author", "input":
				sqlAdd += `VARCHAR(512)`
			case "editor":
				sqlAdd += `TEXT`
			case "comment":
				sqlAdd += `VARCHAR(4) `
				_ = d.CreateComment(tableInfoNew, false)
			default:
				sqlAdd += `VARCHAR(128)`
			}

			sqlAdd += `, `
		}

		sqlAdd = sqlAdd[:len(sqlAdd)-2] + ";"
		sqlAdd = strings.ReplaceAll(sqlAdd, "#DATABASE", DatabaseName)
		sqlAdd = strings.ReplaceAll(sqlAdd, "#TABLE_NAME", tableInfoNew.Table.String)
	}

	sqlRemove := ""
	if len(remove) > 0 {
		for _, c := range remove {
			if c["type"].(string) == "comment" {
				d.DeleteComment(c["column"].(string))
			} else {
				sqlRemove += ` DROP COLUMN ` + c["column"].(string) + `, `
			}
		}

		sqlRemove = sqlRemove[:len(sqlRemove)-2] + ";"
	}

	sqlModify := ``
	sqlCommentRename := ``
	if len(modify) > 0 {
		for _, nc := range modify {
			for _, ocINF := range tableInfoOld.Fields.([]interface{}) {
				oc := ocINF.(map[string]interface{})
				if nc["idx"].(float64) == oc["idx"].(float64) {
					if nc["column"].(string) != oc["column"].(string) {
						if oc["type"].(string) == "comment" {

							sqlCommentRename += `
							USE "#DATABASE"
							EXEC sp_rename "#TABLE_NAME_OLD", "#TABLE_NAME_NEW";
							`

							sqlCommentRename = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
							sqlCommentRename = strings.ReplaceAll(sqlCommentRename, "#TABLE_NAME_OLD", oc["column"].(string)+"_COMMENT")
							sqlCommentRename = strings.ReplaceAll(sqlCommentRename, "#TABLE_NAME_NEW", nc["column"].(string)+"_COMMENT")
						} else {
							sqlModify += `ALTER TABLE #TABLE_NAME
							EXEC sp_rename "#TABLE_NAME".` + oc["column"].(string) + ` ` + nc["column"].(string) + `, 'COLUMN'`
							if oc["type"].(string) != nc["type"].(string) {
								colType := `ALTER COLUMN ` + nc["column"].(string) + ` `
								switch nc["type"].(string) {
								// cusom-tablelist
								case "text":
									colType += `TEXT`
								case "number":
									colType += `INTEGER`
								case "real":
									colType += `DECIMAL(20,20)`

								// cusom-board
								case "title", "author", "input":
									colType += `VARCHAR(512)`
								case "editor":
									colType += `TEXT`
								case "comment":
									colType += `VARCHAR(4)`
									_ = d.CreateComment(tableInfoNew, false)
								default:
									colType += `VARCHAR(128)`
								}
								colType += ` NULL;`
								sqlModify += colType
							}

							sqlModify = strings.ReplaceAll(sqlModify, "#TABLE_NAME", tableInfoNew.Table.String)
						}
					}
					break
				}
			}
		}
		// if strings.Contains(sqlModify, "RENAME COLUMN") {
		// 	sqlModify = sqlModify[:len(sqlModify)-2]
		// }
		// sqlModify += `; `
	}

	if sqlAdd != "" || sqlRemove != "" || sqlModify != "" {
		sql += sqlAdd + sqlRemove + sqlModify
		sql = sql[:len(sql)-2]

		sql += ";"

		sql += sqlCommentRename
	}

	log.Println("Sqlserver/EditCustomBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteBoard - Delete a board table
func (d *Sqlserver) DeleteBoard(tableName string) error {
	sql := `USE "#DATABASE"`
	sql += `
	IF OBJECT_ID('#TABLE_NAME','U') IS NOT NULL
	DROP TABLE "#TABLE_NAME"
	-- GO
	`

	sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", tableName)

	log.Println("Sqlserver/DeleteBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateComment - Create comment table
func (d *Sqlserver) CreateComment(tableInfo models.Board, recreate bool) error {
	sql := `USE "#DATABASE"`
	if recreate {
		sql += `
		IF OBJECT_ID('#TABLE_NAME','U') IS NOT NULL
		DROP TABLE "#TABLE_NAME"
		-- GO
		`

		sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
		sql = strings.ReplaceAll(sql, "#TABLE_NAME", tableInfo.Table.String+"_COMMENT")
		log.Println("Sqlserver/CreateComment: ", sql)

		_, err := Dbo.Exec(sql)
		if err != nil {
			return err
		}
	}

	sql = `USE "#DATABASE"`
	sql += `
	IF OBJECT_ID(N'#TABLE_NAME', N'U') IS NULL
	CREATE TABLE "#TABLE_NAME" (
		IDX BIGINT NOT NULL IDENTITY PRIMARY KEY,
		BOARD_IDX INT NOT NULL,
		CONTENT TEXT NULL DEFAULT NULL,
		IS_MEMBER VARCHAR(2) NULL DEFAULT NULL,
		WRITER_IDX VARCHAR(11) NULL DEFAULT NULL,
		WRITER_NAME VARCHAR(64) NULL DEFAULT NULL,
		WRITER_PASSWORD VARCHAR(128) NULL DEFAULT NULL,
		REG_DTTM BIGINT NULL DEFAULT NULL,
	)
	--GO`

	sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", tableInfo.Table.String+"_COMMENT")

	log.Println("Sqlserver/CreateComment: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteComment - Delete a comment table
func (d *Sqlserver) DeleteComment(tableName string) error {
	sql := `DROP TABLE IF EXISTS ` + "`#TABLE_NAME`" + `;`
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", DatabaseName+"`.`"+tableName+"_COMMENT")

	log.Println("Sqlserver/DeleteComment: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// EditUserTableFields - Edit user table schema
func (d *Sqlserver) EditUserTableFields(fieldsInfoOld []models.UserColumn, fieldsInfoNew []models.UserColumn, notUse []string) error {
	_, _, modify := diffUserTableFields(fieldsInfoOld, fieldsInfoNew)
	log.Println("User fields Modify: ", modify)

	sql := ``

	if len(modify) > 0 {
		sql = `USE "#DATABASE"
		`
		sqlModify := ``
		for _, nm := range modify {
			for _, om := range fieldsInfoOld {
				if nm.Idx.Int64 == om.Idx.Int64 {
					if nm.ColumnName.String != om.ColumnName.String {
						sqlModify += `EXEC sp_rename 'dbo.#TABLE_NAME.` + om.ColumnName.String + `', '` + nm.ColumnName.String + `', 'COLUMN';
						`
					}
					break
				}
			}
		}
		sql += sqlModify
	}

	sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTableName)

	log.Println("Sqlserver/EditUserTableFields: ", sql)

	if sql == "" {
		if len(modify) > 0 {
			return nil
		} else {
			return errors.New("nothing to change")
		}
	}

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserTableFields - Delete user table field
func (d *Sqlserver) DeleteUserTableFields(fieldsInfoRemove []models.UserColumn) error {
	remove := fieldsInfoRemove
	sql := ""

	if len(remove) > 0 {
		sql = `
		USE "#DATABASE"
		ALTER TABLE "#TABLE_NAME"
		`
		sql += `DROP COLUMN`
		for _, r := range remove {
			sql += ` ` + r.ColumnName.String + `, `
		}
		sql = sql[:len(sql)-2]
		sql += sql + `; `
	}

	sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTableName)

	log.Println("Sqlserver/DeleteUserTableFields: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// AddUserTableFields - Add user column
func (d *Sqlserver) AddUserTableFields(fields []models.UserColumn) error {
	sql := ``
	if len(fields) > 0 {
		sql = `USE "#DATABASE"
		ALTER TABLE "#TABLE_NAME"
			ADD`
		for _, a := range fields {
			sql += "\n" + a.ColumnName.String + ` `
			switch a.Type.String {
			case "text":
				sql += `VARCHAR(128)`
			case "number":
				sql += `INTEGER`
			case "real":
				sql += `REAL`
			}

			sql += `, `
		}

		sql = sql[:len(sql)-2] + ";"
		sql = strings.ReplaceAll(sql, "#DATABASE", DatabaseName)
		sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTableName)

		log.Println("Sqlserver/EditUserTableFields: ", sql)

		_, err := Dbo.Exec(sql)
		if err != nil {
			return err
		}
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
