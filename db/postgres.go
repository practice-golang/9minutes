package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/lib/pq"
	"github.com/practice-golang/9minutes/models"
	"github.com/thoas/go-funk"
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
	sql := `CREATE DATABASE "#DATABASE_NAME";`
	sql = strings.ReplaceAll(sql, "#DATABASE_NAME", DatabaseName)

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
		"IDX" BIGSERIAL PRIMARY KEY,
		"NAME" VARCHAR(128) NULL DEFAULT NULL,
		"CODE" VARCHAR(64) NULL DEFAULT NULL,
		"TYPE" VARCHAR(64) NULL DEFAULT NULL,
		"TABLE" VARCHAR(64) NULL DEFAULT NULL,
		"GRANT_READ" VARCHAR(16) NULL DEFAULT NULL,
		"GRANT_WRITE" VARCHAR(16) NULL DEFAULT NULL,
		"GRANT_COMMENT" VARCHAR(16) NULL DEFAULT NULL,
		"FILE_UPLOAD" VARCHAR(2) NULL DEFAULT NULL,
		"FIELDS" TEXT NULL DEFAULT NULL
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
		"CODE" VARCHAR(64) NULL DEFAULT NULL,
		"TYPE" VARCHAR(64) NULL DEFAULT NULL,
		"COLUMN_NAME" VARCHAR(64) NULL DEFAULT NULL,
		"ORDER" VARCHAR(11) NULL DEFAULT NULL
	);`

	sql = strings.ReplaceAll(sql, "#SCHEMA_NAME", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserFieldTable)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserTable - Create user table
func (d *Postgres) CreateUserTable(recreate bool) error {
	sql := `CREATE SCHEMA IF NOT EXISTS #SCHEMA_NAME;`

	if recreate {
		sql += `DROP TABLE IF EXISTS #TABLE_NAME;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS #TABLE_NAME (
		"IDX" SERIAL PRIMARY KEY,
		"USERNAME" VARCHAR(128) UNIQUE NULL DEFAULT NULL,
		"PASSWORD" VARCHAR(128) NULL DEFAULT NULL,
		"EMAIL" VARCHAR(128) UNIQUE NULL DEFAULT NULL,
		"ADMIN" VARCHAR(2) NULL DEFAULT NULL,
		"APPROVAL" VARCHAR(2) NULL DEFAULT NULL,
		"REG_DTTM" BIGINT NULL DEFAULT NULL,

		CONSTRAINT "USERS_UQ" UNIQUE ("USERNAME", "EMAIL")
	);`

	sql = strings.ReplaceAll(sql, "#SCHEMA_NAME", DatabaseName)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTable)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	// Add temp admin
	sql = `
	INSERT INTO #TABLE_NAME ("USERNAME", "PASSWORD", "EMAIL", "ADMIN", "APPROVAL")
		VALUES ('admin', 'admin', 'admin@please.modify', 'Y', 'Y')
	ON CONFLICT ON CONSTRAINT "USERS_UQ" DO NOTHING;`
	// ON CONFLICT ("USERNAME", "EMAIL") DO NOTHING;`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTable)

	_, err = Dbo.Exec(sql)
	if err != nil {
		log.Println(sql)
		return err
	}

	return nil
}

