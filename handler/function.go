package handler

import (
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"

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
