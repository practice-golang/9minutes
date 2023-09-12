package consts

var (
	ProgramName string = "9minutes"
	SiteName    string = "9mSite"

	BcryptCost       int    = 8
	TableBoards      string = "boards"
	TableUsers       string = "users"
	TableUserColumns string = "user_fields"

	UserGrades = map[string]int{
		"admin":        1000,
		"manager":      2000,
		"regular_user": 3000,
		"pending_user": 4000,
		"guest":        5000,
		"banned_user":  6000,
	}

	/* Message */
	// Not use yet
	// MsgFailedToLogin   = `<html><script>alert('Failed to login');location.href="/login.html"</script></html>`
	// MsgForbidden       = `<html><script>location.href="/auth/logout"</script></html>`
	// MsgAlreadyLoggedin = `<html><script>location.href="/"</script></html>`

	MsgPasswordResetEmail        = `<html><script>alert('Password reset email have been sent');location.href="/"</script></html>`
	MsgPasswordResetUserNotFound = `<html><script>alert('User not found');location.href="/"</script></html>`
)
