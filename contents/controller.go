package contents

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/practice-golang/9minutes/auth"
	"github.com/practice-golang/9minutes/board"
	"github.com/practice-golang/9minutes/db"
	"github.com/practice-golang/9minutes/models"
	"github.com/practice-golang/9minutes/user"
	"gopkg.in/guregu/null.v4"
)

// GetContentsListBasicBoard - Get contents
func GetContentsListBasicBoard(c echo.Context) error {
	var data interface{}
	var err error

	search, _ := ioutil.ReadAll(c.Request().Body)

	if !user.CheckPermission(c) {
		return c.JSON(http.StatusForbidden, "")
	}

	user := c.Get("user")
	isAdmin := "N"
	username := ""
	if user != nil {
		claims := user.(*jwt.Token).Claims.(*auth.CustomClaims)
		isAdmin = claims.Admin
		username = claims.UserName
	}

	mode := c.QueryParam("mode")

	if mode == "edit" && isAdmin != "Y" {
		var searchModel models.ContentSearch
		err = json.Unmarshal(search, &searchModel)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		// password for anonymous post
		password := false
		for _, s := range searchModel.Keywords {
			if s.WriterPassword.Valid {
				password = true
				break
			}
		}

		if !password {
			searchModel.Keywords = append(searchModel.Keywords, models.ContentsBasicBoardSET{
				WriterName: null.NewString(username, true),
			})
			search, err = json.Marshal(searchModel)
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
		}
	}

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
	var data models.ContentsBasicBoardSET

	err := json.Unmarshal(dataBytes, &dataMap)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	dataJSON, _ := json.Marshal(dataMap["data"])
	err = json.Unmarshal(dataJSON, &data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	isValid := user.CheckPermission(c)

	if !isValid {
		return c.JSON(http.StatusForbidden, map[string]bool{"permission": false})
	}

	isFileUpload := board.CheckUpload(c)
	if isFileUpload {
		filesJSON, _ := json.Marshal(dataMap["files"])
		data.Files = null.NewString(string(filesJSON), true)
	} else {
		data.Files = null.NewString("", false)
	}

	user := c.Get("user")
	if user != nil {
		claims := user.(*jwt.Token).Claims.(*auth.CustomClaims)
		data.WriterIdx = null.NewString(claims.Idx, true)
		data.WriterName = null.NewString(claims.UserName, true)
		data.WriterPassword = null.NewString("", false)
		data.IsMember = null.NewString("Y", true)
	} else {
		data.WriterIdx = null.NewString("-1", true)
		data.IsMember = null.NewString("N", true)
	}

	sqlResult, err := db.InsertContents(data, dataMap["table"].(string))
	if err != nil {
		log.Println("InsertContents: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	log.Println("AddContentsBasicBoard / last-id: ", lastID)

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
	var data models.ContentsBasicBoardSET

	err := json.Unmarshal(dataBytes, &dataMap)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	dataJSON, _ := json.Marshal(dataMap["data"])
	err = json.Unmarshal(dataJSON, &data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	isValid := user.CheckPermission(c)

	if !isValid {
		return c.JSON(http.StatusForbidden, map[string]bool{"permission": false})
	}

	isFileUpload := board.CheckUpload(c)
	if isFileUpload {
		filesJSON, _ := json.Marshal(dataMap["files"])
		data.Files = null.NewString(string(filesJSON), true)
	} else {
		data.Files = null.NewString("", false)
	}

	user := c.Get("user")
	if user != nil {
		claims := user.(*jwt.Token).Claims.(*auth.CustomClaims)
		if data.IsMember.String == "Y" {
			if data.WriterName.String == claims.UserName {
				data.WriterIdx = null.NewString(claims.Idx, true)
				data.WriterName = null.NewString(claims.UserName, true)
			}
		} else if !data.WriterPassword.Valid {
			return c.JSON(http.StatusBadRequest, map[string]string{"msg": "you can not edit: need password"})
		}
	} else if !data.WriterPassword.Valid {
		log.Println("WTF?? ", data.WriterPassword)
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": "you can not edit: need password"})
	}

	sqlResult, err := db.UpdateContents(data, dataMap["table"].(string))
	if err != nil {
		log.Println("UpdateContents: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	// if db.DBType == db.SQLITE && lastID == 0 {
	if lastID == 0 {
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
	havePermission := user.CheckPermission(c)

	if !havePermission {
		return c.JSON(http.StatusForbidden, map[string]bool{"permission": false})
	}

	dataBytes, _ := ioutil.ReadAll(c.Request().Body)

	var dataMap map[string]interface{}
	var data models.ContentsBasicBoardSET

	err := json.Unmarshal(dataBytes, &dataMap)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	dataJSON, _ := json.Marshal(dataMap["data"])
	err = json.Unmarshal(dataJSON, &data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	data.WriterName = null.NewString("", false)
	user := c.Get("user")
	if user != nil && data.IsMember.String == "Y" {
		claims := user.(*jwt.Token).Claims.(*auth.CustomClaims)
		data.WriterName = null.NewString(claims.UserName, true)
		data.WriterPassword = null.NewString("", false)
	}

	sqlResult, err := db.DeleteContents(data, dataMap["table"].(string))
	if err != nil {
		log.Println("DeleteContents: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	// if db.DBType == db.SQLITE && lastID == 0 {
	if lastID == 0 {
		lastID, _ = strconv.ParseInt(data.Idx.String, 10, 64)
	}

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
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
	if !auth.CheckAuth(c) {
		return c.JSON(http.StatusForbidden, "")
	}

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

	isValid := user.CheckPermission(c)

	if !isValid {
		return c.JSON(http.StatusForbidden, map[string]bool{"permission": false})
	}

	userName := ""
	user := c.Get("user")
	if user != nil {
		claims := user.(*jwt.Token).Claims.(*auth.CustomClaims)
		userName = claims.UserName
	}

	isFileUpload := board.CheckUpload(c)

	sqlResult, err := db.InsertContentsMAP(dataBytes, userName, isFileUpload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	result := map[string]string{
		"affected-rows": fmt.Sprint(affRows),
		"last-id":       fmt.Sprint(lastID),
	}

	return c.JSON(http.StatusOK, result)
}

// UpdateContentsListCustomBoard - Update contents
func UpdateContentsListCustomBoard(c echo.Context) error {
	isValid := user.CheckPermission(c)

	if !isValid {
		return c.JSON(http.StatusForbidden, map[string]bool{"permission": false})
	}

	dataBytes, _ := ioutil.ReadAll(c.Request().Body)

	userName := ""
	user := c.Get("user")
	if user != nil {
		claims := user.(*jwt.Token).Claims.(*auth.CustomClaims)
		userName = claims.UserName
	}

	isFileUpload := board.CheckUpload(c)

	sqlResult, err := db.UpdateContentsMAP(dataBytes, userName, isFileUpload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	// if db.DBType == db.SQLITE && lastID == 0 {
	if lastID == 0 {
		var allData map[string]interface{}
		_ = json.Unmarshal(dataBytes, &allData)
		lastID, _ = strconv.ParseInt(fmt.Sprint(allData["data"].(map[string]interface{})["IDX"]), 10, 64)
	}

	result := map[string]string{
		"affected-rows": fmt.Sprint(affRows),
		"last-id":       fmt.Sprint(lastID),
	}

	return c.JSON(http.StatusOK, result)
}

// DeleteContentsListCustomBoard - Delete contents
func DeleteContentsListCustomBoard(c echo.Context) error {
	isValid := user.CheckPermission(c)

	if !isValid {
		return c.JSON(http.StatusForbidden, map[string]bool{"permission": false})
	}

	dataBytes, _ := ioutil.ReadAll(c.Request().Body)

	userName := ""
	user := c.Get("user")
	if user != nil {
		claims := user.(*jwt.Token).Claims.(*auth.CustomClaims)
		userName = claims.UserName
	}
	sqlResult, err := db.DeleteContentsMAP(dataBytes, userName)
	if err != nil {
		log.Println("DeleteContents custom board: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	lastID, _ := sqlResult.LastInsertId()
	affected, _ := sqlResult.RowsAffected()

	result := map[string]string{
		"affected-rows": fmt.Sprint(affected),
		"last-id":       fmt.Sprint(lastID),
	}

	return c.JSON(http.StatusOK, result)
}

// GetContentsTotalPageMAP - Get total page of custom board
func GetContentsTotalPageMAP(c echo.Context) error {
	if !auth.CheckAuth(c) {
		return c.JSON(http.StatusForbidden, "")
	}

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