// CreateBasicBoard - Create board table
func (d *Postgres) CreateBasicBoard(tableInfo models.Board, recreate bool) error {
	sql := ``
	if recreate {
		sql += `DROP TABLE IF EXISTS #TABLE_NAME;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS #TABLE_NAME (
		"IDX" BIGSERIAL PRIMARY KEY,
		"TITLE" VARCHAR(256) NULL DEFAULT NULL,
		"CONTENT" TEXT NULL DEFAULT NULL,
		"IS_MEMBER" VARCHAR(2) NULL DEFAULT NULL,
		"WRITER_IDX" VARCHAR(11) NULL DEFAULT NULL,
		"WRITER_NAME" VARCHAR(64) NULL DEFAULT NULL,
		"WRITER_PASSWORD" VARCHAR(128) NULL DEFAULT NULL,
		"FILES" TEXT NULL DEFAULT NULL,
		"REG_DTTM" BIGINT NULL DEFAULT NULL
	);`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", DatabaseName+`."`+tableInfo.Table.String+`"`)

	log.Println("Postgres/CreateBasicBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateCustomBoard - Create board table
func (d *Postgres) CreateCustomBoard(tableInfo models.Board, fields []models.Field, recreate bool) error {
	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS #TABLE_NAME;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS #TABLE_NAME (
		"IDX" BIGSERIAL PRIMARY KEY,
		"IS_MEMBER" VARCHAR(2) NULL DEFAULT NULL,
		"WRITER_IDX" VARCHAR(11) NULL DEFAULT NULL,
		"WRITER_NAME" VARCHAR(64) NULL DEFAULT NULL,
		"WRITER_PASSWORD" VARCHAR(128) NULL DEFAULT NULL,
		"FILES" TEXT NULL DEFAULT NULL,
		"REG_DTTM" BIGINT NULL DEFAULT NULL`

	if len(fields) > 0 {
		sql += `,`
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
				colType = "INTEGER"
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

			sql += fmt.Sprintf(`%s "%s" %s NULL DEFAULT NULL, `, "\n", f.ColumnName.String, colType)
		}

		if commentCount > 1 {
			return errors.New("available only 1 comment")
		}

		sql = sql[:len(sql)-2]
	}

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", DatabaseName+`."`+tableInfo.Table.String+`"`)

	sql += `);`

	log.Println("Postgres/CreateCustomBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// EditBasicBoard - Create board table
func (d *Postgres) EditBasicBoard(tableInfoOld models.Board, tableInfoNew models.Board) error {
	sql := `ALTER TABLE #TABLE_NAME_OLD RENAME TO #TABLE_NAME_NEW;`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME_OLD", DatabaseName+`."`+tableInfoOld.Table.String+`"`)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME_NEW", `"`+tableInfoNew.Table.String+`"`)

	log.Println("Postgres/EditBasicBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// EditBasicBoard - Create board table
func (d *Postgres) EditCustomBoard(tableInfoOld models.Board, tableInfoNew models.Board) error {
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
		sql = `ALTER TABLE #TABLE_NAME_OLD RENAME TO #TABLE_NAME_NEW;`

		sql = strings.ReplaceAll(sql, "#TABLE_NAME_OLD", DatabaseName+`."`+tableInfoOld.Table.String+`"`)
		sql = strings.ReplaceAll(sql, "#TABLE_NAME_NEW", `"`+tableInfoNew.Table.String+`"`)
	}

	sqlAdd := ""
	if len(add) > 0 {
		for _, c := range add {
			sqlAdd += ` ADD "` + c["column"].(string) + `" `
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
				sqlAdd += `VARCHAR(4)`
				_ = d.CreateComment(tableInfoNew, false)
			default:
				sqlAdd += `VARCHAR(128)`
			}

			sqlAdd += `, `
			// if i < (len(add) - 1) {
			// 	sqlAdd += `, `
			// }
		}
		// sqlAdd += `; `
	}

	sqlRemove := ""
	if len(remove) > 0 {
		for _, c := range remove {
			if c["type"].(string) == "comment" {
				d.DeleteComment(c["column"].(string))
			} else {
				sqlRemove += ` DROP COLUMN "` + c["column"].(string) + `", `
			}
		}
		// if strings.Contains(sqlRemove, "DROP COLUMN") {
		// 	sqlRemove = sqlRemove[:len(sqlRemove)-2]
		// }
		// sql += sqlRemove + `; `
	}

	sqlTypeChange := ""
	sqlRename := ""
	sqlCommentRename := ""
	if len(modify) > 0 {
		for _, nc := range modify {
			for _, ocINF := range tableInfoOld.Fields.([]interface{}) {
				oc := ocINF.(map[string]interface{})
				if nc["idx"].(float64) == oc["idx"].(float64) {
					if nc["column"].(string) != oc["column"].(string) {
						if oc["type"].(string) == "comment" {
							sqlCommentRename += `ALTER TABLE #TABLE_NAME_OLD RENAME TO #TABLE_NAME_NEW; `

							sqlCommentRename = strings.ReplaceAll(sqlCommentRename, "#TABLE_NAME_OLD", DatabaseName+`."`+oc["column"].(string)+`_COMMENT`+`"`)
							sqlCommentRename = strings.ReplaceAll(sqlCommentRename, "#TABLE_NAME_NEW", `"`+nc["column"].(string)+`_COMMENT`+`"`)
						} else {
							sqlRename += `RENAME COLUMN "` + oc["column"].(string) + `" TO "` + nc["column"].(string) + `", `
							sqlTypeChange += `ALTER COLUMN "` + nc["column"].(string) + `" TYPE `
							switch nc["type"].(string) {
							// cusom-tablelist
							case "text":
								sqlTypeChange += `TEXT`
							case "number":
								sqlTypeChange += `INTEGER`
							case "real":
								sqlTypeChange += `DECIMAL(20,20)`

							// cusom-board
							case "title", "author", "input":
								sqlTypeChange += `VARCHAR(512)`
							case "editor":
								sqlTypeChange += `TEXT`
							case "comment":
								sqlTypeChange += `VARCHAR(4)`
								_ = d.CreateComment(tableInfoNew, false)
							default:
								sqlTypeChange += `VARCHAR(128)`
							}

							sqlTypeChange += ", "

							// if i < (len(modify) - 1) {
							// 	sqlModify += `, `
							// }
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

	if sqlAdd != "" {
		sql += `ALTER TABLE #TABLE_NAME_NEW ` + sqlAdd
		sql = sql[:len(sql)-2]

		sql += ";"
	}
	if sqlRemove != "" {
		sql += `ALTER TABLE #TABLE_NAME_NEW ` + sqlRemove
		sql = sql[:len(sql)-2]
		sql += ";"
	}
	if sqlRename != "" {
		sql += `ALTER TABLE #TABLE_NAME_NEW ` + sqlRename
		sql = sql[:len(sql)-2]
		sql += ";"
	}
	if sqlTypeChange != "" {
		sql += `ALTER TABLE #TABLE_NAME_NEW ` + sqlTypeChange
		sql = sql[:len(sql)-2]
		sql += ";"
	}
	if sqlCommentRename != "" {
		sql += sqlCommentRename
	}
	sql = strings.ReplaceAll(sql, "#TABLE_NAME_NEW", DatabaseName+`."`+tableInfoNew.Table.String+`"`)

	log.Println("Postgres/EditCustomBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteBoard - Delete a board table
func (d *Postgres) DeleteBoard(tableName string) error {
	sql := `DROP TABLE IF EXISTS #TABLE_NAME;`
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", DatabaseName+`."`+tableName+`"`)

	log.Println("Postgres/DeleteBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateComment - Create comment table
func (d *Postgres) CreateComment(tableInfo models.Board, recreate bool) error {
	sql := ``
	if recreate {
		sql += `DROP TABLE IF EXISTS #TABLE_NAME;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS #TABLE_NAME (
		"IDX" BIGSERIAL PRIMARY KEY,
		"BOARD_IDX" BIGSERIAL NOT NULL,
		"CONTENT" TEXT NULL DEFAULT NULL,
		"IS_MEMBER" VARCHAR(2) NULL DEFAULT NULL,
		"WRITER_IDX" VARCHAR(11) NULL DEFAULT NULL,
		"WRITER_NAME" VARCHAR(64) NULL DEFAULT NULL,
		"WRITER_PASSWORD" VARCHAR(128) NULL DEFAULT NULL,
		"REG_DTTM" BIGINT NULL DEFAULT NULL
	);`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", DatabaseName+`."`+tableInfo.Table.String+`_COMMENT`+`"`)

	log.Println("Postgres/CreateComment: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteComment - Delete comment table
func (d *Postgres) DeleteComment(tableName string) error {
	sql := `DROP TABLE IF EXISTS #TABLE_NAME;`
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", DatabaseName+`."`+tableName+`_COMMENT`+`"`)

	log.Println("Postgres/DeleteComment: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// EditUserTableFields - Edit user table schema
func (d *Postgres) EditUserTableFields(fieldsInfoOld []models.UserColumn, fieldsInfoNew []models.UserColumn, notUse []string) error {
	add, remove, modify := diffUserTableFields(fieldsInfoOld, fieldsInfoNew)
	log.Println("User fields Add: ", add)
	log.Println("User fields Remove: ", remove)
	log.Println("User fields Modify: ", modify)

	sql := ""
	if len(add) > 0 {
		for _, a := range add {
			sql += `ALTER TABLE #TABLE_NAME `
			sql += ` ADD COLUMN ` + "`" + a.ColumnName.String + "`" + ` `
			switch a.Type.String {
			// cusom-tablelist
			case "text":
				sql += ` VARCHAR(256)`
			case "number":
				sql += ` INTEGER`
			case "real":
				sql += ` DECIMAL(20,20)`
			}

			sql += `; `
		}
	}

	if len(remove) > 0 && !funk.Contains(notUse, "remove") {
		sqlRemove := `ALTER TABLE #TABLE_NAME `
		for _, r := range remove {
			sqlRemove += ` DROP COLUMN ` + "`" + r.ColumnName.String + "`" + `, `
		}
		if strings.Contains(sqlRemove, "DROP COLUMN") {
			sqlRemove = sqlRemove[:len(sqlRemove)-2]
		}
		sql += sqlRemove + `; `
	}

	if len(modify) > 0 {
		sqlModify := ""
		for _, nm := range modify {
			for _, om := range fieldsInfoOld {
				if nm.Idx.Int64 == om.Idx.Int64 {
					if nm.ColumnName.String != om.ColumnName.String {
						sqlModify += `ALTER TABLE #TABLE_NAME `
						sqlModify += ` CHANGE COLUMN ` + "`" + om.ColumnName.String + "`" + ` TO ` + "`" + nm.ColumnName.String + "`" + `; `
					}
					break
				}
			}
		}
		sql += sqlModify
	}

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTable)

	log.Println("Postgres/EditUserTableFields: ", sql)

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

// AddUserTableFields - Add user column
func (d *Postgres) AddUserTableFields(fields []models.UserColumn) error {
	sql := ""
	for _, a := range fields {
		sql += `ALTER TABLE "#TABLE_NAME" ADD COLUMN ` + a.ColumnName.String + ` `
		switch a.Type.String {
		case "text":
			sql += ` TEXT`
		case "number":
			sql += ` INTEGER`
		case "real":
			sql += ` REAL`
		}

		sql += `; `
	}

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTable)

	log.Println("Sqlite/EditUserTableFields: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// SelectColumnNames - Get column names of table
func (d *Postgres) SelectColumnNames(table string) (sql.Result, error) {
	result, err := Dbo.Exec("PRAGMA TABLE_INFO(" + table + ")")
	if err != nil {
		return nil, err
	}

	return result, nil
}
