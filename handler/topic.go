package handler

import (
	"9minutes/config"
	"9minutes/internal/crud"
	"9minutes/model"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/guregu/null.v4"
)

func GetTopicList(boardCode string, queries map[string]string) (model.TopicPageData, error) {
	var err error

	board := BoardListData[boardCode]

	page := 1
	count := config.TopicCountPerPage
	if queries["page"] != "" {
		page, err = strconv.Atoi(queries["page"])
		if err != nil {
			return model.TopicPageData{}, nil
		}
	}
	if queries["count"] != "" {
		count, _ = strconv.Atoi(queries["count"])
	}
	if count < 1 {
		count = config.TopicCountPerPage
	}

	topicListOption := model.TopicListingOption{}
	topicListOption.Search = null.StringFrom(queries["search"])

	topicListOption.Page = null.IntFrom(int64(page - 1))
	topicListOption.ListCount = null.IntFrom(int64(count))

	topicList, err := crud.GetTopicList(board, topicListOption)
	if err != nil {
		return model.TopicPageData{}, nil
	}

	pageList := []int{}
	pageShowGap := 2
	for i := topicList.CurrentPage - pageShowGap; i <= topicList.CurrentPage+pageShowGap; i++ {
		if i > 0 && i <= topicList.CurrentPage+pageShowGap && i <= topicList.TotalPage {
			pageList = append(pageList, i)
		}
	}
	topicList.PageList = pageList

	pageJumpGap := 5
	topicList.JumpPrev = topicList.CurrentPage - pageJumpGap
	topicList.JumpNext = topicList.CurrentPage + pageJumpGap
	if topicList.JumpPrev < 1 {
		topicList.JumpPrev = 1
	}
	if topicList.JumpNext > topicList.TotalPage {
		topicList.JumpNext = topicList.TotalPage
	}

	topicList.ListCount = count

	return topicList, err
}

func GetTopicData(boardCode, idx string) (model.Topic, error) {
	board := BoardListData[boardCode]

	topic, err := crud.GetTopic(board, idx)
	if err != nil {
		return model.Topic{}, err
	}

	topic.Views = null.IntFrom(topic.Views.Int64 + 1)
	err = crud.UpdateTopic(board, topic, "viewcount")
	if err != nil {
		return model.Topic{}, err
	}

	return topic, nil
}

func ListTopicAPI(c *fiber.Ctx) (err error) {
	boardCode, queries := c.Params("board_code"), c.Queries()
	list, err := GetTopicList(boardCode, queries)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(list)
}

func ReadTopicAPI(c *fiber.Ctx) (err error) {
	boardCode, queries := c.Params("board_code"), c.Queries()
	idx := c.Params("idx")

	board := BoardListData[boardCode]

	topic, err := GetTopicData(boardCode, idx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	comments, err := GetCommentList(board, idx, queries)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := make(map[string]interface{})
	result["topic"] = topic
	result["comments"] = comments

	return c.Status(http.StatusOK).JSON(result)
}

// WriteTopicAPI - Write topic API
func WriteTopicAPI(c *fiber.Ctx) (err error) {
	topic := model.Topic{}

	boardCode := c.Params("board_code")
	board := BoardListData[boardCode]

	err = c.BodyParser(&topic)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
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

	accessible := checkBoardAccessible(board.GrantWrite.String, grade)
	if !accessible {
		return c.Status(http.StatusBadRequest).SendString("access denied")
	}

	if userid != "" {
		topic.AuthorIdx = null.IntFrom(useridx)
		topic.AuthorName = null.StringFrom(userid)
	}

	clientIP := c.Context().RemoteIP().String()
	clientIPs := strings.Split(clientIP, ".")
	topic.AuthorIpFull = null.StringFrom(clientIP)
	topic.AuthorIP = null.StringFrom(clientIPs[0] + "." + clientIPs[1])

	now := time.Now().Format("20060102150405")
	topic.RegDate = null.StringFrom(now)
	topic.Views = null.IntFrom(0)

	if strings.TrimSpace(topic.Files.String) != "" {
		var imIndices []int
		imIndicesStr := strings.Split(topic.Files.String, "|")

		for _, imIdxStr := range imIndicesStr {
			imIdx, _ := strconv.ParseInt(imIdxStr, 0, 64)
			imIndices = append(imIndices, int(imIdx))
		}

		fileDatas, _ := crud.GetUploadedFiles(imIndices)
		for _, fileData := range fileDatas {
			filename := fileData.StorageName.String
			fext := filepath.Ext(filename)
			fname := filename[0 : len(filename)-len(fext)]

			fpathFrom := "upload/" + filename
			fnameTo := fname + "_thumb.png"
			fpathTo := "upload/" + fnameTo

			if !CheckFileExtensionIsImage(filename) {
				continue
			}

			if CopyResizeImagePNG(fpathFrom, fpathTo, 800, 800) != nil {
				continue
			}

			topic.TitleImage = null.StringFrom(fnameTo)
			break
		}
	}

	_, err = crud.WriteTopic(board, topic)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]interface{}{"result": "success"}
	return c.Status(http.StatusOK).JSON(result)
}

