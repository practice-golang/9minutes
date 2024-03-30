package handler

import (
	"9minutes/internal/crud"
	"9minutes/model"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/guregu/null.v4"
)

func GetMemberListAPI(c *fiber.Ctx) error {
	queries := c.Queries()

	log.Println(queries)

	return c.Status(http.StatusOK).JSON("result")
}

func AddMemberAPI(c *fiber.Ctx) (err error) {
	var member model.MemberRequest

	err = c.BodyParser(&member)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	now := time.Now().Format("20060102150405")
	member.RegDate = null.StringFrom(now)

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
