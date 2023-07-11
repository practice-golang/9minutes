package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"9minutes/config"
	"9minutes/crud"
	"9minutes/model"

	"gopkg.in/guregu/null.v4"
	// "github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
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

func HealthCheck(c *fiber.Ctx) error {
	return c.SendString("Ok")
}

func HelloParam(c *fiber.Ctx) error {
	if len(c.Params("name")) > 0 {
		return c.Status(http.StatusOK).SendString("Hello " + c.Params("name"))
	} else {
		return c.Status(http.StatusBadRequest).SendString("Missing parameter")
	}
}

// HandleHTML - Handle HTML template layout
func HandleHTML(c *fiber.Ctx) error {
	name := strings.TrimSuffix(c.Path()[1:], "/")
	params := c.Queries()
	templateMap := fiber.Map{}

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	userid := ""
	useridInterface := sess.Get("userid")
	if useridInterface != nil {
		userid = useridInterface.(string)
	}
	grade := ""
	gradeInterface := sess.Get("grade")
	if gradeInterface != nil {
		grade = gradeInterface.(string)
	}

	templateMap["Title"] = "9minutes"
	templateMap["UserId"] = userid
	templateMap["Grade"] = grade

	switch true {
	case name == "":
		name = "index"

		if params["hello"] != "" {
			log.Printf("Hello: %s", params["hello"])
		}
	case strings.HasPrefix(name, "board"):

	case strings.HasPrefix(name, "admin"):
		if userid == "" {
			name = "status/unauthorized"
			break
		}

		switch name {
		case "admin":
			name = "admin/index"
		case "admin/board-list":
			log.Println(name)
		case "admin/user-columns":
			log.Println(name)
		case "admin/user-list":
			log.Println(name)
		}

		name = "admin/index"
	}

	err = c.Render(name, templateMap)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return c.Status(http.StatusNotFound).SendString("Page not Found")
		}
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return nil
}

// // func HandleLogin(c *router.Context) {
// // 	failBody := `<meta http-equiv="refresh" content="2; url=/"></meta>`

// // 	if c.AuthInfo != nil {
// // 		if c.AuthInfo.(model.AuthInfo).Name.Valid {
// // 			userInfo, err := crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
// // 			if err != nil {
// // 				c.Html(http.StatusBadRequest, []byte(failBody+`Server error`))
// // 				return
// // 			}

// // 			if userInfo.Idx.Valid {
// // 				c.Html(http.StatusBadRequest, []byte(failBody+`Already logged in`))
// // 				return
// // 			}
// // 		}
// // 	}

// // 	h, err := LoadHTML(c)
// // 	if err != nil {
// // 		c.Text(http.StatusInternalServerError, err.Error())
// // 		return
// // 	}

// // 	c.Html(http.StatusOK, h)
// // }

// // func HandleSignup(c *router.Context) {
// // 	failBody := `<meta http-equiv="refresh" content="2; url=/"></meta>`

// // 	if c.AuthInfo != nil {
// // 		if c.AuthInfo.(model.AuthInfo).Name.Valid {
// // 			userInfo, err := crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
// // 			if err != nil {
// // 				c.Html(http.StatusBadRequest, []byte(failBody+`Server error`))
// // 				return
// // 			}

// // 			if userInfo.Idx.Valid {
// // 				c.Html(http.StatusBadRequest, []byte(failBody+`Already logged in`))
// // 				return
// // 			}
// // 		}
// // 	}

// // 	h, err := LoadHTML(c)
// // 	if err != nil {
// // 		c.Text(http.StatusInternalServerError, err.Error())
// // 		return
// // 	}

// // 	// Captcha
// // 	captchaID := captcha.New()
// // 	h = bytes.ReplaceAll(h, []byte("$CAPTCHAID$"), []byte(captchaID))

// // 	c.Html(http.StatusOK, h)
// // }

// func HandleAsset(c *router.Context) {
// 	h, err := LoadFile(c)

// 	if err != nil {
// 		logging.Object.Warn().Err(err).Msg("HandleAsset")
// 	}

// 	c.File(http.StatusOK, h)
// }

