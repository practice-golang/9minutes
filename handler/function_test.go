package handler

import (
	"9minutes/fd"
	"9minutes/model"
	"9minutes/router"
	"bytes"
	"embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"

	// "github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/websocket"
)

//go:embed embed_test/*
var fncEMBED embed.FS

func Test_Index(t *testing.T) {
	type args struct{ c *router.Context }
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_Index",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/index.html", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
		},
		{
			name: "Test_Index_embed",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
		},
		{
			name: "Test_Index_notfound",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storeRootBackup := StoreRoot
			embedRootBackup := EmbedRoot
			compareFilePath := "../html/index.html"
			want, err := os.ReadFile(compareFilePath)
			if err != nil {
				t.Error("Reference file htm not found")
			}

			if tt.name == "Test_Index_embed" {
				StoreRoot = "./not-found"
				EmbedRoot = "embed_test"
				router.Content = fncEMBED
				// tt.args.c.URL.Path = "/index.html"
				want = []byte("Hello embedded world")
			}
			if tt.name == "Test_Index_notfound" {
				StoreRoot = "./not-found"
				EmbedRoot = "not-found"
				want = []byte("File not found")
			}
			// HandleHTML(tt.args.c)
			Index(tt.args.c)

			StoreRoot = storeRootBackup
			EmbedRoot = embedRootBackup

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			patternLinkLogout = `\$LinkLogout\$(.*)\n`
			patternLinkAdmin = `\$LinkAdmin\$(.*)\n`
			reLogout = regexp.MustCompile(patternLinkLogout)
			reAdmin = regexp.MustCompile(patternLinkAdmin)
			patternIncludes = `@INCLUDE@(.*)(\n|$)`
			reIncludes = regexp.MustCompile(patternIncludes)

			/* Include */
			// includes := map[string][]byte{}
			m := reIncludes.FindAllSubmatch(want, -1)
			for _, v := range m {
				includeFileName := string(v[1])
				includeDirective := bytes.TrimSpace(v[0])
				includeStoreFilePath := strings.TrimSpace(StoreRoot + "/" + includeFileName)
				includeEmbedFilePath := strings.TrimSpace(EmbedRoot + "/" + includeFileName)

				include := []byte{}
				switch true {
				case fd.CheckFileExists(includeStoreFilePath, false):
					include, _ = os.ReadFile(includeStoreFilePath)
				case fd.CheckFileExists(includeEmbedFilePath, true):
					include, _ = router.Content.ReadFile(includeEmbedFilePath)
				}

				want = bytes.ReplaceAll(want, includeDirective, include)
			}

			patternLinkLogin = `\$LinkLogin\$(.*)\n`
			patternLinkLogout = `\$LinkLogout\$(.*)\n`
			patternLinkAdmin = `\$LinkAdmin\$(.*)\n`
			patternLinkMyPage = `\$LinkMyPage\$(.*)\n`
			patternYouArePending = `\$YouArePending\$(.*)\n`
			reLogin = regexp.MustCompile(patternLinkLogin)
			reLogout = regexp.MustCompile(patternLinkLogout)
			reAdmin = regexp.MustCompile(patternLinkAdmin)
			reMyPage = regexp.MustCompile(patternLinkMyPage)
			reYouArePending = regexp.MustCompile(patternYouArePending)

			want = reLogout.ReplaceAll(want, []byte(""))
			want = reAdmin.ReplaceAll(want, []byte(""))
			want = reMyPage.ReplaceAll(want, []byte(""))
			want = reYouArePending.ReplaceAll(want, []byte(""))

			want = bytes.ReplaceAll(want, []byte("$USERNAME$"), []byte("Guest"))
			want = bytes.ReplaceAll(want, []byte("$LinkLogin$"), []byte(""))
			want = bytes.ReplaceAll(want, []byte("$LinkLogout$"), []byte(""))
			want = bytes.ReplaceAll(want, []byte("$LinkAdmin$"), []byte(""))
			want = bytes.ReplaceAll(want, []byte("$LinkMyPage$"), []byte(""))
			want = reLogout.ReplaceAll(want, []byte(""))
			want = reAdmin.ReplaceAll(want, []byte(""))

			require.Equal(t, string(want), string(data), tt.name+" not equal"+" embed_test / "+tt.args.c.URL.Path)
		})
	}
}

