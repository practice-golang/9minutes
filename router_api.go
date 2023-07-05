package main

import (
	"9minutes/handler"

	"github.com/gofiber/fiber/v2"
)

func setAPIs(a *fiber.App) {
	/* API */
	gapi := a.Group("/api")
	gapi.Get("/health", handler.HealthCheck)
	gapi.Post("/login", handler.LoginAPI)
	gapi.Get("/logout", handler.LogoutAPI)
	gapi.Post("/signup", handler.SignupAPI)

	/* API myinfo */
	gmyinfo := gapi.Group("/myinfo") // Require add session middleware
	gmyinfo.Get("/", handler.GetMyInfo)
	gmyinfo.Put("/", handler.UpdateMyInfo)
	gmyinfo.Delete("/", handler.ResignUser)
}

// func setApiUploader(r *router.App) {
// 	/* API Uploader */
// 	gu := r.Group(`^/api/uploader`, handler.AuthApiSessionMiddleware)
// 	gu.POST(`/file$`, handler.UploadFile)
// 	gu.POST(`/image$`, handler.UploadImage)

// 	// Delete all of files, images, title-image which is(are) uploaded during writing or editing on board, when cancel
// 	gu.DELETE(`/file$`, handler.DeleteFiles)
// }
