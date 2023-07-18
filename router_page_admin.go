package main

import (
	"9minutes/handler"

	"github.com/gofiber/fiber/v2"
)

func setPageAdmin(a *fiber.App) {
	a.Get("/admin/*", handler.HandleHTML)
}
