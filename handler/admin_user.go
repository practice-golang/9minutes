package handler

import (
	"9minutes/consts"
	"9minutes/crud"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func GetUsersList(c *fiber.Ctx) error {
	queries := c.Queries()
	search := queries["search"]

	result, err := crud.GetUsersListMap(search)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(result)
}

func AddUser(c *fiber.Ctx) error {
	var err error

	now := time.Now().Format("20060102150405")

	userData := make(map[string]interface{})

	err = c.BodyParser(&userData)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	password, err := bcrypt.GenerateFromPassword([]byte(userData["password"].(string)), consts.BcryptCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	userData["password"] = string(password)
	userData["reg-dttm"] = now

	err = crud.AddUserMap(userData)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]string{
		"result": "ok",
	}

	return c.Status(http.StatusOK).JSON(result)
}

func UpdateUser(c *fiber.Ctx) error {
	var err error

	userDatas := []map[string]interface{}{}

	err = c.BodyParser(&userDatas)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	for _, userData := range userDatas {
		if _, ok := userData["password"]; ok {
			password, err := bcrypt.GenerateFromPassword([]byte(userData["password"].(string)), consts.BcryptCost)
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
			userData["password"] = string(password)
		}

		err = crud.UpdateUserMap(userData)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}

	result := map[string]string{
		"result": "ok",
	}

	return c.Status(http.StatusOK).JSON(result)
}

func DeleteUser(c *fiber.Ctx) error {
	userDatas := []map[string]interface{}{}

	err := c.BodyParser(&userDatas)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	for _, userData := range userDatas {
		err = crud.ResignUser(userData["id"].(int64))
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}

	result := map[string]string{
		"result": "ok",
	}

	return c.Status(http.StatusOK).JSON(result)
}
