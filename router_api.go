package main

import (
	"9minutes/handler"

	"github.com/gofiber/fiber/v2"
)

func setAPIs(a *fiber.App) {
	/* API */
	gapi := a.Group("/api")
	gapi.Get("/health", handler.HealthCheckAPI)
	gapi.Post("/login", handler.LoginAPI)
	gapi.Get("/logout", handler.LogoutAPI)
	gapi.Post("/signup", handler.SignupAPI)

	/* API myinfo */
	gmyinfo := gapi.Group("/myinfo") // Require add session middleware
	gmyinfo.Get("/", handler.GetMyInfo)
	gmyinfo.Put("/", handler.UpdateMyInfo)
	gmyinfo.Delete("/", handler.ResignUser)
}

func setApiUploader(r *fiber.App) {
	/* API Uploader */
	gupload := r.Group("/api/uploader") // Require add session middleware
	gupload.Post("/", handler.UploadFile)
	gupload.Post("/files-info", handler.FilesInfo)
	gupload.Delete("/", handler.DeleteFiles)

	// gu.POST(`/image$`, handler.UploadImage)

	// Delete all of files, images, title-image which is(are) uploaded during writing or editing on board, when cancel
}
