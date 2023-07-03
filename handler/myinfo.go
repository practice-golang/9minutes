package handler

import (
	"9minutes/consts"
	"9minutes/crud"
	"9minutes/model"
	"9minutes/router"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

func GetMyInfo(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	name := sess.Get("name")
	if name == nil {
		return c.Status(http.StatusForbidden).Send([]byte("Unauthorized"))
	}

	columnsCount, _ := crud.GetUserColumnsCount()

	switch columnsCount {
	case model.UserDataFieldCount:
		user, err := crud.GetUserByName(name.(string))
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
		}

		user.Password = null.NewString("", false)

		return c.Status(http.StatusOK).JSON(user)
	default:
		user, err := crud.GetUserByNameMap(name.(string))
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
		}

		delete(user.(map[string]interface{}), "password")

		return c.Status(http.StatusOK).JSON(user)
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

func ResignUser(c *router.Context) {
	var err error

	if c.AuthInfo == nil {
		c.Text(http.StatusForbidden, "Unauthorized")
		return
	}

	userData, err := crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	err = crud.ResignUser(userData.Idx.Int64)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	result := map[string]string{
		"result": "ok",
	}

	c.Json(http.StatusOK, result)
}
