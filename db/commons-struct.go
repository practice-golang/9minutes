package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/practice-golang/9minutes/config"
	"github.com/practice-golang/9minutes/models"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver"
	"github.com/doug-martin/goqu/v9/exp"
)

// InsertContents - Crud contents / basic-board
func InsertContents(data interface{}, table string) (map[string]string, error) {
	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	log.Println("InsertContents: ", dbType)

	tbl := table
	if dbType != "sqlite3" {
		if dbType == "postgres" {
			tbl = config.DbInfo.Schema + "." + tbl
		} else {
			tbl = DatabaseName + "." + tbl
		}
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Insert(tbl).Rows(data)

	var sql string
	var args []interface{}
	var lastID, affRows int64

	if dbType == "postgres" {
		ds = ds.Returning("IDX")
		sql, args, _ = ds.ToSQL()

		var idx int64
		rows, err := Dbo.Query(sql, args...)
		if err != nil {
			return nil, err
		}

		var idxs []int64
		for rows.Next() {
			_ = rows.Scan(&idx)
			idxs = append(idxs, idx)
		}

		lastID = idx
		affRows = int64(len(idxs))
	} else {
		sql, args, _ = ds.ToSQL()
		sqlResult, err := Dbo.Exec(sql)
		if err != nil {
			return nil, err
		}

		lastID, _ = sqlResult.LastInsertId()
		affRows, _ = sqlResult.RowsAffected()
	}

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	log.Println("InsertContents: ", sql, args)

	return result, nil
}

// UpdateContents - Crud contents / basic-board
func UpdateContents(data interface{}, table string) (sql.Result, error) {
	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	whereEXP := goqu.Ex{
		"IDX": data.(models.ContentsBasicBoardSET).Idx,
	}
	if data.(models.ContentsBasicBoardSET).WriterName.Valid {
		whereEXP["WRITER_NAME"] = data.(models.ContentsBasicBoardSET).WriterName
	}
	if data.(models.ContentsBasicBoardSET).WriterPassword.Valid {
		whereEXP["WRITER_PASSWORD"] = data.(models.ContentsBasicBoardSET).WriterPassword
	}

	tbl := table
	if dbType != "sqlite3" {
		if config.DbInfo.Type == "postgres" {
			tbl = config.DbInfo.Schema + "." + tbl
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

// DeleteContents - cruD contents / basic-board
func DeleteContents(data interface{}, table string) (sql.Result, error) {
	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	whereEXP := goqu.Ex{
		"IDX": data.(models.ContentsBasicBoardSET).Idx,
	}
	if data.(models.ContentsBasicBoardSET).WriterName.Valid {
		whereEXP["WRITER_NAME"] = data.(models.ContentsBasicBoardSET).WriterName
		whereEXP["WRITER_PASSWORD"] = nil
	} else {
		whereEXP["WRITER_PASSWORD"] = data.(models.ContentsBasicBoardSET).WriterPassword
	}

	tbl := table
	if dbType != "sqlite3" {
		if config.DbInfo.Type == "postgres" {
			tbl = config.DbInfo.Schema + "." + tbl
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

// SelectContents - cRud contents / basic-board
func SelectContents(search interface{}) (interface{}, error) {
	var result interface{}
	var err error

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	var searchBytes models.ContentSearch
	err = json.Unmarshal(search.([]byte), &searchBytes)
	if err != nil {
		return nil, err
	}

	dbms := goqu.New(dbType, Dbo)
	var ds *goqu.SelectDataset
	contentResult := []models.ContentsBasicBoardGET{}
	exps := []goqu.Expression{}
	expsAND := []goqu.Expression{}

	keywords := searchBytes.Keywords

	table := searchBytes.Table.String
	if dbType != "sqlite3" {
		if config.DbInfo.Type == "postgres" {
			table = config.DbInfo.Schema + "." + table
		} else {
			table = DatabaseName + "." + table
		}
	}

	if searchBytes.Options.Count.Int64 > 1 {
		// Content list
		ds = dbms.From(table).Select(models.ContentsBasicBoardList{})
	} else {
		// Not allow getting all of list
		if !searchBytes.Options.Count.Valid || searchBytes.Options.Count.Int64 < 1 {
			searchBytes.Options.Count.SetValid(1)
		}

		// Contents
		ds = dbms.From(table).Select(models.ContentsBasicBoardGET{})
	}

	for _, k := range keywords {
		ex := PrepareWhere(k)
		exAND := exp.Ex{}
		if !ex.IsEmpty() {
			for c, v := range ex {
				if c == "IDX" || c == "BOARD_IDX" {
					val := fmt.Sprintf("%s", v)
					exAND[c] = goqu.Op{"eq": val}
					delete(ex, c)
				} else if c == "WRITER_PASSWORD" {
					val := fmt.Sprintf("%s", v)
					exAND[c] = goqu.Op{"eq": val}
					delete(ex, c)
				} else if c == "WRITER_NAME" {
					val := fmt.Sprintf("%s", v)
					exAND[c] = goqu.Op{"eq": val}
					delete(ex, c)
				} else {
					val := fmt.Sprintf("%s%s%s", "%", v, "%")
					ex[c] = goqu.Op{"like": val}
				}
			}
			exps = append(exps, ex.Expression())
			expsAND = append(expsAND, exAND.Expression())
		}
	}
	ds = ds.Where(goqu.Or(exps...))
	ds = ds.Where(goqu.And(expsAND...))

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

	err = ds.ScanStructs(&contentResult)
	if err != nil {
		log.Println("ds: ", err.Error())
		return nil, err
	}
	if contentResult != nil {
		result = contentResult
	}

	return result, nil
}

// SelectContentsCount - data count -> pages = (data count) / (count per page)
func SelectContentsCount(search interface{}) (uint, error) {
	var cnt uint

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	dbms := goqu.New(dbType, Dbo)
	exps := []goqu.Expression{}
	table := ""

	switch search := search.(type) {
	case models.ContentSearch:
		table = search.Table.String
		if dbType != "sqlite3" {
			if config.DbInfo.Type == "postgres" {
				table = config.DbInfo.Schema + "." + table
			} else {
				table = DatabaseName + "." + table
			}
		}

		keywords := search.Keywords
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
	}
	// log.Println("Select count search type: ", reflect.TypeOf(search))
	ds := dbms.From(table).Select(goqu.COUNT("*").As("PAGE_COUNT"))
	ds = ds.Where(goqu.Or(exps...))

	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	ds.ScanVal(&cnt)

	return cnt, nil
}