func Test_HealthCheck(t *testing.T) {
	type args struct{ c *router.Context }
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_HealthCheck",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/api", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HealthCheck(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, []byte("Ok"), data, "HealthCheck not equal Ok")
		})
	}
}

func Test_Hello(t *testing.T) {
	type args struct {
		c    *router.Context
		want []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_Hello_API",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/api/hello", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("Hello world GET"),
			},
		},
		{
			name: "Test_Hello_GET",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/hello", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("Hello world GET"),
			},
		},
		{
			name: "Test_Hello_POST",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("POST", "/hello", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("Hello world POST"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Hello(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, tt.args.want, data, "Hello not equal %v", string(data))
		})
	}
}

func Test_HelloParam(t *testing.T) {
	type args struct {
		c    *router.Context
		want []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_HelloParam",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/hello", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
					Params:         []string{"test_name"},
				},
				want: []byte("Hello test_name"),
			},
		},
		{
			name: "Test_HelloParam_no_params",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/hello", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("Missing parameter"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HelloParam(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, tt.args.want, data, "HelloParam not equal")
		})
	}
}

func Test_GetParam(t *testing.T) {
	type args struct {
		c    *router.Context
		want []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_GetParam",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/get-param?hello=world", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("hello=world\n"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetParam(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, tt.args.want, data, "GetParam not equal")
		})
	}
}

func Test_PostForm(t *testing.T) {
	type args struct {
		c    *router.Context
		want []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_PostForm",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest(http.MethodGet, "/post-form", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("Hello world GET"),
			},
		},
		{
			name: "Test_PostForm",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest(http.MethodPost, "/post-form", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("hello=world\n"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.c.Request.Method == http.MethodPost {
				dat := url.Values{
					"hello": []string{"world"},
				}
				tt.args.c.Request.PostForm = dat
			}

			PostForm(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, tt.args.want, data, "PostForm not equal")
		})
	}
}

func Test_PostJson(t *testing.T) {
	jsonBody := map[string]interface{}{
		"name": "Thomas",
		"age":  "42",
	}
	body, err := json.Marshal(jsonBody)
	if err != nil {
		t.Errorf("WTWTWT expected error to be nil got %v", err)
	}

	type args struct {
		c    *router.Context
		want []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_PostJson",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest(http.MethodPost, "/post-json", bytes.NewReader(body)),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte(`{"name":"Thomas","age":42}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PostJson(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, tt.args.want, data, "PostJson not equal")
		})
	}
}

func Test_HandleAsset(t *testing.T) {
	listRendererJS, _ := os.ReadFile("../html/assets/js/list-renderer.js")
	type args struct {
		c    *router.Context
		want []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_HandleAsset",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/assets/js/list-renderer.js", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: listRendererJS,
			},
		},
		{
			name: "Test_HandleAsset_embed",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/index.html", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("Hello embedded world"),
			},
		},
		{
			name: "Test_HandleAsset_notfound",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/not-found", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				// want: []byte("Not found"),
				want: []byte(""),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storeRootBackup := StoreRoot
			embedRootBackup := EmbedRoot

			if tt.name == "Test_HandleAsset_embed" {
				StoreRoot = "./not-found"
				EmbedRoot = "embed_test"
				router.Content = fncEMBED
				tt.args.c.URL.Path = "/index.html"
			}

			HandleAsset(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()

			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			StoreRoot = storeRootBackup
			EmbedRoot = embedRootBackup

			require.Equal(t, tt.args.want, data, tt.name+" not equal")
		})
	}
}

