package main

import (
	"9minutes/handler"

	"github.com/gofiber/fiber/v2"
)

func setApiAdmin(a *fiber.App) {
	/* API Admin */
	gadmin := a.Group("/api/admin")
	gadmin.Get("/health", handler.HealthCheck)
	// ad.POST(`/signin$`, handler.SigninAPI)

	/* API Admin - User fileds */
	gauserfield := a.Group("/api/admin/user-columns") // required add auth middleware
	gauserfield.Get("/", handler.GetUserColumns)
	gauserfield.Post("/", handler.AddUserColumn)
	gauserfield.Put("/", handler.UpdateUserColumns)
	gauserfield.Delete("/", handler.DeleteUserColumns)

	/* API Admin - Users */
	gauser := a.Group("/api/admin/user") // required add auth middleware
	gauser.Get("/:search?", handler.GetUsersList)
	gauser.Post("/", handler.AddUser)
	gauser.Put("/", handler.UpdateUser)
	gauser.Delete("/", handler.DeleteUser)

	// /* API Admin - Boards */
	// gab := a.Group(`^/api/admin/boards`, handler.AuthApiSessionMiddleware)
	// gab.GET(`/list$`, handler.GetBoards)
	// gab.POST(`/board$`, handler.AddBoard)
	// gab.PUT(`/board$`, handler.UpdateBoard)
	// gab.DELETE(`/board/[^/]+$`, handler.DeleteBoard)
}
