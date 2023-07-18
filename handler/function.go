package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"9minutes/config"
	"9minutes/internal/crud"
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
		switch name {
		case "board":
			// board := params["board"]
			// page := params["page"]

			name = "board/index"
		case "board/write":
			name = "board/index"
		}

	case strings.HasPrefix(name, "admin"):
		if userid == "" {
			name = "status/unauthorized"
			break
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

func ListContentAPI(c *fiber.Ctx) (err error) {
	boardCode, queries := c.Params("board_code"), c.Queries()
	list, err := GetContentsList(boardCode, queries)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(list)
}

func ReadContentAPI(c *fiber.Ctx) (err error) {
	boardCode, queries := c.Params("board_code"), c.Queries()
	idx, err := strconv.ParseInt(c.Params("idx"), 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	content, err := GetContentData(boardCode, int(idx), queries)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	comments, err := GetCommentsList(boardCode, int(idx), queries)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := make(map[string]interface{})
	result["content"] = content
	result["comments"] = comments

	return c.Status(http.StatusOK).JSON(result)
}

// WriteContentAPI - Write content API
func WriteContentAPI(c *fiber.Ctx) (err error) {
	board := model.Board{}
	content := model.Content{}

	board.BoardCode = null.StringFrom(c.Params("board_code"))

	err = c.BodyParser(&content)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	now := time.Now().Format("20060102150405")
	content.RegDate = null.StringFrom(now)
	content.Views = null.IntFrom(0)

	_, err = crud.WriteContent(board, content)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]interface{}{
		"result": "success",
	}

	// c.Json(http.StatusOK, result)
	return c.Status(http.StatusOK).JSON(result)
}

func UpdateContentAPI(c *fiber.Ctx) (err error) {
	var board model.Board
	var content model.Content
	var deleteList model.FilesToDelete

	board.BoardCode = null.StringFrom(c.Params("board_code"))
	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Board was not found")
	}

	idx, _ := strconv.Atoi(c.Params("idx"))
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

func DeleteContentAPI(c *fiber.Ctx) error {
	var board model.Board

	board.BoardCode = null.StringFrom(c.Params("board_code"))

	idx, _ := strconv.Atoi(c.Params("idx"))

	board, err := crud.GetBoardByCode(board)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Board was not found")
	}

	content, err := crud.GetContent(board, fmt.Sprint(idx))
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Content was not found")
	}

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

	contentIdx := int(content.Idx.Int64)
	commentList, err := crud.GetComments(board, contentIdx, listingOptions)
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

	err = crud.DeleteComment(board, fmt.Sprint(boardIdx), fmt.Sprint(commentIdx))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]interface{}{
		"result": "success",
	}

	return c.Status(http.StatusOK).JSON(result)
}
