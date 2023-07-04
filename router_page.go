package main

import (
	"9minutes/handler"

	"github.com/gofiber/fiber/v2"
)

// func setPAGEs(r *router.App) {
// 	/* HTML */
// 	r.Handle(`^/?$`, handler.Index, "GET")
// 	r.Handle(`^/(index|logout).html$`, handler.HandleHTML, "GET")
// 	r.Handle(`assets/.*[css|js|map|woff|woff2]$`, handler.HandleAsset, "GET")

// 	glogin := r.Group(`^/`)
// 	glogin.Use(handler.AuthSessionMiddleware)
// 	glogin.GET(`login.html`, handler.HandleLogin)
// 	glogin.GET(`signup.html`, handler.HandleSignup)
// }

func setPage(a *fiber.App) {
	a.Get("/*", handler.HandleHTML)
}

// func setPageMyPage(r *router.App) {
// 	/* Mypage */
// 	gmypage := r.Group(`^/`)
// 	gmypage.Use(handler.RestrictSessionMiddleware)
// 	gmypage.GET(`mypage.html$`, handler.HandleHTML)
// }
