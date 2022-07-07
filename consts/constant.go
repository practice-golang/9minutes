package consts

var (
	ProgramName string = "9m"

	BcryptCost       int    = 8
	TableBooks       string = "books"
	TableBoards      string = "boards"
	TableUsers       string = "users"
	TableUserColumns string = "user_fields"

	/* Message */
	MsgFailedToLogin   = `<html><script>alert('Failed to login');location.href="/login.html"</script></html>`
	MsgForbidden       = `<html><script>location.href="/auth/logout"</script></html>`
	MsgAlreadyLoggedin = `<html><script>location.href="/"</script></html>`
)
