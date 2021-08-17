package board

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/practice-golang/9minutes/models"

	"github.com/practice-golang/9minutes/db"

	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"
)

// AddBoards - Insert board(s) info
func AddBoards(c echo.Context) error {
	var err error

	data, _ := ioutil.ReadAll(c.Request().Body)
	boards, fields := prepareInsertData(data)

	sqlResult, err := db.InsertData(boards)
	if err != nil {
		log.Println("InsertData: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	failed := []string{}

	for k, b := range boards {
		f := []models.Field{}
		if fields[k] != nil {
			f = fields[k]
		}
		switch b.Type.String {
		case "basic-board":
			err = db.Dbi.CreateBasicBoard(b, false)
			if err != nil {
				log.Println("Add Boards CreateBasicBoard", err)
				failed = append(failed, b.Name.String)
			}
			err = db.Dbi.CreateComment(b, false)
			if err != nil {
				log.Println("Add Boards CreateBasicBoard", err)
				failed = append(failed, b.Name.String)
			}
		case "custom-board":
			err := db.Dbi.CreateCustomBoard(b, f, false)
			if err != nil {
				log.Println("Add Boards CreateCustomBoard", err)
				failed = append(failed, b.Name.String)
			}
		case "custom-tablelist":
			err := db.Dbi.CreateCustomBoard(b, f, false)
			if err != nil {
				log.Println("Add Boards CreateCustomBoard", err)
				failed = append(failed, b.Name.String)
			}
		}
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	result := map[string]string{
		"fails":         strings.Join(failed, "/"),
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
}

// GetBoard - Get a board info
func GetBoard(c echo.Context) error {
	idx := c.Param("idx")

	dataINTF, err := db.SelectData(models.Board{Idx: null.NewString(idx, true)})
	if err != nil {
		log.Println("SelectData: ", err)
	}
	data := prepareSelectData(dataINTF)

	return c.JSON(http.StatusOK, data)
}

// GetBoards - Get all(but limit 10 by db.SelectData) boards info
func GetBoards(c echo.Context) error {
	var err error
	dataINTF, err := db.SelectData(models.Board{})
	if err != nil {
		log.Println("SelectData: ", err)
	}

	data := prepareSelectData(dataINTF)

	return c.JSON(http.StatusOK, data)
}

// GetBoardByCode - Get a board info
func GetBoardByCode(boardName string) []models.Board {
	dataINTF, err := db.SelectData(models.Board{Code: null.NewString(boardName, true)})
	if err != nil {
		log.Fatal("SelectData by code: ", err)
	}
	data := prepareSelectData(dataINTF)

	return data
}

// SearchBoards - Search board(s) or paging
func SearchBoards(c echo.Context) error {
	var search models.BoardSearch
	var err error

	if err := c.Bind(&search); err != nil {
		log.Fatal("Search_SelectData: ", err)
	}

	dataINTF, err := db.SelectData(search)
	if err != nil {
		log.Fatal("Search_SelectData: ", err)
	}

	data := prepareSelectData(dataINTF)

	return c.JSON(http.StatusOK, data)
}

// EditBoard - Modify board
func EditBoard(c echo.Context) error {
	data, _ := ioutil.ReadAll(c.Request().Body)
	board := prepareUpdateData(data)

	boardPreviousINTF, err := db.SelectData(models.Board{Idx: board.Idx})
	if err != nil {
		log.Fatal("EditData/SelectData: ", err)
	}

	boardsPrevious := prepareSelectData(boardPreviousINTF)
	if len(boardsPrevious) > 0 {
		boardPrevious := boardsPrevious[0]
		switch boardPrevious.Type.String {
		case "basic-board":
			boardPrevious.Fields = nil
			board.Fields = nil
			if cmp.Equal(boardPrevious, board) {
				return c.JSON(http.StatusBadRequest, map[string]string{"msg": "Both tables are same"})
			} else if boardPrevious.Table.String != board.Table.String {
				db.Dbi.EditBasicBoard(boardPrevious, board)
			}
		case "custom-board":
			log.Println("custom-board")

			err := db.Dbi.EditCustomBoard(boardPrevious, board)
			if err != nil {
				return c.JSON(http.StatusBadRequest, err.Error())
			}
		case "custom-tablelist":
			log.Println("custom-tablelist")

			err := db.Dbi.EditCustomBoard(boardPrevious, board)
			if err != nil {
				return c.JSON(http.StatusBadRequest, err.Error())
			}
		default:
			log.Println("Edit board: No proper board type")
		}
	}

	sqlResult, err := db.UpdateData(board)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"msg": string(err.Error())})
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
}

// DeleteBoard - Delete a board
func DeleteBoard(c echo.Context) error {
	idx := c.Param("idx")

	// Drop table
	dataINTF, err := db.SelectData(models.Board{Idx: null.NewString(idx, true)})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	data := prepareSelectData(dataINTF)

	if len(data) > 0 {
		err = db.Dbi.DeleteBoard(data[0].Table.String)
	} else {
		err = errors.New("delete table failed")
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Remove table information data
	sqlResult, err := db.DeleteData("IDX", idx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	lastID, _ := sqlResult.LastInsertId()
	affRows, _ := sqlResult.RowsAffected()

	result := map[string]string{
		"last-id":       fmt.Sprint(lastID),
		"affected-rows": fmt.Sprint(affRows),
	}

	return c.JSON(http.StatusOK, result)
}

// GetTotalPage - Get total page
func GetTotalPage(c echo.Context) error {
	var search models.BoardSearch

	if err := c.Bind(&search); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	data, err := db.SelectCount(search)
	if err != nil {
		log.Fatal("SelectCount: ", err)
	}

	countPerPage := uint(1)
	if search.Options.Count.Valid {
		countPerPage = uint(search.Options.Count.Int64)
	}

	pages := uint(math.Ceil(float64(data) / float64(countPerPage)))

	result := map[string]uint{"total-page": pages}

	return c.JSON(http.StatusOK, result)
}
