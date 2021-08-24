package db

import (
	"database/sql"
	"log"

	"github.com/practice-golang/9minutes/models"
	"gopkg.in/guregu/null.v4"

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

	table := ""
	if dbType == "postgres" {
		table = UserFieldTableNoQuotes
	} else {
		table = UserFieldTable
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Insert(table).Rows(data)
	sql, args, _ := ds.ToSQL()
	log.Println("InsertUserField: ", sql, args)

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

	table := ""
	if dbType == "postgres" {
		table = UserFieldTableNoQuotes
	} else {
		table = UserFieldTable
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.From(table).Select(search)

	ex := PrepareWhere(search)
	if !ex.IsEmpty() {
		for c, v := range ex {
			if c == "IDX" || c == "BOARD_IDX" {
				log.Println("SelectUserFields/goqu_ex: ", v)
				ex[c] = goqu.Op{"eq": v}
				// val := fmt.Sprintf("%s", v)
				// ex[c] = goqu.Op{"eq": val}
			} else {
				ex[c] = goqu.Op{"like": v}
				// val := fmt.Sprintf("%s%s%s", "%", v, "%")
				// ex[c] = goqu.Op{"like": val}
			}
		}
		exps = append(exps, ex.Expression())
	}
	ds = ds.Where(goqu.Or(exps...))

	fieldsResult := []models.UserColumn{}

	sql, args, _ := ds.ToSQL()
	log.Println("SelectUserFields: ", sql, args)

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

	table := ""
	if dbType == "postgres" {
		table = UserFieldTableNoQuotes
	} else {
		table = UserFieldTable
	}

	dbms := goqu.New(dbType, Dbo)
	// var ex goqu.Ex
	for _, d := range data.([]models.UserColumn) {
		ex, err := CheckValidAndPrepareWhere(d)
		if err != nil {
			log.Println("UpdateUserFields where: ", err)
			return nil, err
		}

		ds := dbms.Update(table).Set(d).Where(ex)

		sql, args, _ := ds.ToSQL()
		log.Println("UpdateUserFields: ", sql, args)

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

	table := ""
	if dbType == "postgres" {
		table = UserFieldTableNoQuotes
	} else {
		table = UserFieldTable
	}

	whereEXP := goqu.Ex{target: value}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Delete(table).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println("DeleteUserFieldRow: ", sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SelectUserColumnNames - Get column names of table
func SelectUserColumnNames() (interface{}, error) {
	var result interface{}
	var err error
	colResult := []models.UserColumn{}
	colDefault := []models.UserColumn{
		{
			Idx:        null.NewInt(-1, true),
			Name:       null.NewString("Idx", true),
			Code:       null.NewString("idx", true),
			Type:       null.NewString("text", true),
			ColumnName: null.NewString("IDX", true),
			Order:      null.NewInt(1, true),
		},
		{
			Idx:        null.NewInt(-1, true),
			Name:       null.NewString("User name", true),
			Code:       null.NewString("username", true),
			Type:       null.NewString("text", true),
			ColumnName: null.NewString("USERNAME", true),
			Order:      null.NewInt(2, true),
		},
		{
			Idx:        null.NewInt(-1, true),
			Name:       null.NewString("Password", true),
			Code:       null.NewString("password", true),
			Type:       null.NewString("text", true),
			ColumnName: null.NewString("PASSWORD", true),
			Order:      null.NewInt(3, true),
		},
		{
			Idx:        null.NewInt(-1, true),
			Name:       null.NewString("Email", true),
			Code:       null.NewString("email", true),
			Type:       null.NewString("text", true),
			ColumnName: null.NewString("EMAIL", true),
			Order:      null.NewInt(4, true),
		},
		{
			Idx:        null.NewInt(-1, true),
			Name:       null.NewString("Admin", true),
			Code:       null.NewString("admin", true),
			Type:       null.NewString("text", true),
			ColumnName: null.NewString("ADMIN", true),
			Order:      null.NewInt(5, true),
		},
		{
			Idx:        null.NewInt(-1, true),
			Name:       null.NewString("Approval", true),
			Code:       null.NewString("approval", true),
			Type:       null.NewString("text", true),
			ColumnName: null.NewString("APPROVAL", true),
			Order:      null.NewInt(6, true),
		},
		{
			Idx:        null.NewInt(-1, true),
			Name:       null.NewString("Reg datetime", true),
			Code:       null.NewString("reg-dttm", true),
			Type:       null.NewString("text", true),
			ColumnName: null.NewString("REG_DTTM", true),
			Order:      null.NewInt(7, true),
		},
	}

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	table := ""
	if dbType != "sqlite3" {
		if dbType == "postgres" {
			table = UserFieldTableNoQuotes
		} else if dbType == "sqlserver" {
			table = DatabaseName + ".dbo." + UserFieldTableName
		} else {
			table = UserFieldTable
		}
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.From(table).Select([]models.UserColumn{})

	sql, args, _ := ds.ToSQL()
	log.Println("SelectUserColumnNames: ", sql, args)

	err = ds.ScanStructs(&colResult)
	if err != nil {
		log.Println("ds: ", err.Error())
		return nil, err
	}

	colResult = append(colDefault, colResult...)

	result = colResult

	return result, nil
}

// InsertUser - Crud
func InsertUser(data interface{}) (sql.Result, error) {
	rcds := []goqu.Record{}

	if data != nil {
		rcd := goqu.Record{}
		for k, d := range data.(map[string]interface{}) {
			rcd[k] = d
		}
		rcds = append(rcds, rcd)
	}

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	table := ""
	if dbType == "postgres" {
		table = UserTableNoQuotes
	} else {
		table = UserTable
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Insert(table).Rows(rcds)
	sql, args, _ := ds.ToSQL()
	log.Println("InsertUser: ", sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateUser - crUd
func UpdateUser(data interface{}) (sql.Result, error) {
	whereEXP := goqu.Ex{}

	rcd := goqu.Record{}
	if data != nil {
		for k, d := range data.(map[string]interface{}) {
			if k == "IDX" {
				whereEXP[k] = d
				continue
			}
			rcd[k] = d
		}
	}

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	table := ""
	if dbType == "postgres" {
		table = UserTableNoQuotes
	} else {
		table = UserTable
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Update(table).Set(rcd).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println("UpdateUser: ", sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteUser - cruD
func DeleteUser(idx string) (sql.Result, error) {
	whereEXP := goqu.Ex{"IDX": idx}

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	table := ""
	if dbType == "postgres" {
		table = UserTableNoQuotes
	} else {
		table = UserTable
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Delete(table).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println("DeleteUser: ", sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}
