package main

import (
	"9m/handler"
	"9m/router"
)

func setOthers(r *router.App) {
	/* Static */
	r.GET(`^/static/*`, router.StaticServer)
	r.GET(`^/upload/*`, router.UploadServer)
	r.GET(`^/embed/*`, router.EmbedStaticServer)

	/* Websocket - /ws.html */
	r.GET(`^/ws-echo`, handler.HandleWebsocketEcho)
	r.GET(`^/ws-chat`, handler.HandleWebsocketChat)
}
