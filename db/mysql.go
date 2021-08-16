package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/practice-golang/9minutes/models"
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
		for k, f := range fields {
			log.Println(k, f.Name.String, f.Type.String, f.Order.Int64)
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

			default:
				colType = "VARCHAR(128)"
			}

			sql += fmt.Sprintf(`%s		`+"`%s`"+`		%s,`, "\n", f.ColumnName.String, colType)
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
	return nil
}

// EditCustomBoard - Create custom table
func (d *Mysql) EditCustomBoard(tableInfoOld models.Board, tableInfoNew models.Board) error {
	return nil
}

// DeleteBoard - Delete a board table
func (d *Mysql) DeleteBoard(tableName string) error {
	return nil
}

// CreateComment - Create comment table
func (d *Mysql) CreateComment(tableInfo models.Board, recreate bool) error {
	return nil
}

// EditUserTableFields - Edit user table schema
func (d *Mysql) EditUserTableFields(fieldsInfoOld []models.UserColumn, fieldsInfoNew []models.UserColumn, notUse []string) error {
	return nil
}

// DeleteUserTableFields - Delete user table field
func (d *Mysql) DeleteUserTableFields(fieldsInfoRemove []models.UserColumn) error {
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
func (d *Mysql) AddUserTableFields(fields []models.UserColumn) error {
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
func (d *Mysql) SelectColumnNames(table string) (sql.Result, error) {
	result, err := Dbo.Exec("PRAGMA TABLE_INFO(" + table + ")")
	if err != nil {
		return nil, err
	}

	return result, nil
}
