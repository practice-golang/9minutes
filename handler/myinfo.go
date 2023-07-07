package handler

import (
	"9minutes/consts"
	"9minutes/crud"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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

	// columnsCount, _ := crud.GetUserColumnsCount()
	// switch columnsCount {
	// case model.UserDataFieldCount:
	// 	user, err := crud.GetUserByName(name.(string))
	// 	if err != nil {
	// 		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	// 	}

	// 	user.Password = null.NewString("", false)

	// 	return c.Status(http.StatusOK).JSON(user)
	// default:
	// 	user, err := crud.GetUserByNameMap(name.(string))
	// 	if err != nil {
	// 		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	// 	}

	// 	delete(user.(map[string]interface{}), "password")

	// 	return c.Status(http.StatusOK).JSON(user)
	// }

	user, err := crud.GetUserByNameMap(name.(string))
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	delete(user.(map[string]interface{}), "password")

	return c.Status(http.StatusOK).JSON(user)
}

func UpdateMyInfo(c *fiber.Ctx) error {
	var err error

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	name := sess.Get("name")
	if name == nil {
		return c.Status(http.StatusForbidden).Send([]byte("Unauthorized"))
	}

	userDataOldRaw, err := crud.GetUserByNameMap(name.(string))
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}
	userDataOld := userDataOldRaw.(map[string]interface{})

	userDataNEW := make(map[string]interface{})
	err = json.Unmarshal(c.Body(), &userDataNEW)
	if err != nil {
		return c.Status(http.StatusBadRequest).Send([]byte(err.Error()))
	}

	userDataNEW["idx"] = fmt.Sprint(userDataOld["idx"].(int64))

	if _, ok := userDataNEW["password"]; ok {
		err = bcrypt.CompareHashAndPassword([]byte(userDataOld["password"].(string)), []byte(userDataNEW["old-password"].(string)))
		if err != nil {
			return c.Status(http.StatusBadRequest).Send([]byte("wrong password"))
		}

		password, err := bcrypt.GenerateFromPassword([]byte(userDataNEW["password"].(string)), consts.BcryptCost)
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
		}

		userDataNEW["password"] = string(password)
		delete(userDataNEW, "old-password")
	}

	err = crud.UpdateUserMap(userDataNEW)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	result := map[string]string{
		"result": "ok",
	}

	return c.Status(http.StatusOK).JSON(result)
}

func ResignUser(c *fiber.Ctx) error {
	var err error

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	name := sess.Get("name")
	if name == nil {
		return c.Status(http.StatusForbidden).Send([]byte("Unauthorized"))
	}

	userDataRaw, err := crud.GetUserByNameMap(name.(string))
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}
	userData := userDataRaw.(map[string]interface{})

	err = crud.ResignUser(userData["idx"].(int64))
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	result := map[string]string{
		"result": "ok",
	}

	return c.Status(http.StatusOK).JSON(result)
}
