package main

import (
	"github.com/gofiber/fiber/v2"
)

// setStaticFiles - Set static files
func setStaticFiles(a *fiber.App) {
	a.Static("/favicon.png", HtmlPath+"/favicon.png")
	a.Static("/files", FilesPath)
	a.Static("/upload", UploadPath)
	a.Static("/assets/", HtmlPath+"/assets")
	a.Static("/admin/_app/", HtmlPath+"/admin/_app")
}
