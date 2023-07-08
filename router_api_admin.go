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

	/* API Admin - Boards */
	gaboard := a.Group("/api/admin/boards") // required add auth middleware
	gaboard.Get("/", handler.GetBoards)
	gaboard.Post("/", handler.AddBoard)
	gaboard.Put("/", handler.UpdateBoard)
	gaboard.Delete("/", handler.DeleteBoard)
}
