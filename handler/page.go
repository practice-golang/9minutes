package handler

import (
	"9minutes/consts"
	"9minutes/internal/crud"
	"9minutes/model"
	"html"
	"net/http"
	"strings"

	"gopkg.in/guregu/null.v4"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/microcosm-cc/bluemonday"
)

var bm = bluemonday.UGCPolicy()
var indexPaths = []string{"", "admin", "board"}

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

func appendIndexToRoutePath(name string) string {
	for _, p := range indexPaths {
		if name == p {
			name += "/index"
			name = strings.TrimPrefix(name, "/")
			break
		}
	}

	return name
}

// HandleHTML - Handle HTML template layout
func HandleHTML(c *fiber.Ctx) error {
	routePath := strings.TrimSuffix(c.Path()[1:], "/")
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

	routePath = appendIndexToRoutePath(routePath)
	routePaths := strings.Split(routePath, "/")

	switch routePaths[0] {
	case "board":
		boardCode := queries["board_code"]

		switch routePath {
		case "board/read":
			posting, err := GetPostingData(boardCode, queries["idx"])
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}

			posting.Content = null.StringFrom(html.UnescapeString(posting.Content.String))
			templateMap["Posting"] = posting

			comments, err := GetCommentList(boardCode, queries["idx"], map[string]string{"page": "0"})
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
			templateMap["Comments"] = comments
		case "board/edit":
			posting, err := GetPostingData(boardCode, queries["idx"])
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}

			posting.Content = null.StringFrom(html.UnescapeString(posting.Content.String))
			templateMap["Posting"] = posting
		case "board/write":
			if boardCode == "" {
				return c.Status(http.StatusBadRequest).SendString("no board code")
			}
		default:
			routePath = "status/unauthorized"
		}
	case "mypage":
		if userid == "" {
			routePath = "status/unauthorized"
			break
		}
	case "admin":
		if userid == "" {
			routePath = "status/unauthorized"
			break
		}

		usergrade, err := GetSessionUserGrade(c)
		if err != nil {
			return err
		}

		if usergrade != "admin" {
			routePath = "status/unauthorized"
			break
		}
	}

	err = c.Render(routePath, templateMap)
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
