package handler

import (
	"encoding/json"
	"net/http"

	"9minutes/model"

	"github.com/gofiber/fiber/v2"
)

func LoginAPI(c *fiber.Ctx) error {
	var err error

	signin := model.SignIn{}
	err = json.Unmarshal(c.Request().Body(), &signin)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	sess.Set("name", signin.Name.String)
	sess.Set("ip", c.IP())
	sess.Set("user-agent", c.Get("User-Agent"))
	sess.Set("duration", 60*60*24*7)

	err = sess.Save()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	result := map[string]string{
		"msg": "Signin success",
	}

	return c.Status(http.StatusOK).JSON(result)
}
