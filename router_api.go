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
	gmi.DELETE(`(/?)$`, handler.ResignUser)

	/* API File & Directory */
	g.POST(`/dir/list$`, handler.HandleGetDir)

	/* Captcha */
	g.GET(`/captcha/[^/]+\.png$`, handler.GetCaptchaImage)
	g.PATCH(`/captcha$`, handler.RenewCaptcha)
}

func setApiLogin(r *router.App) {
	/* Login, Logout */
	r.POST(`^/login`, handler.Login)
	r.GET(`^/logout`, handler.Logout)
	r.POST(`/api/signup`, handler.Signup)

	/* User verification - Should be moved at next time */
	r.GET(`/verify`, handler.UserVerification)
}

func setApiUploader(r *router.App) {
	/* API Uploader */
	gu := r.Group(`^/api/uploader`, handler.AuthApiSessionMiddleware)
	gu.POST(`/file$`, handler.UploadFile)
	gu.POST(`/image$`, handler.UploadImage)

	// Delete all of files, images, title-image which is(are) uploaded during writing or editing on board, when cancel
	gu.DELETE(`/file$`, handler.DeleteFiles)
}
