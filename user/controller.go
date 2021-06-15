package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/practice-golang/9minutes/auth"
	"github.com/practice-golang/9minutes/db"
	"github.com/practice-golang/9minutes/models"
	"gopkg.in/guregu/null.v4"

	"github.com/golang-jwt/jwt"
	_ "github.com/golang-jwt/jwt"
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

// DeleteUserFields - Delete user custom field
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

// GetUserColumns - Get column names of USER table
func GetUserColumns(c echo.Context) error {
	sqlResult, err := db.SelectUserColumnNames()
	if err != nil {
		log.Println("GetUserColumns: ", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, sqlResult)
}

// AddUser - Add user
func AddUser(c echo.Context) error {

	return c.JSON(http.StatusOK, "test")
}

// GetUsers - Get a user fields
func GetUsers(c echo.Context) error {
	var err error
	var data interface{}

	search, _ := ioutil.ReadAll(c.Request().Body)

	// data, err = db.SelectUsers(search)
	data, err = db.SelectContentsMAP(search)
	if err != nil {
		log.Println("GetUsers: ", err)
	}

	return c.JSON(http.StatusOK, data)
}

// Login - login
func Login(c echo.Context) error {
	var err error
	var data interface{}

	search, _ := ioutil.ReadAll(c.Request().Body)

	data, err = db.SelectContentsMAP(search)
	if err != nil {
		log.Println("GetUsers: ", err)
	}

	result := map[string]string{}

	if len(data.([]map[string]interface{})) != 1 {
		return c.JSON(http.StatusNoContent, "")
	}

	result["token"], err = auth.PrepareToken(data.([]map[string]interface{})[0])
	if err != nil {
		log.Println("GetUsers token: ", err)
	}

	return c.JSON(http.StatusOK, result)
}

// VerifyToken - Validate token
func VerifyToken(c echo.Context) error {
	result := map[string]string{"msg": "OK"}
	return c.JSON(http.StatusOK, result)
}

func ReissueToken(c echo.Context) error {
	var err error
	var info *jwt.Token
	result := map[string]string{}

	headers := c.Request().Header

	if tokenString, ok := headers["Authorization"]; ok {
		tokens := strings.Split(tokenString[0], " ")

		info, err = jwt.Parse(tokens[1], func(token *jwt.Token) (interface{}, error) {
			return auth.JwtKey, nil
		})
		if err != nil {
			log.Println("jwt parse error: ", err)
		}
		claimMAP := info.Claims.(jwt.MapClaims)

		var claimsJSON []byte
		claimsJSON, err = json.Marshal(claimMAP)
		if err != nil {
			log.Println("json marshal error: ", err)
		}
		var claims auth.CustomClaims
		err = json.Unmarshal(claimsJSON, &claims)
		if err != nil {
			log.Println("json unmarshal error: ", err)
		}

		now := time.Now().Unix()
		delta := int64(time.Unix(now, 0).Sub(time.Unix(claims.RefreshUntil, 0)).Hours() / 24)

		if delta < 0 {
			// access_token
			// claims.ExpiresAt = time.Now().Add(time.Second * 1 * 60).Unix()
			claims.ExpiresAt = time.Now().Add(time.Hour * 1 * 24).Unix()
			// refresh
			if delta > -7 {
				// claims.RefreshUntil = time.Now().Add(time.Second * 2 * 60).Unix()
				claims.RefreshUntil = time.Now().Add(time.Hour * 30 * 24).Unix()
			}
			refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			result["token"], err = refreshToken.SignedString(auth.JwtKey)
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
		} else {
			return c.JSON(http.StatusUnauthorized, errors.New("issue not allowed"))
		}
	}

	return c.JSON(http.StatusOK, result)
}
