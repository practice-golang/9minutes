package handler

import (
	"html"
	"net/http"
	"strings"

	"gopkg.in/guregu/null.v4"

	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
)

var bm = bluemonday.UGCPolicy()
var indexPaths = []string{"", "admin", "board"}
var boardActions = []string{"list", "read", "write", "edit"}

func HealthCheckAPI(c *fiber.Ctx) error {
	return c.SendString("Ok")
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
				routePath = "status/unauthorized"
				break
			}

			list, err := GetPostingList(boardCode, queries)
			if err != nil {
				return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
			}
			templateMap["BoardCode"] = boardCode
			templateMap["Data"] = list

			if board.BoardType.String == "gallery" {
				routePath = strings.Replace(routePath, "board/list", "board/gallery", 1)
			}
		case "board/read":
			accessible := checkBoardAccessible(board.GrantRead.String, grade)
			if !accessible {
				routePath = "status/unauthorized"
				break
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
				routePath = "status/unauthorized"
				break
			}

			if boardCode == "" {
				return c.Status(http.StatusBadRequest).SendString("no board code")
			}
		case "board/edit":
			accessible := checkBoardAccessible(board.GrantWrite.String, grade)
			if !accessible {
				routePath = "status/unauthorized"
				break
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