// func HandleWebsocketEcho(c *router.Context) {
// 	wsock.WebSocketEcho(c.ResponseWriter, c.Request)
// }

// func HandleWebsocketChat(c *router.Context) {
// 	wsock.WebSocketChat(c.ResponseWriter, c.Request)
// }

// func HandleGetDir(c *router.Context) {
// 	path := model.FilePath{}
// 	result := model.FileList{}

// 	err := json.NewDecoder(c.Body).Decode(&path)
// 	if err != nil {
// 		c.Text(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	f, err := os.Stat(path.Path.String)
// 	if err != nil {
// 		c.Text(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	dir := path.Path.String
// 	if f.IsDir() {
// 		dir = path.Path.String + "/"
// 	}

// 	dir = filepath.Dir(dir)
// 	absPath, err := filepath.Abs(dir)
// 	if err != nil {
// 		c.Text(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	sort := 0
// 	switch path.Sort.String {
// 	case "name":
// 		sort = fd.NAME
// 	case "size":
// 		sort = fd.SIZE
// 	case "time":
// 		sort = fd.TIME
// 	default:
// 		sort = fd.NAME
// 	}

// 	order := 0
// 	switch path.Order.String {
// 	case "asc":
// 		order = fd.ASC
// 	case "desc":
// 		order = fd.DESC
// 	default:
// 		order = fd.ASC
// 	}

// 	result.Path = null.StringFrom(dir)
// 	result.FullPath = null.StringFrom(absPath)

// 	files, err := fd.Dir(absPath, sort, order)
// 	if err != nil {
// 		c.Text(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	for _, file := range files {
// 		f, _ := file.Info()
// 		result.Files = append(result.Files, model.FileInfo{
// 			Name:     null.StringFrom(file.Name()),
// 			Size:     null.IntFrom(f.Size()),
// 			DateTime: null.StringFrom(f.ModTime().Format("2006-01-02 15:04:05")),
// 			IsDir:    null.BoolFrom(file.IsDir()),
// 		})
// 	}

// 	c.Json(http.StatusOK, result)
// }

func HandleContentList(c *fiber.Ctx) error {
	var err error

	code := ""
	search := ""
	searchParam := ""
	board := model.Board{}
	queries := c.Queries()

	if queries["code"] != "" {
		code = queries["code"]

		board.BoardCode = null.StringFrom(code)
		board, err = crud.GetBoardByCode(board)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
	}

	// var userInfo model.UserData
	// userGrade := 999

	// switch c.AuthInfo {
	// case nil:
	// 	userGrade = config.UserGrades.IndexOf("guest")
	// default:
	// 	userInfo, err = crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
	// 	if err != nil {
	// 		c.Text(http.StatusInternalServerError, err.Error())
	// 		return
	// 	}

	// 	userGrade = config.UserGrades.IndexOf(userInfo.Grade.String)
	// }

	// accessGrade := config.UserGrades.IndexOf(board.GrantRead.String)

	// if accessGrade < userGrade {
	// 	c.Text(http.StatusForbidden, "Forbidden")
	// 	return
	// }

	if queries["search"] != "" {
		search = queries["search"]
		searchParam = "&search=" + search
	}

	listingOptions := model.ContentListingOptions{}
	listingOptions.Search = null.StringFrom(search)

	listingOptions.Page = null.IntFrom(1)
	listingOptions.ListCount = null.IntFrom(int64(config.ContentsCountPerPage))

	if queries["count"] != "" {
		countPerPage, err := strconv.Atoi(queries["count"])
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}

		listingOptions.ListCount = null.IntFrom(int64(countPerPage))
	}

	if queries["page"] != "" {
		page := queries["page"]
		pageNum, err := strconv.Atoi(page)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}

		listingOptions.Page = null.IntFrom(int64(pageNum))
	}

	listingOptions.Page.Int64--

	switch board.BoardType.String {
	case "board":
		c.Context().URI().SetPath("/board/list.html")
	case "gallery":
		c.Context().URI().SetPath("/gallery/list.html")
	}

	h, err := LoadHTML(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if !board.Idx.Valid {
		return c.Status(http.StatusNotFound).SendString("Board not found")
	}

	h = bytes.ReplaceAll(h, []byte("$CODE$"), []byte(code))
	h = bytes.ReplaceAll(h, []byte("$SEARCH$"), []byte(searchParam))

	list, err := crud.GetContentList(board, listingOptions)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
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
			"is_last_index": func(idx int) bool {
				return idx == len(list.PageList)-1
			},
		},
	)

	tmpl, err = tmpl.Parse(string(h))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// err = tmpl.Execute(c.ResponseWriter, list)
	// if err != nil {
	// 	c.Text(http.StatusInternalServerError, err.Error())
	// 	return err
	// }
	// c.ResponseWriter.WriteHeader(http.StatusOK)
	// // c.Html(http.StatusOK, h)

	return c.Render("list", fiber.Map{})
}

