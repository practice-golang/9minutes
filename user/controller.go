package user

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/practice-golang/9minutes/db"
	"github.com/practice-golang/9minutes/models"
)

// AddFields - Insert user(s) optional table fields
func AddUserFields(c echo.Context) error {
	// var err error

	data, _ := ioutil.ReadAll(c.Request().Body)
	result := map[string]string{
		"msg":  "test",
		"data": string(data),
	}

	return c.JSON(http.StatusOK, result)
}

// GetUserFields - Get a user fields
func GetUserFields(c echo.Context) error {
	var err error

	dataINTF, err := db.SelectUserFields(models.UserField{})
	if err != nil {
		log.Println("SelectUserFields: ", err)
	}
	// data := prepareSelectData(dataINTF)

	// return c.JSON(http.StatusOK, data)
	return c.JSON(http.StatusOK, dataINTF)
}
