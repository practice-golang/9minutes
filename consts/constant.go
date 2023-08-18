package consts

var (
	ProgramName string = "9minutes"

	BcryptCost       int    = 8
	TableBoards      string = "boards"
	TableUsers       string = "users"
	TableUserColumns string = "user_fields"

	UserGrades = []string{
		"admin",
		"manager",
		"regular_user",
		"pending_user",
		"banned_user",
		"guest",
	}

	/* Message */
	MsgFailedToLogin   = `<html><script>alert('Failed to login');location.href="/login.html"</script></html>`
	MsgForbidden       = `<html><script>location.href="/auth/logout"</script></html>`
	MsgAlreadyLoggedin = `<html><script>location.href="/"</script></html>`

	MsgPasswordResetEmail        = `<html><script>alert('Password reset email have been sent');location.href="/"</script></html>`
	MsgPasswordResetUserNotFound = `<html><script>alert('User not found');location.href="/"</script></html>`
)
