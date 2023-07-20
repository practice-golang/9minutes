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
	gadmin.Get("/health", handler.HealthCheckAPI)

	/* API Admin - User fileds */
	gauserfield := gadmin.Group("/user-columns") // required add auth middleware
	gauserfield.Get("/", handler.GetUserColumnsAPI)
	gauserfield.Post("/", handler.AddUserColumnAPI)
	gauserfield.Put("/", handler.UpdateUserColumnsAPI)
	gauserfield.Delete("/", handler.DeleteUserColumnsAPI)

	/* API Admin - Users */
	gauser := gadmin.Group("/user") // required add auth middleware
	gauser.Get("/", handler.GetUserListAPI)
	gauser.Post("/", handler.AddUserAPI)
	gauser.Put("/", handler.UpdateUserAPI)
	gauser.Delete("/", handler.DeleteUserAPI)

	/* API Admin - Boards */
	gaboard := gadmin.Group("/board") // required add auth middleware
	gaboard.Get("/", handler.GetBoardsAPI)
	gaboard.Post("/", handler.AddBoardAPI)
	gaboard.Put("/", handler.UpdateBoardAPI)
	gaboard.Delete("/", handler.DeleteBoardAPI)
}
