package main // import "9minutes"

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"9minutes/config"
	"9minutes/internal/email"

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

func exportStaticEmbed() error {
	exportPath := "."

	err := os.MkdirAll(exportPath, 0755)
	if err != nil {
		return err
	}

	err = fs.WalkDir(EmbedHTML, "static", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			filePath := filepath.Join(exportPath, path)
			err := os.MkdirAll(filepath.Dir(filePath), 0755)
			if err != nil {
				return err
			}

			srcFile, err := EmbedHTML.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			dstFile, err := os.Create(filePath)
			if err != nil {
				return err
			}
			defer dstFile.Close()

			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) > 1 {
		flagGet := flag.String("get", "html|dkim", "Get html, dkim files")

		flag.Usage = func() {
			flagSet := flag.CommandLine
			fmt.Printf("Usage of %s:\n", "9m")
			fmt.Printf("  %-19sRun server\n", "without options")

			order := []string{"get"}
			for _, name := range order {
				flag := flagSet.Lookup(name)
				fmt.Printf("  -%-18s%s\n", flag.Name+" "+flag.Value.String(), flag.Usage)
			}
		}

		flag.Parse()

		switch *flagGet {
		case "html":
			exportStaticEmbed()
			fmt.Println("done to export html files")
		case "dkim":
			email.GenerateKeys()
			fmt.Println("done to generate dkim keys")
		}

		os.Exit(0)
	}

	setupMain()

	println("Listen", ListeningAddress)
	log.Fatal(app.Listen(ListeningAddress))
}
