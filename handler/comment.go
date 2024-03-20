package handler

import (
	"9minutes/config"
	"9minutes/internal/crud"
	"9minutes/model"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/guregu/null.v4"
)

func GetCommentList(board model.Board, topicIdx string, queries map[string]string) (model.CommentPageData, error) {
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

	commentListOption := model.CommentListingOptions{}
	commentListOption.Search = null.StringFrom(commentSearch)

	commentListOption.Page = null.IntFrom(int64(page - 1))
	commentListOption.ListCount = null.IntFrom(int64(count))

	comments, err := crud.GetComments(board, topicIdx, commentListOption)
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

func GetCommentsAPI(c *fiber.Ctx) (err error) {
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
	topicIdx := c.Params("topic_idx")

	board := BoardListData[boardCode]
	accessible := checkBoardAccessible(board.GrantRead.String, grade)
	if !accessible {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"status": 403, "message": "forbidden"})
	}

	comments, err := GetCommentList(board, topicIdx, queries)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(comments)
}

func GetComment(c *fiber.Ctx) (err error) {
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
	topicIdx := c.Params("topic_idx")
	commentIdx := c.Params("comment_idx")

	board := BoardListData[boardCode]
	accessible := checkBoardAccessible(board.GrantRead.String, grade)
	if !accessible {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"status": 403, "message": "forbidden"})
	}

	comment, err := crud.GetComment(board, topicIdx, commentIdx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(comment)
}

func WriteCommentAPI(c *fiber.Ctx) (err error) {
	comment := model.Comment{}

	boardCode := c.Params("board_code")
	board := BoardListData[boardCode]

	err = c.BodyParser(&comment)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

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

	if !checkBoardAccessible(board.GrantComment.String, grade) {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"status": 403, "message": "forbidden"})
	}

	if userid != "" {
		comment.AuthorIdx = null.IntFrom(useridx)
		comment.AuthorName = null.StringFrom(userid)
	}

	clientIP := c.Context().RemoteIP().String()
	clientIPs := strings.Split(clientIP, ".")
	comment.AuthorIpFull = null.StringFrom(clientIP)
	comment.AuthorIP = null.StringFrom(clientIPs[0] + "." + clientIPs[1])

	topicIdx, err := strconv.ParseInt(c.Params("topic_idx"), 0, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
	comment.TopicIdx = null.IntFrom(topicIdx)

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

	result := map[string]interface{}{"result": "success"}
	return c.Status(http.StatusOK).JSON(result)
}

func UpdateCommentAPI(c *fiber.Ctx) error {
	boardCode := c.Params("board_code")
	topicIdx := c.Params("topic_idx")
	commentIdx := c.Params("comment_idx")

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

	board := BoardListData[boardCode]

	accessible := checkBoardAccessible(board.GrantComment.String, grade)
	if !accessible {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"status": 403, "message": "forbidden"})
	}

	commentPrev, err := crud.GetComment(board, topicIdx, commentIdx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if commentPrev.AuthorIdx.Int64 != useridx || grade == "admin" {
		result := map[string]interface{}{"result": "fail", "msg": "user is not author"}
		return c.Status(http.StatusBadRequest).JSON(result)
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

	err = crud.UpdateComment(board, comment, fmt.Sprint(topicIdx))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]interface{}{"result": "success"}
	return c.Status(http.StatusOK).JSON(result)
}

func DeleteCommentAPI(c *fiber.Ctx) error {
	boardCode := c.Params("board_code")
	topicIdx := c.Params("topic_idx")
	commentIdx := c.Params("comment_idx")
	board := BoardListData[boardCode]

	if strings.TrimSpace(topicIdx) == "" {
		result := map[string]interface{}{"result": "fail", "msg": "empty topic index"}
		return c.Status(http.StatusBadRequest).JSON(result)
	}
	if strings.TrimSpace(commentIdx) == "" {
		result := map[string]interface{}{"result": "fail", "msg": "empty comment index"}
		return c.Status(http.StatusBadRequest).JSON(result)
	}

	comment, err := crud.GetComment(board, topicIdx, commentIdx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
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

	accessible := checkBoardAccessible(board.GrantComment.String, grade)
	if !accessible {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"status": 403, "message": "forbidden"})
	}

	switch true {
	case useridx < 0 || userid == "" || comment.AuthorIdx.Int64 < 0:
		deletePassword := ""
		deletePasswords := c.GetReqHeaders()["Delete-Password"]
		if len(deletePasswords) > 0 {
			deletePassword = deletePasswords[0]
		}

		if comment.EditPassword.String != deletePassword {
			result := map[string]interface{}{"result": "fail", "msg": "incorrect password"}
			return c.Status(http.StatusBadRequest).JSON(result)
		}
	case grade != "admin" && comment.AuthorIdx.Int64 != useridx:
		result := map[string]interface{}{"result": "fail", "msg": "user is not author"}
		return c.Status(http.StatusBadRequest).JSON(result)
	}

	uploadIndices := strings.Split(comment.Files.String, "|")
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

	err = crud.DeleteComment(board, fmt.Sprint(topicIdx), fmt.Sprint(commentIdx))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]interface{}{"result": "success"}
	return c.Status(http.StatusOK).JSON(result)
}
