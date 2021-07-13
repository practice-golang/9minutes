package comments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/practice-golang/9minutes/auth"
	"github.com/practice-golang/9minutes/db"
	"github.com/practice-golang/9minutes/models"
	"github.com/practice-golang/9minutes/user"
	"gopkg.in/guregu/null.v4"
)

// GetComments - Get comments
func GetComments(c echo.Context) error {
	var data interface{}
	var err error

	search, _ := ioutil.ReadAll(c.Request().Body)

	data, err = db.SelectComments(search)
	if err != nil {
		log.Println("Get comments: ", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, data)
}

// AddComment - Add comment
func AddComment(c echo.Context) error {
	var err error
	var dataMap map[string]interface{}
	var data models.CommentSET

	dataBytes, _ := ioutil.ReadAll(c.Request().Body)

	err = json.Unmarshal(dataBytes, &dataMap)
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

	user := c.Get("user")
	if user != nil {
		data.IsMember = null.NewString("Y", true)
		claims := user.(*jwt.Token).Claims.(*auth.CustomClaims)
		data.WriterName = null.NewString(claims.UserName, true)
		data.WriterPassword = null.NewString("", false)
	} else if data.WriterPassword.Valid && data.WriterPassword.String == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": "Need password for guest"})
	}

	sqlResult, err := db.InsertComment(data, dataMap["table"].(string)+"_COMMENT")
	if err != nil {
		log.Println("InsertComments: ", err)
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

// EditComment - Edit comment
func EditComment(c echo.Context) error {
	var err error
	var dataMap map[string]interface{}
	var data models.CommentSET

	dataBytes, _ := ioutil.ReadAll(c.Request().Body)

	err = json.Unmarshal(dataBytes, &dataMap)
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

	user := c.Get("user")
	if user != nil {
		claims := user.(*jwt.Token).Claims.(*auth.CustomClaims)
		data.WriterName = null.NewString(claims.UserName, true)
	}
	sqlResult, err := db.UpdateComment(data, dataMap["table"].(string)+"_COMMENT")
	if err != nil {
		log.Println("UpdateComments: ", err)
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

// DeleteComment - Delete comment
func DeleteComment(c echo.Context) error {
	var err error
	var dataMap map[string]interface{}
	var data models.CommentSET

	dataBytes, _ := ioutil.ReadAll(c.Request().Body)

	err = json.Unmarshal(dataBytes, &dataMap)
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

	data.WriterName = null.NewString("", false)
	user := c.Get("user")
	if user != nil && data.IsMember.String == "Y" {
		claims := user.(*jwt.Token).Claims.(*auth.CustomClaims)
		data.WriterName = null.NewString(claims.UserName, true)
		data.WriterPassword = null.NewString("", false)
	}

	sqlResult, err := db.DeleteComment(data, dataMap["table"].(string)+"_COMMENT")
	if err != nil {
		log.Println("UpdateComments: ", err)
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
