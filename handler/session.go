package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store

func NewSessionStore() {
	store = session.New(session.Config{
		CookieSecure:   true,
		CookieHTTPOnly: false,
		// CookieSameSite: "Strict",
		CookieSameSite: "None", // For cross-origin
	})
}

func GetSessionUserGrade(c *fiber.Ctx) (string, error) {
	var userGrade string
	var userGradeInterface interface{}

	sess, err := store.Get(c)
	if err != nil {
		return userGrade, err
	}

	userGradeInterface = sess.Get("grade")
	if userGradeInterface != nil {
		userGrade = userGradeInterface.(string)
	}

	return userGrade, err
}
