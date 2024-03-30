package handler

import (
	"9minutes/internal/crud"
	"9minutes/model"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/guregu/null.v4"
)

func GetMemberListAPI(c *fiber.Ctx) (err error) {
	boardCode := c.Params("board_code")
	board := BoardListData[boardCode]

	members, err := crud.GetMemberList(board.Idx.Int64)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(members)
}

func AddMemberAPI(c *fiber.Ctx) (err error) {
	boardCode := c.Params("board_code")
	userIdxSTR := c.Params("user_idx")

	board := BoardListData[boardCode]
	userIDX, err := strconv.ParseInt(userIdxSTR, 10, 64)
	grade := "member"
	now := time.Now().Format("20060102150405")
	if !board.Idx.Valid || board.Idx.Int64 < 1 {
		return c.Status(http.StatusBadRequest).Send([]byte("board code is wrong"))
	}
	if err != nil {
		return c.Status(http.StatusBadRequest).Send([]byte(err.Error()))
	}
	if strings.TrimSpace(userIdxSTR) == "" || !board.Idx.Valid || board.Idx.Int64 < 1 {
		return c.Status(http.StatusBadRequest).Send([]byte("user idx is wrong"))
	}

	member := model.Member{
		BoardIdx: null.IntFrom(board.Idx.Int64),
		UserIdx:  null.IntFrom(userIDX),
		Grade:    null.StringFrom(grade),
		RegDate:  null.StringFrom(now),
	}

	_, idx, err := crud.AddMember(member)
	if err != nil {
		log.Println("AddMemberAPI:", err)
	}

	idxSTR := strconv.FormatInt(idx, 10)
	result := map[string]string{
		"result": "ok",
		"idx":    idxSTR,
	}

	return c.Status(http.StatusOK).JSON(result)
}

func UpdateMemberAPI(c *fiber.Ctx) (err error) {
	datas := []map[string]interface{}{}
	datasSuccess := []map[string]interface{}{}
	datasFailed := []map[string]interface{}{}

	err = c.BodyParser(&datas)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	result := map[string]interface{}{"result": "ok"}
	if len(datasFailed) > 0 {
		result["result"] = "failed"
		result["failed"] = datasFailed
		result["success"] = datasSuccess
	}

	return c.Status(http.StatusOK).JSON(result)
}

func DeleteMemberAPI(c *fiber.Ctx) (err error) {
	datas := []map[string]interface{}{}
	datasSuccess := []map[string]interface{}{}
	datasFailed := []map[string]interface{}{}

	err = c.BodyParser(&datas)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	result := map[string]interface{}{"result": "ok"}
	if len(datasFailed) > 0 {
		result["result"] = "failed"
		result["failed"] = datasFailed
		result["success"] = datasSuccess
	}

	return c.Status(http.StatusOK).JSON(result)

}
