package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"9minutes/config"
	"9minutes/handler"
	"9minutes/internal/db"
	"9minutes/internal/email"
	"9minutes/model"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"

	"gopkg.in/ini.v1"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/sijms/go-ora/v2"
	_ "modernc.org/sqlite"
)

func setupINI() {
	iniPath := "config.ini"

	cfg, err := ini.Load(iniPath)
	if err != nil {
		f, err := os.Create(iniPath)
		if err != nil {
			log.Fatal("Create INI: ", err)
		}
		defer f.Close()

		_, err = f.WriteString(sampleINI + "\n")
		if err != nil {
			log.Fatal("Create INI: ", err)
		}

		fmt.Println("config.ini is created")
		fmt.Println("Please modify config.ini then run again")

		os.Exit(1)
	}

	if cfg != nil {
		if cfg.Section("").HasKey("SITE_NAME") {
			config.SiteName = cfg.Section("").Key("SITE_NAME").String()
		}

		if cfg.Section("server").HasKey("ADDRESS") {
			ListeningIP = cfg.Section("server").Key("ADDRESS").String()
		}
		if cfg.Section("server").HasKey("PORT") {
			ListeningPort = cfg.Section("server").Key("PORT").String()
		}

		if cfg.Section("dirpaths").HasKey("STATIC_PATH") {
			StaticPath = cfg.Section("dirpaths").Key("STATIC_PATH").String()
		}
		if cfg.Section("dirpaths").HasKey("UPLOAD_PATH") {
			UploadPath = cfg.Section("dirpaths").Key("UPLOAD_PATH").String()
		}
		if cfg.Section("dirpaths").HasKey("HTML_PATH") {
			handler.StoreRoot = cfg.Section("dirpaths").Key("HTML_PATH").String()
		}

		// if cfg.Section("session").HasKey("STORE_TYPE") {
		// 	switch cfg.Section("session").Key("STORE_TYPE").String() {
		// 	case "etcd":
		// 		sessionStoreInfo.StoreType = auth.ETCD
		// 	case "redis":
		// 		sessionStoreInfo.StoreType = auth.REDIS
		// 	default:
		// 		sessionStoreInfo.StoreType = auth.MEMSTORE
		// 	}

		// 	if sessionStoreInfo.StoreType != auth.MEMSTORE {
		// 		if cfg.Section("session").HasKey("ADDRESS") {
		// 			sessionStoreInfo.Address = cfg.Section("session").Key("ADDRESS").String()
		// 		}
		// 		if cfg.Section("session").HasKey("PORT") {
		// 			sessionStoreInfo.Address = cfg.Section("session").Key("PORT").String()
		// 		}
		// 	}
		// }

		if cfg.Section("database").HasKey("DBTYPE") {
			switch cfg.Section("database").Key("DBTYPE").String() {
			case "mysql":
				db.Info = config.DatabaseInfoMySQL
			case "postgres":
				db.Info = config.DatabaseInfoPgPublic
			case "sqlserver":
				db.Info = config.DatabaseInfoSqlServer
			case "oracle":
				db.Info = config.DatabaseInfoOracle
			default:
				db.Info = config.DatabaseInfoSQLite
				if cfg.Section("database").HasKey("FILENAME") {
					db.Info.FilePath = cfg.Section("database").Key("FILENAME").String()
				}
			}

			if db.Info.DatabaseType != model.SQLITE {
				if cfg.Section("database").HasKey("ADDRESS") {
					db.Info.Addr = cfg.Section("database").Key("ADDRESS").String()
				}
				if cfg.Section("database").HasKey("PORT") {
					db.Info.Port = cfg.Section("database").Key("PORT").String()
				}
				if cfg.Section("database").HasKey("USER") {
					db.Info.GrantID = cfg.Section("database").Key("USER").String()
				}
				if cfg.Section("database").HasKey("PASSWORD") {
					db.Info.GrantPassword = cfg.Section("database").Key("PASSWORD").String()
				}
				if cfg.Section("database").HasKey("DATABASE") {
					db.Info.DatabaseName = cfg.Section("database").Key("DATABASE").String()
				}

				if db.Info.DatabaseType == model.SQLSERVER || db.Info.DatabaseType == model.POSTGRES {
					if cfg.Section("database").HasKey("SCHEMA") {
						db.Info.SchemaName = cfg.Section("database").Key("SCHEMA").String()
					}
				}
				if db.Info.DatabaseType == model.ORACLE {
					if cfg.Section("database").HasKey("FILEPATH") {
						db.Info.FilePath = cfg.Section("database").Key("FILEPATH").String()
					}
				}
			}
		}

		if cfg.Section("email").HasKey("USE_EMAIL") {
			email.Info.UseEmail = cfg.Section("email").Key("USE_EMAIL").MustBool(false)
		}
		if cfg.Section("email").HasKey("DOMAIN") {
			email.Info.Domain = cfg.Section("email").Key("DOMAIN").String()
		}
		if cfg.Section("email").HasKey("SEND_DIRECT") {
			email.Info.SendDirect = cfg.Section("email").Key("SEND_DIRECT").MustBool(false)
		}
		if cfg.Section("email").HasKey("DKIM_PATH") {
			if email.Info.UseEmail && email.Info.SendDirect {
				dkimKey, err := os.ReadFile(cfg.Section("email").Key("DKIM_PATH").String())
				if err != nil {
					panic("check dkim path. " + err.Error())
				}
				email.Info.Service.KeyDKIM = string(dkimKey)
			}
		}
		if cfg.Section("email").HasKey("FROM_ADDRESS") {
			email.Info.SenderInfo.Email = cfg.Section("email").Key("FROM_ADDRESS").String()
		}
		if cfg.Section("email").HasKey("FROM_NAME") {
			email.Info.SenderInfo.Name = cfg.Section("email").Key("FROM_NAME").String()
		}
		if cfg.Section("email").HasKey("SERVER") {
			email.Info.Service.Host = cfg.Section("email").Key("SERVER").String()
		}
		if cfg.Section("email").HasKey("PORT") {
			email.Info.Service.Port = cfg.Section("email").Key("PORT").String()
		}
		if cfg.Section("email").HasKey("USER") {
			email.Info.Service.ID = cfg.Section("email").Key("USER").String()
		}
		if cfg.Section("email").HasKey("PASSWORD") {
			email.Info.Service.Password = cfg.Section("email").Key("PASSWORD").String()
		}
	}
}

