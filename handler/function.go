package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"html/template"

	"9minutes/config"
	"9minutes/crud"
	"9minutes/fd"
	"9minutes/logging"
	"9minutes/model"
	"9minutes/router"
	"9minutes/wsock"

	"gopkg.in/guregu/null.v4"
	// "github.com/goccy/go-json"

	"github.com/microcosm-cc/bluemonday"
)

var (
	patternLinkLogin     = `\$LinkLogin\$(.*)\n`
	patternLinkLogout    = `\$LinkLogout\$(.*)\n`
	patternLinkAdmin     = `\$LinkAdmin\$(.*)\n`
	patternLinkMyPage    = `\$LinkMyPage\$(.*)\n`
	patternYouArePending = `\$YouArePending\$(.*)\n`
	reLogin              = regexp.MustCompile(patternLinkLogin)
	reLogout             = regexp.MustCompile(patternLinkLogout)
	reAdmin              = regexp.MustCompile(patternLinkAdmin)
	reMyPage             = regexp.MustCompile(patternLinkMyPage)
	reYouArePending      = regexp.MustCompile(patternYouArePending)

	patternIncludes = `@INCLUDE@(.*)(\n|$)`
	reIncludes      = regexp.MustCompile(patternIncludes)
)

var bm = bluemonday.UGCPolicy()

// var bm := bluemonday.NewPolicy()
// bm.AllowStandardURLs()
// bm.AllowAttrs("href").OnElements("a")
// bm.AllowElements([]string{"p", "br", "pre", "code"}...)

func Index(c *router.Context) {
	c.URL.Path = "/index.html"
	HandleHTML(c)
}

func AdminIndex(c *router.Context) {
	c.URL.Path = "/admin/index.html"
	HandleHTML(c)
}

func HealthCheck(c *router.Context) {
	c.Text(http.StatusOK, "Ok")
}

func Hello(c *router.Context) {
	switch c.Method {
	case http.MethodGet:
		c.Text(http.StatusOK, "Hello world GET")
	case http.MethodPost:
		c.Text(http.StatusOK, "Hello world POST")
	}
}

func HelloParam(c *router.Context) {
	if len(c.Params) > 0 {
		c.Text(http.StatusOK, "Hello "+c.Params[0])
	} else {
		c.Text(http.StatusBadRequest, "Missing parameter")
	}
}

func GetParam(c *router.Context) {
	result := ""

	params := c.URL.Query()

	for k := range c.URL.Query() {
		result += k + "=" + params.Get(k) + "\n"
	}

	c.Text(http.StatusOK, result)
}

func PostForm(c *router.Context) {
	result := ""

	switch c.Method {
	case http.MethodGet:
		result = "Hello world GET"
	case http.MethodPost:
		c.ParseForm()
		for k := range c.Form {
			result += k + "=" + c.FormValue(k) + "\n"
		}
	}

	c.Text(http.StatusOK, result)
}

func PostJson(c *router.Context) {
	user := model.UserInfo{}

	err := json.NewDecoder(c.Body).Decode(&user)
	if err != nil {
		c.Text(http.StatusBadRequest, err.Error())
		return
	}

	c.Json(http.StatusOK, user)
}

func HandleHTML(c *router.Context) {
	h, err := LoadHTML(c)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	c.Html(http.StatusOK, h)
}

func HandleLogin(c *router.Context) {
	failBody := `<meta http-equiv="refresh" content="2; url=/"></meta>`

	if c.AuthInfo != nil {
		if c.AuthInfo.(model.AuthInfo).Name.Valid {
			userInfo, err := crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
			if err != nil {
				c.Html(http.StatusBadRequest, []byte(failBody+`Server error`))
				return
			}

			if userInfo.Idx.Valid {
				c.Html(http.StatusBadRequest, []byte(failBody+`Already logged in`))
				return
			}
		}
	}

	h, err := LoadHTML(c)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	c.Html(http.StatusOK, h)
}