func Test_WebsocketEcho(t *testing.T) {
	t.Run("WebsocketEcho", func(t *testing.T) {
		serverHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := &router.Context{Request: r, ResponseWriter: w}
			HandleWebsocketEcho(c)
		})
		s := httptest.NewServer(serverHandler)
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")

		w, err := websocket.Dial(u, "", s.URL)
		if err != nil {
			t.Errorf("Dial error %v", err)
		}

		msg := []byte("Hello")

		i, err := w.Write(msg)
		if err != nil {
			t.Errorf("Write error %v", err)
		}

		require.Equal(t, len(msg), i, "WebsocketEcho not equal")
		log.Println(i)
	})
}

func Test_WebsocketChat(t *testing.T) {
	t.Run("WebsocketChat", func(t *testing.T) {
		serverHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := &router.Context{Request: r, ResponseWriter: w}
			HandleWebsocketChat(c)
		})
		s := httptest.NewServer(serverHandler)
		defer s.Close()

		u := "ws" + strings.TrimPrefix(s.URL, "http")

		w, err := websocket.Dial(u, "", s.URL)
		if err != nil {
			t.Errorf("Dial error %v", err)
		}

		msg := []byte("Hello")

		i, err := w.Write(msg)
		if err != nil {
			t.Errorf("Write error %v", err)
		}

		require.Equal(t, len(msg), i, "WebsocketChat not equal")
	})
}

func TestHandleGetDir(t *testing.T) {
	type args struct {
		c        *router.Context
		jsonBody map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Test_GetDir_name_desc",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest(http.MethodPost, "/api/dir/list", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				jsonBody: map[string]interface{}{"path": "..", "sort": "name", "order": "desc"},
			},
			want: []byte(`
			{"path":"..",
			"full-path":"C:\\Users",
			"files":[
				{"name":"setup.go","size":3546,"datetime":"2022-01-11 04:47:53","isdir":false},
				{"name":"router_page_content.go","size":3445,"datetime":"2022-01-11 04:55:55","isdir":false},
				{"name":"router_page_admin.go","size":3445,"datetime":"2022-01-11 04:55:55","isdir":false},
				{"name":"router_page.go","size":3445,"datetime":"2022-01-11 04:55:55","isdir":false},
				{"name":"router_others.go","size":1388,"datetime":"2022-01-11 04:47:53","isdir":false},
				{"name":"router_notuse.go","size":1388,"datetime":"2022-01-11 04:47:53","isdir":false},
				{"name":"router_api_board.go","size":1388,"datetime":"2022-01-11 04:47:53","isdir":false},
				{"name":"router_api_admin.go","size":1388,"datetime":"2022-01-11 04:47:53","isdir":false},
				{"name":"router_api.go","size":1388,"datetime":"2022-01-11 04:47:53","isdir":false},
				{"name":"requests_user_fields.http","size":1048,"datetime":"2022-01-11 04:47:53","isdir":false},
				{"name":"requests_board.http","size":1048,"datetime":"2022-01-11 04:47:53","isdir":false},
				{"name":"requests.http","size":1048,"datetime":"2022-01-11 04:47:53","isdir":false},
				{"name":"main_test.go","size":1048,"datetime":"2022-01-11 04:47:53","isdir":false},
				{"name":"main.go","size":1048,"datetime":"2022-01-11 04:47:53","isdir":false},
				{"name":"go.sum","size":22179,"datetime":"2022-01-11 04:42:54","isdir":false},
				{"name":"go.mod","size":2200,"datetime":"2022-01-11 04:47:46","isdir":false},
				{"name":"delete_all_container.cmd","size":224,"datetime":"2022-01-10 20:06:41","isdir":false}
			]}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.args.jsonBody)
			tt.args.c.Request = httptest.NewRequest(http.MethodPost, "/api/dir/list", bytes.NewReader(body))
			HandleGetDir(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			var want model.FileList
			err = json.Unmarshal(tt.want, &want)
			if err != nil {
				t.Error(err)
			}

			var got model.FileList
			err = json.Unmarshal(data, &got)
			if err != nil {
				t.Error(err)
			}

			require.Equal(t, want.Path, got.Path, "GetDir not equal")
			for i, v := range want.Files {
				require.Equal(t, v.Name, got.Files[i].Name, "GetDir not equal")
				if v.Name.String == "go.mod" {
					break
				}
			}
		})
	}
}