// setupConfig - setup configurations from environment variables or ini
func setupConfig() {
	db.Info = config.DatabaseInfoSQLite
	// email.Info = config.EmailServerSMTP
	email.Info = config.EmailServerDirect

	envPORT := os.Getenv("PORT")
	envDBMS := os.Getenv("DATABASE_TYPE")
	if envPORT != "" {
		ListeningIP = "0.0.0.0"
		ListeningPort = envPORT

		StaticPath = "static"
		UploadPath = "upload"
		handler.StoreRoot = "static/html"

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

func setupDB() {
	var err error

	if db.SetupDB() != nil {
		log.Fatal("SetupDB:", err)
	}

	if db.Obj.CreateDB() != nil {
		log.Fatal("CreateDB:", err)
	}

	if db.Obj.CreateBoardTable() != nil {
		log.Fatal("CreateBoardTable:", err)
	}
	if db.Obj.CreateUploadTable() != nil {
		log.Fatal("CreateUploadTable:", err)
	}

	if db.Obj.CreateUserTable() != nil {
		log.Fatal("CreateUserTable:", err)
	}
	if db.Obj.CreateUserVerificationTable() != nil {
		log.Fatal("CreateUserVerificationTable:", err)
	}
}

func loadBoardDatas() {
	handler.LoadBoardListData()
}

func loadUserColumnDatas() {
	handler.LoadUserColumnDatas()
}

func setupRouter() {
	var engine *html.Engine

	if IsStaticEmbed {
		htmlRoot, err := fs.Sub(EmbedHTML, "static/html")
		if err != nil {
			log.Fatal(err)
		}
		engine = html.NewFileSystem(http.FS(htmlRoot), ".html")
	} else {
		engine = html.New("./static/html", ".html")
	}

	engine.AddFunc("unescape", unEscape)
	engine.AddFunc("format_date", formatDate)
	engine.AddFunc("js_array", jsArray)

	// engine.Debug(true)

	cfg := fiber.Config{
		AppName:               "9minutes",
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		Views:                 engine,
	}
	app = fiber.New(cfg)
	app.Use(recover.New())
	app.Use(cors.New())

	setApiAdmin(app)    // API Admin
	setApiBoard(app)    // API Board
	setApiUploader(app) // API Uploader
	setAPIs(app)        // API

	setStaticFiles(app) // Files, Assets
	setPage(app)        // HTML templates

	loadBoardDatas()      // Board list
	loadUserColumnDatas() // User column list
}
