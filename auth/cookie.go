package auth

import (
	"net/http"
)

func SetCookieHeader(w http.ResponseWriter, token string, duration int64) {
	// SameSite: http.SameSiteNoneMode,
	session := http.Cookie{
		Name:     "token",
		Value:    token,
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   int(duration),
		HttpOnly: true,
	}

	w.Header().Add("Set-Cookie", session.String())
}

func ExpireCookie(w http.ResponseWriter) {
	session := http.Cookie{
		Name:     "token",
		MaxAge:   0,
		HttpOnly: true,
	}

	w.Header().Add("Set-Cookie", session.String())
}

// GetCookie - Not use
// func GetCookie(r *http.Request) {
// 	token, err := r.Cookie("token")
// 	if err != nil {
// 		log.Println("GetCookie:", err)
// 	}
// 	cookies := r.Cookies()

// 	log.Println(token)
// 	log.Println(cookies)
// }
