package main // import "github.com/practice-golang/9minutes"

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gopkg.in/ini.v1"

	_ "modernc.org/sqlite"

	"github.com/practice-golang/9minutes/auth"
	"github.com/practice-golang/9minutes/board"
	"github.com/practice-golang/9minutes/comments"
	"github.com/practice-golang/9minutes/config"
	"github.com/practice-golang/9minutes/contents"
	"github.com/practice-golang/9minutes/db"
	"github.com/practice-golang/9minutes/user"
)

var (
	//go:embed static
	staticPATH embed.FS
	//go:embed templates
	templatePATH embed.FS
	//go:embed samples/9minutes.ini
	sampleINI string

	jwtKey = []byte("9minutes")
)

func setupDB() error {
	var err error
	info := config.DbInfo

	switch config.DbInfo.Type {
	case "sqlite":
		db.DBType = db.SQLITE
		db.Dsn = info.Filename
	case "mysql":
		db.DBType = db.MYSQL
		db.Dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/",
			info.User, info.Password, info.Server, info.Port)
		db.DatabaseName = info.Database
		db.BoardManagerTable = db.DatabaseName + "." + db.BoardManagerTable
	case "postgres":
		db.DBType = db.POSTGRES

		// DB creation
		if info.Database != "postgres" {
			db.Dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
				info.Server, info.Port, info.User, info.Password)
			db.Dbi, err = db.InitDB(db.DBType)
			if err != nil {
				log.Fatal("InitDB - CreateDB: ", err)
			}
			err = db.Dbi.CreateDB()
			if err != nil {
				log.Println("Create DB (ignore this if msg is already exists): ", err)
			}
			db.Dbo.Close()
		}

		db.Dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			info.Server, info.Port, info.User, info.Password, info.Database)

		db.DatabaseName = info.Schema
		db.BoardManagerTable = db.DatabaseName + "." + db.BoardManagerTable
	case "sqlserver":
		db.DBType = db.SQLSERVER
		db.Dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			info.User, info.Password, info.Server, info.Port, info.Database)
		db.DatabaseName = info.Database
		db.BoardManagerTable = db.DatabaseName + ".dbo." + db.BoardManagerTable
	default:
		log.Fatal("nothing to support DB")
	}

	db.Dbi, err = db.InitDB(db.DBType)
	if err != nil {
		log.Fatal("InitDB: ", err)
	}

	recreate := false
	err = db.Dbi.CreateBoardManagerTable(recreate)
	if err != nil {
		log.Fatal("Create Board manager Table: ", err)
	}
	err = db.Dbi.CreateUserFieldTable(recreate)
	if err != nil {
		log.Fatal("Create User field Table: ", err)
	}
	err = db.Dbi.CreateUserTable(recreate)
	if err != nil {
		log.Fatal("Create User Table: ", err)
	}

	return err
}