func HandleReadContent(c *fiber.Ctx) error {
	h, err := LoadHTML(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	queries := c.Queries()
	if queries["code"] != "" {
		code := queries["code"]
		h = bytes.ReplaceAll(h, []byte("$CODE$"), []byte(code))

		board := model.Board{}
		board.BoardCode = null.StringFrom(queries["code"])
		board, err = crud.GetBoardByCode(board)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		// var userInfo model.UserData
		// userGrade := 999

		// switch c.AuthInfo {
		// case nil:
		// 	userGrade = config.UserGrades.IndexOf("guest")
		// default:
		// 	userInfo, err = crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
		// 	if err != nil {
		// 		c.Text(http.StatusInternalServerError, err.Error())
		// 		return
		// 	}

		// 	userGrade = config.UserGrades.IndexOf(userInfo.Grade.String)
		// }

		// accessGrade := config.UserGrades.IndexOf(board.GrantRead.String)

		// if accessGrade < userGrade {
		// 	c.Text(http.StatusForbidden, "Forbidden")
		// 	return
		// }

		idx := queries["idx"]
		content, err := crud.GetContent(board, idx)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		content.Views = null.IntFrom(content.Views.Int64 + 1)
		err = crud.UpdateContent(board, content, "viewcount")
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		listingOptions := model.CommentListingOptions{}
		listingOptions.Search = null.StringFrom(queries["search"])

		listingOptions.Page = null.IntFrom(-1)
		listingOptions.ListCount = null.IntFrom(int64(config.CommentCountPerPage))

		comments, err := crud.GetComments(board, content, listingOptions)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
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

	// c.Html(http.StatusOK, h)

	return c.Render("read", fiber.Map{})
}

// HandleWriteContent - Write content page
func HandleWriteContent(c *fiber.Ctx) error {
	h, err := LoadHTML(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	queries := c.Queries()
	if queries["code"] != "" {
		code := queries["code"]
		h = bytes.ReplaceAll(h, []byte("$CODE$"), []byte(code))

		board := model.Board{}
		board.BoardCode = null.StringFrom(queries["code"])
		board, err = crud.GetBoardByCode(board)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		// var userInfo model.UserData
		// userGrade := 999

		// switch c.AuthInfo {
		// case nil:
		// 	userGrade = config.UserGrades.IndexOf("guest")
		// default:
		// 	userInfo, err = crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
		// 	if err != nil {
		// 		c.Text(http.StatusInternalServerError, err.Error())
		// 		return
		// 	}

		// 	userGrade = config.UserGrades.IndexOf(userInfo.Grade.String)
		// }

		// accessGrade := config.UserGrades.IndexOf(board.GrantWrite.String)

		// if accessGrade < userGrade {
		// 	c.Text(http.StatusForbidden, "Forbidden")
		// 	return
		// }
	}

	// c.Html(http.StatusOK, h)
	return c.Render("write", fiber.Map{})
}

func HandleEditContent(c *fiber.Ctx) error {
	h, err := LoadHTML(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	queries := c.Queries()
	if queries["code"] != "" {

		code := queries["code"]
		h = bytes.ReplaceAll(h, []byte("$CODE$"), []byte(code))

		board := model.Board{}
		board.BoardCode = null.StringFrom(queries["code"])
		board, err = crud.GetBoardByCode(board)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		// idx := queries["idx"]
		// content, err := crud.GetContent(board, idx)
		// if err != nil {
		// 	return c.Status(http.StatusInternalServerError).SendString(err.Error())
		// }

		// userInfo, err := crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
		// if err != nil {
		// 	c.Html(http.StatusOK, []byte(`<meta http-equiv="refresh" content="0; url=/"></meta>`))
		// 	return
		// }

		// if content.AuthorIdx.Int64 != userInfo.Idx.Int64 && userInfo.Grade.String != "admin" && userInfo.Grade.String != "manager" {
		// 	c.Html(http.StatusOK, []byte(`<meta http-equiv="refresh" content="0; url=/"></meta>`))
		// 	return
		// }

		// h = bytes.ReplaceAll(h, []byte("$IDX$"), []byte(idx))
		// h = bytes.ReplaceAll(h, []byte("$TITLE$"), []byte(content.Title.String))
		// h = bytes.ReplaceAll(h, []byte("$TITLE_IMAGE$"), []byte(content.TitleImage.String))
		// h = bytes.ReplaceAll(h, []byte("$AUTHOR_IDX$"), []byte(fmt.Sprint(content.AuthorIdx.Int64)))
		// h = bytes.ReplaceAll(h, []byte("$VIEWS$"), []byte(fmt.Sprint(content.Views.Int64)))
		// h = bytes.ReplaceAll(h, []byte("$CONTENT$"), []byte(content.Content.String))
		// h = bytes.ReplaceAll(h, []byte("$FILE_LIST$"), []byte(content.Files.String))
	}

	// c.Html(http.StatusOK, h)
	return c.Render("edit", fiber.Map{})
}

// WriteContent - Write content API
func WriteContent(c *fiber.Ctx) error {
	var err error

	board := model.Board{}
	content := model.Content{}

	uri := strings.Split(c.Context().URI().String(), "/")
	code := null.StringFrom(uri[len(uri)-1])
	board.BoardCode = code

	c.BodyParser(&content)

	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// userInfo, err := crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
	// if err != nil {
	// 	c.Text(http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// content.AuthorIdx = userInfo.Idx

	now := time.Now().Format("20060102150405")
	content.RegDate = null.StringFrom(now)

	content.Views = null.IntFrom(0)

	r, err := crud.WriteContent(board, content)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	boardIDX := board.Idx.Int64
	postIDX, _ := r.LastInsertId()

	files := strings.Split(content.Files.String, "?")

	for _, f := range files {
		if f == "" {
			continue
		}
		files := strings.Split(f, "/")

		filename := files[0]
		storename := files[1]

		crud.UpdateUploadedFile(boardIDX, postIDX, filename, storename)
	}

	result := map[string]interface{}{
		"result": "success",
	}

	// c.Json(http.StatusOK, result)
	return c.Status(http.StatusOK).JSON(result)
}

func UpdateContent(c *fiber.Ctx) error {
	var board model.Board
	var content model.Content
	var deleteList model.FilesToDelete

	uri := strings.Split(c.Context().URI().String(), "/")

	code := uri[len(uri)-2]
	board.BoardCode = null.StringFrom(code)

	board, err := crud.GetBoardByCode(board)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Board was not found")
	}

	// userInfo, err := crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
	// if err != nil {
	// 	c.Text(http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// content.AuthorIdx = userInfo.Idx

	idx, _ := strconv.Atoi(uri[len(uri)-1])
	content.Idx = null.IntFrom(int64(idx))

	rbody := c.Body()

	err = json.Unmarshal(rbody, &content)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	err = crud.UpdateContent(board, content, "update")
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	err = json.Unmarshal(rbody, &deleteList)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	files := strings.Split(content.Files.String, "?")
	for _, f := range files {
		if f == "" {
			continue
		}
		files := strings.Split(f, "/")

		filename := files[0]
		storename := files[1]

		crud.UpdateUploadedFile(board.Idx.Int64, content.Idx.Int64, filename, storename)
	}
	for _, f := range deleteList.DeleteFiles {
		crud.UpdateUploadedFile(board.Idx.Int64, content.Idx.Int64, f.FileName.String, f.StoreName.String)
	}

	// for _, f := range deleteList.DeleteFiles {
	// 	filepath := router.UploadPath + "/" + f.StoreName.String
	// 	err = crud.DeleteUploadedFile(board.Idx.Int64, content.Idx.Int64, f.FileName.String, f.StoreName.String)
	// 	if err != nil {
	// 		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	// 	}
	// 	DeleteUploadFile(filepath)
	// }

	result := map[string]interface{}{
		"result": "success",
	}

	return c.Status(http.StatusOK).JSON(result)
}

func DeleteContent(c *fiber.Ctx) error {
	var board model.Board

	uri := strings.Split(c.Context().URI().String(), "/")

	code := uri[len(uri)-2]
	board.BoardCode = null.StringFrom(code)

	idx, _ := strconv.Atoi(uri[len(uri)-1])

	board, err := crud.GetBoardByCode(board)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Board was not found")
	}

	content, err := crud.GetContent(board, fmt.Sprint(idx))
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Content was not found")
	}

	// userInfo, err := crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
	// if err != nil {
	// 	c.Text(http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// if content.AuthorIdx.Int64 != userInfo.Idx.Int64 && userInfo.Grade.String != "admin" && userInfo.Grade.String != "manager" {
	// 	c.Text(http.StatusForbidden, "forbidden")
	// 	return
	// }

	deleteFiles := strings.Split(content.Files.String, "?")
	deleteList := model.FilesToDelete{}
	for _, df := range deleteFiles {
		if !strings.Contains(df, "/") {
			continue
		}

		deleteFile := model.File{}

		dfs := strings.Split(df, "/")
		deleteFile.FileName = null.StringFrom(dfs[0])
		deleteFile.StoreName = null.StringFrom(dfs[1])

		deleteList.DeleteFiles = append(deleteList.DeleteFiles, deleteFile)
	}

	// for _, f := range deleteList.DeleteFiles {
	// 	filepath := router.UploadPath + "/" + f.StoreName.String
	// 	err = crud.DeleteUploadedFile(board.Idx.Int64, content.Idx.Int64, f.FileName.String, f.StoreName.String)
	// 	if err != nil {
	// 		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	// 	}
	// 	DeleteUploadFile(filepath)
	// }

	err = crud.DeleteContent(board, fmt.Sprint(idx))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	err = crud.DeleteComments(board, fmt.Sprint(idx))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]interface{}{
		"result": "success",
	}

	return c.Status(http.StatusOK).JSON(result)
}

