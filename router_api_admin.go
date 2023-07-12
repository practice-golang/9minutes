package main

import (
	"9minutes/handler"

	"github.com/gofiber/fiber/v2"
)

func checkAdmin(c *fiber.Ctx) error {
	usergrade, err := handler.GetSessionUserGrade(c)
	if err != nil {
		return err
	}

	if usergrade != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"status":  403,
			"message": "forbidden",
		})
	}

	return c.Next()
}

func setApiAdmin(a *fiber.App) {

	/* API Admin */
	gadmin := a.Group("/api/admin")
	gadmin.Use(checkAdmin)
	gadmin.Get("/health", handler.HealthCheck)

	/* API Admin - User fileds */
	gauserfield := gadmin.Group("/user-columns") // required add auth middleware
	gauserfield.Get("/", handler.GetUserColumns)
	gauserfield.Post("/", handler.AddUserColumn)
	gauserfield.Put("/", handler.UpdateUserColumns)
	gauserfield.Delete("/", handler.DeleteUserColumns)

	/* API Admin - Users */
	gauser := a.Group("/api/admin/user") // required add auth middleware
	gauser.Get("/", handler.GetUserList)
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