func HandleSignup(c *router.Context) {
	failBody := `<meta http-equiv="refresh" content="2; url=/"></meta>`

	if c.AuthInfo != nil {
		if c.AuthInfo.(model.AuthInfo).Name.Valid {
			userInfo, err := crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
			if err != nil {
				c.Html(http.StatusBadRequest, []byte(failBody+`Server error`))
				return
			}

			if userInfo.Idx.Valid {
				c.Html(http.StatusBadRequest, []byte(failBody+`Already logged in`))
				return
			}
		}
	}

	h, err := LoadHTML(c)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	c.Html(http.StatusOK, h)
}

func HandleAsset(c *router.Context) {
	h, err := LoadFile(c)

	if err != nil {
		logging.Object.Warn().Err(err).Msg("HandleAsset")
	}

	c.File(http.StatusOK, h)
}

func HandleWebsocketEcho(c *router.Context) {
	wsock.WebSocketEcho(c.ResponseWriter, c.Request)
}

func HandleWebsocketChat(c *router.Context) {
	wsock.WebSocketChat(c.ResponseWriter, c.Request)
}

func HandleGetDir(c *router.Context) {
	path := model.FilePath{}
	result := model.FileList{}

	err := json.NewDecoder(c.Body).Decode(&path)
	if err != nil {
		c.Text(http.StatusBadRequest, err.Error())
		return
	}

	f, err := os.Stat(path.Path.String)
	if err != nil {
		c.Text(http.StatusBadRequest, err.Error())
		return
	}

	dir := path.Path.String
	if f.IsDir() {
		dir = path.Path.String + "/"
	}

	dir = filepath.Dir(dir)
	absPath, err := filepath.Abs(dir)
	if err != nil {
		c.Text(http.StatusBadRequest, err.Error())
		return
	}

	sort := 0
	switch path.Sort.String {
	case "name":
		sort = fd.NAME
	case "size":
		sort = fd.SIZE
	case "time":
		sort = fd.TIME
	default:
		sort = fd.NAME
	}

	order := 0
	switch path.Order.String {
	case "asc":
		order = fd.ASC
	case "desc":
		order = fd.DESC
	default:
		order = fd.ASC
	}

	result.Path = null.StringFrom(dir)
	result.FullPath = null.StringFrom(absPath)

	files, err := fd.Dir(absPath, sort, order)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	for _, file := range files {
		result.Files = append(result.Files, model.FileInfo{
			Name:     null.StringFrom(file.Name()),
			Size:     null.IntFrom(file.Size()),
			DateTime: null.StringFrom(file.ModTime().Format("2006-01-02 15:04:05")),
			IsDir:    null.BoolFrom(file.IsDir()),
		})
	}

	c.Json(http.StatusOK, result)
}

func HandleContentList(c *router.Context) {
	var err error

	code := ""
	board := model.Board{}
	queries := c.URL.Query()

	if queries.Get("code") != "" {
		code = queries.Get("code")

		board.BoardCode = null.StringFrom(code)
		board, err = crud.GetBoardByCode(board)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
		}
	}

	var userInfo model.UserData
	userGrade := 999

	switch c.AuthInfo {
	case nil:
		userGrade = config.UserGrades.IndexOf("guest")
	default:
		userInfo, err = crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
			return
		}

		userGrade = config.UserGrades.IndexOf(userInfo.Grade.String)
	}

	accessGrade := config.UserGrades.IndexOf(board.GrantRead.String)

	if accessGrade < userGrade {
		c.Text(http.StatusForbidden, "Forbidden")
		return
	}

	listingOptions := model.ContentListingOptions{}
	listingOptions.Search = null.StringFrom(queries.Get("search"))

	listingOptions.Page = null.IntFrom(1)
	listingOptions.ListCount = null.IntFrom(int64(config.ContentsCountPerPage))

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

	if !board.Idx.Valid {
		c.Text(http.StatusInternalServerError, "Board not found")
		return
	}

	h = bytes.ReplaceAll(h, []byte("$CODE$"), []byte(code))

	list, err := crud.GetContentList(board, listingOptions)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
	}

	listJSON, _ := json.Marshal(list)
	h = bytes.ReplaceAll(h, []byte("$CONTENT_LIST$"), listJSON)

	c.Html(http.StatusOK, h)
}

