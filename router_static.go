package main

import (
	"github.com/gofiber/fiber/v2"
)

// setStatic - Set static
func setStatic(a *fiber.App) {
	a.Static("/files", HtmlPath)
}
