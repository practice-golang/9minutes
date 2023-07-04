package main

import (
	"9minutes/handler"

	"github.com/gofiber/fiber/v2"
)

func setApiBoard(a *fiber.App) {
	/* API Board */
	// gabrd := a.Group("/api/board", handler.AuthApiSessionMiddleware)
	gabrd := a.Group("/api/board") // Require add session middleware
	gabrd.Post("/content/:board", handler.WriteContent)
	gabrd.Put("/content/:board/:idx", handler.UpdateContent)
	gabrd.Delete("/content/:board/:idx", handler.DeleteContent)
	gabrd.Get("/comment/:board/:idx", handler.GetComments)
	gabrd.Post("/comment/:board/:idx", handler.WriteComment)
	gabrd.Delete("/comment/:board/:commentidx", handler.DeleteComment)
}