func UpdateTopicAPI(c *fiber.Ctx) (err error) {
	var topic model.Topic
	var deleteList model.FilesToDelete

	boardCode := c.Params("board_code")
	board := BoardListData[boardCode]

	// idx, _ := strconv.Atoi(c.Params("idx"))
	idx := c.Params("idx")
	if strings.TrimSpace(idx) == "" {
		result := map[string]interface{}{"result": "fail", "msg": "empty index"}
		return c.Status(http.StatusBadRequest).JSON(result)
	}

	idxNUM, _ := strconv.Atoi(idx)
	topic.Idx = null.IntFrom(int64(idxNUM))

	rbody := c.Body()

	topicPrev, err := crud.GetTopic(board, idx)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
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

	err = json.Unmarshal(rbody, &topic)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	switch true {
	case topicPrev.AuthorIdx.Int64 < 0:
		if topic.EditPassword.String != topicPrev.EditPassword.String {
			result := map[string]interface{}{"result": "fail", "msg": "incorrect password"}
			return c.Status(http.StatusBadRequest).JSON(result)
		}
	case grade != "admin" && topicPrev.AuthorIdx.Int64 != useridx:
		result := map[string]interface{}{"result": "fail", "msg": "user is not author"}
		return c.Status(http.StatusBadRequest).JSON(result)
	}

	if strings.TrimSpace(topic.Files.String) != "" {
		var imIndices []int
		imIndicesStr := strings.Split(topic.Files.String, "|")

		for _, imIdxStr := range imIndicesStr {
			imIdx, _ := strconv.ParseInt(imIdxStr, 0, 64)
			imIndices = append(imIndices, int(imIdx))
		}

		fileDatas, _ := crud.GetUploadedFiles(imIndices)
		for _, fileData := range fileDatas {
			filename := fileData.StorageName.String
			fext := filepath.Ext(filename)
			fname := filename[0 : len(filename)-len(fext)]

			fpathFrom := "upload/" + filename
			fnameTo := fname + "_thumb.png"
			fpathTo := "upload/" + fnameTo

			if !CheckFileExtensionIsImage(filename) {
				continue
			}

			if CopyResizeImagePNG(fpathFrom, fpathTo, 800, 800) != nil {
				continue
			}

			topic.TitleImage = null.StringFrom(fnameTo)
			break
		}
	}

	err = crud.UpdateTopic(board, topic, "update")
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	err = json.Unmarshal(rbody, &deleteList)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	DeleteUploadFile("upload/" + topicPrev.TitleImage.String)

	result := map[string]interface{}{"result": "success"}
	return c.Status(http.StatusOK).JSON(result)
}

func DeleteTopicAPI(c *fiber.Ctx) error {
	boardCode := c.Params("board_code")
	board := BoardListData[boardCode]

	idx := c.Params("idx")
	if strings.TrimSpace(idx) == "" {
		result := map[string]interface{}{"result": "fail", "msg": "empty index"}
		return c.Status(http.StatusBadRequest).JSON(result)
	}

	topic, err := crud.GetTopic(board, idx)
	if err != nil {
		result := map[string]interface{}{"result": "fail", "msg": err.Error()}
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

	switch true {
	case useridx < 0 || userid == "" || topic.AuthorIdx.Int64 < 0:
		deletePassword := ""
		deletePasswords := c.GetReqHeaders()["Delete-Password"]
		if len(deletePasswords) > 0 {
			deletePassword = deletePasswords[0]
		}

		if topic.EditPassword.String != deletePassword {
			result := map[string]interface{}{"result": "fail", "msg": "incorrect password"}
			return c.Status(http.StatusBadRequest).JSON(result)
		}
	case grade != "admin" && topic.AuthorIdx.Int64 != useridx:
		result := map[string]interface{}{"result": "fail", "msg": "user is not author"}
		return c.Status(http.StatusBadRequest).JSON(result)
	}

	uploadIndices := strings.Split(topic.Files.String, "|")
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

	err = crud.DeleteTopic(board, idx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	err = crud.DeleteComments(board, idx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	DeleteUploadFile("upload/" + topic.TitleImage.String)

	result := map[string]interface{}{"result": "success"}
	return c.Status(http.StatusOK).JSON(result)
}
