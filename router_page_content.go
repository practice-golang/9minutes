package main

import (
	"9minutes/handler"
	"9minutes/router"
)

func setPageContent(r *router.App) {
	/* Content - Read */
	gbrdr := r.Group(`^/board`)
	gbrdr.Use(handler.AuthSessionMiddleware)
	// gbrdr.GET(`/(list|gallery).html(\?[^\?]+)?$`, handler.HandleContentList)
	gbrdr.GET(`(\?[^\?]+)?$`, handler.HandleContentList)
	gbrdr.GET(`/read.html(\?[^\?]+)?$`, handler.HandleReadContent)

	/* Content - Edit, Write */
	gbew := r.Group(`^/board`)
	gbew.Use(handler.RestrictSessionMiddleware)
	gbew.GET(`/write.html(\?[^\?]+)?$`, handler.HandleWriteContent)
	gbew.GET(`/edit.html(\?[^\?]+)?$`, handler.HandleEditContent)
}