func HandleContentListTmpl(c *router.Context) {
	var err error

	code := ""
	board := model.Board{}
	queries := c.URL.Query()

	if queries.Get("code") != "" {
		code = queries.Get("code")

		board.BoardCode = null.StringFrom(code)
		board, err = crud.GetBoardByCode(board)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
		}
	}

	var userInfo model.UserData
	userGrade := 999

	switch c.AuthInfo {
	case nil:
		userGrade = config.UserGrades.IndexOf("guest")
	default:
		userInfo, err = crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
			return
		}

		userGrade = config.UserGrades.IndexOf(userInfo.Grade.String)
	}

	accessGrade := config.UserGrades.IndexOf(board.GrantRead.String)

	if accessGrade < userGrade {
		c.Text(http.StatusForbidden, "Forbidden")
		return
	}

	listingOptions := model.ContentListingOptions{}
	listingOptions.Search = null.StringFrom(queries.Get("search"))

	listingOptions.Page = null.IntFrom(1)
	listingOptions.ListCount = null.IntFrom(int64(config.ContentsCountPerPage))

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

	if !board.Idx.Valid {
		c.Text(http.StatusInternalServerError, "Board not found")
		return
	}

	h = bytes.ReplaceAll(h, []byte("$CODE$"), []byte(code))

	list, err := crud.GetContentList(board, listingOptions)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
	}

	for i := list.CurrentPage - 1; i < (list.CurrentPage + 2); i++ {
		if i < 1 || i > list.TotalPage {
			continue
		}
		if i > list.TotalPage {
			break
		}

		list.PageList = append(list.PageList, i)
	}

	listJSON, _ := json.Marshal(list)
	h = bytes.ReplaceAll(h, []byte("$CONTENT_LIST$"), listJSON)

	tmpl := template.New("list")
	tmpl = tmpl.Funcs(
		template.FuncMap{
			"jump_to_before": func(page int) string {
				jumpPage := page - 5
				if jumpPage < 1 {
					jumpPage = 1
				}
				result := fmt.Sprint(jumpPage)
				return result
			},
			"jump_to_after": func(page int) string {
				jumpPage := page + 5
				if jumpPage > list.TotalPage {
					jumpPage = list.TotalPage
				}
				result := fmt.Sprint(jumpPage)
				return result
			},
		},
	)

	tmpl, err = tmpl.Parse(string(h))
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	err = tmpl.Execute(c.ResponseWriter, list)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}
	c.ResponseWriter.WriteHeader(http.StatusOK)
}

