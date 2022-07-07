package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SetCookieHeader(t *testing.T) {
	recorder := httptest.NewRecorder()

	http.SetCookie(recorder, &http.Cookie{Name: "test", Value: "test_value"})

	request := &http.Request{Header: http.Header{"Cookie": recorder.HeaderMap["Set-Cookie"]}}
	cookie, err := request.Cookie("test")

	require.NoError(t, err, "failed to read 'test' Cookie: %v", err)
	require.Equal(t, cookie.Value, "test_value")
}

func Test_ExpireCookie(t *testing.T) {
	w := httptest.NewRecorder()
	ExpireCookie(w)
	r := w.Result()
	defer r.Body.Close()

	data := r.Cookies()
	require.Equal(t, data, []*http.Cookie{{Name: "token", MaxAge: 0, HttpOnly: true, Raw: "token=; HttpOnly"}})
}
