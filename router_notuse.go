package main

import (
	"9minutes/handler"
	"9minutes/router"
)

func setRouterNotUse(r *router.App) {
	/* Group */
	gh := r.Group(`^/hello`)
	gh.Handle(`$`, handler.Hello, "GET", "POST")
	gh.GET(`/([\p{L}\d_]+)$`, handler.HelloParam)

	/* Middleware */
	gm := r.Group(``, handler.HelloMiddleware)
	gm.GET(`^/hi/([\p{L}\d_]+)$`, handler.HelloParam)

	/* Restricted - Cookie */
	r.POST(`^/signin$`, handler.Signin)

	// gr := r.Group(``, handler.AuthMiddleware)
	gr := r.Group(``, handler.RestrictSessionMiddleware)
	gr.GET(`^/restricted$`, handler.RestrictedHello)

	/* Restricted - Header */
	gapi := r.Group(`^/api`, handler.AuthApiMiddleware)
	gapi.GET(`/restricted$`, handler.RestrictedHello)

	/* Other GET, POST */
	r.GET(`^/get-param$`, handler.GetParam)
	r.Handle(`^/post-form$`, handler.PostForm, "GET", "POST")
	r.Handle(`^/post-json$`, handler.PostJson, router.AllMethods...)

}
