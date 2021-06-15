package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/tidwall/gjson"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver"
)

// InsertContentsMAP - Crud / custom-board
func InsertContentsMAP(data interface{}) (sql.Result, error) {
	rcds := []goqu.Record{}
	var allData map[string]interface{}

	_ = json.Unmarshal(data.([]byte), &allData)

	jsonData, _ := json.Marshal(allData["data"].(map[string]interface{}))
	log.Println(string(jsonData))

	if data != nil {
		rcd := goqu.Record{}
		for k, d := range allData["data"].(map[string]interface{}) {
			log.Println(k, d)
			rcd[k] = d
		}
		rcds = append(rcds, rcd)
	}

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Insert(allData["table"]).Rows(rcds)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SelectContentsMAP - cRud contents MAP / custom-board, custom-tablelist
func SelectContentsMAP(search interface{}) (interface{}, error) {
	var result []map[string]interface{}
	var err error

	offset := uint(0)
	count := listCount

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	// log.Println(string(search.([]byte)))
	jsonBody, ok := gjson.Parse(string(search.([]byte))).Value().(map[string]interface{})
	if !ok {
		log.Println("Cannot parse jsonBody")
	}

	colNames := []interface{}{}
	log.Println("contentsMAP: ", jsonBody["columns"])
	for _, c := range jsonBody["columns"].([]interface{}) {
		colNames = append(colNames, c.(string))
	}

	dbms := goqu.New(dbType, Dbo)
	exps := []goqu.Expression{}
	ex := goqu.Ex{}

	switch jsonBody["keywords"].(type) {
	case ([]interface{}):
		keywords := jsonBody["keywords"].([]interface{})
		for _, keywordOBJ := range keywords {
			ex := goqu.Ex{}
			for k, d := range keywordOBJ.(map[string]interface{}) {
				val := ""
				if k == "IDX" {
					val = fmt.Sprintf("%s", d)
					ex[k] = goqu.Op{"eq": val}
				} else {
					val = fmt.Sprintf("%s%s%s", "%", d, "%")
					ex[k] = goqu.Op{"like": val}
				}
			}
			exps = append(exps, ex.Expression())
		}
	case (map[string]interface{}):
		keywords := jsonBody["keywords"].(map[string]interface{})
		for k, v := range keywords {
			val := fmt.Sprintf("%s", v)
			ex[k] = goqu.Op{"eq": val}
		}
	}

	exps = append(exps, ex.Expression())

	ds := dbms.From(jsonBody["table"].(string)).Select(colNames...)
	ds = ds.Where(goqu.Or(exps...))

	orderDirection := goqu.C(OrderScope).Asc()
	cnt := uint(count)

	if options, ok := jsonBody["options"].(map[string]interface{}); ok {
		// Order direction
		if opt, ok := options["order"]; ok && opt == "desc" {
			orderDirection = goqu.C(OrderScope).Desc()
		}

		// Rows count
		if optCount, ok := options["count"]; ok {
			var optcntINT int
			switch val := optCount.(type) {
			case string:
				optcntINT, _ = (strconv.Atoi(val))
			case float64:
				optcntINT = int(val)
			}

			cnt = uint(optcntINT)
		}

		// Paging
		if page, ok := options["page"]; ok {
			var pageINT int
			switch val := page.(type) {
			case string:
				pageINT, _ = (strconv.Atoi(val))
			case float64:
				pageINT = int(val)
			}

			offset = uint(pageINT)
		}
	}

	ds = ds.Order(orderDirection)
	ds = ds.Limit(cnt)
	ds = ds.Offset(offset * cnt)

	sql, args, _ := ds.ToSQL()
	log.Println("SelectContentsMAP: ", sql, args)

	rows, err := Dbo.Query(sql, args...)
	if err != nil {
		log.Println("rowsMAP: ", err.Error())
		return nil, err
	}
	// cols, _ := rows.Columns()
	cols := colNames

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		colPtrs := make([]interface{}, len(cols))
		for i := range columns {
			colPtrs[i] = &columns[i]
		}

		if err := rows.Scan(colPtrs...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := colPtrs[i].(*interface{})
			switch (*val).(type) {
			case string:
				m[colName.(string)] = (*val).(string)
			case uint8:
				m[colName.(string)] = (*val).(uint8)
			case []uint8:
				// string or double or integer or numeric
				m[colName.(string)] = string([]byte((*val).([]uint8)))
			case int64:
				m[colName.(string)] = (*val).(int64)
			case float64:
				m[colName.(string)] = (*val).(float64)
			}
		}

		result = append(result, m)
	}

	return result, err
}

