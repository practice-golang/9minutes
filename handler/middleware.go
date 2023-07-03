package handler

import (
	"log"
	"net/http"

	"9minutes/auth"
	"9minutes/router"
)

// HelloMiddleware - middleware for test
func HelloMiddleware(next router.Handler) router.Handler {
	return func(c *router.Context) {
		c.Text(http.StatusInternalServerError, "Middle ware test error")
		// next(c)
	}
}

// AdminOrManagerMiddleware - Check if user is admin or manager
func AdminOrManagerMiddleware(next router.Handler) router.Handler {
	return func(c *router.Context) {
		// claim, err := auth.GetCookieSession(c)
		// if err != nil {
		// 	auth.ExpireCookie(c.ResponseWriter)
		// 	c.Html(http.StatusOK, []byte(`<meta http-equiv="refresh" content="0; url=/"></meta>`))

		// 	return
		// }

		// userInfo, err := crud.GetUserByName(claim.Name.String)
		// if err != nil {
		// 	auth.ExpireCookie(c.ResponseWriter)
		// 	c.Html(http.StatusOK, []byte(`<meta http-equiv="refresh" content="0; url=/"></meta>`))

		// 	return
		// }

		// if userInfo.Grade.String != "admin" && userInfo.Grade.String != "manager" {
		// 	auth.ExpireCookie(c.ResponseWriter)
		// 	c.Html(http.StatusOK, []byte(`<meta http-equiv="refresh" content="0; url=/"></meta>`))

		// 	return
		// }

		// c.AuthInfo = claim

		next(c)
	}
}

// RemoveTrailingSlashMiddleware - remove trailing slash
func RemoveTrailingSlashMiddleware(next router.Handler) router.Handler {
	return func(c *router.Context) {
		if c.URL.Path[len(c.URL.Path)-1:] == "/" {
			http.Redirect(c.ResponseWriter, c.Request, c.URL.Path[:len(c.URL.Path)-1], http.StatusMovedPermanently)
		}

		next(c)
	}
}

// AuthMiddleware - Check if user is logged in for API
func AuthMiddleware(next router.Handler) router.Handler {
	return func(c *router.Context) {
		claim, err := auth.GetClaim(*c.Request, "cookie")
		if err != nil {
			auth.ExpireCookie(c.ResponseWriter)

			c.Text(http.StatusUnauthorized, "Auth error")

			return
		}

		c.AuthInfo = claim

		next(c)
	}
}

// AuthSessionMiddleware - Check if user is logged in for session
func AuthSessionMiddleware(next router.Handler) router.Handler {
	return func(c *router.Context) {
		// claim, err := auth.GetCookieSession(c)
		// if err != nil {
		// 	auth.ExpireCookie(c.ResponseWriter)

		// 	next(c)

		// 	return
		// }

		// c.AuthInfo = claim

		next(c)
	}
}

// RestrictSessionMiddleware - Check if user is logged in for session
func RestrictSessionMiddleware(next router.Handler) router.Handler {
	return func(c *router.Context) {
		// claim, err := auth.GetCookieSession(c)
		// if err != nil {
		// 	auth.ExpireCookie(c.ResponseWriter)

		// 	// c.Text(http.StatusUnauthorized, "Auth error")
		// 	c.Html(http.StatusOK, []byte(`<meta http-equiv="refresh" content="0; url=/login.html"></meta>`))

		// 	return
		// }

		// c.AuthInfo = claim

		next(c)
	}
}

func AuthApiSessionMiddleware(next router.Handler) router.Handler {
	return func(c *router.Context) {
		// claim, err := auth.GetCookieSession(c)
		// if err != nil {
		// 	auth.ExpireCookie(c.ResponseWriter)
		// 	// c.Text(http.StatusUnauthorized, "Auth error")

		// 	next(c)

		// 	return
		// }

		// c.AuthInfo = claim

		next(c)
	}
}

func AuthApiMiddleware(next router.Handler) router.Handler {
	return func(c *router.Context) {
		c.Request.Header.Get("Authorization")
		claim, err := auth.GetClaim(*c.Request, "header")
		if err != nil {
			log.Println("AuthApiMiddleware:", err)
			c.Text(http.StatusUnauthorized, "Auth error")

			return
		}

		c.AuthInfo = claim

		next(c)
	}
}
