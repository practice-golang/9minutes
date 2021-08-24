package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/practice-golang/9minutes/models"
	"github.com/thoas/go-funk"
)

type Mysql struct{ Dsn string }

// initDB - Prepare DB
func (d *Mysql) initDB() (*sql.DB, error) {
	var err error

	Dbo, err = sql.Open("mysql", d.Dsn)
	if err != nil {
		return nil, err
	}

	return Dbo, nil
}

func (d *Mysql) CreateDB() error {
	return nil
}

// CreateBoardManagerTable - Create board manager table
func (d *Mysql) CreateBoardManagerTable(recreate bool) error {
	sql := `CREATE DATABASE IF NOT EXISTS #DATABASE_NAME;`
	sql = strings.ReplaceAll(sql, "#DATABASE_NAME", DatabaseName)
	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	sql = ""
	if recreate {
		sql += `DROP TABLE IF EXISTS #TABLE_NAME;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS #TABLE_NAME (
		IDX INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
		NAME VARCHAR(128) NULL DEFAULT NULL,
		CODE VARCHAR(64) NULL DEFAULT NULL,
		TYPE VARCHAR(64) NULL DEFAULT NULL,
		` + "`TABLE`" + ` VARCHAR(64) NULL DEFAULT NULL,
		GRANT_READ VARCHAR(16) NULL DEFAULT NULL,
		GRANT_WRITE VARCHAR(16) NULL DEFAULT NULL,
		GRANT_COMMENT VARCHAR(16) NULL DEFAULT NULL,
		FILE_UPLOAD VARCHAR(2) NULL DEFAULT NULL,
		FIELDS TEXT NULL DEFAULT NULL,

		PRIMARY KEY (IDX),
		INDEX IDX (IDX)
	)
	COLLATE='utf8_general_ci'
	ENGINE=InnoDB;`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", BoardManagerTable)

	_, err = Dbo.Exec(sql)
	if err != nil {
		log.Println(sql)
		return err
	}

	return nil
}

// CreateUserFieldTable - Create user manager table
func (d *Mysql) CreateUserFieldTable(recreate bool) error {
	sql := `CREATE DATABASE IF NOT EXISTS #DATABASE_NAME;`
	sql = strings.ReplaceAll(sql, "#DATABASE_NAME", DatabaseName)
	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	sql = ""
	if recreate {
		sql += `DROP TABLE IF EXISTS #TABLE_NAME;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS #TABLE_NAME (
		IDX INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
		NAME VARCHAR(128) NULL DEFAULT NULL,
		CODE VARCHAR(64) NULL DEFAULT NULL,
		TYPE VARCHAR(64) NULL DEFAULT NULL,
		COLUMN_NAME VARCHAR(64) NULL DEFAULT NULL,
		` + "`ORDER`" + ` INT(11) NULL DEFAULT NULL,

		PRIMARY KEY (IDX),
		INDEX IDX (IDX)
	)
	COLLATE='utf8_general_ci'
	ENGINE=InnoDB;`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserFieldTable)

	_, err = Dbo.Exec(sql)
	if err != nil {
		log.Println(sql)
		return err
	}

	return nil
}

// CreateUserTable - Create user table
func (d *Mysql) CreateUserTable(recreate bool) error {
	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS #TABLE_NAME;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS #TABLE_NAME (
		IDX INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
		USERNAME VARCHAR(128) NULL DEFAULT NULL,
		PASSWORD VARCHAR(128) NULL DEFAULT NULL,
		EMAIL VARCHAR(128) NULL DEFAULT NULL,
		ADMIN VARCHAR(2) NULL DEFAULT NULL,
		APPROVAL VARCHAR(2) NULL DEFAULT NULL,
		REG_DTTM INT(14) UNSIGNED NULL DEFAULT NULL,

		PRIMARY KEY(IDX),
		UNIQUE INDEX USERNAME (USERNAME),
		UNIQUE INDEX EMAIL (EMAIL),
		INDEX IDX (IDX)
	)
	COLLATE='utf8_general_ci'
	ENGINE=InnoDB;`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTable)

	_, err := Dbo.Exec(sql)
	if err != nil {
		log.Println(sql)
		return err
	}

	// Add temp admin
	sql = `
	INSERT IGNORE INTO #TABLE_NAME (USERNAME, ` + "`PASSWORD`" + `, EMAIL, ` + "`ADMIN`" + `, APPROVAL)
		VALUES ("admin", "admin", "admin@please.modify", "Y", "Y");`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTable)

	_, err = Dbo.Exec(sql)
	if err != nil {
		log.Println(sql)
		return err
	}

	return nil
}

// CreateBasicBoard - Create board table
func (d *Mysql) CreateBasicBoard(tableInfo models.Board, recreate bool) error {
	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS #TABLE_NAME;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS #TABLE_NAME (
		IDX INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
		TITLE VARCHAR(256) NULL DEFAULT NULL,
		CONTENT TEXT NULL DEFAULT NULL,
		IS_MEMBER VARCHAR(2) NULL DEFAULT NULL,
		WRITER_IDX VARCHAR(11) NULL DEFAULT NULL,
		WRITER_NAME VARCHAR(64) NULL DEFAULT NULL,
		WRITER_PASSWORD VARCHAR(128) NULL DEFAULT NULL,
		FILES TEXT NULL DEFAULT NULL,
		REG_DTTM INT(14) UNSIGNED NULL DEFAULT NULL,

		PRIMARY KEY(IDX),
		INDEX IDX (IDX)
	);`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", DatabaseName+"."+tableInfo.Table.String)

	log.Println("MySQL/CreateBasicBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// CreateCustomBoard - Create board table
func (d *Mysql) CreateCustomBoard(tableInfo models.Board, fields []models.Field, recreate bool) error {
	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS #TABLE_NAME;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS #TABLE_NAME (
		IDX INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
		IS_MEMBER VARCHAR(2) NULL DEFAULT NULL,
		WRITER_IDX VARCHAR(11) NULL DEFAULT NULL,
		WRITER_NAME VARCHAR(64) NULL DEFAULT NULL,
		WRITER_PASSWORD VARCHAR(128) NULL DEFAULT NULL,
		FILES TEXT NULL DEFAULT NULL,
		REG_DTTM INT(14) UNSIGNED NULL DEFAULT NULL,`

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
				colType = "INT(16)"
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

			sql += fmt.Sprintf(`%s		`+"`%s`"+`		%s,`, "\n", f.ColumnName.String, colType)
		}

		if commentCount > 1 {
			return errors.New("available only 1 comment")
		}
	}

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", DatabaseName+"."+tableInfo.Table.String)

	sql += `
		PRIMARY KEY(IDX),
		INDEX IDX (IDX)
	);`

	log.Println("MySQL/CreateCustomBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// EditBasicBoard - Create board table
func (d *Mysql) EditBasicBoard(tableInfoOld models.Board, tableInfoNew models.Board) error {
	// sql := `RENAME TABLE ` + "`#TABLE_NAME_OLD`" + ` TO ` + "`#TABLE_NAME_NEW`" + ` ;`
	sql := `ALTER TABLE ` + "`#TABLE_NAME_OLD`" + ` RENAME TO ` + "`#TABLE_NAME_NEW`" + `;`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME_OLD", DatabaseName+"`.`"+tableInfoOld.Table.String)
	sql = strings.ReplaceAll(sql, "#TABLE_NAME_NEW", DatabaseName+"`.`"+tableInfoNew.Table.String)

	log.Println("MySQL/EditBasicBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// EditCustomBoard - Create custom table
func (d *Mysql) EditCustomBoard(tableInfoOld models.Board, tableInfoNew models.Board) error {
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

	commentTableCNT := 0
	sqlCommentExist := `SELECT COUNT(TABLE_NAME) FROM information_schema.tables WHERE table_schema = "#DATABASE" AND TABLE_NAME = "#TABLE_NAME";`
	sqlCommentExist = strings.ReplaceAll(sqlCommentExist, "#DATABASE", DatabaseName)
	sqlCommentExist = strings.ReplaceAll(sqlCommentExist, "#TABLE_NAME", tableInfoOld.Table.String+"_COMMENT")
	row := Dbo.QueryRow(sqlCommentExist)
	err := row.Scan(&commentTableCNT)
	if err != nil {
		return err
	}

	sql := ""
	sqlTableRename := ""
	sqlCommentRename := ""
	if tableInfoOld.Table.String != tableInfoNew.Table.String {
		sqlTableRename = "ALTER TABLE `#TABLE_NAME_OLD` RENAME TO `#TABLE_NAME_NEW`;"

		sqlTableRename = strings.ReplaceAll(sqlTableRename, "#TABLE_NAME_OLD", DatabaseName+"`.`"+tableInfoOld.Table.String)
		sqlTableRename = strings.ReplaceAll(sqlTableRename, "#TABLE_NAME_NEW", DatabaseName+"`.`"+tableInfoNew.Table.String)

		if commentTableCNT > 0 {
			sqlCommentRename += "ALTER TABLE `#TABLE_NAME_OLD` RENAME TO `#TABLE_NAME_NEW`;"

			sqlCommentRename = strings.ReplaceAll(sqlCommentRename, "#TABLE_NAME_OLD", DatabaseName+"`.`"+tableInfoOld.Table.String+"_COMMENT`")
			sqlCommentRename = strings.ReplaceAll(sqlCommentRename, "#TABLE_NAME_NEW", DatabaseName+"`.`"+tableInfoNew.Table.String+"_COMMENT`")
		}
	}

	sqlAdd := ""
	if len(add) > 0 {
		for _, c := range add {
			sqlAdd += " ADD COLUMN `" + c["column"].(string) + "` "
			switch c["type"].(string) {
			// cusom-tablelist
			case "text":
				sqlAdd += `TEXT`
			case "number":
				sqlAdd += `INT(16)`
			case "real":
				sqlAdd += `DECIMAL(20,20)`

			// cusom-board
			case "title", "author", "input":
				sqlAdd += `VARCHAR(512)`
			case "editor":
				sqlAdd += `TEXT`
			case "comment":
				sqlAdd += `VARCHAR(4)`
				if commentTableCNT == 0 {
					_ = d.CreateComment(tableInfoNew, false)
				}
			default:
				sqlAdd += `VARCHAR(128)`
			}

			sqlAdd += ` NULL, `
			// if i < (len(add) - 1) {
			// 	sqlAdd += `, `
			// }
		}
		// sqlAdd += `; `
	}

	sqlRemove := ""
	if len(remove) > 0 {
		for _, c := range remove {
			sqlRemove += " DROP COLUMN `" + c["column"].(string) + "`, "

			if c["type"].(string) == "comment" {
				commentCount = 0
				sqlCommentRename = ""
				d.DeleteComment(tableInfoOld.Table.String)
			}
		}
		// if strings.Contains(sqlRemove, "DROP COLUMN") {
		// 	sqlRemove = sqlRemove[:len(sqlRemove)-2]
		// }
		// sql += sqlRemove + `; `
	}

	sqlModify := ""
	if len(modify) > 0 {
		for _, nc := range modify {
			for _, ocINF := range tableInfoOld.Fields.([]interface{}) {
				oc := ocINF.(map[string]interface{})
				if nc["idx"].(float64) == oc["idx"].(float64) {
					if nc["column"].(string) != oc["column"].(string) {
						if oc["type"].(string) == "comment" {
							if commentCount > 0 {
								sqlModify += "CHANGE COLUMN `" + oc["column"].(string) + "` `" + nc["column"].(string) + "` VARCHAR(4) NULL, "
							}
						} else {
							sqlModify += "CHANGE COLUMN `" + oc["column"].(string) + "` `" + nc["column"].(string) + "` "
							switch nc["type"].(string) {
							// cusom-tablelist
							case "text":
								sqlModify += `TEXT`
							case "number":
								sqlModify += `INT(16)`
							case "real":
								sqlModify += `DECIMAL(20,20)`

							// cusom-board
							case "title", "author", "input":
								sqlModify += `VARCHAR(512)`
							case "editor":
								sqlModify += `TEXT`
							case "comment":
								sqlModify += `VARCHAR(4)`
								_ = d.CreateComment(tableInfoNew, false)
							default:
								sqlModify += `VARCHAR(128)`
							}

							sqlModify += " NULL"
							sqlModify += ", "
						}
					}
					break
				}
			}
		}
	}

	if sqlTableRename != "" {
		log.Println("MySQL/EditCustomBoard: ", sql)

		_, err := Dbo.Exec(sqlTableRename)
		if err != nil {
			return err
		}
	}
	if sqlCommentRename != "" {
		log.Println("MySQL/EditCustomBoard: ", sql)

		_, err := Dbo.Exec(sqlCommentRename)
		if err != nil {
			return err
		}
	}
	if sqlAdd != "" {
		sql = `
		ALTER TABLE ` + "`#TABLE_NAME_NEW`" + "\n" + sqlAdd
		sql = sql[:len(sql)-2] + ";"

		sql = strings.ReplaceAll(sql, "#TABLE_NAME_NEW", DatabaseName+"`.`"+tableInfoNew.Table.String)

		log.Println("MySQL/EditCustomBoard: ", sql)

		_, err := Dbo.Exec(sql)
		if err != nil {
			return err
		}
	}
	if sqlRemove != "" {
		sql = `
		ALTER TABLE ` + "`#TABLE_NAME_NEW`" + "\n" + sqlRemove
		sql = sql[:len(sql)-2] + ";"

		sql = strings.ReplaceAll(sql, "#TABLE_NAME_NEW", DatabaseName+"`.`"+tableInfoNew.Table.String)

		log.Println("MySQL/EditCustomBoard: ", sql)

		_, err := Dbo.Exec(sql)
		if err != nil {
			return err
		}
	}
	if sqlModify != "" {
		sql = `
		ALTER TABLE ` + "`#TABLE_NAME_NEW`" + "\n" + sqlModify
		sql = sql[:len(sql)-2] + ";"

		sql += sqlCommentRename

		sql = strings.ReplaceAll(sql, "#TABLE_NAME_NEW", DatabaseName+"`.`"+tableInfoNew.Table.String)

		log.Println("MySQL/EditCustomBoard: ", sql)

		_, err := Dbo.Exec(sql)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteBoard - Delete a board table
func (d *Mysql) DeleteBoard(tableName string) error {
	sql := `DROP TABLE IF EXISTS ` + "`#TABLE_NAME`" + `;`
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", DatabaseName+"`.`"+tableName)

	log.Println("MySQL/DeleteBoard: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	_ = d.DeleteComment(tableName)

	return nil
}

// CreateComment - Create comment table
func (d *Mysql) CreateComment(tableInfo models.Board, recreate bool) error {
	sql := ""
	if recreate {
		sql += `DROP TABLE IF EXISTS #TABLE_NAME;`
	}
	sql += `
	CREATE TABLE IF NOT EXISTS #TABLE_NAME (
		IDX INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
		BOARD_IDX INT(11) UNSIGNED NOT NULL,
		CONTENT TEXT NULL DEFAULT NULL,
		IS_MEMBER VARCHAR(2) NULL DEFAULT NULL,
		WRITER_IDX VARCHAR(11) NULL DEFAULT NULL,
		WRITER_NAME VARCHAR(64) NULL DEFAULT NULL,
		WRITER_PASSWORD VARCHAR(128) NULL DEFAULT NULL,
		REG_DTTM INT(14) UNSIGNED NULL DEFAULT NULL,

		PRIMARY KEY(IDX),
		INDEX IDX (IDX)
	);`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", DatabaseName+"."+tableInfo.Table.String+"_COMMENT")

	log.Println("MySQL/CreateComment: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// Not use yet
// EditComment - Edit a comment table
func (d *Mysql) EditComment(tableInfoOld models.Board, tableInfoNew models.Board) error {
	sql := `ALTER TABLE ` + "`#TABLE_NAME_OLD`" + ` RENAME TO ` + "`#TABLE_NAME_NEW`" + `;`

	sql = strings.ReplaceAll(sql, "#TABLE_NAME_OLD", DatabaseName+"`.`"+tableInfoOld.Table.String+"_COMMENT")
	sql = strings.ReplaceAll(sql, "#TABLE_NAME_NEW", DatabaseName+"`.`"+tableInfoNew.Table.String+"_COMMENT")

	log.Println("MySQL/EditComment: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// DeleteComment - Delete a comment table
func (d *Mysql) DeleteComment(tableName string) error {
	sql := `DROP TABLE IF EXISTS ` + "`#TABLE_NAME`" + `;`
	sql = strings.ReplaceAll(sql, "#TABLE_NAME", DatabaseName+"`.`"+tableName+"_COMMENT")

	log.Println("MySQL/DeleteComment: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// EditUserTableFields - Edit user table schema
func (d *Mysql) EditUserTableFields(fieldsInfoOld []models.UserColumn, fieldsInfoNew []models.UserColumn, notUse []string) error {
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
				sql += ` INT(16)`
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

	log.Println("MySQL/EditUserTableFields: ", sql)

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
func (d *Mysql) DeleteUserTableFields(fieldsInfoRemove []models.UserColumn) error {
	remove := fieldsInfoRemove
	sql := ""

	if len(remove) > 0 {
		sqlRemove := "ALTER TABLE #TABLE_NAME"
		for _, r := range remove {
			sqlRemove += ` DROP COLUMN ` + r.ColumnName.String + `, `
		}
		if strings.Contains(sqlRemove, "DROP COLUMN") {
			sqlRemove = sqlRemove[:len(sqlRemove)-2]
		}
		sql += sqlRemove + `;`
	}

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTable)

	log.Println("MySQL/DeleteUserTableFields: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// AddUserTableFields - Add user column
func (d *Mysql) AddUserTableFields(fields []models.UserColumn) error {
	sql := ""
	for _, a := range fields {
		sql += "ALTER TABLE #TABLE_NAME ADD COLUMN `" + a.ColumnName.String + "`"
		switch a.Type.String {
		case "text":
			sql += ` VARCHAR(256)`
		case "number":
			sql += ` INT(16)`
		case "real":
			sql += ` DECIMAL(20,20)`
		}

		sql += " NULL"
		sql += `;`
	}

	sql = strings.ReplaceAll(sql, "#TABLE_NAME", UserTable)

	log.Println("MySQL/EditUserTableFields: ", sql)

	_, err := Dbo.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// SelectColumnNames - Get column names of table
func (d *Mysql) SelectColumnNames(table string) (sql.Result, error) {
	result, err := Dbo.Exec("PRAGMA TABLE_INFO(" + table + ")")
	if err != nil {
		return nil, err
	}

	return result, nil
}
