package contents

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/practice-golang/9minutes/db"
	"github.com/practice-golang/9minutes/models"
)

// GetContentsListBasicBoard - Get contents
func GetContentsListBasicBoard(c echo.Context) error {
	var data interface{}
	var err error

	search, _ := ioutil.ReadAll(c.Request().Body)

	data, err = db.SelectContents(search)
	if err != nil {
		log.Println("Get contents: ", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, data)
}

// AddContentsBasicBoard - Add contents
func AddContentsBasicBoard(c echo.Context) error {
	dataBytes, _ := ioutil.ReadAll(c.Request().Body)

	var dataMap map[string]interface{}
	var data models.ContentsBasicBoard

	err := json.Unmarshal(dataBytes, &dataMap)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	dataJSON, _ := json.Marshal(dataMap["data"])
	err = json.Unmarshal(dataJSON, &data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	sqlResult, err := db.InsertContents(data, dataMap["table"].(string))
	if err != nil {
		log.Println("InsertContents: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
}

// UpdateContentsBasicBoard - Update contents
func UpdateContentsBasicBoard(c echo.Context) error {
	dataBytes, _ := ioutil.ReadAll(c.Request().Body)

	var dataMap map[string]interface{}
	var data models.ContentsBasicBoard

	err := json.Unmarshal(dataBytes, &dataMap)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	dataJSON, _ := json.Marshal(dataMap["data"])
	err = json.Unmarshal(dataJSON, &data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	sqlResult, err := db.UpdateContents(data, dataMap["table"].(string))
	if err != nil {
		log.Println("InsertContents: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	if db.DBType == db.SQLITE && lastID == 0 {
		lastID, _ = strconv.ParseInt(data.Idx.String, 10, 64)
	}

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
}

// DeleteContentsBasicBoard - Delete contents
func DeleteContentsBasicBoard(c echo.Context) error {
	dataBytes, _ := ioutil.ReadAll(c.Request().Body)

	var dataMap map[string]interface{}
	var data models.ContentsBasicBoard

	err := json.Unmarshal(dataBytes, &dataMap)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	dataJSON, _ := json.Marshal(dataMap["data"])
	err = json.Unmarshal(dataJSON, &data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	sqlResult, err := db.DeleteContents(data, dataMap["table"].(string))
	if err != nil {
		log.Println("InsertContents: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	return c.JSON(http.StatusOK, sqlResult)
}

// GetContentsTotalPage - Get total page of basic board
func GetContentsTotalPage(c echo.Context) error {
	var search models.ContentSearch

	if err := c.Bind(&search); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	data, err := db.SelectContentsCount(search)
	if err != nil {
		log.Fatal("SelectCount: ", err)
	}

	countPerPage := uint(1)
	if search.Options.Count.Valid {
		countPerPage = uint(search.Options.Count.Int64)
	}

	pages := uint(math.Ceil(float64(data) / float64(countPerPage)))
	// log.Println("search: ", search.Keywords, data, pages)

	result := map[string]uint{"total-page": pages}

	return c.JSON(http.StatusOK, result)
}

// GetContentsListCustomBoard - Get contents
func GetContentsListCustomBoard(c echo.Context) error {
	var data interface{}
	var err error

	search, _ := ioutil.ReadAll(c.Request().Body)

	data, err = db.SelectContentsMAP(search)
	if err != nil {
		log.Println("Get contents: ", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, data)
}

// AddContentsListCustomBoard - Add contents
func AddContentsListCustomBoard(c echo.Context) error {
	dataBytes, _ := ioutil.ReadAll(c.Request().Body)

	sqlResult, err := db.InsertContentsMAP(dataBytes)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
}

// UpdateContentsListCustomBoard - Update contents
func UpdateContentsListCustomBoard(c echo.Context) error {
	dataBytes, _ := ioutil.ReadAll(c.Request().Body)

	sqlResult, err := db.UpdateContentsMAP(dataBytes)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	if db.DBType == db.SQLITE && lastID == 0 {
		var allData map[string]interface{}
		_ = json.Unmarshal(dataBytes, &allData)
		lastID, _ = strconv.ParseInt(fmt.Sprint(allData["data"].(map[string]interface{})["IDX"]), 10, 64)
	}

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
}

// DeleteContentsListCustomBoard - Delete contents
func DeleteContentsListCustomBoard(c echo.Context) error {
	dataBytes, _ := ioutil.ReadAll(c.Request().Body)

	sqlResult, err := db.DeleteContentsMAP(dataBytes)
	if err != nil {
		log.Println("DeleteContents custom board: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	return c.JSON(http.StatusOK, sqlResult)
}

// GetContentsTotalPageMAP - Get total page of custom board
func GetContentsTotalPageMAP(c echo.Context) error {
	search, _ := ioutil.ReadAll(c.Request().Body)

	data, count, err := db.SelectContentsCountMAP(search)
	if err != nil {
		log.Fatal("SelectCount: ", err)
	}

	// countPerPage := uint(1)
	countPerPage := count

	pages := uint(math.Ceil(float64(data) / float64(countPerPage)))
	// log.Println("search: ", string(search), data, pages)

	result := map[string]uint{"total-page": pages}

	return c.JSON(http.StatusOK, result)
}
