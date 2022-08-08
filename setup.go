package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"9minutes/auth"
	"9minutes/config"
	"9minutes/db"
	"9minutes/email"
	"9minutes/fd"
	"9minutes/handler"
	"9minutes/logging"
	"9minutes/router"
	"9minutes/wsock"

	"github.com/alexedwards/scs/etcdstore"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/gomodule/redigo/redis"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/rs/cors"
	"gopkg.in/ini.v1"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
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

		if cfg.Section("session").HasKey("STORE_TYPE") {
			switch cfg.Section("session").Key("STORE_TYPE").String() {
			case "etcd":
				sessionStoreInfo.StoreType = auth.ETCD
			case "redis":
				sessionStoreInfo.StoreType = auth.REDIS
			default:
				sessionStoreInfo.StoreType = auth.MEMSTORE
			}

			if sessionStoreInfo.StoreType != auth.MEMSTORE {
				if cfg.Section("session").HasKey("ADDRESS") {
					sessionStoreInfo.Address = cfg.Section("session").Key("ADDRESS").String()
				}
				if cfg.Section("session").HasKey("PORT") {
					sessionStoreInfo.Address = cfg.Section("session").Key("PORT").String()
				}
			}
		}

		if cfg.Section("database").HasKey("DBTYPE") {
			switch cfg.Section("database").Key("DBTYPE").String() {
			case "mysql":
				db.Info = config.DatabaseInfoMySQL
			case "postgres":
				db.Info = config.DatabaseInfoPgPublic
			case "sqlserver":
				db.Info = config.DatabaseInfoSqlServer
			default:
				db.Info = config.DatabaseInfoSQLite
				if cfg.Section("database").HasKey("FILENAME") {
					db.Info.FilePath = cfg.Section("database").Key("FILENAME").String()
				}
			}

			if db.Info.DatabaseType != db.SQLITE {
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

				if db.Info.DatabaseType == db.SQLSERVER || db.Info.DatabaseType == db.POSTGRES {
					if cfg.Section("database").HasKey("SCHEMA") {
						db.Info.SchemaName = cfg.Section("database").Key("SCHEMA").String()
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
			dkimKey, err := os.ReadFile(cfg.Section("email").Key("DKIM_PATH").String())
			if err != nil {
				panic(err)
			}
			email.Info.Service.KeyDKIM = string(dkimKey)
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

func setupSession() {
	auth.SessionManager = scs.New()

	switch sessionStoreInfo.StoreType {
	case auth.ETCD:
		addr := config.StoreInfoETCD.Address + ":" + config.StoreInfoETCD.Port
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{addr},
			DialTimeout: 5 * time.Second,
		})

		if err != nil {
			log.Fatal(err)
		}

		auth.SessionManager.Store = etcdstore.New(cli)
	case auth.REDIS:
		addr := config.StoreInfoRedis.Address + ":" + config.StoreInfoRedis.Port
		pool := &redis.Pool{
			MaxIdle: 10,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", addr)
			},
		}

		auth.SessionManager.Store = redisstore.New(pool)
	default:
		auth.SessionManager.Store = memstore.New()
	}

	auth.SessionManager.Lifetime = 3 * time.Hour
	auth.SessionManager.IdleTimeout = 20 * time.Minute
	auth.SessionManager.Cookie.Name = "session_id"

	// auth.SessionManager.Cookie.Domain = "example.com"
	// auth.SessionManager.Cookie.HttpOnly = true
	// auth.SessionManager.Cookie.Path = "/example/"
	// auth.SessionManager.Cookie.Persist = true
	// auth.SessionManager.Cookie.SameSite = http.SameSiteStrictMode
	// auth.SessionManager.Cookie.Secure = true
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

	// Not use
	// err = db.Obj.CreateTable()
	// if err != nil {
	// 	log.Fatal("CreateTable:", err)
	// }

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

func setupKey() {
	auth.Secret = "practice-golang/9m secret"

	privKeyExist := fd.CheckFileExists(auth.JwtPrivateKeyFileName, false)
	pubKeyExist := fd.CheckFileExists(auth.JwtPublicKeyFileName, false)
	if privKeyExist && pubKeyExist {
		auth.LoadRsaKeys()
	} else {
		auth.GenerateRsaKeys()
		auth.SaveRsaKeys()
	}

	err := auth.GenerateKeySet()
	if err != nil {
		panic(err)
	}
}

func setupLogger() {
	logging.SetupLogger()

	go func() {
		now := time.Now()
		zone, i := now.Zone()
		nextDay := now.AddDate(0, 0, 1).In(time.FixedZone(zone, i))
		nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location())
		restTimeNextDay := time.Until(nextDay)
		time.Sleep(restTimeNextDay)
		for {
			if time.Now().Format("15") == "00" {
				logging.RenewLogger()
				time.Sleep(24 * time.Hour)
			} else {
				time.Sleep(time.Second)
			}
		}
	}()
}

func setupRouter() {
	router.StaticPath = StaticPath
	router.UploadPath = UploadPath
	router.Content = Content
	router.EmbedStatic = EmbedStatic

	router.SetupStaticServer()

	r := router.New()

	setPAGEs(r)        // HTML, Assets, Login/Signup
	setPageAdmin(r)    // Admin
	setPageContent(r)  // Content
	setPageMyPage(r)   // MyPage
	setPageHTMLs(r)    // HTML for both user and anonymous
	setApiBoard(r)     // API Board
	setApiUploader(r)  // API Uploader
	setApiLogin(r)     // API Login, Logout, Signup
	setApiAdmin(r)     // API Admin
	setAPIs(r)         // API
	setOthers(r)       // Others
	setRouterNotUse(r) // Not use, should be removed at future

	ServerHandler = auth.SessionManager.LoadAndSave(cors.Default().Handler(r))
	// ServerHandler = cors.Default().Handler(r)
	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"http://"+listen},
	// 	AllowedMethods:   []string{"GET"},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowCredentials: true,
	// 	Debug:            false,
	// })
	// ServerHandler := c.Handler(r)

}

func doSetup() {
	_ = os.Mkdir(StaticPath, os.ModePerm)
	_ = os.Mkdir(UploadPath, os.ModePerm)
	_ = os.Mkdir(config.HtmlPath, os.ModePerm)

	setupSession()
	setupDB()
	setupKey()
	setupLogger()
	setupRouter()

	wsock.InitWebSocketChat()
}
