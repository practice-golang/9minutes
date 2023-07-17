package main

import (
	"github.com/gofiber/fiber/v2"
)

// setStaticFiles - Set static files
func setStaticFiles(a *fiber.App) {
	a.Static("/files", FilesPath)
}

// setStaticAssets - Set static js, css
func setStaticAssets(a *fiber.App) {
	a.Static("/assets/", HtmlPath+"/assets")
	a.Static("/admin/_app/", HtmlPath+"/admin/_app")
	a.Static("/board/_app/", HtmlPath+"/board/_app")
}
