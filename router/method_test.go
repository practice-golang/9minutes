package router

import (
	"embed"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed embed_test/*
var fncStatic embed.FS

func Test_Router(t *testing.T) {
	type args struct{ c *Context }

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Test_Router_api",
			args: args{
				c: &Context{
					Request:        httptest.NewRequest("GET", "/api", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
			want: []byte("Ok"),
		},
		{
			name: "Test_Router_root",
			args: args{
				c: &Context{
					Request:        httptest.NewRequest("GET", "/", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
			want: []byte(`"Ok"`),
		},
		{
			name: "Test_Router_not_found",
			args: args{
				c: &Context{
					Request:        httptest.NewRequest("GET", "/not-found", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
			want: []byte("Not found"),
		},
		{
			name: "Test_Router_method_not_allowed",
			args: args{
				c: &Context{
					Request:        httptest.NewRequest("POST", "/", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
			want: []byte("Method not allowed"),
		},
		{
			name: "Test_Router_html",
			args: args{
				c: &Context{
					Request:        httptest.NewRequest("GET", "/index.html", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
			want: []byte("Ok HTML"),
		},
		{
			name: "Test_Router_file",
			args: args{
				c: &Context{
					Request:        httptest.NewRequest("GET", "/assets/js/not-exist-but-show-ok.js", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
			want: []byte("Ok File"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			g := r.Group(`^/api`)
			g.Handle(`/?$`, func(c *Context) { c.Text(http.StatusOK, "Ok") }, "GET")
			r.Handle(`^/?$`, func(c *Context) { c.Json(http.StatusOK, "Ok") }, "GET")
			r.Handle(`^/index.html?$`, func(c *Context) { c.Html(http.StatusOK, []byte("Ok HTML")) }, "GET")
			r.Handle(`^/assets/.*[css|js|map|woff|woff2]$`, func(c *Context) { c.File(http.StatusOK, []byte("Ok File")) }, "GET")
			r.ServeHTTP(tt.args.c.ResponseWriter, tt.args.c.Request)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, tt.want, data, "not equal")
		})
	}
}

func Test_Static(t *testing.T) {
	type args struct{ c *Context }
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Test_Static_embed",
			args: args{
				c: &Context{
					Request:        httptest.NewRequest("GET", "/embed/static_test.txt", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
			want: []byte("Hello static"),
		},
		{
			name: "Test_Static_storage",
			args: args{
				c: &Context{
					Request:        httptest.NewRequest("GET", "/static/hello.md", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
			want: []byte("# `Hello world!`"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StaticPath = "../static"
			EmbedPath = "embed_test"
			EmbedStatic = fncStatic

			SetupStaticServer()

			r := New()
			r.Handle("^/embed/.*$", EmbedStaticServer, "GET")
			r.Handle("^/static/.*$", StaticServer, "GET")
			r.ServeHTTP(tt.args.c.ResponseWriter, tt.args.c.Request)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, tt.want, data, "not equal")
		})
	}
}
