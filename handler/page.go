package handler

import (
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
var boardActions = []string{"list", "read", "write", "edit"}

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
	if userid == "" {
		grade = "guest"
	}

	templateMap["Title"] = "9minutes" // Todo: Remove or change to site title
	templateMap["UserId"] = userid
	templateMap["Grade"] = grade

	// boardListingOptions := model.BoardListingOptions{}
	// boardListingOptions.Page = null.IntFrom(0)
	// boardListingOptions.ListCount = null.IntFrom(9999)
	// boardList, _ := crud.GetBoards(boardListingOptions)
	// templateMap["BoardPageData"] = boardList
	templateMap["BoardList"] = BoardListData

	templateMap["PendingUser"] = true
	if grade != "user_hold" {
		templateMap["PendingUser"] = false
	}

	routePath = appendIndexToRoutePath(routePath)
	routePaths := strings.Split(routePath, "/")

	switch routePaths[0] {
	case "board":
		boardActionExist := checkBoardActionExist(routePaths[1])
		if !boardActionExist {
			return c.Status(http.StatusForbidden).SendString("Forbidden")
		}

		boardCode := queries["board_code"]
		board := BoardListData[boardCode]

		switch routePath {
		case "board/list":
			accessible := checkBoardAccessible(board.GrantRead.String, grade)
			if !accessible {
				return c.Status(http.StatusForbidden).SendString("Forbidden")
			}

			list, err := GetPostingList(boardCode, queries)
			if err != nil {
				return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
			}
			templateMap["BoardCode"] = boardCode
			templateMap["PostingList"] = list
		case "board/read":
			accessible := checkBoardAccessible(board.GrantRead.String, grade)
			if !accessible {
				return c.Status(http.StatusForbidden).SendString("Forbidden")
			}

			posting, err := GetPostingData(boardCode, queries["idx"])
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}

			posting.Content = null.StringFrom(html.UnescapeString(posting.Content.String))
			templateMap["Posting"] = posting

			comments, err := GetCommentList(board, queries["idx"], map[string]string{"page": "0"})
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
			templateMap["Comments"] = comments
		case "board/write":
			accessible := checkBoardAccessible(board.GrantWrite.String, grade)
			if !accessible {
				return c.Status(http.StatusForbidden).SendString("Forbidden")
			}

			if boardCode == "" {
				return c.Status(http.StatusBadRequest).SendString("no board code")
			}
		case "board/edit":
			accessible := checkBoardAccessible(board.GrantWrite.String, grade)
			if !accessible {
				return c.Status(http.StatusForbidden).SendString("Forbidden")
			}

			posting, err := GetPostingData(boardCode, queries["idx"])
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}

			posting.Content = null.StringFrom(html.UnescapeString(posting.Content.String))
			templateMap["Posting"] = posting
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

		if grade != "admin" {
			routePath = "status/unauthorized"
			break
		}

		routePath = "admin/index"
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
