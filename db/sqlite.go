package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/practice-golang/9minutes/models"
	"github.com/thoas/go-funk"
	// _ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
)

type Sqlite struct{ Dsn string }

// initDB - Prepare DB
func (d *Sqlite) initDB() (*sql.DB, error) {
	var err error

	Dbo, err = sql.Open("sqlite", d.Dsn)
	if err != nil {
		return nil, err
	}

	return Dbo, nil
}

func (d *Sqlite) CreateDB() error {
	return nil
}

// CreateBoardManagerTable - Create board manager table
func (d *Sqlite) CreateBoardManagerTable(recreate bool) error {
	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS "#TABLE_NAME";`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS "#TABLE_NAME" (
		"IDX"				INTEGER,
		"NAME"				TEXT,
		"CODE"				TEXT,
		"TYPE"				TEXT,
		"TABLE"				TEXT UNIQUE,
		"GRANT_READ"		TEXT,
		"GRANT_WRITE"		TEXT,
		"FIELDS"			TEXT,
		PRIMARY KEY("IDX" AUTOINCREMENT)
	);`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", BoardManagerTable)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserFieldTable - Create user manager table
func (d *Sqlite) CreateUserFieldTable(recreate bool) error {
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

// CreateUserTable - Create user table
func (d *Sqlite) CreateUserTable(recreate bool) error {
	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS "#TABLE_NAME";`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS "#TABLE_NAME" (
		"IDX"			INTEGER,
		"NAME"			TEXT UNIQUE,
		"PASSWORD"		TEXT,
		"EMAIL"			TEXT UNIQUE,
		"ADMIN"			TEXT,
		"APPROVAL"		TEXT,
		"REG_DTTM"		TEXT,
		PRIMARY KEY("IDX" AUTOINCREMENT)
	);`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTable)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateBasicBoard - Create board table
func (d *Sqlite) CreateBasicBoard(tableInfo models.Board, recreate bool) error {
	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS "#TABLE_NAME";`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS "#TABLE_NAME" (
		"IDX"				INTEGER,
		"TITLE"				TEXT,
		"CONTENT"			TEXT,
		"WRITER_IDX"		TEXT,
		"WRITER_NAME"		TEXT,
		"WRITER_PASSWORD"	TEXT,
		"REG_DTTM"			TEXT,
		PRIMARY KEY("IDX" AUTOINCREMENT)
	);`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", tableInfo.Table.String)

	log.Println("Sqlite/CreateBasicBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateCustomBoard - Create board table
func (d *Sqlite) CreateCustomBoard(tableInfo models.Board, fields []models.Field, recreate bool) error {
	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS "#TABLE_NAME";`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS "#TABLE_NAME" (
		"IDX"		INTEGER,
	`

	if len(fields) > 0 {
		for k, f := range fields {
			log.Println(k, f.Name.String, f.Type.String, f.Order.Int64)
			colType := ""
			switch f.Type.String {
			// cusom-tablelist
			case "text":
				colType = "TEXT"
			case "number":
				colType = "INTEGER"
			case "real", "double":
				colType = "REAL"

			// cusom-board
			case "title", "author", "input":
				colType = "TEXT"
			case "editor":
				colType = "TEXT"

			default:
				colType = "TEXT"
			}

			sql += fmt.Sprintf(`%s		"%s"		%s,`, "\n", f.ColumnName.String, colType)
		}
	}

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", tableInfo.Table.String)

	sql += `
		PRIMARY KEY("IDX" AUTOINCREMENT)
	);
	`

	log.Println("Sqlite/CreateCustomBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// EditBasicBoard - Create board table
func (d *Sqlite) EditBasicBoard(tableInfoOld models.Board, tableInfoNew models.Board) error {
	sql := `ALTER TABLE "#TABLE_NAME_OLD" RENAME TO "#TABLE_NAME_NEW";`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME_OLD", tableInfoOld.Table.String)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME_NEW", tableInfoNew.Table.String)

	log.Println("Sqlite/EditBasicBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func diffCustomBoardFields(old, new []map[string]interface{}) (add, remove, modify []map[string]interface{}) {
	var diff []map[string]interface{}

	for i := 0; i < 2; i++ {
		for _, s1 := range old {
			found := false
			for _, s2 := range new {
				if s1["idx"] == s2["idx"] {
					if i == 0 && !cmp.Equal(s1, s2) {
						modify = append(modify, s2)
					}

					found = true
					break
				}
			}

			if !found {
				diff = append(diff, s1)
			}
		}

		if i == 0 {
			remove = diff
			old, new = new, old
		} else {
			add = diff
		}

		diff = []map[string]interface{}{}
	}

	return
}

// EditCustomBoard - Create custom table
func (d *Sqlite) EditCustomBoard(tableInfoOld models.Board, tableInfoNew models.Board) error {
	// var err error
	// log.Println(tableInfoOld.Table.String, tableInfoNew.Table.String)
	// log.Println(tableInfoOld.Fields, tableInfoNew.Fields)

	var newFieldITF []map[string]interface{}
	_ = json.Unmarshal([]byte(tableInfoNew.Fields.(string)), &newFieldITF)

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
		sql = `ALTER TABLE "#TABLE_NAME_OLD" RENAME TO "#TABLE_NAME_NEW"; `

		sql = strings.ReplaceAll(sql, "#TABLE_NAME_OLD", tableInfoOld.Table.String)
		sql = strings.ReplaceAll(sql, "#TABLE_NAME_NEW", tableInfoNew.Table.String)
	}

	if len(add) > 0 {
		sql += `ALTER TABLE "#TABLE_NAME_NEW" `
		for i, c := range add {
			sql += ` ADD COLUMN ` + c["column"].(string) + ` `
			switch c["type"].(string) {
			// cusom-tablelist
			case "text":
				sql += ` TEXT`
			case "number":
				sql += ` INTEGER`
			case "real":
				sql += ` REAL`

			// cusom-board
			case "title", "author", "input":
				sql += ` TEXT`
			case "editor":
				sql += ` TEXT`
			case "comment":
				_ = d.CreateComment(tableInfoNew, false)
			}

			if i < (len(add) - 1) {
				sql += `, `
			}
		}
		sql += `; `
	}

	if len(remove) > 0 {
		sqlRemove := `ALTER TABLE "#TABLE_NAME_NEW" `
		for _, c := range remove {
			if c["type"].(string) == "comment" {
				d.DeleteComment(c["column"].(string))
			} else {
				sqlRemove += ` DROP COLUMN ` + c["column"].(string) + `, `
			}
		}
		if strings.Contains(sqlRemove, "DROP COLUMN") {
			sqlRemove = sqlRemove[:len(sqlRemove)-2]
		}
		sql += sqlRemove + `; `
	}

	if len(modify) > 0 {
		sqlModify := `ALTER TABLE "#TABLE_NAME_NEW" `
		sqlCommentRename := ``
		for _, nc := range modify {
			for _, ocINF := range tableInfoOld.Fields.([]interface{}) {
				oc := ocINF.(map[string]interface{})
				if nc["idx"].(float64) == oc["idx"].(float64) {
					if nc["column"].(string) != oc["column"].(string) {
						if oc["type"].(string) == "comment" {
							sqlCommentRename += `ALTER TABLE "#TABLE_NAME_OLD" RENAME TO "#TABLE_NAME_NEW"; `

							sqlCommentRename = strings.ReplaceAll(sqlCommentRename, "#TABLE_NAME_OLD", oc["column"].(string)+"_COMMENT")
							sqlCommentRename = strings.ReplaceAll(sqlCommentRename, "#TABLE_NAME_NEW", nc["column"].(string)+"_COMMENT")
						} else {
							sqlModify += ` RENAME COLUMN ` + oc["column"].(string) + ` TO ` + nc["column"].(string) + `, `
						}
					}
					break
				}
			}
		}
		if strings.Contains(sqlModify, "RENAME COLUMN") {
			sqlModify = sqlModify[:len(sqlModify)-2]
		}
		sql += sqlModify + `; `
		sql += sqlCommentRename
	}

	sql = strings.ReplaceAll(sql, "#TABLE_NAME_NEW", tableInfoNew.Table.String)

	log.Println("Sqlite/EditCustomBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteBoard - Delete a board table
func (d *Sqlite) DeleteBoard(tableName string) error {
	sql := `DROP TABLE IF EXISTS "#TABLE_NAME";`
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", tableName)

	log.Println("Sqlite/DeleteBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateComment - Create comment table
func (d *Sqlite) CreateComment(tableInfo models.Board, recreate bool) error {
	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS "#TABLE_NAME";`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS "#TABLE_NAME" (
		"IDX"				INTEGER,
		"BOARD_IDX"			INTEGER,
		"CONTENT"			TEXT,
		"WRITER_IDX"		TEXT,
		"WRITER_NAME"		TEXT,
		"WRITER_PASSWORD"	TEXT,
		"REG_DTTM"			TEXT,
		PRIMARY KEY("IDX" AUTOINCREMENT)
	);`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", tableInfo.Table.String+"_COMMENT")

	log.Println("Sqlite/CreateComment: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteComment - Delete a comment table
func (d *Sqlite) DeleteComment(tableName string) error {
	sql := `DROP TABLE IF EXISTS "#TABLE_NAME";`
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", tableName+"_COMMENT")

	log.Println("Sqlite/DeleteComment: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// AddUserTableFields - Add user column
func (d *Sqlite) AddUserTableFields(fields []models.UserColumn) error {
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

	log.Println("Sqlite/AddUserTableFields: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// EditUserTableFields - Edit user table schema
func (d *Sqlite) EditUserTableFields(fieldsInfoOld []models.UserColumn, fieldsInfoNew []models.UserColumn, notUse []string) error {
	add, remove, modify := diffUserTableFields(fieldsInfoOld, fieldsInfoNew)
	log.Println("User fields Add: ", add)
	log.Println("User fields Remove: ", remove)
	log.Println("User fields Modify: ", modify)

	sql := ""
	if len(add) > 0 {
		for _, a := range add {
			sql += `ALTER TABLE "#TABLE_NAME" `
			sql += ` ADD COLUMN "` + a.ColumnName.String + `" `
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
	}

	if len(remove) > 0 && !funk.Contains(notUse, "remove") {
		sqlRemove := `ALTER TABLE "#TABLE_NAME" `
		for _, r := range remove {
			sqlRemove += ` DROP COLUMN ` + r.ColumnName.String + `, `
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
						sqlModify += `ALTER TABLE "#TABLE_NAME" `
						sqlModify += ` RENAME COLUMN "` + om.ColumnName.String + `" TO "` + nm.ColumnName.String + `"; `
					}
					break
				}
			}
		}
		sql += sqlModify
	}

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTable)

	log.Println("Sqlite/EditUserTableFields: ", sql)

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

func diffUserTableFields(fieldsInfoOld, fieldsInfoNew []models.UserColumn) (add, remove, modify []models.UserColumn) {
	var diff []models.UserColumn

	old := fieldsInfoOld
	new := fieldsInfoNew

	for i := 0; i < 2; i++ {
		for _, s1 := range old {
			found := false
			for _, s2 := range new {
				if s1.Idx == s2.Idx {
					if i == 0 && !cmp.Equal(s1, s2) {
						modify = append(modify, s2)
					}

					found = true
					break
				}
			}

			if !found {
				diff = append(diff, s1)
			}
		}

		if i == 0 {
			remove = diff
			old, new = new, old
		} else {
			add = diff
		}

		diff = []models.UserColumn{}
	}

	return
}

// DeleteUserTableFields - Delete user table field
func (d *Sqlite) DeleteUserTableFields(fieldsInfoRemove []models.UserColumn) error {
	remove := fieldsInfoRemove
	sql := ""

	if len(remove) > 0 {
		for _, r := range remove {
			sql += `ALTER TABLE "#TABLE_NAME" DROP COLUMN "` + r.ColumnName.String + `";`
		}
	}

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTable)

	log.Println("Sqlite/DeleteUserTableFields: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// SelectColumnNames - Get column names of table
func (d *Sqlite) SelectColumnNames(table string) (sql.Result, error) {
	result, err := Dbo.Query("PRAGMA TABLE_INFO(" + table + ")")
	if err != nil {
		return nil, err
	}

	log.Println(result)

	return nil, nil
}
