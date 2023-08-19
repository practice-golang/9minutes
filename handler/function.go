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

	"9minutes/config"
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

	list, err := GetPostingList(boardCode, queries)
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
		case "board/read", "board/edit":
			boardCode := queries["board_code"]
			idx := queries["idx"]

			posting, err := GetPostingData(boardCode, idx)
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
			posting.Content = null.StringFrom(html.UnescapeString(posting.Content.String))
			templateMap["Posting"] = posting

			comments, err := GetCommentsList(boardCode, idx, map[string]string{})
			if err != nil {
				return err
			}
			templateMap["Comments"] = comments

		case "board/write":
			boardCode := queries["board_code"]

			if boardCode == "" {
				return c.Status(http.StatusBadRequest).SendString("no board code")
			}

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

func ListPostingAPI(c *fiber.Ctx) (err error) {
	boardCode, queries := c.Params("board_code"), c.Queries()
	list, err := GetPostingList(boardCode, queries)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(list)
}

func ReadPostingAPI(c *fiber.Ctx) (err error) {
	boardCode, queries := c.Params("board_code"), c.Queries()
	idx := c.Params("idx")

	posting, err := GetPostingData(boardCode, idx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	comments, err := GetCommentsList(boardCode, idx, queries)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := make(map[string]interface{})
	result["posting"] = posting
	result["comments"] = comments

	return c.Status(http.StatusOK).JSON(result)
}

// WritePostingAPI - Write posting API
func WritePostingAPI(c *fiber.Ctx) (err error) {
	board := model.Board{}
	posting := model.Posting{}

	board.BoardCode = null.StringFrom(c.Params("board_code"))

	err = c.BodyParser(&posting)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	now := time.Now().Format("20060102150405")
	posting.RegDate = null.StringFrom(now)
	posting.Views = null.IntFrom(0)

	_, err = crud.WritePosting(board, posting)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]interface{}{
		"result": "success",
	}

	// c.Json(http.StatusOK, result)
	return c.Status(http.StatusOK).JSON(result)
}

func UpdatePostingAPI(c *fiber.Ctx) (err error) {
	var board model.Board
	var posting model.Posting
	var deleteList model.FilesToDelete

	board.BoardCode = null.StringFrom(c.Params("board_code"))
	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Board was not found")
	}

	idx, _ := strconv.Atoi(c.Params("idx"))
	posting.Idx = null.IntFrom(int64(idx))

	rbody := c.Body()

	err = json.Unmarshal(rbody, &posting)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	err = crud.UpdatePosting(board, posting, "update")
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	err = json.Unmarshal(rbody, &deleteList)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	result := map[string]interface{}{
		"result": "success",
	}

	return c.Status(http.StatusOK).JSON(result)
}

func DeletePostingAPI(c *fiber.Ctx) error {
	var board model.Board

	board.BoardCode = null.StringFrom(c.Params("board_code"))

	idx, _ := strconv.Atoi(c.Params("idx"))

	board, err := crud.GetBoardByCode(board)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Board was not found")
	}

	posting, err := crud.GetPosting(board, fmt.Sprint(idx))
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Posting was not found")
	}

	uploadIndices := strings.Split(posting.Files.String, "|")
	for _, f := range uploadIndices {
		if f == "" {
			continue
		}
		fidx, err := strconv.Atoi(f)
		if err != nil {
			continue
		}

		fdata, err := crud.GetUploadedFile(fidx)
		if err != nil {
			continue
		}

		err = crud.DeleteUploadedFile(int64(fidx))
		if err != nil {
			continue
		}

		filepath := config.UploadPath + "/" + fdata.StorageName.String
		DeleteUploadFile(filepath)
	}

	err = crud.DeletePosting(board, fmt.Sprint(idx))
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
	postingIdx := c.Params("idx")
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	comments, err := GetCommentsList(boardCode, postingIdx, queries)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(comments)
}

func WriteComment(c *fiber.Ctx) error {
	board := model.Board{BoardCode: null.StringFrom(c.Params("board_code"))}
	comment := model.Comment{}

	postingIdx, err := strconv.ParseInt(c.Params("posting_idx"), 0, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
	comment.BoardIdx = null.IntFrom(postingIdx)

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
