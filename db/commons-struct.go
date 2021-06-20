package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/practice-golang/9minutes/models"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver"
)

// InsertContents - Crud contents / basic-board
func InsertContents(data interface{}, table string) (sql.Result, error) {
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

// UpdateContents - Crud contents / basic-board
func UpdateContents(data interface{}, table string) (sql.Result, error) {
	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	whereEXP := goqu.Ex{
		"IDX":         data.(models.ContentsBasicBoard).Idx,
		"WRITER_NAME": data.(models.ContentsBasicBoard).WriterName,
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Update(table).Set(data).Where(whereEXP)
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

	whereEXP := goqu.Ex{"idx": data.(models.ContentsBasicBoard).Idx}

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
	contentResult := []models.ContentsBasicBoard{}
	exps := []goqu.Expression{}

	keywords := searchBytes.Keywords
	if searchBytes.Options.Count.Int64 > 1 {
		// Content list
		ds = dbms.From(searchBytes.Table.String).Select(models.ContentsListBasicBoard{})
	} else {
		// Not allow getting all of list
		if !searchBytes.Options.Count.Valid || searchBytes.Options.Count.Int64 < 1 {
			searchBytes.Options.Count.SetValid(1)
		}

		// Contents
		ds = dbms.From(searchBytes.Table.String).Select(models.ContentsBasicBoard{})
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
		keywords := search.Keywords
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
	}
	// log.Println("Select count search type: ", reflect.TypeOf(search))
	ds := dbms.From(table).Select(goqu.COUNT("*").As("PAGE_COUNT"))
	ds = ds.Where(goqu.Or(exps...))

	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	ds.ScanVal(&cnt)

	return cnt, nil
}
