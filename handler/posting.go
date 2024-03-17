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

func GetPostingList(boardCode string, queries map[string]string) (model.PostingPageData, error) {
	var err error

	board := BoardListData[boardCode]

	page := 1
	count := config.PostingListCountPerPage
	if queries["page"] != "" {
		page, err = strconv.Atoi(queries["page"])
		if err != nil {
			return model.PostingPageData{}, nil
		}
	}
	if queries["count"] != "" {
		count, err = strconv.Atoi(queries["count"])
		if err != nil {
			return model.PostingPageData{}, nil
		}
	}

	listingOptions := model.PostingListingOptions{}
	listingOptions.Search = null.StringFrom(queries["search"])

	listingOptions.Page = null.IntFrom(int64(page - 1))
	listingOptions.ListCount = null.IntFrom(int64(count))

	list, err := crud.GetPostingList(board, listingOptions)
	if err != nil {
		return model.PostingPageData{}, nil
	}

	pageList := []int{}
	pageShowGap := 2
	for i := list.CurrentPage - pageShowGap; i <= list.CurrentPage+pageShowGap; i++ {
		if i > 0 && i <= list.CurrentPage+pageShowGap && i <= list.TotalPage {
			pageList = append(pageList, i)
		}
	}
	list.PageList = pageList

	pageJumpGap := 5
	list.JumpPrev = list.CurrentPage - pageJumpGap
	list.JumpNext = list.CurrentPage + pageJumpGap
	if list.JumpPrev < 1 {
		list.JumpPrev = 1
	}
	if list.JumpNext > list.TotalPage {
		list.JumpNext = list.TotalPage
	}

	list.ListCount = count

	return list, err
}

func GetPostingData(boardCode, idx string) (model.Posting, error) {
	board := BoardListData[boardCode]

	posting, err := crud.GetPosting(board, idx)
	if err != nil {
		return model.Posting{}, err
	}

	posting.Views = null.IntFrom(posting.Views.Int64 + 1)
	err = crud.UpdatePosting(board, posting, "viewcount")
	if err != nil {
		return model.Posting{}, err
	}

	return posting, nil
}

func ListPostingAPI(c *fiber.Ctx) (err error) {
	boardCode, queries := c.Params("board_code"), c.Queries()
	list, err := GetPostingList(boardCode, queries)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(list)
}

func ReadPostingAPI(c *fiber.Ctx) (err error) {
	boardCode, queries := c.Params("board_code"), c.Queries()
	idx := c.Params("idx")

	board := BoardListData[boardCode]

	posting, err := GetPostingData(boardCode, idx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	comments, err := GetCommentList(board, idx, queries)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := make(map[string]interface{})
	result["posting"] = posting
	result["comments"] = comments

	return c.Status(http.StatusOK).JSON(result)
}

// WritePostingAPI - Write posting API
func WritePostingAPI(c *fiber.Ctx) (err error) {
	posting := model.Posting{}

	err = c.BodyParser(&posting)
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
	if userid == "" {
		return c.Status(http.StatusBadRequest).SendString("userid is empty")
	}

	posting.AuthorIdx = null.IntFrom(useridx)
	posting.AuthorName = null.StringFrom(userid)

	now := time.Now().Format("20060102150405")
	posting.RegDate = null.StringFrom(now)
	posting.Views = null.IntFrom(0)

	boardCode := c.Params("board_code")
	board := BoardListData[boardCode]

	if strings.TrimSpace(posting.Files.String) != "" {
		var imIndices []int
		imIndicesStr := strings.Split(posting.Files.String, "|")

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

			posting.TitleImage = null.StringFrom(fnameTo)
			break
		}
	}

	_, err = crud.WritePosting(board, posting)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]interface{}{
		"result": "success",
	}

	return c.Status(http.StatusOK).JSON(result)
}

func UpdatePostingAPI(c *fiber.Ctx) (err error) {
	var posting model.Posting
	var deleteList model.FilesToDelete

	boardCode := c.Params("board_code")
	board := BoardListData[boardCode]

	idx, _ := strconv.Atoi(c.Params("idx"))
	posting.Idx = null.IntFrom(int64(idx))

	rbody := c.Body()

	postingPrev, err := crud.GetPosting(board, c.Params("idx"))
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	err = json.Unmarshal(rbody, &posting)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	if strings.TrimSpace(posting.Files.String) != "" {
		var imIndices []int
		imIndicesStr := strings.Split(posting.Files.String, "|")

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

			posting.TitleImage = null.StringFrom(fnameTo)
			break
		}
	}

	err = crud.UpdatePosting(board, posting, "update")
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	err = json.Unmarshal(rbody, &deleteList)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	DeleteUploadFile("upload/" + postingPrev.TitleImage.String)

	result := map[string]interface{}{
		"result": "success",
	}

	return c.Status(http.StatusOK).JSON(result)
}

func DeletePostingAPI(c *fiber.Ctx) error {
	boardCode := c.Params("board_code")
	board := BoardListData[boardCode]

	// idx, _ := strconv.Atoi(c.Params("idx"))
	idx := c.Params("idx")
	if strings.TrimSpace(idx) == "" {
		result := map[string]interface{}{
			"result": "fail",
			"msg":    "index is empty",
		}
		return c.Status(http.StatusBadRequest).JSON(result)
	}

	posting, err := crud.GetPosting(board, idx)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Posting was not found")
	}

	uploadIndices := strings.Split(posting.Files.String, "|")
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

	err = crud.DeletePosting(board, idx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	err = crud.DeleteComments(board, idx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	DeleteUploadFile("upload/" + posting.TitleImage.String)

	result := map[string]interface{}{
		"result": "success",
	}

	return c.Status(http.StatusOK).JSON(result)
}
