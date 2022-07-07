package main

import (
	"9minutes/handler"
	"9minutes/router"
)

func setPageAdmin(r *router.App) {
	/* Admin */
	gadmin := r.Group(`^/admin`)
	gadmin.Use(handler.RestrictSessionMiddleware, handler.AdminOrManagerMiddleware, handler.RemoveTrailingSlashMiddleware)
	gadmin.GET(`(/?)$`, handler.AdminIndex)
	gadmin.GET(`/boards-list.html(\?[^\?]+)?$`, handler.HandleBoardList)
	gadmin.GET(`/users-list.html(\?[^\?]+)?$`, handler.HandleUserList)
	gadmin.GET(`/[^/]+.html$`, handler.HandleHTML)
}
