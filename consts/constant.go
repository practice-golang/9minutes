package consts

import (
	"9minutes/model"
)

var (
	ProgramName string = "9minutes"
	SiteName    string = "9mSite"

	BcryptCost       int    = 8
	TableBoards      string = "boards"
	TableUsers       string = "users"
	TableUserColumns string = "user_fields"

	UserGrades = map[string]model.UserGrade{
		"admin":         {Name: "Admin", Code: "admin", Point: 100},
		"board_manager": {Name: "Board manager", Code: "board_manager", Point: 200},
		"board_member":  {Name: "Board member", Code: "board_member", Point: 300},
		"regular_user":  {Name: "Regular user", Code: "regular_user", Point: 400},
		"guest":         {Name: "Guest", Code: "guest", Point: 500},
		"pending_user":  {Name: "Pending user", Code: "pending_user", Point: 600},
		"expired_user":  {Name: "Expired user", Code: "expired_user", Point: 700},
		"resigned_user": {Name: "Resigned user", Code: "resigned_user", Point: 800},
		"banned_user":   {Name: "Banned user", Code: "banned_user", Point: 900},
	}

	/* Message */
	// Not use yet
	// MsgFailedToLogin   = `<html><script>alert('Failed to login');location.href="/login.html"</script></html>`
	// MsgForbidden       = `<html><script>location.href="/auth/logout"</script></html>`
	// MsgAlreadyLoggedin = `<html><script>location.href="/"</script></html>`

	MsgPasswordResetEmail        = `<html><script>alert('Password reset email have been sent');location.href="/"</script></html>`
	MsgPasswordResetUserNotFound = `<html><script>alert('User not found');location.href="/"</script></html>`
)
