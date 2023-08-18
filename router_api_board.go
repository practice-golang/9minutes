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

	/* API Posting */
	gbrd.Get("/:board_code", handler.ListPostingAPI)
	gbrd.Get("/:board_code/posting/:idx", handler.ReadPostingAPI)
	gbrd.Post("/:board_code/posting", handler.WritePostingAPI)
	gbrd.Put("/:board_code/posting/:idx", handler.UpdatePostingAPI)
	gbrd.Delete("/:board_code/posting/:idx", handler.DeletePostingAPI)

	/* API Comment */
	gbrd.Get("/:board_code/:idx/comment", handler.GetComments)
	gbrd.Post("/:board_code/:posting_idx/comment", handler.WriteComment)
	gbrd.Delete("/:board_code/:idx/comment/:comment_idx", handler.DeleteComment)
}
