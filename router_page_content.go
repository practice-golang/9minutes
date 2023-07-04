package main

import (
	"9minutes/handler"

	"github.com/gofiber/fiber/v2"
)

func setPageContent(a *fiber.App) {
	/* Content - Read */
	gbread := a.Group(`^/board`) // Require add session middleware
	gbread.Get("/", handler.HandleContentList)
	gbread.Get("/read.html", handler.HandleReadContent)

	/* Content - Edit, Write */
	gbwrite := a.Group("/board") // Require add session middleware
	gbwrite.Get("/write.html", handler.HandleWriteContent)
	gbwrite.Get("/edit.html", handler.HandleEditContent)
}
