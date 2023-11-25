package handler

import (
	"9minutes/consts"
	"9minutes/internal/crud"
	"9minutes/model"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func RestrictedHello(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	userid := sess.Get("userid")
	if userid == nil {
		return c.Status(http.StatusForbidden).Send([]byte("Unauthorized"))
	}

	return c.Status(http.StatusOK).Send([]byte("Hello " + userid.(string)))
}

func GetUserColumnsAPI(c *fiber.Ctx) error {
	// result, err := crud.GetUserColumnsList()
	// if err != nil {
	// 	return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	// }
	result := UserColumnsALL

	return c.Status(http.StatusOK).JSON(result)
}

func AddUserColumnAPI(c *fiber.Ctx) error {
	var userColumn model.UserColumn

	err := c.BodyParser(&userColumn)
	if err != nil {
		return c.Status(http.StatusBadRequest).Send([]byte(err.Error()))
	}

	err = crud.AddUserColumn(userColumn)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	result := map[string]string{
		"result": "ok",
	}

	SetUserColumnsALL()

	return c.Status(http.StatusOK).JSON(result)
}

func UpdateUserColumnsAPI(c *fiber.Ctx) error {
	var userColumns []model.UserColumn

	err := c.BodyParser(&userColumns)
	if err != nil {
		return c.Status(http.StatusBadRequest).Send([]byte(err.Error()))
	}

	for _, userColumn := range userColumns {
		err = crud.UpdateUserColumn(userColumn)
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
		}
	}

	result := map[string]string{
		"result": "ok",
	}

	SetUserColumnsALL()

	return c.Status(http.StatusOK).JSON(result)
}

func DeleteUserColumnsAPI(c *fiber.Ctx) error {
	var userColumns []model.UserColumn

	err := c.BodyParser(&userColumns)
	if err != nil {
		return c.Status(http.StatusBadRequest).Send([]byte(err.Error()))
	}

	for _, userColumn := range userColumns {
		err := crud.DeleteUserColumn(userColumn)
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
		}
	}

	result := map[string]string{
		"result": "ok",
	}

	SetUserColumnsALL()

	return c.Status(http.StatusOK).JSON(result)
}

func GetUserGrades(c *fiber.Ctx) error {
	result := consts.UserGrades

	return c.Status(http.StatusOK).JSON(result)
}
