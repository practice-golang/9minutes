package handler

import (
	"9minutes/auth"
	"9minutes/consts"
	"9minutes/crud"
	"9minutes/db"
	"9minutes/model"
	"9minutes/np"
	"9minutes/router"
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/blockloop/scan"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

func Login(c *router.Context) {
	failBody := `<meta http-equiv="refresh" content="2; url=/"></meta>`

	destination := c.FormValue("destination")
	if destination == "" {
		destination = "/"
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		c.Html(http.StatusBadRequest, []byte(failBody+`Missing parameter`))
		return
	}

	var users []model.UserData
	table := db.GetFullTableName(db.Info.UserTable)
	dbtype := db.GetDatabaseTypeString()

	column := np.CreateString(model.UserData{}, dbtype, "", false)
	where := np.CreateString(map[string]interface{}{"USERNAME": nil}, dbtype, "", false)

	sql := `SELECT
	` + column.Names + `
	FROM ` + table + `
	WHERE ` + where.Names + `='` + username + `'`

	rows, err := db.Con.Query(sql)
	if err != nil {
		c.Html(http.StatusBadRequest, []byte(failBody+`DB error or User may not exists`))
		return
	}
	defer rows.Close()

	err = scan.Rows(&users, rows)
	if err != nil {
		c.Html(http.StatusBadRequest, []byte(failBody+`User may not exists`))
		return
	}

	if len(users) == 0 {
		c.Html(http.StatusBadRequest, []byte(failBody+`User not exists`))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(users[0].Password.String), []byte(password))
	if err != nil {
		c.Html(http.StatusBadRequest, []byte(failBody+`Check password`))
		return
	}

	authinfo := model.AuthInfo{
		Name:     null.NewString(username, true),
		IpAddr:   null.NewString(c.RemoteAddr, true),
		Platform: null.NewString("", true),
		Duration: null.NewInt(60*60*24*7, true),
		// Duration: null.NewInt(10, true), // 10 seconds test
	}

	// auth.SetupCookieToken(c.ResponseWriter, authinfo)
	auth.SetCookieSession(c, authinfo)

	// c.Html(http.StatusOK, []byte(`<meta http-equiv="refresh" content="0; url=/"></meta>`))
	c.Html(http.StatusOK, []byte(`<meta http-equiv="refresh" content="0; url=`+destination+`"></meta>`))
}

// Logout - Expire cookie
func Logout(c *router.Context) {
	auth.ExpireCookie(c.ResponseWriter)
	auth.DestroyCookieSession(c)

	// authinfo := model.AuthInfo{}
	// if c.AuthInfo != nil {
	// 	authinfo = c.AuthInfo.(model.AuthInfo)
	// }
	// c.Text(http.StatusOK, "Good bye "+authinfo.Name.String)
	c.Html(http.StatusOK, []byte(`<meta http-equiv="refresh" content="0; url=/"></meta>`))
}

// Signup - Create new user
func Signup(c *router.Context) {
	var err error

	now := time.Now().Format("20060102150405")
	columnsCount, _ := crud.GetUserColumnsCount()

	switch columnsCount {
	case model.UserDataFieldCount:
		var userData model.UserData

		err = json.NewDecoder(c.Body).Decode(&userData)
		if err != nil {
			c.Text(http.StatusBadRequest, err.Error())
			return
		}

		password, err := bcrypt.GenerateFromPassword([]byte(userData.Password.String), consts.BcryptCost)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
			return
		}
		userData.Password = null.StringFrom(string(password))

		userData.RegDTTM = null.StringFrom(now)

		userData.Grade = null.StringFrom("pending_user")
		userData.Approval = null.StringFrom("N")

		err = crud.AddUser(userData)
	default:
		userData := make(map[string]interface{})

		err = json.NewDecoder(c.Body).Decode(&userData)
		if err != nil {
			c.Text(http.StatusBadRequest, err.Error())
			return
		}

		password, err := bcrypt.GenerateFromPassword([]byte(userData["password"].(string)), consts.BcryptCost)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
			return
		}
		userData["password"] = string(password)

		userData["reg-dttm"] = now

		userData["grade"] = "pending_user"
		userData["approval"] = "N"

		err = crud.AddUserMap(userData)
	}

	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	result := map[string]string{
		"result": "ok",
	}

	c.Json(http.StatusOK, result)
}

func HandleUserList(c *router.Context) {
	// Use struct with default columns or map with default and user defined columns
	columnsCount, _ := crud.GetUserColumnsCount()
	// columnsCount, _ := db.Obj.GetColumnCount(db.Info.UserTable)

	// queries := c.URL.Query()
	// search := queries.Get("search")

	var err error
	queries := c.URL.Query()

	listingOptions := model.UserListingOptions{}
	listingOptions.Search = null.StringFrom(queries.Get("search"))

	listingOptions.Page = null.IntFrom(1)
	listingOptions.ListCount = null.IntFrom(100)

	if queries.Get("count") != "" {
		countPerPage, err := strconv.Atoi(queries.Get("count"))
		if err != nil {
			c.Text(http.StatusBadRequest, err.Error())
			return
		}

		listingOptions.ListCount = null.IntFrom(int64(countPerPage))
	}

	if queries.Get("page") != "" {
		page := queries.Get("page")
		pageNum, err := strconv.Atoi(page)
		if err != nil {
			c.Text(http.StatusBadRequest, err.Error())
			return
		}

		listingOptions.Page = null.IntFrom(int64(pageNum))
	}

	listingOptions.Page.Int64--

	h, err := LoadHTML(c)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	var listJSON []byte

	switch columnsCount {
	case model.UserDataFieldCount:
		// result, err := crud.GetUsersList(search)
		// if err != nil {
		// 	c.Text(http.StatusInternalServerError, err.Error())
		// 	return
		// }

		// c.Json(http.StatusOK, result)

		list, err := crud.GetUsers(listingOptions)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
		}

		listJSON, _ = json.Marshal(list)
	default:
		// result, err := crud.GetUsersListMap(search)
		// if err != nil {
		// 	c.Text(http.StatusInternalServerError, err.Error())
		// 	return
		// }

		// c.Json(http.StatusOK, result)

		list, err := crud.GetUsersMap(listingOptions)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
		}

		listJSON, _ = json.Marshal(list)
	}

	h = bytes.ReplaceAll(h, []byte("$USER_LIST$"), listJSON)

	c.Html(http.StatusOK, h)
}
