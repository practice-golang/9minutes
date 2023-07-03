package main

import (
	"9minutes/handler"

	"github.com/gofiber/fiber/v2"
)

func setApiAdmin(a *fiber.App) {
	/* API Admin */
	ad := a.Group("/api/admin")
	ad.Get("/health", handler.HealthCheck)
	// ad.POST(`/signin$`, handler.SigninAPI)

	// /* API Admin restricted - Header */
	// gapiad := a.Group(`^/api/admin`, handler.AuthApiSessionMiddleware)
	// gapiad.GET(`/restricted$`, handler.RestrictedApiAdminHello)

	// /* API Admin - User fileds */
	// gauf := a.Group(`^/api/admin/user-columns`, handler.AuthApiSessionMiddleware)
	// gauf.GET(`/list$`, handler.GetUserColumns)
	// gauf.POST(`/column$`, handler.AddUserColumn)
	// gauf.PUT(`/column$`, handler.UpdateUserColumn)
	// gauf.DELETE(`/column/[^/]+$`, handler.DeleteUserColumn)

	// /* API Admin - Users */
	// gau := a.Group(`^/api/admin/users`, handler.AuthApiSessionMiddleware)
	// gau.GET(`/list(\?[^\?]+)?$`, handler.GetUsersList)
	// gau.POST(`/user$`, handler.AddUser)
	// gau.PUT(`/user$`, handler.UpdateUser)
	// gau.DELETE(`/user/[^/]+$`, handler.DeleteUser)

	// /* API Admin - Boards */
	// gab := a.Group(`^/api/admin/boards`, handler.AuthApiSessionMiddleware)
	// gab.GET(`/list$`, handler.GetBoards)
	// gab.POST(`/board$`, handler.AddBoard)
	// gab.PUT(`/board$`, handler.UpdateBoard)
	// gab.DELETE(`/board/[^/]+$`, handler.DeleteBoard)
}
