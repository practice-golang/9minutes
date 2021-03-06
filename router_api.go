package main

import (
	"9minutes/handler"
	"9minutes/router"
)

func setAPIs(r *router.App) {
	/* API */
	g := r.Group(`^/api`)
	g.GET(`/?$`, handler.HealthCheck)
	g.GET(`/hello$`, handler.Hello)
	g.POST(`/signin$`, handler.SigninAPI)

	/* API myinfo */
	gmi := r.Group(`^/api/myinfo`, handler.AuthApiSessionMiddleware)
	gmi.GET(`(/?)$`, handler.GetMyInfo)
	gmi.PUT(`(/?)$`, handler.UpdateMyInfo)

	/* API File & Directory */
	g.POST(`/dir/list$`, handler.HandleGetDir)
}

func setApiLogin(r *router.App) {
	/* Login, Logout */
	r.POST(`^/login`, handler.Login)
	r.GET(`^/logout`, handler.Logout)
	r.POST(`/api/signup`, handler.Signup)
}

func setApiUploader(r *router.App) {
	/* API Uploader */
	gu := r.Group(`^/api/uploader`, handler.AuthApiSessionMiddleware)
	gu.POST(`/file$`, handler.UploadFile)
	gu.POST(`/image$`, handler.UploadImage)
}
