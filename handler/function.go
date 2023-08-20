package handler

import (
	"9minutes/internal/crud"
	"9minutes/model"
	"html"
	"log"
	"net/http"
	"strings"

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

	boardListingOptions := model.BoardListingOptions{}
	boardListingOptions.Page = null.IntFrom(0)
	boardListingOptions.ListCount = null.IntFrom(9999)
	boardList, _ := crud.GetBoards(boardListingOptions)
	templateMap["BoardPageData"] = boardList

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

	boardListingOptions := model.BoardListingOptions{}
	boardListingOptions.Page = null.IntFrom(0)
	boardListingOptions.ListCount = null.IntFrom(9999)
	boardList, _ := crud.GetBoards(boardListingOptions)
	templateMap["BoardPageData"] = boardList

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

			comments, err := GetCommentList(boardCode, idx, map[string]string{"page": "0"})
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
