package auth

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"9m/model"

	"gopkg.in/guregu/null.v4"
)

func Test_SetupCookieToken(t *testing.T) {
	type args struct {
		w        http.ResponseWriter
		authinfo model.AuthInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "SetupCookieToken",
			args: args{
				w: http.ResponseWriter(httptest.NewRecorder()),
				authinfo: model.AuthInfo{
					Name:     null.StringFrom("test_name"),
					IpAddr:   null.StringFrom("192.168.1.1"),
					Platform: null.StringFrom("test_platform"),
					Duration: null.IntFrom(3600),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GenerateRsaKeys()
			if err != nil {
				t.Errorf("GenerateKeys() error = %v", err)
				return
			}
			err = GenerateKeySet()
			if err != nil {
				t.Errorf("GenerateKeySet() error = %v", err)
				return
			}

			err = SetupCookieToken(tt.args.w, tt.args.authinfo)
			if err != nil {
				t.Errorf("SetupCookieToken() error = %v", err)
			}
		})
	}
}

func Test_GetClaim(t *testing.T) {
	type args struct {
		r_cookie    http.Request
		r_header    http.Request
		authinfo    model.AuthInfo
		from_cookie string
		from_header string
	}
	tests := []struct {
		name    string
		args    args
		want    model.AuthInfo
		wantErr bool
	}{
		{
			name: "GetClaim",
			args: args{
				r_cookie: *httptest.NewRequest(http.MethodGet, "/test", nil),
				r_header: *httptest.NewRequest(http.MethodGet, "/test", nil),
				authinfo: model.AuthInfo{
					Name:     null.StringFrom("test_name"),
					IpAddr:   null.StringFrom("192.168.0.1"),
					Platform: null.StringFrom("test_platform"),
					Duration: null.IntFrom(3600),
				},
				from_cookie: "cookie",
				from_header: "header",
			},
			want: model.AuthInfo{
				Name:     null.StringFrom("test_name"),
				IpAddr:   null.StringFrom("192.168.0.1"),
				Platform: null.StringFrom("test_platform"),
				Duration: null.IntFrom(3600),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GenerateRsaKeys()
			if err != nil {
				t.Errorf("GenerateKeys() error = %v", err)
				return
			}
			err = GenerateKeySet()
			if err != nil {
				t.Errorf("GenerateKey() error = %v", err)
				return
			}
			token, err := GenerateToken(tt.args.authinfo)
			if err != nil {
				t.Errorf("GenerateToken() error = %v", err)
				return
			}

			tt.args.r_cookie.AddCookie(&http.Cookie{
				Name:  "token",
				Value: token,
			})
			tt.args.r_header.Header.Add("Authorization", "Bearer "+token)

			gotCookie, err := GetClaim(tt.args.r_cookie, tt.args.from_cookie)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClaim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCookie, tt.want) {
				t.Errorf("GetClaim()\nerror = %v\nwants = %v", gotCookie, tt.want)
			}

			tt.args.r_cookie.Header.Del("Cookie")
			gotCookieErr, err := GetClaim(tt.args.r_cookie, tt.args.from_cookie)
			if !reflect.DeepEqual(err.Error(), "http: named cookie not present") {
				t.Errorf("GetClaim()\nerror = %v\nwants = %v", gotCookieErr, tt.want)
			}

			gotHeader, err := GetClaim(tt.args.r_header, tt.args.from_header)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClaim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHeader, tt.want) {
				t.Errorf("GetClaim()\nerror = %v\nwants = %v", gotHeader, tt.want)
			}

			tt.args.r_header.Header.Set("Authorization", "Bearer "+token+"causeError")
			gotHeaderErr1, _ := GetClaim(tt.args.r_header, tt.args.from_header)
			if !reflect.DeepEqual(gotHeaderErr1, model.AuthInfo{}) {
				t.Errorf("GetClaim()\nerror = %v\nwants = %v", gotHeaderErr1, tt.want)
			}

			tt.args.r_header.Header.Del("Authorization")
			gotHeaderErr, err := GetClaim(tt.args.r_header, tt.args.from_header)
			if !reflect.DeepEqual(err.Error(), "bearer not found") {
				t.Errorf("GetClaim()\nerror = %v\nwants = %v", gotHeaderErr, tt.want)
			}
		})
	}
}
