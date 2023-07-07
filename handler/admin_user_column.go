package handler

import (
	"9minutes/crud"
	"9minutes/model"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func RestrictedHello(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	name := sess.Get("name")
	if name == nil {
		return c.Status(http.StatusForbidden).Send([]byte("Unauthorized"))
	}

	return c.Status(http.StatusOK).Send([]byte("Hello " + name.(string)))
}

func GetUserColumns(c *fiber.Ctx) error {
	result, err := crud.GetUserColumnsList()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(result)
}

func AddUserColumn(c *fiber.Ctx) error {
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

	return c.Status(http.StatusOK).JSON(result)
}

func UpdateUserColumns(c *fiber.Ctx) error {
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

	return c.Status(http.StatusOK).JSON(result)
}

func DeleteUserColumns(c *fiber.Ctx) error {
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

	return c.Status(http.StatusOK).JSON(result)
}