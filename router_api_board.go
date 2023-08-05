package main

import (
	"9minutes/handler"

	"github.com/gofiber/fiber/v2"
)

func setApiBoard(a *fiber.App) {
	/* API Board */
	gbrd := a.Group("/api/board") // Require add session middleware

	/* API Board list */
	gbrd.Get("/list", handler.BoardListAPI)

	/* API Content */
	gbrd.Get("/:board_code", handler.ListContentAPI)
	gbrd.Get("/:board_code/content/:idx", handler.ReadContentAPI)
	gbrd.Post("/:board_code/content", handler.WriteContentAPI)
	gbrd.Put("/:board_code/content/:idx", handler.UpdateContentAPI)
	gbrd.Delete("/:board_code/content/:idx", handler.DeleteContentAPI)

	/* API Comment */
	gbrd.Get("/:board_code/:idx/comment", handler.GetComments)
	gbrd.Post("/:board_code/:content_idx/comment", handler.WriteComment)
	gbrd.Delete("/:board_code/:idx/comment/:comment_idx", handler.DeleteComment)
}