func GetComments(c *fiber.Ctx) error {
	var err error
	var board model.Board

	uri := strings.Split(c.Context().URI().String(), "/")

	code := uri[len(uri)-2]
	board.BoardCode = null.StringFrom(code)
	queries := c.Queries()

	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Board was not found")
	}

	// var userInfo model.UserData
	// userGrade := 999

	// switch c.AuthInfo {
	// case nil:
	// 	userGrade = config.UserGrades.IndexOf("guest")
	// default:
	// 	userInfo, err = crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
	// 	if err != nil {
	// 		c.Text(http.StatusInternalServerError, err.Error())
	// 		return err
	// 	}

	// 	userGrade = config.UserGrades.IndexOf(userInfo.Grade.String)
	// }

	// accessGrade := config.UserGrades.IndexOf(board.GrantRead.String)

	// if accessGrade < userGrade {
	// 	c.Text(http.StatusForbidden, "Forbidden")
	// 	return err
	// }

	idx, _ := strconv.Atoi(uri[len(uri)-1])
	content, err := crud.GetContent(board, fmt.Sprint(idx))
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Content was not found")
	}

	listingOptions := model.CommentListingOptions{}
	listingOptions.Search = null.StringFrom(queries["search"])

	listingOptions.Page = null.IntFrom(1)
	listingOptions.ListCount = null.IntFrom(int64(config.CommentCountPerPage))

	if queries["count"] != "" {
		countPerPage, err := strconv.Atoi(queries["count"])
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}

		listingOptions.ListCount = null.IntFrom(int64(countPerPage))
	}

	if queries["page"] != "" {
		page := queries["page"]
		pageNum, err := strconv.Atoi(page)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}

		listingOptions.Page = null.IntFrom(int64(pageNum))
	}

	listingOptions.Page.Int64--

	commentList, err := crud.GetComments(board, content, listingOptions)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(commentList)
}

