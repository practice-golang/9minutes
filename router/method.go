package router

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
	"regexp"
	"strings"

	"9m/logging"
)

var StaticServer Handler
var UploadServer Handler
var EmbedStaticServer Handler

func New(middleware ...Middleware) *App {
	app := &App{
		DefaultRoute: func(c *Context) {
			c.Text(http.StatusNotFound, "Not found")
		},
		MethodNotAllowed: func(c *Context) {
			c.Text(http.StatusNotFound, "Method not allowed")
		},
		Middlewares: middleware,
	}

	return app
}

func (a *App) Use(middlewares ...Middleware) {
	a.Middlewares = append(a.Middlewares, middlewares...)
}

func (a *App) Group(prefix string, middleware ...Middleware) *RouteGroup {
	group := &RouteGroup{
		App:         a,
		Prefix:      prefix,
		Middlewares: middleware,
	}

	return group
}

func (g *RouteGroup) Use(middlewares ...Middleware) {
	g.Middlewares = append(g.Middlewares, middlewares...)
}

func (g *RouteGroup) Handle(pattern string, handler Handler, methods ...string) {
	appMiddlewares := g.App.Middlewares
	g.App.Middlewares = append(g.App.Middlewares, g.Middlewares...)
	g.App.Handle(g.Prefix+pattern, handler, methods...)
	g.App.Middlewares = appMiddlewares
}

// Aliases - Bypass to Handle of Group
func (g *RouteGroup) GET(pattern string, handler Handler)     { g.Handle(pattern, handler, "GET") }
func (g *RouteGroup) HEAD(pattern string, handler Handler)    { g.Handle(pattern, handler, "HEAD") }
func (g *RouteGroup) POST(pattern string, handler Handler)    { g.Handle(pattern, handler, "POST") }
func (g *RouteGroup) PUT(pattern string, handler Handler)     { g.Handle(pattern, handler, "PUT") }
func (g *RouteGroup) DELETE(pattern string, handler Handler)  { g.Handle(pattern, handler, "DELETE") }
func (g *RouteGroup) CONNECT(pattern string, handler Handler) { g.Handle(pattern, handler, "CONNECT") }
func (g *RouteGroup) OPTIONS(pattern string, handler Handler) { g.Handle(pattern, handler, "OPTIONS") }
func (g *RouteGroup) TRACE(pattern string, handler Handler)   { g.Handle(pattern, handler, "TRACE") }
func (g *RouteGroup) PATCH(pattern string, handler Handler)   { g.Handle(pattern, handler, "PATCH") }

func (a *App) Handle(pattern string, handler Handler, methods ...string) {
	re := regexp.MustCompile(pattern)
	m := Methods{}

	for _, method := range methods {
		switch method {
		case "*":
			for _, method := range AllMethods {
				m[method] = true
			}
		default:
			m[strings.ToUpper(method)] = true
		}
	}

	route := Route{Pattern: re, Handler: handler, Methods: m, Middlewares: a.Middlewares}

	a.Routes = append(a.Routes, route)
}

// Aliases - Bypass to Handle of App
func (a *App) GET(pattern string, handler Handler)     { a.Handle(pattern, handler, "GET") }
func (a *App) HEAD(pattern string, handler Handler)    { a.Handle(pattern, handler, "HEAD") }
func (a *App) POST(pattern string, handler Handler)    { a.Handle(pattern, handler, "POST") }
func (a *App) PUT(pattern string, handler Handler)     { a.Handle(pattern, handler, "PUT") }
func (a *App) DELETE(pattern string, handler Handler)  { a.Handle(pattern, handler, "DELETE") }
func (a *App) CONNECT(pattern string, handler Handler) { a.Handle(pattern, handler, "CONNECT") }
func (a *App) OPTIONS(pattern string, handler Handler) { a.Handle(pattern, handler, "OPTIONS") }
func (a *App) TRACE(pattern string, handler Handler)   { a.Handle(pattern, handler, "TRACE") }
func (a *App) PATCH(pattern string, handler Handler)   { a.Handle(pattern, handler, "PATCH") }

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{Request: r, ResponseWriter: w}

	b, _ := ioutil.ReadAll(c.Body)
	c.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	logger := logging.Object.Log()
	if json.Valid(b) {
		bc := new(bytes.Buffer)
		json.Compact(bc, b)
		logger = logger.RawJSON("body", bc.Bytes())
	} else {
		logger = logger.Fields(map[string]interface{}{"body": bytes.ReplaceAll(b, []byte("\n"), []byte(""))})
	}

	logger.Timestamp().
		Str("method", c.Method).
		Str("path", c.URL.Path).
		Str("remote", c.RemoteAddr).
		Str("user-agent", c.UserAgent()).
		Fields(map[string]interface{}{"header": c.Request.Header}).
		Send()

	methodNotAllowed := false
	for _, rt := range a.Routes {
		matches := rt.Pattern.FindStringSubmatch(c.URL.Path)

		if len(matches) > 0 {
			if !rt.Methods[c.Method] {
				methodNotAllowed = true
				continue
			}

			if len(matches) > 1 {
				c.Params = matches[1:]
			}

			for _, m := range rt.Middlewares {
				rt.Handler = m(rt.Handler)
			}

			rt.Handler(c)
			return
		}
	}

	if methodNotAllowed {
		a.MethodNotAllowed(c)
	} else {
		a.DefaultRoute(c)
	}
}

func (c *Context) Text(code int, body string) {
	c.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	c.WriteHeader(code)

	c.ResponseWriter.Write([]byte(body))
}

func (c *Context) Json(code int, body interface{}) {
	c.ResponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	c.WriteHeader(code)

	result, err := json.Marshal(body)
	if err != nil {
		log.Println("Json error:", err)
		return
	}

	c.ResponseWriter.Write(result)
}

func (c *Context) Html(code int, body []byte) {
	c.ResponseWriter.Header().Set("Content-Type", "text/html")
	c.WriteHeader(code)

	c.ResponseWriter.Write(body)
}

func (c *Context) File(code int, body []byte) {
	c.ResponseWriter.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(c.URL.Path)))
	c.WriteHeader(code)

	c.ResponseWriter.Write(body)
}

// SetupStaticServer - Serving internal `embedded` static files
func SetupStaticServer() {
	var err error

	EmbedContent, err = fs.Sub(fs.FS(EmbedStatic), EmbedPath)
	if err != nil {
		logging.Object.Warn().Err(err).Msg("SetupStatic")
	}

	e := http.StripPrefix("/embed/", http.FileServer(http.FS(EmbedContent))) // embed storage
	EmbedStaticServer = func(c *Context) { e.ServeHTTP(c.ResponseWriter, c.Request) }

	s := http.StripPrefix("/static/", http.FileServer(http.Dir(StaticPath))) // real storage
	StaticServer = func(c *Context) { s.ServeHTTP(c.ResponseWriter, c.Request) }

	u := http.StripPrefix("/upload/", http.FileServer(http.Dir(UploadPath))) // real storage
	UploadServer = func(c *Context) { u.ServeHTTP(c.ResponseWriter, c.Request) }
}