func dumpHandler(c echo.Context, reqBody, resBody []byte) {
	header := time.Now().Format("2006-01-02 15:04:05") + " - "
	body := string(reqBody)
	body = strings.Replace(body, "\r\n", "", -1)
	body = strings.Replace(body, "\n", "", -1)
	data := header + body + "\n"

	f, err := os.OpenFile(
		"request-body.log",
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		os.FileMode(0777),
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	if _, err = f.WriteString(data); err != nil {
		log.Println(err)
		return
	}
}

// Admin page
func adminTemplateHandler(c echo.Context) error {
	if !auth.IsAdmin(c) {
		return c.JSON(http.StatusForbidden, "")
	}

	status := http.StatusOK
	sPath := fs.FS(staticPATH)

	var content []byte
	var err error
	var contents string

	filePath := c.Request().URL.RequestURI()
	// Remove slash
	if filePath[0] == '/' {
		filePath = filePath[1:]
	}

	content, err = fs.ReadFile(sPath, filePath)
	if err != nil {
		return c.HTML(http.StatusNotFound, "Not found")
	}

	contents = string(content)

	return c.HTML(status, contents)
}

// Choose board template
func boardTemplateHandler(c echo.Context) error {
	status := http.StatusOK
	code := c.QueryParam("code")
	mode := c.QueryParam("mode") // write
	tPath := fs.FS(templatePATH)

	log.Println("boardTemplateHandler mode: ", mode)

	boardInfos := board.GetBoardByCode(code)
	boardType := ""
	tableName := ""
	boardCode := ""
	if len(boardInfos) > 0 {
		boardType = boardInfos[0].Type.String
		tableName = boardInfos[0].Table.String
		boardCode = boardInfos[0].Code.String
	}

	var content []byte
	var err error
	var contents string

	switch boardType {
	case "basic-board":
		switch mode {
		case "write", "edit":
			content, err = fs.ReadFile(tPath, "templates/basic-board-writer.html")
		case "read":
			content, err = fs.ReadFile(tPath, "templates/basic-board-reader.html")
		default:
			content, err = fs.ReadFile(tPath, "templates/basic-board-list.html")
		}

		contents = string(content)
	case "custom-board":
		switch mode {
		case "write", "edit":
			content, err = fs.ReadFile(tPath, "templates/custom-board-writer.html")
		case "read":
			content, err = fs.ReadFile(tPath, "templates/custom-board-reader.html")
		default:
			content, err = fs.ReadFile(tPath, "templates/custom-board-list.html")
		}

		columnsInterface := []map[string]interface{}{
			{
				"idx":    0,
				"name":   "Idx",
				"column": "IDX",
				"json":   "idx",
				"type":   "number",
				"order":  1,
			},
		}
		for _, f := range boardInfos[0].Fields.([]interface{}) {
			columnsInterface = append(columnsInterface, f.(map[string]interface{}))
		}

		columns, _ := json.Marshal(columnsInterface)
		contents = strings.ReplaceAll(string(content), "'##__COLUMNS__##'", string(columns))
	case "custom-tablelist":
		content, err = fs.ReadFile(tPath, "templates/custom-tablelist.html")
		columnsInterface := []map[string]interface{}{
			{
				"idx":    0,
				"name":   "Idx",
				"column": "IDX",
				"json":   "idx",
				"type":   "number",
				"order":  1,
			},
		}
		for _, f := range boardInfos[0].Fields.([]interface{}) {
			columnsInterface = append(columnsInterface, f.(map[string]interface{}))
		}

		columns, _ := json.Marshal(columnsInterface)
		contents = strings.ReplaceAll(string(content), "'##__COLUMNS__##'", string(columns))
	default:
		return c.HTML(http.StatusNotFound, "404 Not found")
	}

	contents = strings.ReplaceAll(contents, "'##__TABLE_NAME__##'", tableName)
	contents = strings.ReplaceAll(contents, "'##__BOARD_CODE__##'", boardCode)

	if err != nil {
		log.Println("template READ: ", boardType, err)
	}

	return c.HTML(status, contents)
}

func setupServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(
		middleware.CORS(),
		middleware.Recover(),
	)

	staticRoot := "/static"
	// routeTargetFilename := "$1"
	rewriteTargetFilename := "page-loader"

	jwtConfigNoResponse := middleware.JWTConfig{
		Claims:     &auth.CustomClaims{},
		SigningKey: jwtKey,
		ErrorHandlerWithContext: func(e error, c echo.Context) error {
			_ = user.CheckPermission(c)

			return c.String(http.StatusOK, "")
		},
	}

	jwtConfigRestricted := middleware.JWTConfig{
		Claims:     &auth.CustomClaims{},
		SigningKey: jwtKey,
		ErrorHandlerWithContext: func(e error, c echo.Context) error {
			result := map[string]string{"msg": e.Error()}

			return c.JSON(http.StatusUnauthorized, result)
		},
	}

	jwtConfigBoard := middleware.JWTConfig{
		Skipper: func(c echo.Context) bool {
			code := c.QueryParam("code")
			mode := c.QueryParam("mode") // read, write

			boardInfos := board.GetBoardByCode(code)
			if len(boardInfos) == 0 {
				return false
			}

			switch true {
			case (mode == "write" && boardInfos[0].GrantWrite.String == "all") ||
				(mode != "write" && boardInfos[0].GrantRead.String == "all"):
				return true
			default:
				return false
			}
		},
		Claims:     &auth.CustomClaims{},
		SigningKey: jwtKey,
		ErrorHandlerWithContext: func(e error, c echo.Context) error {
			result := map[string]string{"msg": e.Error()}

			return c.JSON(http.StatusUnauthorized, result)
		},
	}

	contentHandler := echo.WrapHandler(http.FileServer(http.FS(staticPATH)))
	contentRewriteAdmin := middleware.RewriteWithConfig(middleware.RewriteConfig{
		RegexRules: map[*regexp.Regexp]string{
			regexp.MustCompile(`^/admin/([^\?]+)(\?(.*)|)`): staticRoot + "/" + rewriteTargetFilename + ".html",
		},
	})

	e.GET("/admin/*", contentHandler, contentRewriteAdmin)

	contentRewriteBody := middleware.RewriteWithConfig(middleware.RewriteConfig{
		RegexRules: map[*regexp.Regexp]string{
			regexp.MustCompile(`^/page/([^\?]+)(\?(.*)|)`): staticRoot + "/pages/$1.html",
		},
	})
	bd := e.Group("/page")
	bd.Use(middleware.JWTWithConfig(jwtConfigRestricted))
	// bd.GET("/*", contentHandler, contentRewriteBody)
	bd.GET("/*", adminTemplateHandler, contentRewriteBody)

	contentRewriteUsers := middleware.RewriteWithConfig(middleware.RewriteConfig{
		RegexRules: map[*regexp.Regexp]string{
			regexp.MustCompile(`^/users/([^\?]+)(\?(.*)|)`): staticRoot + "/users/$1.html",
		},
	})
	e.GET("/users/*", contentHandler, contentRewriteUsers)

	contentRewriteAssets := middleware.RewriteWithConfig(middleware.RewriteConfig{
		RegexRules: map[*regexp.Regexp]string{
			regexp.MustCompile(`^/assets/([^\?]+)(\?(.*)|)`): staticRoot + "/assets/$1",
		},
	})
	e.GET("/assets/*", contentHandler, contentRewriteAssets)

	contentRewrite := middleware.Rewrite(map[string]string{"/*": staticRoot + "/"})
	e.GET("/*", contentHandler, contentRewrite)

	e.GET("/board", boardTemplateHandler)

	a := e.Group("/api/admin")
	a.Use(middleware.JWTWithConfig(jwtConfigRestricted))
	a.GET("/board/:idx", board.GetBoard)
	a.GET("/boards", board.GetBoards)
	a.POST("/boards", board.SearchBoards)
	a.PUT("/boards", board.AddBoards)
	a.PATCH("/board", board.EditBoard)
	a.DELETE("/board/:idx", board.DeleteBoard)
	a.POST("/total-page", board.GetTotalPage)

	a.GET("/user-fields", user.GetUserFields)
	a.PUT("/user-fields", user.AddUserFields)
	a.PATCH("/user-fields", user.EditUserFields)
	a.DELETE("/user-fields/:idx", user.DeleteUserFields)

	a.GET("/user-columns", user.GetUserColumns)
	a.POST("/users", user.GetUsers)
	a.PUT("/users", user.AddUser)
	a.PATCH("/users", user.EditUser)
	a.DELETE("/users/:idx", user.DeleteUser)

	u := e.Group("/api/user")
	u.POST("/login", user.Login)
	u.GET("/token", user.ReissueToken)

	ua := e.Group("/api/user")
	ua.Use(middleware.JWTWithConfig(jwtConfigNoResponse))
	ua.POST("/token/verify", user.VerifyToken)
	ua.GET("/permission", user.CheckPermission)

	bb := e.Group("/api/basic-board")
	bb.Use(middleware.JWTWithConfig(jwtConfigBoard))
	bb.POST("/contents", contents.GetContentsListBasicBoard)
	bb.PUT("/contents", contents.AddContentsBasicBoard)
	bb.PATCH("/contents", contents.UpdateContentsBasicBoard)
	bb.DELETE("/contents", contents.DeleteContentsBasicBoard)
	bb.POST("/total-page", contents.GetContentsTotalPage)

	cb := e.Group("/api/custom-board")
	cb.Use(middleware.JWTWithConfig(jwtConfigBoard))
	cb.POST("/contents-list", contents.GetContentsListCustomBoard)
	cb.PUT("/contents-list", contents.AddContentsListCustomBoard)
	cb.PATCH("/contents-list", contents.UpdateContentsListCustomBoard)
	cb.DELETE("/contents-list", contents.DeleteContentsListCustomBoard)
	cb.POST("/total-page", contents.GetContentsTotalPageMAP)

	cm := e.Group("/api/comments")
	cm.POST("", comments.GetComments)
	cm.PUT("", comments.AddComments)

	return e
}

