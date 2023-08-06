package handler

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"9minutes/internal/crud"
	"9minutes/model"

	"gopkg.in/guregu/null.v4"

	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
)

var bm = bluemonday.UGCPolicy()

func HealthCheckAPI(c *fiber.Ctx) error {
	return c.SendString("Ok")
}

func HelloParam(c *fiber.Ctx) error {
	if len(c.Params("name")) > 0 {
		return c.Status(http.StatusOK).SendString("Hello " + c.Params("name"))
	} else {
		return c.Status(http.StatusBadRequest).SendString("Missing parameter")
	}
}

func HandleBoardHTML(c *fiber.Ctx) error {
	name := strings.TrimSuffix(c.Path()[1:], "/")
	queries := c.Queries()
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

	boardCode := queries["board_code"]
	page := queries["page"]

	if boardCode == "" {
		return c.Status(http.StatusBadRequest).SendString("missing parameter - board")
	}
	if page == "" {
		page = "1"
	}

	list, err := GetContentsList(boardCode, queries)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	templateMap["BoardCode"] = boardCode
	templateMap["BoardList"] = list

	err = c.Render(name, templateMap)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return c.Status(http.StatusNotFound).SendString("Page not Found")
		}
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return nil
}

// HandleHTML - Handle HTML template layout
func HandleHTML(c *fiber.Ctx) error {
	name := strings.TrimSuffix(c.Path()[1:], "/")
	queries := c.Queries()
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

		if queries["hello"] != "" {
			log.Printf("Hello: %s", queries["hello"])
		}
	case strings.HasPrefix(name, "board"):
		switch name {
		case "board":
			name = "board/index"
		case "board/read":
			boardCode := queries["board_code"]
			idx := queries["idx"]

			content, err := GetContentData(boardCode, idx)
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
			content.Content = null.StringFrom(html.UnescapeString(content.Content.String))
			templateMap["Content"] = content

		case "board/write":
			name = "board/write"
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

func BoardListAPI(c *fiber.Ctx) (err error) {
	queries := c.Queries()

	listingOptions := model.BoardListingOptions{}
	listingOptions.Search = null.StringFrom(queries["search"])

	listingOptions.Page = null.IntFrom(1)
	listingOptions.ListCount = null.IntFrom(10)

	if queries["page"] != "" {
		page := queries["page"]
		pageNum, err := strconv.Atoi(page)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}

		if queries["list-count"] != "" {
			countPerPage, err := strconv.Atoi(queries["list-count"])
			if err != nil {
				return c.Status(http.StatusBadRequest).SendString(err.Error())
			}

			listingOptions.ListCount = null.IntFrom(int64(countPerPage))
		}

		listingOptions.Page = null.IntFrom(int64(pageNum))
	}

	listingOptions.Page.Int64--

	result, err := crud.GetBoards(listingOptions)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(result)
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
	idx := c.Params("idx")

	content, err := GetContentData(boardCode, idx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	comments, err := GetCommentsList(boardCode, idx, queries)
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

func GetComments(c *fiber.Ctx) (err error) {
	boardCode, queries := c.Params("board_code"), c.Queries()
	contentIdx := c.Params("idx")
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	comments, err := GetCommentsList(boardCode, contentIdx, queries)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(comments)
}

func WriteComment(c *fiber.Ctx) error {
	board := model.Board{BoardCode: null.StringFrom(c.Params("board_code"))}
	comment := model.Comment{}

	contentIdx, err := strconv.ParseInt(c.Params("content_idx"), 0, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
	comment.BoardIdx = null.IntFrom(contentIdx)

	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("board not found")
	}

	err = c.BodyParser(&comment)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	comment.AuthorIdx = null.IntFrom(-1)

	// Check session
	sess, err := store.Get(c)
	useridInterface := sess.Get("idx")
	if err == nil && useridInterface != nil {
		userIDX, err := strconv.Atoi(useridInterface.(string))
		if err == nil {
			comment.AuthorIdx = null.IntFrom(int64(userIDX))
		}
	}

	comment.Content = null.StringFrom(bm.Sanitize(comment.Content.String))
	if comment.Content.String == "" {
		return c.Status(http.StatusBadRequest).SendString("comment is empty")
	}

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
