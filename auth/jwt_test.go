package auth

import (
	"reflect"
	"testing"

	"9minutes/model"

	"gopkg.in/guregu/null.v4"
)

func Test_GenerateAndParseToken(t *testing.T) {
	authinfo := model.AuthInfo{
		Name:     null.StringFrom("test_name"),
		IpAddr:   null.StringFrom("192.168.0.1"),
		Os:       null.StringFrom("test_platform"),
		Duration: null.IntFrom(3600),
	}

	t.Run("GenerateToken", func(t *testing.T) {
		err := GenerateRsaKeys()
		if err != nil {
			t.Errorf("GenerateKey() error = %v", err)
			return
		}

		err = GenerateKeySet()
		if err != nil {
			t.Errorf("GenerateKey() error = %v", err)
			return
		}

		payloadSTR, err := GenerateToken(authinfo)
		if err != nil {
			t.Errorf("GenerateToken() error = %v", err)
			return
		}

		_, got, err := ParseToken(payloadSTR)
		if err != nil {
			t.Errorf("ParseToken() error = %v", err)
			return
		}

		if !reflect.DeepEqual(got, authinfo) {
			t.Errorf("ParseToken()\ngot = %v\nwant = %v", got, authinfo)
		}
	})
}
