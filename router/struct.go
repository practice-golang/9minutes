package router

import (
	"net/http"
	"regexp"
)

type Handler func(*Context)
type Middleware func(Handler) Handler

type Methods map[string]bool

type Route struct {
	Pattern     *regexp.Regexp
	Handler     Handler
	Methods     Methods
	Middlewares []Middleware
}

type App struct {
	Routes           []Route
	DefaultRoute     Handler
	MethodNotAllowed Handler
	Middlewares      []Middleware
}

type RouteGroup struct {
	App         *App
	Prefix      string
	Middlewares []Middleware
}

type Context struct {
	http.ResponseWriter
	*http.Request
	Params   []string
	AuthInfo interface{}
}
