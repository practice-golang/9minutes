package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/practice-golang/9minutes/config"
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

	table := ""
	if dbType == "postgres" {
		table = BoardManagerTableNoQuotes
	} else {
		table = BoardManagerTable
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Insert(table).Rows(data)
	sql, args, _ := ds.ToSQL()
	log.Println("InsertData: ", sql, args)

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

	table := ""
	if dbType == "postgres" {
		table = BoardManagerTableNoQuotes
	} else {
		table = BoardManagerTable
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.From(table).Select(search)

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
					if c == "IDX" || c == "BOARD_IDX" {
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

	table := ""
	if dbType == "postgres" {
		table = BoardManagerTableNoQuotes
	} else {
		table = BoardManagerTable
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Update(table).Set(data).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println("DB UpdateData: ", sql, args)

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

	table := ""
	if dbType == "postgres" {
		table = BoardManagerTableNoQuotes
	} else {
		table = BoardManagerTable
	}

	whereEXP := goqu.Ex{target: value}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Delete(table).Where(whereEXP)
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

	table := ""
	if dbType == "postgres" {
		table = BoardManagerTableNoQuotes
	} else {
		table = BoardManagerTable
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.From(table).Select(goqu.COUNT("*").As("PAGE_COUNT"))

	switch search := search.(type) {
	case models.BoardSearch:
		keywords := search.Keywords
		exps := []goqu.Expression{}
		for _, k := range keywords {
			ex := PrepareWhere(k)
			if !ex.IsEmpty() {
				for c, v := range ex {
					if c == "IDX" || c == "BOARD_IDX" {
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
	commentResult := []models.CommentList{}
	exps := []goqu.Expression{}

	keywords := searchBytes.Keywords

	table := searchBytes.Table.String + "_COMMENT"
	if dbType != "sqlite3" {
		if dbType == "postgres" {
			table = config.DbInfo.Schema + "." + table
		} else if dbType == "sqlserver" {
			table = DatabaseName + ".dbo." + table
		} else {
			table = DatabaseName + "." + table
		}
	}
	if searchBytes.Options.Count.Int64 > 1 {
		// Comment list
		ds = dbms.From(table).Select(models.CommentList{})
	} else {
		// Not allow getting all of list
		if !searchBytes.Options.Count.Valid || searchBytes.Options.Count.Int64 < 1 {
			searchBytes.Options.Count.SetValid(1)
		}

		// Comment
		ds = dbms.From(table).Select(models.CommentList{})
	}

	for _, k := range keywords {
		ex := PrepareWhere(k)
		if !ex.IsEmpty() {
			for c, v := range ex {
				if c == "IDX" || c == "BOARD_IDX" {
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

	tbl := table
	if dbType != "sqlite3" {
		if dbType == "postgres" {
			tbl = config.DbInfo.Schema + "." + tbl
		} else if dbType == "sqlserver" {
			tbl = DatabaseName + ".dbo." + tbl
		} else {
			tbl = DatabaseName + "." + tbl
		}
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Insert(tbl).Rows(data)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateComment - crUd comment
func UpdateComment(data interface{}, table string) (sql.Result, error) {
	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	whereEXP := goqu.Ex{
		"IDX": data.(models.CommentSET).Idx,
	}
	if data.(models.CommentSET).WriterName.Valid {
		whereEXP["WRITER_NAME"] = data.(models.CommentSET).WriterName
	}
	whereEXP["WRITER_PASSWORD"] = goqu.Op{"eq": nil}
	if data.(models.CommentSET).WriterPassword.Valid {
		whereEXP["WRITER_PASSWORD"] = data.(models.CommentSET).WriterPassword
	}

	tbl := table
	if dbType != "sqlite3" {
		if dbType == "postgres" {
			tbl = config.DbInfo.Schema + "." + tbl
		} else if dbType == "sqlserver" {
			tbl = DatabaseName + ".dbo." + tbl
		} else {
			tbl = DatabaseName + "." + tbl
		}
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Update(tbl).Set(data).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteComment - cruD comment
func DeleteComment(data interface{}, table string) (sql.Result, error) {
	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	whereEXP := goqu.Ex{
		"IDX": data.(models.CommentSET).Idx,
	}
	if data.(models.CommentSET).WriterName.Valid {
		whereEXP["WRITER_NAME"] = data.(models.CommentSET).WriterName
		whereEXP["WRITER_PASSWORD"] = nil
	} else {
		whereEXP["WRITER_PASSWORD"] = data.(models.CommentSET).WriterPassword
	}

	tbl := table
	if dbType != "sqlite3" {
		if dbType == "postgres" {
			tbl = config.DbInfo.Schema + "." + tbl
		} else if dbType == "sqlserver" {
			tbl = DatabaseName + ".dbo." + tbl
		} else {
			tbl = DatabaseName + "." + tbl
		}
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Delete(tbl).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}
