package handler

import (
	"9minutes/config"
	"9minutes/consts"
	"html"
	"log"
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

	clientIP := c.Context().RemoteIP()
	log.Println("Client IP:", clientIP)

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

	// Todo: create template func and remove this
	boardGradeRanks := []int{}
	for _, d := range BoardListData {
		boardGradeRanks = append(boardGradeRanks, consts.BoardGrades[d.GrantRead.String].Rank)
	}

	templateMap["UserGrades"] = consts.UserGrades
	templateMap["BoardGrades"] = consts.BoardGrades

	templateMap["UserInfo"] = userInfo
	templateMap["BoardGradeRanks"] = boardGradeRanks

	templateMap["SiteName"] = "9minutes" // Todo: Remove or change to site title
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

			topic, err := GetTopicData(boardCode, queries["idx"])
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}

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
