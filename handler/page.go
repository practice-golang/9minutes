package handler

import (
	"9minutes/consts"
	"9minutes/internal/crud"
	"9minutes/model"
	"html"
	"log"
	"net/http"
	"strings"

	"gopkg.in/guregu/null.v4"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/microcosm-cc/bluemonday"
)

var bm = bluemonday.UGCPolicy()

func HealthCheckAPI(c *fiber.Ctx) error {
	return c.SendString("Ok")
}

func getSessionValue(sess *session.Session, key string) (result string) {
	value := sess.Get(key)
	if value != nil {
		result = value.(string)
	}

	return result
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

	userid := getSessionValue(sess, "userid")
	grade := getSessionValue(sess, "grade")

	templateMap["Title"] = "9minutes" // Todo: Remove or change to site title
	templateMap["UserId"] = userid
	templateMap["Grade"] = grade

	boardListingOptions := model.BoardListingOptions{}
	boardListingOptions.Page = null.IntFrom(0)
	boardListingOptions.ListCount = null.IntFrom(9999)
	// boardList, _ := crud.GetBoards(boardListingOptions)
	// templateMap["BoardPageData"] = boardList
	templateMap["BoardList"] = BoardListALL

	templateMap["PendingUser"] = true
	if grade != "pending_user" {
		templateMap["PendingUser"] = false
	}

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

			if name == "board/read" {
				comments, err := GetCommentList(boardCode, idx, map[string]string{"page": "0"})
				if err != nil {
					return err
				}
				templateMap["Comments"] = comments
			}

		case "board/write":
			boardCode := queries["board_code"]

			if boardCode == "" {
				return c.Status(http.StatusBadRequest).SendString("no board code")
			}

			name = "board/write"
		}
	case strings.HasPrefix(name, "mypage"):
		if userid == "" {
			name = "status/unauthorized"
			break
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

func HandlePostingList(c *fiber.Ctx) error {
	name := strings.TrimSuffix(c.Path()[1:], "/")
	queries := c.Queries()
	templateMap := fiber.Map{}

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	userid := getSessionValue(sess, "userid")
	grade := getSessionValue(sess, "grade")

	templateMap["Title"] = "9minutes" // Todo: Remove or change to site title
	templateMap["UserId"] = userid
	templateMap["Grade"] = grade

	boardCode := queries["board_code"]
	if boardCode == "" {
		return c.Status(http.StatusBadRequest).SendString("missing parameter - board")
	}
	page := queries["page"]
	if page == "" {
		page = "1"
	}

	boardListingOptions := model.BoardListingOptions{}
	boardListingOptions.Page = null.IntFrom(0)
	boardListingOptions.ListCount = null.IntFrom(9999)
	// boardList, _ := crud.GetBoards(boardListingOptions)
	// templateMap["BoardPageData"] = boardList
	templateMap["BoardList"] = BoardListALL

	currentBoard, _ := crud.GetBoardByCode(model.Board{BoardCode: null.StringFrom(boardCode)})

	templateMap["Accessible"] = false
	if consts.UserGrades[grade].Point >= consts.UserGradesForGrant[currentBoard.GrantRead.String].Point {
		templateMap["Accessible"] = true
	}

	templateMap["PendingUser"] = true
	if grade != "pending_user" {
		templateMap["PendingUser"] = false
	}

	list, err := GetPostingList(boardCode, queries)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	templateMap["BoardCode"] = boardCode
	templateMap["PostingList"] = list

	err = c.Render(name, templateMap)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return c.Status(http.StatusNotFound).SendString("Page not Found")
		}

		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return nil
}