func WriteComment(c *fiber.Ctx) error {
	board := model.Board{}
	comment := model.Comment{}

	uri := strings.Split(c.Context().URI().String(), "/")
	boardCode := uri[len(uri)-2]
	board.BoardCode = null.StringFrom(boardCode)

	board, err := crud.GetBoardByCode(board)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Board was not found")
	}

	// var userInfo model.UserData
	// userGrade := 999

	// switch c.AuthInfo {
	// case nil:
	// 	userGrade = config.UserGrades.IndexOf("guest")
	// default:
	// 	userInfo, err = crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
	// 	if err != nil {
	// 		c.Text(http.StatusInternalServerError, err.Error())
	// 		return
	// 	}

	// 	userGrade = config.UserGrades.IndexOf(userInfo.Grade.String)
	// }

	// accessGrade := config.UserGrades.IndexOf(board.GrantComment.String)

	// if accessGrade < userGrade {
	// 	c.Text(http.StatusForbidden, "Forbidden")
	// 	return
	// }

	contentIdx, _ := strconv.Atoi(uri[len(uri)-1])

	err = c.BodyParser(comment)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	comment.Content = null.StringFrom(bm.Sanitize(comment.Content.String))
	if comment.Content.String == "" {
		return c.Status(http.StatusBadRequest).SendString("comment is empty")
	}

	comment.BoardIdx = null.IntFrom(int64(contentIdx))

	// userInfo, err = crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
	// if err != nil {
	// 	c.Text(http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// comment.AuthorIdx = userInfo.Idx

	now := time.Now().Format("20060102150405")
	comment.RegDate = null.StringFrom(now)

	err = crud.WriteComment(board, comment)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]interface{}{
		"result": "success",
	}

	return c.Status(http.StatusOK).JSON(result)
}

