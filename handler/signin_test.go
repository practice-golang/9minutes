package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"9m/auth"
	"9m/router"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/stretchr/testify/require"
)

func Test_Signin(t *testing.T) {
	type args struct {
		c        *router.Context
		want     []byte
		run_func router.Handler
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_Signin",
			args: args{
				c: &router.Context{
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
					Request:        httptest.NewRequest("POST", "/signin", bytes.NewBuffer([]byte(`{"name": "test_user","password": "12345"}`))),
					Params:         []string{},
					AuthInfo:       nil,
				},
				want:     []byte(`"Signin success"`),
				run_func: Signin,
			},
		},
	}

	auth.SessionManager = scs.New()
	auth.SessionManager.Store = memstore.New()
	auth.SessionManager.Lifetime = 3 * time.Hour
	auth.SessionManager.IdleTimeout = 20 * time.Minute
	auth.SessionManager.Cookie.Name = "session_id"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newCTX, err := auth.SessionManager.Load(tt.args.c.Context(), "")
			if err != nil {
				t.Errorf("error loading from session manager: %v", err)
			}

			tt.args.c.Request = tt.args.c.WithContext(newCTX)
			tt.args.run_func(tt.args.c)

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

func Test_SigninAPI(t *testing.T) {
	type args struct {
		c        *router.Context
		want     []byte
		run_func router.Handler
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_SigninAPI",
			args: args{
				c: &router.Context{
					Request: httptest.NewRequest(
						"GET", "/signin",
						bytes.NewBuffer([]byte(`{"name": "test_user","password": "12345"}`))),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
				},
				want:     []byte("Signin success"),
				run_func: SigninAPI,
			},
		},
	}

	auth.GenerateRsaKeys()
	auth.SaveRsaKeys()
	err := auth.GenerateKeySet()
	if err != nil {
		panic(err)
	}

	defer os.Remove("private.key")
	defer os.Remove("public.key")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.run_func(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			result := make(map[string]string)
			err = json.Unmarshal(data, &result)
			if err != nil {
				t.Errorf("expected error to parse data %v", err)
			}

			require.Equal(t, tt.args.want, []byte(result["msg"]), "Not equal")
		})
	}
}
