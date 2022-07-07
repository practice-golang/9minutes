package handler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"9m/auth"
	"9m/model"
	"9m/router"

	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func Test_HelloMiddleware(t *testing.T) {
	type args struct {
		c    *router.Context
		want []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_HelloMiddleware",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/api", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte("Middle ware test error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := router.New(HelloMiddleware)
			h.Handle("/api", func(c *router.Context) {}, "GET")
			h.ServeHTTP(tt.args.c.ResponseWriter, tt.args.c.Request)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, tt.args.want, data, "intend error")
		})
	}
}

func Test_AuthMiddleware(t *testing.T) {
	type args struct {
		c        *router.Context
		authinfo model.AuthInfo
		want     []byte
		run_func router.Middleware
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_AuthApiMiddleware",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/api", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte(""),
				authinfo: model.AuthInfo{
					Name:     null.StringFrom("test_name"),
					IpAddr:   null.StringFrom("192.168.0.1"),
					Platform: null.StringFrom("test_platform"),
					Duration: null.IntFrom(3600),
				},
				run_func: AuthApiMiddleware,
			},
		},
		{
			name: "Test_AuthApiMiddleware_error",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/api", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want:     []byte("Auth error"),
				run_func: AuthApiMiddleware,
			},
		},
		{
			name: "Test_AuthMiddleware",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/api", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want: []byte(""),
				authinfo: model.AuthInfo{
					Name:     null.StringFrom("test_name"),
					IpAddr:   null.StringFrom("192.168.0.1"),
					Platform: null.StringFrom("test_platform"),
					Duration: null.IntFrom(3600),
				},
				run_func: AuthMiddleware,
			},
		},
		{
			name: "Test_AuthMiddleware_error",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/api", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want:     []byte("Auth error"),
				run_func: AuthMiddleware,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Test_AuthApiMiddleware" {
				err := auth.GenerateRsaKeys()
				if err != nil {
					t.Errorf("GenerateKeys() error = %v", err)
					return
				}
				err = auth.GenerateKeySet()
				if err != nil {
					t.Errorf("GenerateKey() error = %v", err)
					return
				}
				token, err := auth.GenerateToken(tt.args.authinfo)
				if err != nil {
					t.Errorf("GenerateToken() error = %v", err)
					return
				}

				tt.args.c.Request.Header.Add("Authorization", "Bearer "+token)
			} else if tt.name == "Test_AuthMiddleware" {
				err := auth.GenerateKeySet()
				if err != nil {
					t.Errorf("GenerateKey() error = %v", err)
					return
				}
				token, err := auth.GenerateToken(tt.args.authinfo)
				if err != nil {
					t.Errorf("GenerateToken() error = %v", err)
					return
				}

				tt.args.c.Request.AddCookie(&http.Cookie{
					Name:  "token",
					Value: token,
				})
			}

			h := router.New(tt.args.run_func)
			h.Handle("/api", func(c *router.Context) {}, "GET")
			h.ServeHTTP(tt.args.c.ResponseWriter, tt.args.c.Request)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, tt.args.want, data, "Not equal")
		})
	}
}
