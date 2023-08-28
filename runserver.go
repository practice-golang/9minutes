package server // import "9minutes"

import (
	"log"
	"os"

	"9minutes/handler"
)

func RunServer() {
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
