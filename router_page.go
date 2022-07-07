package main

import (
	"9minutes/handler"
	"9minutes/router"
)

func setPAGEs(r *router.App) {
	/* HTML */
	r.Handle(`^/?$`, handler.Index, "GET")
	r.Handle(`^/(index|logout).html$`, handler.HandleHTML, "GET")
	r.Handle(`assets/.*[css|js|map|woff|woff2]$`, handler.HandleAsset, "GET")

	glogin := r.Group(`^/`)
	glogin.Use(handler.AuthSessionMiddleware)
	glogin.GET(`login.html`, handler.HandleLogin)
	glogin.GET(`signup.html`, handler.HandleSignup)
}

func setPageHTMLs(r *router.App) {
	/* Both user and anonymous */
	gl := r.Group(`^/`)
	gl.Use(handler.AuthSessionMiddleware, handler.RemoveTrailingSlashMiddleware)
	gl.GET(`[^/]+.html$`, handler.HandleHTML)
	gl.GET(`assets/css/[^/]+.html$`, handler.HandleHTML)
}

func setPageMyPage(r *router.App) {
	/* Mypage */
	gmypage := r.Group(`^/`)
	gmypage.Use(handler.RestrictSessionMiddleware)
	gmypage.GET(`mypage.html$`, handler.HandleHTML)
}
