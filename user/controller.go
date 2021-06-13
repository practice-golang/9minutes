package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/practice-golang/9minutes/db"
	"github.com/practice-golang/9minutes/models"
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

	sqlResult, err := db.InsertUserField(data)
	if err != nil {
		log.Println("AddUserFields: ", err)
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
	var data models.UserColumn

	dataJSON, _ := ioutil.ReadAll(c.Request().Body)

	err = json.Unmarshal(dataJSON, &data)
	if err != nil {
		log.Println("AddUserFields: ", err)
	}

	previousFieldsINTF, _ := db.SelectUserFields(models.UserColumn{})
	previousFields := previousFieldsINTF.([]models.UserColumn)

	_ = db.Dbi.EditUserTableFields(previousFields, []models.UserColumn{data})

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
