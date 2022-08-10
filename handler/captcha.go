package handler

import (
	"9minutes/router"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dchest/captcha"
)

func GetCaptchaImage(c *router.Context) {
	captchaID := strings.TrimSuffix(filepath.Base(c.URL.Path), filepath.Ext(c.URL.Path))

	if c.URL.Query().Has("reload") {
		captcha.Reload(captchaID)
	}

	c.ResponseWriter.Header().Set("Content-Type", "image/png")
	err := captcha.WriteImage(c.ResponseWriter, captchaID, captcha.StdWidth, captcha.StdHeight)
	if err != nil {
		log.Println(err)
	}
}

func RenewCaptcha(c *router.Context) {
	captchaID := captcha.New()

	c.Json(http.StatusOK, map[string]interface{}{"captcha-id": captchaID})
}

func VerifyCaptcha(c *router.Context) {
	result := captcha.VerifyString(c.FormValue("captcha-id"), c.FormValue("captcha-answer"))

	if result {
		c.Json(http.StatusOK, map[string]interface{}{"result": "success"})
	} else {
		c.Json(http.StatusOK, map[string]interface{}{"result": "fail"})
	}
}
