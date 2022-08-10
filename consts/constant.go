package consts

var (
	ProgramName string = "9minutes"

	// TableBooks       string = "books"
	BcryptCost       int    = 8
	TableBoards      string = "boards"
	TableUsers       string = "users"
	TableUserColumns string = "user_fields"

	/* Message */
	MsgFailedToLogin   = `<html><script>alert('Failed to login');location.href="/login.html"</script></html>`
	MsgForbidden       = `<html><script>location.href="/auth/logout"</script></html>`
	MsgAlreadyLoggedin = `<html><script>location.href="/"</script></html>`
)
