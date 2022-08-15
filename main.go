package main // import "9minutes"

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"9minutes/config"
	"9minutes/db"
	"9minutes/email"
	"9minutes/handler"
	"9minutes/logging"
)

//go:embed 9minutes.ini
var sampleINI string

//go:embed html/*
var Content embed.FS

//go:embed embed/*
var EmbedStatic embed.FS

var StaticPath = config.StaticPath
var UploadPath = config.UploadPath

var sessionStoreInfo = config.StoreInfoMemory

var (
	ListeningIP   string = "localhost"
	ListeningPort string = "4416"
	ServerHandler http.Handler

	ListeningAddress string
)

func firstRun() {
	// db.Info = config.DatabaseInfoSQLite
	// db.Info = config.DatabaseInfoMySQL
	// db.Info = config.DatabaseInfoPgPublic
	// db.Info = config.DatabaseInfoSqlServer

	// db.Info = config.DatabaseInfoOracle
	// db.InfoOracleAdmin = config.DatabaseInfoOracleSystem
	db.Info = config.DatabaseInfoOracleCloud
	db.InfoOracleAdmin = config.DatabaseInfoOracleCloudAdmin

	email.Info = config.EmailServerDirect
	// email.Info = config.EmailServerSMTP

	envPORT := os.Getenv("PORT")
	envDBMS := os.Getenv("DATABASE_TYPE")
	if envPORT != "" {
		ListeningIP = "0.0.0.0"
		ListeningPort = envPORT

		StaticPath = "./static"
		UploadPath = "./upload"
		handler.StoreRoot = "./html"

		envAddress := os.Getenv("DATABASE_ADDRESS")
		envDbPort := os.Getenv("DATABASE_PORT")
		envProtocol := os.Getenv("DATABASE_PROTOCOL")
		envDbName := os.Getenv("DATABASE_NAME")
		envDbID := os.Getenv("DATABASE_ID")
		envDbPassword := os.Getenv("DATABASE_PASSWORD")

		switch envDBMS {
		case "mysql":
			db.Info = config.DatabaseInfoMySQL
			db.Info.Addr = envAddress
			db.Info.Port = envDbPort
			db.Info.Protocol = envProtocol
			db.Info.DatabaseName = envDbName
			db.Info.GrantID = envDbID
			db.Info.GrantPassword = envDbPassword
		case "postgres":
			db.Info = config.DatabaseInfoPgPublic
		case "sqlserver":
			db.Info = config.DatabaseInfoSqlServer
		default:
			db.Info = config.DatabaseInfoMySQL
		}
	} else {
		setupINI()
	}

	ListeningAddress = ListeningIP + ":" + ListeningPort
}

func writeEmbedToDir(dir string) {
	rootDir, err := Content.ReadDir(dir)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	for _, f := range rootDir {
		if dir+"/"+f.Name() == "html/admin" {
			continue
		}

		if f.IsDir() {
			os.MkdirAll(dir+"/"+f.Name(), os.ModePerm)
			writeEmbedToDir(dir + "/" + f.Name())
		} else {
			sf, err := os.Create(dir + "/" + f.Name())
			if err != nil {
				log.Fatal("failed to create file for embedded html: ", err)
			}
			defer sf.Close()

			ef, err := Content.ReadFile(dir + "/" + f.Name())
			if err != nil {
				log.Fatal("failed to read embedded html: ", err)
			}

			_, err = sf.Write(ef)
			if err != nil {
				log.Fatal("failed to write file from embedded html: ", err)
			}
		}
	}
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
			writeEmbedToDir("html")
			fmt.Println("done to export html files")
		case "dkim":
			email.GenerateKeys()
			fmt.Println("done to generate dkim keys")
		}

		os.Exit(0)
	}

	firstRun()
	doSetup()

	logging.Object.Log().Timestamp().Str("listen", ListeningIP+"\n").Send()
	println("Listen", ListeningAddress)

	err := http.ListenAndServe(ListeningAddress, ServerHandler)
	if err != nil {
		logging.Object.Warn().Err(err).Timestamp().Msg("Server start failed")
	}
}