func HandleReadContent(c *router.Context) {
	h, err := LoadHTML(c)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	queries := c.URL.Query()
	if queries.Get("code") != "" {
		code := queries.Get("code")
		h = bytes.ReplaceAll(h, []byte("$CODE$"), []byte(code))

		board := model.Board{}
		board.BoardCode = null.StringFrom(queries.Get("code"))
		board, err = crud.GetBoardByCode(board)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
		}

		var userInfo model.UserData
		userGrade := 999

		switch c.AuthInfo {
		case nil:
			userGrade = config.UserGrades.IndexOf("guest")
		default:
			userInfo, err = crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
			if err != nil {
				c.Text(http.StatusInternalServerError, err.Error())
				return
			}

			userGrade = config.UserGrades.IndexOf(userInfo.Grade.String)
		}

		accessGrade := config.UserGrades.IndexOf(board.GrantRead.String)

		if accessGrade < userGrade {
			c.Text(http.StatusForbidden, "Forbidden")
			return
		}

		idx := queries.Get("idx")
		content, err := crud.GetContent(board, idx)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
		}

		content.Views = null.IntFrom(content.Views.Int64 + 1)
		err = crud.UpdateContent(board, content, "viewcount")
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
		}

		listingOptions := model.CommentListingOptions{}
		listingOptions.Search = null.StringFrom(queries.Get("search"))

		listingOptions.Page = null.IntFrom(0)
		listingOptions.ListCount = null.IntFrom(int64(config.CommentCountPerPage))

		comments, err := crud.GetComments(board, content, listingOptions)
		commentsJSON, _ := json.Marshal(comments)

		h = bytes.ReplaceAll(h, []byte("$IDX$"), []byte(idx))
		h = bytes.ReplaceAll(h, []byte("$TITLE$"), []byte(content.Title.String))
		h = bytes.ReplaceAll(h, []byte("$AUTHOR_IDX$"), []byte(fmt.Sprint(content.AuthorIdx.Int64)))
		h = bytes.ReplaceAll(h, []byte("$AUTHOR_NAME$"), []byte(fmt.Sprint(content.AuthorName.String)))
		h = bytes.ReplaceAll(h, []byte("$VIEWS$"), []byte(fmt.Sprint(content.Views.Int64)))
		h = bytes.ReplaceAll(h, []byte("$CONTENT$"), []byte(content.Content.String))
		h = bytes.ReplaceAll(h, []byte("$FILE_LIST$"), []byte(content.Files.String))

		h = bytes.ReplaceAll(h, []byte("$COMMENTS$"), commentsJSON)
	}

	c.Html(http.StatusOK, h)
}

// HandleWriteContent - Write content page
func HandleWriteContent(c *router.Context) {
	h, err := LoadHTML(c)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	queries := c.URL.Query()
	if queries.Get("code") != "" {
		code := queries.Get("code")
		h = bytes.ReplaceAll(h, []byte("$CODE$"), []byte(code))

		board := model.Board{}
		board.BoardCode = null.StringFrom(queries.Get("code"))
		board, err = crud.GetBoardByCode(board)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
		}

		var userInfo model.UserData
		userGrade := 999

		switch c.AuthInfo {
		case nil:
			userGrade = config.UserGrades.IndexOf("guest")
		default:
			userInfo, err = crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
			if err != nil {
				c.Text(http.StatusInternalServerError, err.Error())
				return
			}

			userGrade = config.UserGrades.IndexOf(userInfo.Grade.String)
		}

		accessGrade := config.UserGrades.IndexOf(board.GrantWrite.String)

		if accessGrade < userGrade {
			c.Text(http.StatusForbidden, "Forbidden")
			return
		}
	}

	c.Html(http.StatusOK, h)
}

func HandleEditContent(c *router.Context) {
	h, err := LoadHTML(c)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	queries := c.URL.Query()
	if queries.Get("code") != "" {

		code := queries.Get("code")
		h = bytes.ReplaceAll(h, []byte("$CODE$"), []byte(code))

		board := model.Board{}
		board.BoardCode = null.StringFrom(queries.Get("code"))
		board, err = crud.GetBoardByCode(board)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
		}

		idx := queries.Get("idx")
		content, err := crud.GetContent(board, idx)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
		}

		userInfo, err := crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
		if err != nil {
			c.Html(http.StatusOK, []byte(`<meta http-equiv="refresh" content="0; url=/"></meta>`))
			return
		}

		if content.AuthorIdx.Int64 != userInfo.Idx.Int64 && userInfo.Grade.String != "admin" && userInfo.Grade.String != "manager" {
			c.Html(http.StatusOK, []byte(`<meta http-equiv="refresh" content="0; url=/"></meta>`))
			return
		}

		h = bytes.ReplaceAll(h, []byte("$IDX$"), []byte(idx))
		h = bytes.ReplaceAll(h, []byte("$TITLE$"), []byte(content.Title.String))
		h = bytes.ReplaceAll(h, []byte("$TITLE_IMAGE$"), []byte(content.TitleImage.String))
		h = bytes.ReplaceAll(h, []byte("$AUTHOR_IDX$"), []byte(fmt.Sprint(content.AuthorIdx.Int64)))
		h = bytes.ReplaceAll(h, []byte("$VIEWS$"), []byte(fmt.Sprint(content.Views.Int64)))
		h = bytes.ReplaceAll(h, []byte("$CONTENT$"), []byte(content.Content.String))
		h = bytes.ReplaceAll(h, []byte("$FILE_LIST$"), []byte(content.Files.String))
	}

	c.Html(http.StatusOK, h)
}

