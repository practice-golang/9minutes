package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
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
	iniPath := "9minutes.ini"

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

		fmt.Println("9minutes.ini is created")
		fmt.Println("Please modify 9minutes.ini then run again")

		os.Exit(1)
	}

	if cfg != nil {
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
				db.InfoOracleAdmin = config.DatabaseInfoOracleSystem
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

					if cfg.Section("database_admin").HasKey("ADDRESS") {
						db.InfoOracleAdmin.Addr = cfg.Section("database_admin").Key("ADDRESS").String()
					}
					if cfg.Section("database_admin").HasKey("PORT") {
						db.InfoOracleAdmin.Port = cfg.Section("database_admin").Key("PORT").String()
					}
					if cfg.Section("database_admin").HasKey("USER") {
						db.InfoOracleAdmin.GrantID = cfg.Section("database_admin").Key("USER").String()
					}
					if cfg.Section("database_admin").HasKey("PASSWORD") {
						db.InfoOracleAdmin.GrantPassword = cfg.Section("database_admin").Key("PASSWORD").String()
					}
					if cfg.Section("database_admin").HasKey("DATABASE") {
						db.InfoOracleAdmin.DatabaseName = cfg.Section("database_admin").Key("DATABASE").String()
					}
					if cfg.Section("database_admin").HasKey("FILEPATH") {
						db.InfoOracleAdmin.FilePath = cfg.Section("database_admin").Key("FILEPATH").String()
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
					panic(err)
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

func setupDB() {
	var err error

	err = db.SetupDB()
	if err != nil {
		log.Fatal("SetupDB:", err)
	}

	err = db.Obj.CreateDB()
	if err != nil {
		log.Fatal("CreateDB:", err)
	}

	err = db.Obj.CreateBoardTable()
	if err != nil {
		log.Fatal("CreateBoardTable:", err)
	}
	err = db.Obj.CreateUploadTable()
	if err != nil {
		log.Fatal("CreateUploadTable:", err)
	}

	err = db.Obj.CreateUserTable()
	if err != nil {
		log.Fatal("CreateUserTable:", err)
	}
	err = db.Obj.CreateUserVerificationTable()
	if err != nil {
		log.Fatal("CreateUserVerificationTable:", err)
	}
}

func setupRouter() {
	engine := html.New("./static/html", ".html")
	engine.AddFunc("unescape", func(s string) template.HTML {
		return template.HTML(s)
	})

	// engine.Debug(true)
	cfg := fiber.Config{
		AppName:               "9minutes",
		DisableStartupMessage: false,
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
}

func doSetup() {
	// _ = os.Mkdir(StaticPath, os.ModePerm)
	// _ = os.Mkdir(config.HtmlPath, os.ModePerm)
	_ = os.Mkdir(UploadPath, os.ModePerm)

	handler.NewSessionStore()

	// setupSession()
	setupDB()
	// setupKey()
	setupRouter()
}
