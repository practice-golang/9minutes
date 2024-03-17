package handler

import (
	"9minutes/config"
	"9minutes/internal/crud"
	"9minutes/model"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/guregu/null.v4"
)

func GetCommentList(board model.Board, postingIDX string, queries map[string]string) (model.CommentPageData, error) {
	var err error

	page := 1
	count := config.CommentCountPerPage
	if queries["page"] != "" {
		page, err = strconv.Atoi(queries["page"])
		if err != nil {
			return model.CommentPageData{}, nil
		}
	}
	if queries["count"] != "" {
		count, err = strconv.Atoi(queries["count"])
		if err != nil {
			return model.CommentPageData{}, nil
		}
	}

	commentSearch := queries["search"]

	commentOptions := model.CommentListingOptions{}
	commentOptions.Search = null.StringFrom(commentSearch)

	commentOptions.Page = null.IntFrom(int64(page - 1))
	commentOptions.ListCount = null.IntFrom(int64(count))

	comments, err := crud.GetComments(board, postingIDX, commentOptions)
	if err != nil {
		return model.CommentPageData{}, err
	}

	pageList := []int{}
	pageShowGap := 2
	for i := comments.CurrentPage - pageShowGap; i <= comments.CurrentPage+pageShowGap; i++ {
		if i > 0 && i <= comments.CurrentPage+pageShowGap && i <= comments.TotalPage {
			pageList = append(pageList, i)
		}
	}
	comments.PageList = pageList

	pageJumpGap := 5
	comments.JumpPrev = comments.CurrentPage - pageJumpGap
	comments.JumpNext = comments.CurrentPage + pageJumpGap
	if comments.JumpPrev < 1 {
		comments.JumpPrev = 1
	}
	if comments.JumpNext > comments.TotalPage {
		comments.JumpNext = comments.TotalPage
	}

	return comments, nil
}

func GetComments(c *fiber.Ctx) (err error) {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	userid := getSessionValue(sess, "userid")
	grade := getSessionValue(sess, "grade")
	if userid == "" {
		grade = "guest"
	}

	boardCode, queries := c.Params("board_code"), c.Queries()
	postingIdx := c.Params("posting_idx")

	board := BoardListData[boardCode]
	accessible := checkBoardAccessible(board.GrantRead.String, grade)
	if !accessible {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"status": 403, "message": "forbidden"})
	}

	comments, err := GetCommentList(board, postingIdx, queries)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(comments)
}

func WriteComment(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	userid := getSessionValue(sess, "userid")
	grade := getSessionValue(sess, "grade")
	if userid == "" {
		grade = "guest"
	}

	comment := model.Comment{}
	postingIdx, err := strconv.ParseInt(c.Params("posting_idx"), 0, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
	comment.PostingIdx = null.IntFrom(postingIdx)

	err = c.BodyParser(&comment)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	boardCode := c.Params("board_code")
	board := BoardListData[boardCode]
	if !checkBoardAccessible(board.GrantComment.String, grade) {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"status": 403, "message": "forbidden"})
	}

	useridx := int64(-1)
	useridxInterface := sess.Get("idx")
	if useridxInterface != nil {
		useridx = useridxInterface.(int64)
	}

	comment.AuthorIdx = null.IntFrom(useridx)
	comment.AuthorName = null.StringFrom(userid)

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

func UpdateComment(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	userid := getSessionValue(sess, "userid")
	grade := getSessionValue(sess, "grade")
	if userid == "" {
		grade = "guest"
	}

	boardCode := c.Params("board_code")
	board := BoardListData[boardCode]
	postingIdx := c.Params("posting_idx")
	commentIdx := c.Params("comment_idx")

	accessible := checkBoardAccessible(board.GrantComment.String, grade)
	if !accessible {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"status": 403, "message": "forbidden"})
	}

	comment := model.Comment{}
	err = c.BodyParser(&comment)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	commentIdxINT, err := strconv.Atoi(commentIdx)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("comment index is not correct")
	}
	comment.Idx = null.IntFrom(int64(commentIdxINT))

	err = crud.UpdateComment(board, comment, fmt.Sprint(postingIdx))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]interface{}{
		"result": "success",
	}

	return c.Status(http.StatusOK).JSON(result)
}

func DeleteComment(c *fiber.Ctx) error {
	boardCode := c.Params("board_code")
	postingIdx := c.Params("posting_idx")
	commentIdx := c.Params("comment_idx")

	board := BoardListData[boardCode]

	err := crud.DeleteComment(board, fmt.Sprint(postingIdx), fmt.Sprint(commentIdx))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]interface{}{
		"result": "success",
	}

	return c.Status(http.StatusOK).JSON(result)
}
