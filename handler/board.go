package handler

import (
	"9minutes/crud"
	"9minutes/db"
	"9minutes/model"
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/guregu/null.v4"
)

func HandleBoardList(c *fiber.Ctx) error {
	var err error

	// queries := c.URL.Query()
	queries := c.Queries()

	listingOptions := model.BoardListingOptions{}
	listingOptions.Search = null.StringFrom(queries["search"])

	listingOptions.Page = null.IntFrom(1)
	listingOptions.ListCount = null.IntFrom(100)

	if queries["count"] != "" {
		countPerPage, err := strconv.Atoi(queries["count"])
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}

		listingOptions.ListCount = null.IntFrom(int64(countPerPage))
	}

	if queries["page"] != "" {
		page := queries["page"]
		pageNum, err := strconv.Atoi(page)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}

		listingOptions.Page = null.IntFrom(int64(pageNum))
	}

	listingOptions.Page.Int64--

	h, err := LoadHTML(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	list, err := crud.GetBoards(listingOptions)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	listJSON, _ := json.Marshal(list)
	h = bytes.ReplaceAll(h, []byte("$BOARD_LIST$"), listJSON)

	return c.Status(http.StatusOK).Send(h)
}

// GetBoards - API Get boards list
func GetBoards(c *fiber.Ctx) error {
	queries := c.Queries()

	listingOptions := model.BoardListingOptions{}
	listingOptions.Search = null.StringFrom(queries["search"])

	listingOptions.Page = null.IntFrom(1)
	listingOptions.ListCount = null.IntFrom(2)

	if queries["page"] != "" {
		page := queries["page"]
		pageNum, err := strconv.Atoi(page)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}

		if queries["count"] != "" {
			countPerPage, err := strconv.Atoi(queries["count"])
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
	boardList := result.BoardList

	// catalog 작업전까지 일단 무시
	if len(boardList) > 0 {
		var fields model.Field
		if boardList[0].Fields != nil {
			err = json.Unmarshal(boardList[0].Fields.([]byte), &fields)
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
			boardList[0].Fields = fields
		}
	}

	return c.Status(http.StatusOK).JSON(result)
}

// AddBoard - Add board
func AddBoard(c *fiber.Ctx) error {
	var board model.Board

	err := c.BodyParser(&board)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	err = crud.AddBoard(board)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	err = db.Obj.CreateBoard(board, false)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	err = db.Obj.CreateComment(board, false)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]string{
		"result": "ok",
	}

	return c.Status(http.StatusOK).JSON(result)
}

func UpdateBoard(c *fiber.Ctx) error {
	var err error

	boardDatas := []map[string]interface{}{}
	boardDatasSucess := []map[string]interface{}{}
	boardDatasFailed := []map[string]interface{}{}

	err = c.BodyParser(&boardDatas)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	for _, boardData := range boardDatas {
		var boardOLD, boardNEW model.Board

		idx, err := strconv.Atoi(boardData["idx"].(string))
		if err != nil {
			responseData := map[string]interface{}{"data": boardData, "error": err.Error()}
			boardDatasFailed = append(boardDatasFailed, responseData)
			continue
		}

		boardNEW = model.Board{
			Idx:          null.IntFrom(int64(idx)),
			BoardName:    null.StringFrom(boardData["board-name"].(string)),
			BoardCode:    null.StringFrom(boardData["board-code"].(string)),
			BoardType:    null.StringFrom(boardData["board-type"].(string)),
			BoardTable:   null.StringFrom(boardData["board-table"].(string)),
			CommentTable: null.StringFrom(boardData["comment-table"].(string)),
			GrantRead:    null.StringFrom(boardData["grant-read"].(string)),
			GrantWrite:   null.StringFrom(boardData["grant-write"].(string)),
			GrantComment: null.StringFrom(boardData["grant-comment"].(string)),
			GrantUpload:  null.StringFrom(boardData["grant-upload"].(string)),
			Fields:       boardData["fields"],
		}

		boardOLD, err = crud.GetBoardByIdx(boardNEW)
		if err != nil {
			responseData := map[string]interface{}{"data": boardData, "error": err.Error()}
			boardDatasFailed = append(boardDatasFailed, responseData)
			continue
		}

		err = crud.UpdateBoard(boardNEW)
		if err != nil {
			responseData := map[string]interface{}{"data": boardData, "error": err.Error()}
			boardDatasFailed = append(boardDatasFailed, responseData)
			continue
		}

		if boardOLD.BoardTable.String != boardNEW.BoardTable.String {
			err = db.Obj.RenameBoard(boardOLD, boardNEW)
			if err != nil {
				responseData := map[string]interface{}{"data": boardData, "error": err.Error()}
				boardDatasFailed = append(boardDatasFailed, responseData)
				continue
			}
		}

		if boardOLD.CommentTable.String != boardNEW.CommentTable.String {
			err = db.Obj.RenameComment(boardOLD, boardNEW)
			if err != nil {
				responseData := map[string]interface{}{"data": boardData, "error": err.Error()}
				boardDatasFailed = append(boardDatasFailed, responseData)
				continue
			}
		}

		responseData := map[string]interface{}{"data": boardData, "error": ""}
		boardDatasSucess = append(boardDatasSucess, responseData)
	}

	result := map[string]interface{}{"result": "ok"}
	if len(boardDatasFailed) > 0 {
		result["result"] = "fail"
		result["failed"] = boardDatasFailed
		result["success"] = boardDatasSucess
	}

	return c.Status(http.StatusOK).JSON(result)
}

func DeleteBoard(c *fiber.Ctx) error {
	boardDatas := []map[string]interface{}{}
	boardDatasSuccess := []map[string]interface{}{}
	boardDatasFailed := []map[string]interface{}{}

	err := c.BodyParser(&boardDatas)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	for _, boardData := range boardDatas {
		var board model.Board

		idx, _ := strconv.Atoi(boardData["idx"].(string))

		board.Idx = null.IntFrom(int64(idx))

		board, err := crud.GetBoardByIdx(board)
		if err != nil {
			responseData := map[string]interface{}{"data": boardData, "error": err.Error()}
			boardDatasFailed = append(boardDatasFailed, responseData)
			continue
		}

		err = crud.DeleteBoard(board)
		if err != nil {
			responseData := map[string]interface{}{"data": boardData, "error": err.Error()}
			boardDatasFailed = append(boardDatasFailed, responseData)
			continue
		}

		err = db.Obj.DeleteBoard(board)
		if err != nil {
			responseData := map[string]interface{}{"data": boardData, "error": err.Error()}
			boardDatasFailed = append(boardDatasFailed, responseData)
			continue
		}

		err = db.Obj.DeleteComment(board)
		if err != nil {
			responseData := map[string]interface{}{"data": boardData, "error": err.Error()}
			boardDatasFailed = append(boardDatasFailed, responseData)
			continue
		}

		responseData := map[string]interface{}{"data": boardData, "error": ""}
		boardDatasSuccess = append(boardDatasSuccess, responseData)
	}

	result := map[string]interface{}{"result": "ok"}
	if len(boardDatasFailed) > 0 {
		result["result"] = "failed"
		result["failed"] = boardDatasFailed
		result["success"] = boardDatasSuccess
	}

	return c.Status(http.StatusOK).JSON(result)
}
