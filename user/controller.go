package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/practice-golang/9minutes/db"
	"github.com/practice-golang/9minutes/models"
	"gopkg.in/guregu/null.v4"
)

// AddFields - Insert user(s) optional table fields
func AddUserFields(c echo.Context) error {
	var err error
	var data []models.UserColumn

	dataJSON, _ := ioutil.ReadAll(c.Request().Body)

	err = json.Unmarshal(dataJSON, &data)
	if err != nil {
		log.Println("AddUserFields: ", err)
	}

	_ = db.Dbi.AddUserTableFields(data)

	sqlResult, err := db.InsertUserField(data)
	if err != nil {
		// log.Println("AddUserFields: ", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
}

// GetUserFields - Get a user fields
func GetUserFields(c echo.Context) error {
	var err error

	dataINTF, err := db.SelectUserFields(models.UserColumn{})
	if err != nil {
		log.Println("SelectUserFields: ", err)
	}

	return c.JSON(http.StatusOK, dataINTF)
}

// EditUserFields - Modify user fields
func EditUserFields(c echo.Context) error {
	var err error
	var data []models.UserColumn

	dataJSON, _ := ioutil.ReadAll(c.Request().Body)

	err = json.Unmarshal(dataJSON, &data)
	if err != nil {
		log.Println("AddUserFields: ", err)
	}

	previousFieldsINTF, _ := db.SelectUserFields(models.UserColumn{})
	previousFields := previousFieldsINTF.([]models.UserColumn)

	err = db.Dbi.EditUserTableFields(previousFields, data, []string{"remove"})
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{"msg": string(err.Error())})
	}

	sqlResult, err := db.UpdateUserFields(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"msg": string(err.Error())})
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
}

// DeleteUserFields - Delete a board
func DeleteUserFields(c echo.Context) error {
	idx := c.Param("idx")
	idxInt, _ := strconv.Atoi(idx)
	idxNullInt := null.NewInt(int64(idxInt), true)

	// Drop table
	data, err := db.SelectUserFields(models.UserColumn{Idx: idxNullInt})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if len(data.([]models.UserColumn)) > 0 {
		db.Dbi.DeleteUserTableFields(data.([]models.UserColumn))
	}

	// Remove field info
	sqlResult, err := db.DeleteUserFieldRow("IDX", fmt.Sprint(idx))
	if err != nil {
		log.Println("DeleteUserFieldRow: ", err)
		result := map[string]string{
			"msg":  err.Error(),
			"desc": "If exist rest of fields rows, please remove manually",
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
}
