package db

import (
	"database/sql"
	"log"

	"github.com/practice-golang/9minutes/models"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver"
)

// InsertUserField - Crud
func InsertUserField(data interface{}) (sql.Result, error) {
	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Insert(UserFieldTable).Rows(data)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SelectUserFields - cRud
func SelectUserFields(search interface{}) (interface{}, error) {
	var result interface{}
	var err error

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.From(UserFieldTable).Select(search)

	boardResult := []models.UserColumn{}

	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	err = ds.ScanStructs(&boardResult)
	if err != nil {
		log.Println("ds: ", err.Error())
		return nil, err
	}
	if boardResult != nil {
		result = boardResult
	}

	return result, nil
}

// UpdateUserFields - crUd
func UpdateUserFields(data interface{}) (sql.Result, error) {
	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	whereEXP, err := CheckValidAndPrepareWhere(data)
	if err != nil {
		return nil, err
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Update(UserFieldTable).Set(data).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}
