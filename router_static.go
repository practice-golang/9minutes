package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

// setStaticFiles - Set static files
func setStaticFiles(a *fiber.App) {
	if IsStaticEmbed {
		configFavicon := filesystem.Config{Root: http.FS(EmbedHTML), PathPrefix: "/static/html"}
		configFiles := filesystem.Config{Root: http.FS(EmbedStatic), PathPrefix: "/static/files"}
		configAssets := filesystem.Config{Root: http.FS(EmbedHTML), PathPrefix: "/static/html/assets"}
		configAdminApp := filesystem.Config{Root: http.FS(EmbedHTML), PathPrefix: "/static/html/admin/_app"}

		a.Use("/favicon.png", filesystem.New(configFavicon))
		a.Use("/files", filesystem.New(configFiles))
		a.Use("/assets", filesystem.New(configAssets))
		a.Use("/admin/_app", filesystem.New(configAdminApp))
	} else {
		a.Static("/favicon.png", HtmlPath+"/favicon.png")
		a.Static("/files", FilesPath)
		a.Static("/assets/", HtmlPath+"/assets")
		a.Static("/admin/_app/", HtmlPath+"/admin/_app")
	}

	a.Static("/upload", UploadPath)
}
