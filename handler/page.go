package handler

import (
	"9minutes/config"
	"9minutes/consts"
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

	useridx := int64(-1)
	useridxInterface := sess.Get("idx")
	if useridxInterface != nil {
		useridx = useridxInterface.(int64)
	}
	userid := getSessionValue(sess, "userid")
	grade := getSessionValue(sess, "grade")
	if userid == "" {
		grade = "guest"
	}

	userInfo := map[string]interface{}{
		"UserIdx":       useridx,
		"UserID":        userid,
		"UserGrade":     grade,
		"UserGradeRank": consts.UserGrades[grade].Rank,
		"PendingUser":   (grade == "user_hold"),
	}

	templateMap["SiteName"] = config.SiteName

	templateMap["UserInfo"] = userInfo
	templateMap["UserGrades"] = consts.UserGrades
	templateMap["BoardGrades"] = consts.BoardGrades

	templateMap["BoardList"] = BoardListData
	templateMap["TopicListCount"] = config.TopicCountPerPage

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

		templateMap["BoardCode"] = boardCode
		templateMap["BoardName"] = board.BoardName.String

		switch routePath {
		case "board/list":
			accessible := checkBoardAccessible(board.GrantRead.String, grade)
			if !accessible {
				routePath = "status/unauthorized"
				break
			}

			list, err := GetTopicList(boardCode, queries)
			if err != nil {
				return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
			}
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

			topic, err := GetTopicData(boardCode, queries["idx"])
			if err != nil {
				routePath = "status/page_not_found"
				break
			}

			topic.EditPassword = null.NewString("", false)

			topic.Content = null.StringFrom(html.UnescapeString(topic.Content.String))
			templateMap["Topic"] = topic

			templateMap["GrantComment"] = false
			if checkBoardAccessible(board.GrantComment.String, grade) {
				templateMap["GrantComment"] = true
			}

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

			topic, err := GetTopicData(boardCode, queries["idx"])
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}

			switch true {
			case useridx < 0 || userid == "" || topic.AuthorIdx.Int64 < 0:
				editPassword := string(c.Request().Header.Cookie("password"))
				c.ClearCookie("password")

				if topic.EditPassword.String != editPassword {
					routePath = "status/access_denied"
				}
			case grade != "admin" && topic.AuthorIdx.Int64 != useridx:
				routePath = "status/access_denied"
			}

			topic.Content = null.StringFrom(html.UnescapeString(topic.Content.String))

			templateMap["Topic"] = topic
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
