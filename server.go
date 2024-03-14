package server // import "9minutes"

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"9minutes/handler"
	"9minutes/internal/email"
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

func parseArgs() {
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
}

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
