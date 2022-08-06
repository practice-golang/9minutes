package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"9minutes/model"
	"9minutes/router"

	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func Test_RestrictedHello(t *testing.T) {
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
			name: "Test_RestrictedHello",
			args: args{
				c: &router.Context{
					Request:        httptest.NewRequest("GET", "/api", nil),
					ResponseWriter: http.ResponseWriter(httptest.NewRecorder()),
					AuthInfo: model.AuthInfo{
						Name:     null.StringFrom("test_name"),
						IpAddr:   null.StringFrom("192.168.0.1"),
						Platform: null.StringFrom("test_platform"),
						Duration: null.IntFrom(3600),
					},
				},
				want:     []byte("Hello test_name"),
				run_func: RestrictedHello,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.run_func(tt.args.c)

			res := tt.args.c.ResponseWriter.(*httptest.ResponseRecorder).Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			require.Equal(t, tt.args.want, data, "Not equal")
		})
	}
}
