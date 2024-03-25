package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store

func NewSessionStore() {
	store = session.New(session.Config{
		CookieSecure:   false,
		CookieHTTPOnly: true,
		// CookieSameSite: "None", // For cross-origin
		// CookieSameSite: "Strict",
		CookieSameSite: "Lax",
	})
}

func getSessionValue(sess *session.Session, key string) string {
	var result string

	value := sess.Get(key)
	if value != nil {
		result = value.(string)
	}

	return result
}

func GetSessionUserGrade(c *fiber.Ctx) (string, error) {
	sess, err := store.Get(c)
	if err != nil {
		return "", err
	}

	result := getSessionValue(sess, "grade")

	return result, nil
}
