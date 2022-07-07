package router

import (
	"embed"
	"io/fs"
)

var (
	StaticPath   string = "../static"
	UploadPath   string = "../upload"
	EmbedPath    string = "embed"
	Content      embed.FS
	EmbedStatic  embed.FS
	EmbedContent fs.FS
	AllMethods   = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
)
