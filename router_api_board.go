package main

import (
	"9minutes/handler"

	"github.com/gofiber/fiber/v2"
)

func setApiBoard(a *fiber.App) {
	/* API Board */
	gbrd := a.Group("/api/board") // Require add session middleware

	/* API Content */
	gbrd.Get("/content/:board_code", handler.ListContentAPI)
	gbrd.Get("/content/:board_code/:idx", handler.ReadContentAPI)
	gbrd.Post("/content/:board_code", handler.WriteContentAPI)
	gbrd.Put("/content/:board_code/:idx", handler.UpdateContent)
	gbrd.Delete("/content/:board_code/:idx", handler.DeleteContent)

	/* API Comment */
	gbrd.Get("/comment/:board/:idx", handler.GetComments)
	gbrd.Post("/comment/:board/:idx", handler.WriteComment)
	gbrd.Delete("/comment/:board/:commentidx", handler.DeleteComment)
}
