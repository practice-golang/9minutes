package handler

import (
	"9m/consts"
	"9m/crud"
	"9m/model"
	"9m/router"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func GetMyInfo(c *router.Context) {
	if c.AuthInfo == nil {
		c.Text(http.StatusForbidden, "Unauthorized")
		return
	}

	columnsCount, _ := crud.GetUserColumnsCount()

	switch columnsCount {
	case model.UserDataFieldCount:
		user, err := crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
			return
		}

		c.Json(http.StatusOK, user)
	default:
		user, err := crud.GetUserByNameMap(c.AuthInfo.(model.AuthInfo).Name.String)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
			return
		}

		c.Json(http.StatusOK, user)
	}
}

func UpdateMyInfo(c *router.Context) {
	var err error

	if c.AuthInfo == nil {
		c.Text(http.StatusForbidden, "Unauthorized")
		return
	}
	userDataOLD, err := crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	userDataNEW := make(map[string]interface{})
	err = json.NewDecoder(c.Body).Decode(&userDataNEW)
	if err != nil {
		c.Text(http.StatusBadRequest, err.Error())
		return
	}

	userDataNEW["idx"] = fmt.Sprint(userDataOLD.Idx.Int64)

	if _, ok := userDataNEW["password"]; ok {
		err = bcrypt.CompareHashAndPassword([]byte(userDataOLD.Password.String), []byte(userDataNEW["old-password"].(string)))
		if err != nil {
			c.Text(http.StatusBadRequest, "wrong password")
			return
		}

		password, err := bcrypt.GenerateFromPassword([]byte(userDataNEW["password"].(string)), consts.BcryptCost)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
			return
		}

		userDataNEW["password"] = string(password)
		delete(userDataNEW, "old-password")
	}

	err = crud.UpdateUserMap(userDataNEW)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	result := map[string]string{
		"result": "ok",
	}

	c.Json(http.StatusOK, result)
}
