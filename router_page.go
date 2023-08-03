package main

import (
	"9minutes/handler"

	"github.com/gofiber/fiber/v2"
)

func setPage(a *fiber.App) {
	a.Get("/board/list", handler.HandleBoardHTML)
	a.Get("/*", handler.HandleHTML)
}
