package handler

import (
	"9minutes/config"
	"9minutes/consts"
	"9minutes/internal/crud"
	"9minutes/model"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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
	} else {
		passwordEncrypted, err := bcrypt.GenerateFromPassword([]byte(comment.EditPassword.String), consts.BcryptCost)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		comment.EditPassword = null.StringFrom(string(passwordEncrypted))
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

	// err = crud.WriteComment(board, comment)
	_, commentIdx, err := crud.WriteComment(board, comment)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	var imIndices []int
	var fdatas []model.StoredFileInfo
	if strings.TrimSpace(comment.Files.String) != "" {
		imIndicesStr := strings.Split(comment.Files.String, "|")

		for _, imIdxStr := range imIndicesStr {
			imIdx, _ := strconv.ParseInt(imIdxStr, 0, 64)
			imIndices = append(imIndices, int(imIdx))
		}

		fdatas, _ = crud.GetUploadedFiles(imIndices)
		for _, fdata := range fdatas {
			crud.SetUploadedFileIndex(fdata.Idx.Int64, topicIdx, commentIdx, "write")
		}
	}

	result := map[string]interface{}{"result": "success"}
	return c.Status(http.StatusOK).JSON(result)
}

func UpdateCommentAPI(c *fiber.Ctx) error {
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

	accessible := checkBoardAccessible(board.GrantComment.String, grade)
	if !accessible {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"status": 403, "message": "forbidden"})
	}

	commentPrev, err := crud.GetComment(board, topicIdx, commentIdx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	switch true {
	case grade != "admin" && commentPrev.AuthorIdx.Int64 < 0:
		editPassword := ""
		editPasswords := c.GetReqHeaders()["Edit-Password"]
		if len(editPasswords) > 0 {
			editPassword = editPasswords[0]
		}

		err = bcrypt.CompareHashAndPassword([]byte(commentPrev.EditPassword.String), []byte(editPassword))
		if err != nil {
			result := map[string]interface{}{"result": "fail", "msg": "incorrect password"}
			return c.Status(http.StatusBadRequest).JSON(result)
		}
	case grade != "admin" && commentPrev.AuthorIdx.Int64 != useridx:
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

	comment.Content = null.StringFrom(bm.Sanitize(comment.Content.String))
	if comment.Content.String == "" {
		return c.Status(http.StatusBadRequest).SendString("comment is empty")
	}

	err = crud.UpdateComment(board, comment, fmt.Sprint(topicIdx))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	var imIndices []int
	if strings.TrimSpace(comment.Files.String) != "" {
		imIndicesStr := strings.Split(comment.Files.String, "|")
		for _, imIdxStr := range imIndicesStr {
			imIdx, _ := strconv.ParseInt(imIdxStr, 0, 64)
			imIndices = append(imIndices, int(imIdx))
		}

		fdatas, _ := crud.GetUploadedFiles(imIndices)
		for _, fdata := range fdatas {
			crud.SetUploadedFileIndex(fdata.Idx.Int64, commentPrev.TopicIdx.Int64, commentPrev.Idx.Int64, "update")
		}
	}

	imIndices = []int{}
	if strings.TrimSpace(comment.DeleteFiles.String) != "" {
		imIndicesStr := strings.Split(comment.DeleteFiles.String, "|")
		for _, imIdxStr := range imIndicesStr {
			imIdx, _ := strconv.ParseInt(imIdxStr, 0, 64)
			imIndices = append(imIndices, int(imIdx))
		}

		fdatas, _ := crud.GetUploadedFiles(imIndices)
		for _, fdata := range fdatas {
			crud.DeleteUploadedFile(fdata.Idx.Int64, commentPrev.TopicIdx.Int64, commentPrev.Idx.Int64)
			if fdata.StorageName.Valid && fdata.StorageName.String != "" {
				DeleteUploadFile(config.UploadPath + "/" + fdata.StorageName.String)
			}
		}
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

	comment, err := crud.GetComment(board, topicIdx, commentIdx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	switch true {
	case grade != "admin" && comment.AuthorIdx.Int64 < 0:
		deletePassword := ""
		deletePasswords := c.GetReqHeaders()["Delete-Password"]
		if len(deletePasswords) > 0 {
			deletePassword = deletePasswords[0]
		}

		err = bcrypt.CompareHashAndPassword([]byte(comment.EditPassword.String), []byte(deletePassword))
		if err != nil {
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

		topicIDX, _ := strconv.ParseInt(fdata.TopicIdx.String, 10, 64)
		commentIDX, _ := strconv.ParseInt(fdata.CommentIdx.String, 10, 64)
		err = crud.DeleteUploadedFile(int64(fidx), topicIDX, commentIDX)
		if err != nil {
			continue
		}

		if fdata.StorageName.Valid && fdata.StorageName.String != "" {
			filepath := config.UploadPath + "/" + fdata.StorageName.String
			DeleteUploadFile(filepath)
		} else {
			log.Println("Empty filename:", fdata)
		}
	}

	err = crud.DeleteComment(board, fmt.Sprint(topicIdx), fmt.Sprint(commentIdx))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]interface{}{"result": "success"}
	return c.Status(http.StatusOK).JSON(result)
}