func DeleteComment(c *fiber.Ctx) error {
	var board model.Board

	uri := strings.Split(c.Context().URI().String(), "/")

	code := uri[len(uri)-3]
	board.BoardCode = null.StringFrom(code)

	boardIdx, _ := strconv.Atoi(uri[len(uri)-2])
	commentIdx, _ := strconv.Atoi(uri[len(uri)-1])

	board, err := crud.GetBoardByCode(board)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Board was not found")
	}

	// var userInfo model.UserData
	// userGrade := 999

	// switch c.AuthInfo {
	// case nil:
	// 	userGrade = config.UserGrades.IndexOf("guest")
	// default:
	// 	userInfo, err = crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
	// 	if err != nil {
	// 		c.Text(http.StatusInternalServerError, err.Error())
	// 		return
	// 	}

	// 	userGrade = config.UserGrades.IndexOf(userInfo.Grade.String)
	// }

	// accessGrade := config.UserGrades.IndexOf(board.GrantComment.String)

	// if accessGrade < userGrade {
	// 	c.Text(http.StatusForbidden, "Forbidden")
	// 	return
	// }

	// comment, err := crud.GetComment(board, fmt.Sprint(boardIdx), fmt.Sprint(commentIdx))
	// if err != nil {
	// 	return c.Status(http.StatusNotFound).SendString("Comment was not found")
	// }

	// userInfo, err = crud.GetUserByName(c.AuthInfo.(model.AuthInfo).Name.String)
	// if err != nil {
	// 	c.Text(http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// if comment.AuthorIdx.Int64 != userInfo.Idx.Int64 && userInfo.Grade.String != "admin" && userInfo.Grade.String != "manager" {
	// 	c.Text(http.StatusForbidden, "forbidden")
	// 	return
	// }

	err = crud.DeleteComment(board, fmt.Sprint(boardIdx), fmt.Sprint(commentIdx))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]interface{}{
		"result": "success",
	}

	return c.Status(http.StatusOK).JSON(result)
}