// WriteContent - Write content API
func WriteContent(c *router.Context) {
	board := model.Board{}
	content := model.Content{}

	uri := strings.Split(c.URL.Path, "/")
	code := null.StringFrom(uri[len(uri)-1])
	board.BoardCode = code

	err := json.NewDecoder(c.Body).Decode(&content)
	if err != nil {
		c.Text(http.StatusBadRequest, err.Error())
		return
	}

	board, err = crud.GetBoardByCode(board)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	userInfo, err := crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}
	content.AuthorIdx = userInfo.Idx

	now := time.Now().Format("20060102150405")
	content.RegDTTM = null.StringFrom(now)

	content.Views = null.IntFrom(0)

	err = crud.WriteContent(board, content)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	result := map[string]interface{}{
		"result": "success",
	}

	c.Json(http.StatusOK, result)
}

func UpdateContent(c *router.Context) {
	var board model.Board
	var content model.Content

	uri := strings.Split(c.URL.Path, "/")

	code := uri[len(uri)-2]
	board.BoardCode = null.StringFrom(code)

	board, err := crud.GetBoardByCode(board)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	userInfo, err := crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}
	content.AuthorIdx = userInfo.Idx

	idx, _ := strconv.Atoi(uri[len(uri)-1])
	content.Idx = null.IntFrom(int64(idx))

	err = json.NewDecoder(c.Body).Decode(&content)
	if err != nil {
		c.Text(http.StatusBadRequest, err.Error())
		return
	}

	err = crud.UpdateContent(board, content, "update")
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	result := map[string]interface{}{
		"result": "success",
	}

	c.Json(http.StatusOK, result)
}

func DeleteContent(c *router.Context) {
	var board model.Board

	uri := strings.Split(c.URL.Path, "/")

	code := uri[len(uri)-2]
	board.BoardCode = null.StringFrom(code)

	idx, _ := strconv.Atoi(uri[len(uri)-1])

	board, err := crud.GetBoardByCode(board)
	if err != nil {
		c.Text(http.StatusNotFound, "Board was not found")
	}

	content, err := crud.GetContent(board, fmt.Sprint(idx))
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	if content.AuthorIdx.Int64 != userInfo.Idx.Int64 && userInfo.Grade.String != "admin" && userInfo.Grade.String != "manager" {
		c.Text(http.StatusForbidden, "forbidden")
		return
	}

	err = crud.DeleteContent(board, fmt.Sprint(idx))
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	err = crud.DeleteComments(board, fmt.Sprint(idx))
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	result := map[string]interface{}{
		"result": "success",
	}

	c.Json(http.StatusOK, result)
}

