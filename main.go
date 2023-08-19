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

//go:embed all:static/html/*
var EmbedHTML embed.FS

//go:embed all:static/*
var EmbedStatic embed.FS

var (
	StaticPath = config.StaticPath
	HtmlPath   = config.HtmlPath
	FilesPath  = config.FilesPath
	UploadPath = config.UploadPath

	ListeningIP      string = "localhost"
	ListeningPort    string = "4416"
	ListeningAddress string = ListeningIP + ":" + ListeningPort

	app *fiber.App
)

func main() {
	parseArgs()

	_ = os.Mkdir(UploadPath, os.ModePerm)
	setupConfig()
	handler.NewSessionStore()
	setupDB()
	setupRouter()

	println("Listen", ListeningAddress)
	log.Fatal(app.Listen(ListeningAddress))
}
