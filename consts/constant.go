package consts

import (
	"9minutes/model"
)

var (
	ProgramName string = "9minutes"
	SiteName    string = "9mSite"

	BcryptCost       int    = 8
	TableBoards      string = "boards"
	TableUploads     string = "uploads"
	TableUsers       string = "users"
	TableUserColumns string = "user_fields"
	TableMembers     string = "members"

	UserGrades = map[string]model.UserGrade{
		"admin":       {Name: "Admin", Code: "admin", Rank: 1},
		"user_active": {Name: "Active user", Code: "user_active", Rank: 200},
		"guest":       {Name: "Guest", Code: "guest", Rank: 500},
		"user_hold":   {Name: "Holding user", Code: "user_hold", Rank: 600},
		"user_quit":   {Name: "Quit user", Code: "user_quit", Rank: 700},
		"user_banned": {Name: "Banned user", Code: "user_banned", Rank: 1000},
	}

	BoardGrades = map[string]model.UserGrade{
		"admin":         {Name: "Admin", Code: "admin", Rank: 1},
		"board_manager": {Name: "Board manager", Code: "board_manager", Rank: 100},
		"board_member":  {Name: "Board member", Code: "board_member", Rank: 150},
		"user_active":   {Name: "Active user", Code: "user_active", Rank: 200},
		"guest":         {Name: "Guest", Code: "guest", Rank: 500},
		"user_hold":     {Name: "Holding user", Code: "user_hold", Rank: 600},
		"user_quit":     {Name: "Quit user", Code: "user_quit", Rank: 700},
		"user_banned":   {Name: "Banned user", Code: "user_banned", Rank: 1000},
	}

	/* Message */
	// Not use yet
	// MsgFailedToLogin   = `<html><script>alert('Failed to login');location.href="/login.html"</script></html>`
	// MsgForbidden       = `<html><script>location.href="/auth/logout"</script></html>`
	// MsgAlreadyLoggedin = `<html><script>location.href="/"</script></html>`

	MsgPasswordResetEmail        = `<html><script>alert('Password reset email have been sent');location.href="/"</script></html>`
	MsgPasswordResetUserNotFound = `<html><script>alert('User not found');location.href="/"</script></html>`
)
