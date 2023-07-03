package main

import (
	"github.com/gofiber/fiber/v2"
)

func setRouterNotUse(a *fiber.App) {
	// /* Group */
	// gh := a.Group("/hello")
	// gh.GET(`/([\p{L}\d_]+)$`, handler.HelloParam)

	// /* Middleware */
	// gm := a.Group(``, handler.HelloMiddleware)
	// gm.GET(`^/hi/([\p{L}\d_]+)$`, handler.HelloParam)

	// /* Restricted - Cookie */
	// a.POST(`^/signin$`, handler.Signin)

	// // gr := r.Group(``, handler.AuthMiddleware)
	// gr := a.Group(``, handler.RestrictSessionMiddleware)
	// gr.GET(`^/restricted$`, handler.RestrictedHello)

	// /* Restricted - Header */
	// gapi := a.Group(`^/api`, handler.AuthApiMiddleware)
	// gapi.GET(`/restricted$`, handler.RestrictedHello)

	// /* Other GET, POST */
	// a.GET(`^/get-param$`, handler.GetParam)
	// a.Handle(`^/post-form$`, handler.PostForm, "GET", "POST")
	// a.Handle(`^/post-json$`, handler.PostJson, router.AllMethods...)
}
