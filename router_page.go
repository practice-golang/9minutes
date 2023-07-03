package main

import (
	"9minutes/handler"
	"9minutes/router"

	"github.com/gofiber/fiber/v2"
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

func setPageHTMLs(a *fiber.App) {
	// /* Both user and anonymous */
	// gl := a.Group("/")
	// gl.Use(handler.AuthSessionMiddleware, handler.RemoveTrailingSlashMiddleware)
	// gl.Get(`[^/]+.html$`, handler.HandleHTML)
	// gl.Get(`assets/css/[^/]+.html$`, handler.HandleHTML)
	a.Static("/", StaticPath)
}

func setPageMyPage(r *router.App) {
	/* Mypage */
	gmypage := r.Group(`^/`)
	gmypage.Use(handler.RestrictSessionMiddleware)
	gmypage.GET(`mypage.html$`, handler.HandleHTML)
}
