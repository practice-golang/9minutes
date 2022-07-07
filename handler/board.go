package handler

import (
	"9minutes/crud"
	"9minutes/db"
	"9minutes/model"
	"9minutes/router"
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/guregu/null.v4"
)

func HandleBoardList(c *router.Context) {
	var err error

	queries := c.URL.Query()

	listingOptions := model.BoardListingOptions{}
	listingOptions.Search = null.StringFrom(queries.Get("search"))

	listingOptions.Page = null.IntFrom(1)
	listingOptions.ListCount = null.IntFrom(100)

	if queries.Get("count") != "" {
		countPerPage, err := strconv.Atoi(queries.Get("count"))
		if err != nil {
			c.Text(http.StatusBadRequest, err.Error())
			return
		}

		listingOptions.ListCount = null.IntFrom(int64(countPerPage))
	}

	if queries.Get("page") != "" {
		page := queries.Get("page")
		pageNum, err := strconv.Atoi(page)
		if err != nil {
			c.Text(http.StatusBadRequest, err.Error())
			return
		}

		listingOptions.Page = null.IntFrom(int64(pageNum))
	}

	listingOptions.Page.Int64--

	h, err := LoadHTML(c)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
		return
	}

	list, err := crud.GetBoards(listingOptions)
	if err != nil {
		c.Text(http.StatusInternalServerError, err.Error())
	}

	listJSON, _ := json.Marshal(list)
	h = bytes.ReplaceAll(h, []byte("$BOARD_LIST$"), listJSON)

	c.Html(http.StatusOK, h)
}

// GetBoards - API Get boards list
func GetBoards(c *router.Context) {
	queries := c.URL.Query()

	listingOptions := model.BoardListingOptions{}
	listingOptions.Search = null.StringFrom(queries.Get("search"))

	listingOptions.Page = null.IntFrom(1)
	listingOptions.ListCount = null.IntFrom(2)

	if queries.Get("page") != "" {
		page := queries.Get("page")
		pageNum, err := strconv.Atoi(page)
		if err != nil {
			c.Text(http.StatusBadRequest, err.Error())
			return
		}

		if queries.Get("count") != "" {
			countPerPage, err := strconv.Atoi(queries.Get("count"))
			if err != nil {
				c.Text(http.StatusBadRequest, err.Error())
				return
			}

			listingOptions.ListCount = null.IntFrom(int64(countPerPage))
		}

		listingOptions.Page = null.IntFrom(int64(pageNum))
	}

	listingOptions.Page.Int64--

	result, err := crud.GetBoards(listingOptions)
	boardList := result.BoardList

	// catalog 작업전까지 일단 무시
	if len(boardList) > 0 {
		var fields model.Field
		if boardList[0].Fields != nil {
			err = json.Unmarshal(boardList[0].Fields.([]byte), &fields)
			if err != nil {
				c.Text(http.StatusOK, err.Error())
				return
			}
			boardList[0].Fields = fields
		}
	}

	c.Json(http.StatusOK, result)
}

// AddBoard - Add board
func AddBoard(c *router.Context) {
	var board model.Board

	err := json.NewDecoder(c.Body).Decode(&board)
	if err != nil {
		c.Text(http.StatusOK, err.Error())
		return
	}

	err = crud.AddBoard(board)
	if err != nil {
		c.Text(http.StatusOK, err.Error())
		return
	}

	err = db.Obj.CreateBoard(board, false)
	if err != nil {
		c.Text(http.StatusOK, err.Error())
		return
	}

	err = db.Obj.CreateComment(board, false)
	if err != nil {
		c.Text(http.StatusOK, err.Error())
		return
	}

	result := map[string]string{
		"result": "ok",
	}

	c.Json(http.StatusOK, result)
}

func UpdateBoard(c *router.Context) {
	var boardOLD, boardNEW model.Board

	err := json.NewDecoder(c.Body).Decode(&boardNEW)
	if err != nil {
		c.Text(http.StatusOK, err.Error())
		return
	}

	boardOLD, err = crud.GetBoardByIdx(boardNEW)
	if err != nil {
		c.Text(http.StatusOK, err.Error())
		return
	}

	err = crud.UpdateBoard(boardNEW)
	if err != nil {
		c.Text(http.StatusOK, err.Error())
		return
	}

	if boardOLD.BoardTable.String != boardNEW.BoardTable.String {
		err = db.Obj.RenameBoard(boardOLD, boardNEW)
		if err != nil {
			c.Text(http.StatusOK, err.Error())
			return
		}
	}

	if boardOLD.CommentTable.String != boardNEW.CommentTable.String {
		err = db.Obj.RenameComment(boardOLD, boardNEW)
		if err != nil {
			c.Text(http.StatusOK, err.Error())
			return
		}
	}

	result := map[string]string{
		"result": "ok",
	}

	c.Json(http.StatusOK, result)
}

func DeleteBoard(c *router.Context) {
	var board model.Board

	uri := strings.Split(c.URL.Path, "/")
	idx, _ := strconv.Atoi(uri[len(uri)-1])

	board.Idx = null.IntFrom(int64(idx))

	board, err := crud.GetBoardByIdx(board)
	if err != nil {
		c.Text(http.StatusOK, err.Error())
		return
	}

	err = crud.DeleteBoard(board)
	if err != nil {
		c.Text(http.StatusOK, err.Error())
		return
	}

	err = db.Obj.DeleteBoard(board)
	if err != nil {
		c.Text(http.StatusOK, err.Error())
		return
	}

	err = db.Obj.DeleteComment(board)
	if err != nil {
		c.Text(http.StatusOK, err.Error())
		return
	}

	result := map[string]string{
		"result": "ok",
	}

	c.Json(http.StatusOK, result)
}
