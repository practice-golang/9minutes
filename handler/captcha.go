package handler

import (
	"9minutes/router"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dchest/captcha"
	"github.com/gofiber/fiber/v2"
)

func GetCaptchaImage(c *fiber.Ctx) error {
	var err error
	captchaID := strings.TrimSuffix(filepath.Base(c.Path()), filepath.Ext(c.Path()))

	if c.Query("reload") != "" {
		captcha.Reload(captchaID)
		return nil
	}

	c.Response().Header.Set("Content-Type", "image/png")
	err = captcha.WriteImage(c.Response().BodyWriter(), captchaID, captcha.StdWidth, captcha.StdHeight)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func RenewCaptcha(c *fiber.Ctx) error {
	captchaID := captcha.New()
	return c.Status(http.StatusOK).JSON(map[string]interface{}{"captcha-id": captchaID})
}

func VerifyCaptcha(c *router.Context) {
	result := captcha.VerifyString(c.FormValue("captcha-id"), c.FormValue("captcha-answer"))

	if result {
		c.Json(http.StatusOK, map[string]interface{}{"result": "success"})
	} else {
		c.Json(http.StatusOK, map[string]interface{}{"result": "fail"})
	}
}
