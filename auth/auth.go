package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"9minutes/model"
	"9minutes/router"

	"gopkg.in/guregu/null.v4"
)

func SetupCookieToken(w http.ResponseWriter, authinfo model.AuthInfo) error {
	token, err := GenerateToken(authinfo)
	if err != nil {
		log.Println("SetupCookieToken:", err)
		return err
	}

	SetCookieHeader(w, token, authinfo.Duration.Int64)

	return nil
}

func SetCookieSession(c *router.Context, authinfo model.AuthInfo) error {
	SessionManager.Put(c.Context(), "userid", authinfo.Name.String)
	SessionManager.Put(c.Context(), "ip", authinfo.IpAddr.String)
	SessionManager.Put(c.Context(), "platform", authinfo.Platform.String)

	return nil
}

func GetCookieSession(c *router.Context) (model.AuthInfo, error) {
	var result model.AuthInfo

	if SessionManager == nil || !SessionManager.Exists(c.Context(), "userid") {
		return result, errors.New("userid is empty")
	}

	result = model.AuthInfo{
		Name:     null.NewString(SessionManager.GetString(c.Context(), "userid"), true),
		IpAddr:   null.NewString(SessionManager.GetString(c.Context(), "ip"), true),
		Platform: null.NewString(SessionManager.GetString(c.Context(), "platform"), true),
	}

	return result, nil
}

func DestroyCookieSession(c *router.Context) error {
	err := SessionManager.Destroy(c.Context())

	return err
}

func GetClaim(r http.Request, from string) (model.AuthInfo, error) {
	var result model.AuthInfo
	var dataCookie *http.Cookie
	var dataHeader string
	var token string
	var err error

	switch from {
	case "cookie":
		dataCookie, err = r.Cookie("token")
		if err != nil {
			// log.Println("GetCookie cookie:", err)
			return result, err
		}

		token = dataCookie.Value
	case "header":
		dataHeader = r.Header.Get("Authorization")
		dataHeaders := strings.Split(dataHeader, " ") // Bearer token
		if dataHeaders[0] != "Bearer" {
			// log.Println("GetCookie cookie:", "Bearer not found")
			return result, errors.New("bearer not found")
		}

		token = dataHeaders[1]
	}

	_, result, err = ParseToken(token)
	if err != nil {
		// log.Println("GetCookie parse token:", err)
		return model.AuthInfo{}, err
	}

	return result, nil
}
