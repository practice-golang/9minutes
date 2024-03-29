package handler

import (
	"9minutes/model"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetMemberListAPI(c *fiber.Ctx) error {
	queries := c.Queries()

	log.Println(queries)

	return c.Status(http.StatusOK).JSON("result")
}

func AddMemberAPI(c *fiber.Ctx) (err error) {
	var member model.MemberRequest
	now := time.Now().Format("20060102150405")

	// data := make(map[string]interface{})
	err = c.BodyParser(&member)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	log.Println(now)

	result := map[string]string{"result": "ok"}

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
