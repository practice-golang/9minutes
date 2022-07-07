package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"9minutes/auth"
	"9minutes/model"
	"9minutes/router"

	"gopkg.in/guregu/null.v4"
)

func Signin(c *router.Context) {
	var err error

	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Json(http.StatusBadRequest, err.Error())
	}

	// log.Println(c.Request.RemoteAddr)
	// log.Println(c.Request.UserAgent())

	signin := model.SignIn{}
	err = json.Unmarshal(b, &signin)
	if err != nil {
		c.Json(http.StatusInternalServerError, err.Error())
	}

	authinfo := model.AuthInfo{
		Name:     null.NewString(signin.Name.String, true),
		IpAddr:   null.NewString(c.RemoteAddr, true),
		Platform: null.NewString("", true),
		Duration: null.NewInt(60*60*24*7, true),
		// Duration: null.NewInt(10, true), // 10 seconds test
	}

	// auth.SetupCookieToken(c.ResponseWriter, authinfo)
	auth.SetCookieSession(c, authinfo)

	c.Json(http.StatusOK, "Signin success")
}

func SigninAPI(c *router.Context) {
	var err error

	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Json(http.StatusBadRequest, err.Error())
	}

	// log.Println(c.Request.RemoteAddr)
	// log.Println(c.Request.UserAgent())

	signin := model.SignIn{}
	err = json.Unmarshal(b, &signin)
	if err != nil {
		c.Json(http.StatusInternalServerError, err.Error())
	}

	authinfo := model.AuthInfo{
		Name:     null.NewString(signin.Name.String, true),
		IpAddr:   null.NewString(c.RemoteAddr, true),
		Platform: null.NewString("", true),
		Duration: null.NewInt(60*60*24*7, true),
		// Duration: null.NewInt(10, true), // 10 seconds test
	}

	token, err := auth.GenerateToken(authinfo)
	if err != nil {
		c.Json(http.StatusInternalServerError, err.Error())
	}

	result := map[string]string{
		"token": token,
		"msg":   "Signin success",
	}

	c.Json(http.StatusOK, result)
}
