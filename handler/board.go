package handler

import (
	"9minutes/consts"
	"9minutes/internal/crud"
	"9minutes/internal/db"
	"9minutes/model"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/guregu/null.v4"
)

// GetBoardsAPI - API Get boards list
func GetBoardsAPI(c *fiber.Ctx) error {
	queries := c.Queries()

	boardListOptions := model.BoardListingOption{}
	boardListOptions.Search = null.StringFrom(queries["search"])

	boardListOptions.Page = null.IntFrom(1)
	boardListOptions.ListCount = null.IntFrom(20)

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

			boardListOptions.ListCount = null.IntFrom(int64(countPerPage))
		}

		boardListOptions.Page = null.IntFrom(int64(pageNum))
	}

	boardListOptions.Page.Int64--

	result, err := crud.GetBoards(boardListOptions)
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

// AddBoardAPI - Add board
func AddBoardAPI(c *fiber.Ctx) error {
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

	LoadBoardListData()

	return c.Status(http.StatusOK).JSON(result)
}

func UpdateBoardAPI(c *fiber.Ctx) error {
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

	LoadBoardListData()

	return c.Status(http.StatusOK).JSON(result)
}

func DeleteBoardAPI(c *fiber.Ctx) error {
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

	LoadBoardListData()

	return c.Status(http.StatusOK).JSON(result)
}

func BoardListAPI(c *fiber.Ctx) (err error) {
	queries := c.Queries()

	boardListOption := model.BoardListingOption{}
	boardListOption.Search = null.StringFrom(queries["search"])

	boardListOption.Page = null.IntFrom(1)
	boardListOption.ListCount = null.IntFrom(10)

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

			boardListOption.ListCount = null.IntFrom(int64(countPerPage))
		}

		boardListOption.Page = null.IntFrom(int64(pageNum))
	}

	boardListOption.Page.Int64--

	result, err := crud.GetBoards(boardListOption)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(result)
}

func checkBoardActionExist(action string) bool {
	result := false
	for _, a := range boardActions {
		if a == action {
			result = true
			break
		}
	}

	return result
}

func checkBoardAccessible(boardGradeKey string, userGradeKey string) bool {
	boardRank := consts.BoardGrades[boardGradeKey].Rank
	userRank := consts.UserGrades[userGradeKey].Rank

	result := false
	if userRank <= boardRank {
		result = true
	}

	return result
}

func LoadBoardListData() map[string]model.Board {
	BoardListData = map[string]model.Board{}

	boardListOption := model.BoardListingOption{}
	boardListOption.Page = null.IntFrom(0)
	boardListOption.ListCount = null.IntFrom(99999)
	boardList, _ := crud.GetBoards(boardListOption)

	for _, b := range boardList.BoardList {
		BoardListData[b.BoardCode.String] = b
	}

	return BoardListData
}
