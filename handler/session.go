package handler

import "github.com/gofiber/fiber/v2/middleware/session"

var store *session.Store

func NewSessionStore() {
	store = session.New(session.Config{
		CookieSecure:   false,
		CookieHTTPOnly: false,
		CookieSameSite: "None", // For cross-origin
	})
}
