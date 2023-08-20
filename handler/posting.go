package handler

import (
	"9minutes/config"
	"9minutes/internal/crud"
	"9minutes/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/guregu/null.v4"
)

func GetPostingList(boardCODE string, queries map[string]string) (model.PostingPageData, error) {
	var err error

	board := model.Board{BoardCode: null.StringFrom(boardCODE)}
	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return model.PostingPageData{}, err
	}

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

	return list, err
}

func GetPostingData(boardCode, idx string) (model.Posting, error) {
	var err error

	board := model.Board{BoardCode: null.StringFrom(boardCode)}
	posting := model.Posting{}

	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return model.Posting{}, err
	}

	posting, err = crud.GetPosting(board, idx)
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

	posting, err := GetPostingData(boardCode, idx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	comments, err := GetCommentList(boardCode, idx, queries)
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
	board := model.Board{}
	posting := model.Posting{}

	board.BoardCode = null.StringFrom(c.Params("board_code"))

	err = c.BodyParser(&posting)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	now := time.Now().Format("20060102150405")
	posting.RegDate = null.StringFrom(now)
	posting.Views = null.IntFrom(0)

	_, err = crud.WritePosting(board, posting)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]interface{}{
		"result": "success",
	}

	// c.Json(http.StatusOK, result)
	return c.Status(http.StatusOK).JSON(result)
}

func UpdatePostingAPI(c *fiber.Ctx) (err error) {
	var board model.Board
	var posting model.Posting
	var deleteList model.FilesToDelete

	board.BoardCode = null.StringFrom(c.Params("board_code"))
	board, err = crud.GetBoardByCode(board)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Board was not found")
	}

	idx, _ := strconv.Atoi(c.Params("idx"))
	posting.Idx = null.IntFrom(int64(idx))

	rbody := c.Body()

	err = json.Unmarshal(rbody, &posting)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	err = crud.UpdatePosting(board, posting, "update")
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	err = json.Unmarshal(rbody, &deleteList)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	result := map[string]interface{}{
		"result": "success",
	}

	return c.Status(http.StatusOK).JSON(result)
}

func DeletePostingAPI(c *fiber.Ctx) error {
	var board model.Board

	board.BoardCode = null.StringFrom(c.Params("board_code"))

	idx, _ := strconv.Atoi(c.Params("idx"))

	board, err := crud.GetBoardByCode(board)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Board was not found")
	}

	posting, err := crud.GetPosting(board, fmt.Sprint(idx))
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

	err = crud.DeletePosting(board, fmt.Sprint(idx))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	err = crud.DeleteComments(board, fmt.Sprint(idx))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]interface{}{
		"result": "success",
	}

	return c.Status(http.StatusOK).JSON(result)
}
