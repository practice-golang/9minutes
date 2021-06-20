package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/practice-golang/9minutes/models"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver"
)

func getDialect() (string, error) {
	var dbType string

	switch DBType {
	case SQLITE:
		dbType = "sqlite3"
	case MYSQL:
		dbType = "mysql"
	case POSTGRES:
		dbType = "postgres"
	case SQLSERVER:
		dbType = "sqlserver"
	default:
		return dbType, errors.New("nothing to support DB")
	}

	return dbType, nil
}

// InsertData - Crud
func InsertData(data interface{}) (sql.Result, error) {
	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Insert(BoardManagerTable).Rows(data)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SelectData - cRud
func SelectData(search interface{}) (interface{}, error) {
	var result interface{}
	var err error

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.From(BoardManagerTable).Select(search)

	boardResult := []models.Board{}

	switch search := search.(type) {
	case models.Board:
		exps := PrepareWhere(search)
		if !exps.IsEmpty() {
			ds = ds.Where(exps)
		}
		// cnt := listCount
		cnt := uint(10)
		ds = ds.Limit(cnt)

	case models.BoardSearch:
		keywords := search.Keywords
		exps := []goqu.Expression{}
		for _, k := range keywords {
			ex := PrepareWhere(k)
			if !ex.IsEmpty() {
				for c, v := range ex {
					if c == "IDX" {
						val := fmt.Sprintf("%s", v)
						ex[c] = goqu.Op{"eq": val}
					} else {
						val := fmt.Sprintf("%s%s%s", "%", v, "%")
						ex[c] = goqu.Op{"like": val}
					}
				}
				exps = append(exps, ex.Expression())
			}
		}
		ds = ds.Where(goqu.Or(exps...))

		orderDirection := goqu.C(OrderScope).Asc()
		if search.Options.Order.String == "desc" {
			orderDirection = goqu.C(OrderScope).Desc()
		}
		ds = ds.Order(orderDirection)

		cnt := listCount
		if search.Options.Count.Valid {
			cnt = uint(search.Options.Count.Int64)
		}
		ds = ds.Limit(cnt)

		offset := uint(0)
		if search.Options.Page.Valid {
			offset = uint(search.Options.Page.Int64)
		}
		ds = ds.Offset(offset * cnt)
	}

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

// UpdateData - crUd
func UpdateData(data interface{}) (sql.Result, error) {
	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	whereEXP, err := CheckValidAndPrepareWhere(data)
	if err != nil {
		return nil, err
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Update(BoardManagerTable).Set(data).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteData - cruD
func DeleteData(target, value string) (sql.Result, error) {
	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	whereEXP := goqu.Ex{target: value}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Delete(BoardManagerTable).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SelectCount - data count -> pages = (data count) / (count per page)
func SelectCount(search interface{}) (uint, error) {
	var cnt uint

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.From(BoardManagerTable).Select(goqu.COUNT("*").As("PAGE_COUNT"))

	switch search := search.(type) {
	case models.BoardSearch:
		keywords := search.Keywords
		exps := []goqu.Expression{}
		for _, k := range keywords {
			ex := PrepareWhere(k)
			if !ex.IsEmpty() {
				for c, v := range ex {
					if c == "IDX" {
						val := fmt.Sprintf("%s", v)
						ex[c] = goqu.Op{"eq": val}
					} else {
						val := fmt.Sprintf("%s%s%s", "%", v, "%")
						ex[c] = goqu.Op{"like": val}
					}
				}
				exps = append(exps, ex.Expression())
			}
		}
		ds = ds.Where(goqu.Or(exps...))
	}

	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	ds.ScanVal(&cnt)

	return cnt, nil
}

/* Comments */

// SelectComments - cRud
func SelectComments(search interface{}) (interface{}, error) {
	var result interface{}
	var err error

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	var searchBytes models.CommentSearch
	err = json.Unmarshal(search.([]byte), &searchBytes)
	if err != nil {
		return nil, err
	}

	dbms := goqu.New(dbType, Dbo)
	var ds *goqu.SelectDataset
	commentResult := []models.Comments{}
	exps := []goqu.Expression{}

	keywords := searchBytes.Keywords

	table := searchBytes.Table.String + "_COMMENT"
	if searchBytes.Options.Count.Int64 > 1 {
		// Comment list
		ds = dbms.From(table).Select(models.Comments{})
	} else {
		// Not allow getting all of list
		if !searchBytes.Options.Count.Valid || searchBytes.Options.Count.Int64 < 1 {
			searchBytes.Options.Count.SetValid(1)
		}

		// Comment
		ds = dbms.From(table).Select(models.Comments{})
	}

	for _, k := range keywords {
		ex := PrepareWhere(k)
		if !ex.IsEmpty() {
			for c, v := range ex {
				if c == "IDX" {
					val := fmt.Sprintf("%s", v)
					ex[c] = goqu.Op{"eq": val}
				} else {
					val := fmt.Sprintf("%s%s%s", "%", v, "%")
					ex[c] = goqu.Op{"like": val}
				}
			}
			exps = append(exps, ex.Expression())
		}
	}

	ds = ds.Where(goqu.Or(exps...))

	orderDirection := goqu.C(OrderScope).Asc()
	if searchBytes.Options.Order.String == "desc" {
		orderDirection = goqu.C(OrderScope).Desc()
	}
	ds = ds.Order(orderDirection)

	cnt := listCount
	if searchBytes.Options.Count.Valid {
		cnt = uint(searchBytes.Options.Count.Int64)
	}
	ds = ds.Limit(cnt)

	offset := uint(0)
	if searchBytes.Options.Page.Valid {
		offset = uint(searchBytes.Options.Page.Int64)
	}
	ds = ds.Offset(offset * cnt)

	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	err = ds.ScanStructs(&commentResult)
	if err != nil {
		log.Println("ds: ", err.Error())
		return nil, err
	}
	if commentResult != nil {
		result = commentResult
	}

	return result, nil
}

// InsertComment - Crud comment
func InsertComment(data interface{}, table string) (sql.Result, error) {
	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Insert(table).Rows(data)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}
