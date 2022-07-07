package main

import (
	"9minutes/handler"
	"9minutes/router"
)

func setApiAdmin(r *router.App) {
	/* API Admin */
	a := r.Group(`^/api/admin`)
	a.GET(`/?$`, handler.HealthCheck)
	a.POST(`/signin$`, handler.SigninAPI)

	/* API Admin restricted - Header */
	gapiadmin := r.Group(`^/api/admin`, handler.AuthApiSessionMiddleware)
	gapiadmin.GET(`/restricted$`, handler.RestrictedApiAdminHello)

	/* API Admin - User fileds */
	gauf := r.Group(`^/api/admin/user-columns`, handler.AuthApiSessionMiddleware)
	gauf.GET(`/list$`, handler.GetUserColumns)
	gauf.POST(`/column$`, handler.AddUserColumn)
	gauf.PUT(`/column$`, handler.UpdateUserColumn)
	gauf.DELETE(`/column/[^/]+$`, handler.DeleteUserColumn)

	/* API Admin - Users */
	gau := r.Group(`^/api/admin/users`, handler.AuthApiSessionMiddleware)
	gau.GET(`/list(\?[^\?]+)?$`, handler.GetUsersList)
	gau.POST(`/user$`, handler.AddUser)
	gau.PUT(`/user$`, handler.UpdateUser)
	gau.DELETE(`/user/[^/]+$`, handler.DeleteUser)

	/* API Admin - Boards */
	gab := r.Group(`^/api/admin/boards`, handler.AuthApiSessionMiddleware)
	gab.GET(`/list$`, handler.GetBoards)
	gab.POST(`/board$`, handler.AddBoard)
	gab.PUT(`/board$`, handler.UpdateBoard)
	gab.DELETE(`/board/[^/]+$`, handler.DeleteBoard)
}
