package server

import (
	"9minutes/config"
	"embed"

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