func GetComments(c *router.Context) {
	var err error
	var board model.Board

	uri := strings.Split(c.URL.Path, "/")

	code := uri[len(uri)-2]
	board.BoardCode = null.StringFrom(code)
	queries := c.URL.Query()

	board, err = crud.GetBoardByCode(board)
	if err != nil {
		c.Text(http.StatusNotFound, "Board was not found")
	}

	var userInfo model.UserData
	userGrade := 999

	switch c.AuthInfo {
	case nil:
		userGrade = config.UserGrades.IndexOf("guest")
	default:
		userInfo, err = crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
			return
		}

		userGrade = config.UserGrades.IndexOf(userInfo.Grade.String)
	}

	accessGrade := config.UserGrades.IndexOf(board.GrantRead.String)

	if accessGrade < userGrade {
		log.Println(c.AuthInfo, accessGrade, userGrade)
		c.Text(http.StatusForbidden, "Forbidden")
		return
	}

	idx, _ := strconv.Atoi(uri[len(uri)-1])
	content, err := crud.GetContent(board, fmt.Sprint(idx))
	if err != nil {
		c.Text(http.StatusNotFound, "Content was not found")
	}

	listingOptions := model.CommentListingOptions{}
	listingOptions.Search = null.StringFrom(queries.Get("search"))

	listingOptions.Page = null.IntFrom(1)
	listingOptions.ListCount = null.IntFrom(int64(config.CommentCountPerPage))

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

	commentList, err := crud.GetComments(board, content, listingOptions)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
	}

	c.Json(http.StatusOK, commentList)
}

func WriteComment(c *router.Context) {
	board := model.Board{}
	comment := model.Comment{}

	uri := strings.Split(c.URL.Path, "/")
	boardCode := uri[len(uri)-2]
	board.BoardCode = null.StringFrom(boardCode)

	board, err := crud.GetBoardByCode(board)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	var userInfo model.UserData
	userGrade := 999

	switch c.AuthInfo {
	case nil:
		userGrade = config.UserGrades.IndexOf("guest")
	default:
		userInfo, err = crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
		if err != nil {
			c.Text(http.StatusInternalServerError, err.Error())
			return
		}

		userGrade = config.UserGrades.IndexOf(userInfo.Grade.String)
	}

	accessGrade := config.UserGrades.IndexOf(board.GrantComment.String)

	if accessGrade < userGrade {
		c.Text(http.StatusForbidden, "Forbidden")
		return
	}

	contentIdx, _ := strconv.Atoi(uri[len(uri)-1])

	err = json.NewDecoder(c.Body).Decode(&comment)
	if err != nil {
		c.Text(http.StatusBadRequest, err.Error())
		return
	}

	comment.Content = null.StringFrom(bm.Sanitize(comment.Content.String))
	if comment.Content.String == "" {
		c.Text(http.StatusBadRequest, "Comment is empty")
		return
	}

	comment.BoardIdx = null.IntFrom(int64(contentIdx))

	userInfo, err = crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}
	comment.AuthorIdx = userInfo.Idx

	now := time.Now().Format("20060102150405")
	comment.RegDTTM = null.StringFrom(now)

	err = crud.WriteComment(board, comment)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	result := map[string]interface{}{
		"result": "success",
	}

	c.Json(http.StatusOK, result)
}

func DeleteComment(c *router.Context) {
	var board model.Board

	uri := strings.Split(c.URL.Path, "/")

	code := uri[len(uri)-3]
	board.BoardCode = null.StringFrom(code)

	boardIdx, _ := strconv.Atoi(uri[len(uri)-2])
	commentIdx, _ := strconv.Atoi(uri[len(uri)-1])

	board, err := crud.GetBoardByCode(board)
	if err != nil {
		c.Text(http.StatusNotFound, "Board was not found")
	}

	comment, err := crud.GetComment(board, fmt.Sprint(boardIdx), fmt.Sprint(commentIdx))
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	if comment.AuthorIdx.Int64 != userInfo.Idx.Int64 && userInfo.Grade.String != "admin" && userInfo.Grade.String != "manager" {
		c.Text(http.StatusForbidden, "forbidden")
		return
	}

	err = crud.DeleteComment(board, fmt.Sprint(boardIdx), fmt.Sprint(commentIdx))
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	result := map[string]interface{}{
		"result": "success",
	}

	c.Json(http.StatusOK, result)
}