func main() {
	var err error
	listenAddress := "127.0.0.1"
	listenPort := "2510"
	cfg, err := ini.Load("9minutes.ini")

	if err != nil {
		log.Print("Fail to read ini. ")

		f, err := os.Create("9minutes.ini")
		if err != nil {
			log.Fatal("Create INI: ", err)
		}
		defer f.Close()

		_, err = f.WriteString(sampleINI + "\n")
		if err != nil {
			log.Fatal("Create INI: ", err)
		}

		log.Println("9minutes.ini is created")
	}

	// Rewrite path regexp check
	// r := regexp.MustCompile(`^/admin/([^\?]+)(\?(.*)?)`)
	// out := r.FindStringSubmatch(("/admin/board?idx=123"))
	// log.Println(out)

	if cfg != nil {
		config.DbInfo.Type = cfg.Section("database").Key("DBTYPE").String()
		config.DbInfo.Server = cfg.Section("database").Key("ADDRESS").String()
		config.DbInfo.Port, _ = cfg.Section("database").Key("PORT").Int()
		config.DbInfo.User = cfg.Section("database").Key("USER").String()
		config.DbInfo.Password = cfg.Section("database").Key("PASSWORD").String()
		config.DbInfo.Database = cfg.Section("database").Key("DATABASE").String()
		config.DbInfo.Schema = cfg.Section("database").Key("SCHEMA").String()
		config.DbInfo.Filename = cfg.Section("database").Key("FILENAME").String()

		if cfg.Section("server").HasKey("ADDRESS") {
			listenAddress = cfg.Section("server").Key("ADDRESS").String()
		}
		if cfg.Section("server").HasKey("PORT") {
			listenPort = cfg.Section("server").Key("PORT").String()
		}
	}

	var fileConnectionLog *os.File

	// sql where target
	db.UpdateScope = []string{"idx", "IDX"} // UPDATE ... WHERE IDX=?
	db.IgnoreScope = []string{}             // Ignore if nil or null
	db.OrderScope = "IDX"

	err = setupDB()
	if err != nil {
		log.Fatal("Setup DB: ", err)
	}

	auth.JwtKey = jwtKey
	e := setupServer()

	fileConnectionLog, err = os.OpenFile(
		"connection.log",
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		os.FileMode(0777),
	)
	if err != nil {
		log.Fatalln("Connection log: ", err)
	}
	defer fileConnectionLog.Close()

	e.Use(
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: `${time_rfc3339} - remote_ip:${remote_ip}, host:${host}, ` +
				`method:${method}, uri:${uri},status:${status}, error:${error}, ` +
				`${header:Authorization}, query:${query:property}, form:${form}, ` + "\n",
			Output: fileConnectionLog,
		}),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"*"},
			AllowHeaders:     []string{"Authorization", "Content-Type"},
			AllowCredentials: true,
			AllowMethods: []string{
				echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE,
				echo.HEAD, echo.OPTIONS,
			},
		}),
	)

	e.Use(middleware.BodyDump(dumpHandler))

	// e.Logger.Fatal(e.Start("127.0.0.1:2918"))
	e.Logger.Fatal(e.Start(listenAddress + ":" + listenPort))
}
