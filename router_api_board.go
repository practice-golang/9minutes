package main

import (
	"9minutes/handler"
	"9minutes/router"
)

func setApiBoard(r *router.App) {
	/* API Board */
	gabrd := r.Group(`^/api/board`, handler.AuthApiSessionMiddleware)
	gabrd.POST(`/content/[^/]+$`, handler.WriteContent)
	gabrd.PUT(`/content/[^/]+/[^/]+$`, handler.UpdateContent)
	gabrd.DELETE(`/content/[^/]+/[^/]+$`, handler.DeleteContent)
	gabrd.GET(`/comment/[^/]+/[^/]+(\?[^\?]+)?$`, handler.GetComments)
	gabrd.POST(`/comment/[^/]+/[^/]+$`, handler.WriteComment)
	gabrd.DELETE(`/comment/[^/]+/[^/]+/[^/]+$`, handler.DeleteComment)
}
