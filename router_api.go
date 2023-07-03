package main

import (
	"9minutes/handler"
	"9minutes/router"

	"github.com/gofiber/fiber/v2"
)

func setAPIs(a *fiber.App) {
	/* API */
	g := a.Group("/api")
	g.Get("/health", handler.HealthCheck)
	g.Post("/login", handler.LoginAPI)

	// /* API myinfo */
	// gmi := a.Group("/api/myinfo", handler.AuthApiSessionMiddleware)
	gmi := a.Group("/api/myinfo") // Require add session middleware
	gmi.Get("/", handler.GetMyInfo)
	// gmi.PUT(`(/?)$`, handler.UpdateMyInfo)
	// gmi.DELETE(`(/?)$`, handler.ResignUser)

	// /* API File & Directory */
	// g.POST(`/dir/list$`, handler.HandleGetDir)

	// /* Captcha */
	// g.GET(`/captcha/[^/]+\.png$`, handler.GetCaptchaImage)
	// g.PATCH(`/captcha$`, handler.RenewCaptcha)
}

func setApiLogin(a *fiber.App) {
	/* Login, Logout */
	a.Post("/login", handler.Login)
	// a.GET(`^/logout`, handler.Logout)
	// a.POST(`/api/signup`, handler.Signup)

	// /* Reset password */
	// a.POST(`/password-reset$`, handler.ResetPassword)

	// /* User verification - Should be moved at next time */
	// a.GET(`/verify`, handler.UserVerification)
}

func setApiUploader(r *router.App) {
	/* API Uploader */
	gu := r.Group(`^/api/uploader`, handler.AuthApiSessionMiddleware)
	gu.POST(`/file$`, handler.UploadFile)
	gu.POST(`/image$`, handler.UploadImage)

	// Delete all of files, images, title-image which is(are) uploaded during writing or editing on board, when cancel
	gu.DELETE(`/file$`, handler.DeleteFiles)
}