// UpdateContentsMAP - crUd / custom-board
func UpdateContentsMAP(data interface{}) (sql.Result, error) {
	var err error
	var allData map[string]interface{}
	whereEXP := goqu.Ex{}

	_ = json.Unmarshal(data.([]byte), &allData)

	rcd := goqu.Record{}
	if allData["data"] != nil {
		for k, d := range allData["data"].(map[string]interface{}) {
			if k == "IDX" {
				whereEXP["IDX"] = d
			}
			rcd[k] = d
		}
	}

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Update(allData["table"]).Set(rcd).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteContentsMAP - cruD contents / custom-board
func DeleteContentsMAP(data interface{}) (sql.Result, error) {
	var err error
	var allData map[string]interface{}
	whereEXP := goqu.Ex{}

	_ = json.Unmarshal(data.([]byte), &allData)

	rcd := goqu.Record{}
	if allData["data"] != nil {
		for k, d := range allData["data"].(map[string]interface{}) {
			if k == "IDX" {
				whereEXP["IDX"] = d
			}
			rcd[k] = d
		}
	}

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	dbms := goqu.New(dbType, Dbo)
	ds := dbms.Delete(allData["table"]).Where(whereEXP)
	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	result, err := Dbo.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SelectContentsCountMAP - data count -> pages = (data count) / (count per page)
func SelectContentsCountMAP(search interface{}) (uint, uint, error) {
	var result uint

	offset := uint(0)
	count := listCount

	dbType, err := getDialect()
	if err != nil {
		log.Println("ERR Select DBType: ", err)
	}

	// log.Println(string(search.([]byte)))
	jsonBody, ok := gjson.Parse(string(search.([]byte))).Value().(map[string]interface{})
	if !ok {
		log.Println("Cannot parse jsonBody")
	}

	log.Println("contentsMAP search body: ", jsonBody)
	log.Println("contentsMAP: ", jsonBody["columns"])

	dbms := goqu.New(dbType, Dbo)
	exps := []goqu.Expression{}

	keywords := jsonBody["keywords"].([]interface{})
	for _, keywordOBJ := range keywords {
		ex := goqu.Ex{}
		for k, d := range keywordOBJ.(map[string]interface{}) {
			val := fmt.Sprintf("%s%s%s", "%", d, "%")
			ex[k] = goqu.Op{"like": val}
		}
		exps = append(exps, ex.Expression())
	}
	log.Println("kworkds", keywords)

	ds := dbms.From(jsonBody["table"].(string)).Select(goqu.COUNT("*").As("PAGE_COUNT"))
	ds = ds.Where(goqu.Or(exps...))

	orderDirection := goqu.C(OrderScope).Asc()
	options := jsonBody["options"].(map[string]interface{})
	if opt, ok := options["order"]; ok && opt == "desc" {
		orderDirection = goqu.C(OrderScope).Desc()
	}
	ds = ds.Order(orderDirection)

	cnt := uint(count)
	if optCount, ok := options["count"]; ok {
		var optcntINT int
		switch val := optCount.(type) {
		case string:
			optcntINT, _ = (strconv.Atoi(val))
		case float64:
			optcntINT = int(val)
		}

		cnt = uint(optcntINT)
	}
	ds = ds.Limit(cnt)

	// 페이징
	if page, ok := options["page"]; ok {
		var pageINT int
		switch val := page.(type) {
		case string:
			pageINT, _ = (strconv.Atoi(val))
		case float64:
			pageINT = int(val)
		}

		offset = uint(pageINT)
	}
	ds = ds.Offset(offset * cnt)

	sql, args, _ := ds.ToSQL()
	log.Println(sql, args)

	ds.ScanVal(&result)

	return result, cnt, nil
}
