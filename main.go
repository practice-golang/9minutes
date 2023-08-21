package main // import "9minutes"

import (
	"embed"
	"log"
	"os"

	"9minutes/config"
	"9minutes/handler"

	"github.com/gofiber/fiber/v2"
)

//go:embed 9minutes.ini
var sampleINI string

//go:embed static/html/favicon.png
var Favicon embed.FS

//go:embed all:static/html/*
var EmbedHTML embed.FS

//go:embed all:static/*
var EmbedStatic embed.FS

var (
	PathPrefix string = ""
	StaticPath string = config.StaticPath
	HtmlPath   string = config.HtmlPath
	FilesPath  string = config.FilesPath
	UploadPath string = config.UploadPath

	ListeningIP      string = "localhost"
	ListeningPort    string = "4416"
	ListeningAddress string = ListeningIP + ":" + ListeningPort

	IsStaticEmbed bool = false
)

var app *fiber.App

func main() {
	parseArgs()

	if _, err := os.Stat(HtmlPath); os.IsNotExist(err) {
		IsStaticEmbed = true
	}

	_ = os.Mkdir(UploadPath, os.ModePerm)
	setupConfig()
	handler.NewSessionStore()
	setupDB()
	setupRouter()

	println("Listen", ListeningAddress)
	log.Fatal(app.Listen(ListeningAddress))
}
