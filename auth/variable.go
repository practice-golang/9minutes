package auth

import "github.com/alexedwards/scs/v2"

var (
	JwtPrivateKeyFileName = "private.key"
	JwtPublicKeyFileName  = "public.key"

	SessionManager *scs.SessionManager
)
