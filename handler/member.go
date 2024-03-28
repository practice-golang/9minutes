package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetMemberListAPI(c *fiber.Ctx) error {
	queries := c.Queries()

	log.Println(queries)

	return nil
}
