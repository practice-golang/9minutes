package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/practice-golang/9minutes/models"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver"
)

// InsertUserField - Crud
func InsertUserField(data []models.UserColumn) (sql.Result, error) {
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

	exps := []goqu.Expression{}

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.From(UserFieldTable).Select(search)

	ex := PrepareWhere(search)
	if !ex.IsEmpty() {
		for c, v := range ex {
			val := fmt.Sprintf("%d", v)
			ex[c] = goqu.Op{"eq": val}
		}
		exps = append(exps, ex.Expression())
	}
	ds = ds.Where(goqu.Or(exps...))

	fieldsResult := []models.UserColumn{}

	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	err = ds.ScanStructs(&fieldsResult)
	if err != nil {
		log.Println("ds: ", err.Error())
		return nil, err
	}
	if fieldsResult != nil {
		result = fieldsResult
	}

	return result, nil
}

// UpdateUserFields - crUd
func UpdateUserFields(data interface{}) (sql.Result, error) {
	var result sql.Result

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	log.Println("WTF??????")

	dbms := goqu.New(dbType, Dbo)
	// var ex goqu.Ex
	for _, d := range data.([]models.UserColumn) {
		ex, err := CheckValidAndPrepareWhere(d)
		if err != nil {
			log.Println("UpdateUserFields where: ", err)
			return nil, err
		}

		ds := dbms.Update(UserFieldTable).Set(d).Where(ex)

		sql, args, _ := ds.ToSQL()
		log.Println(sql, args)

		result, err = Dbo.Exec(sql)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func DeleteUserFieldRow(target, value string) (sql.Result, error) {
	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	whereEXP := goqu.Ex{target: value}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Delete(UserFieldTable).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}
