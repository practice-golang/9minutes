package auth

import "github.com/alexedwards/scs/v2"

const (
	_        = iota
	MEMSTORE // memstore
	REDIS    // redis
	ETCD     // etcd

	// Not use
	// MYSQL // MySQL
	// PG // Postgres
	// SQLSERVER // MS SQL server
	// SQLITE3 // SQLite3
)

type (
	// Not yet support auth
	SessionStoreInfo struct {
		StoreType int
		Address   string
		Port      string
	}
)

var (
	JwtPrivateKeyFileName = "private.key"
	JwtPublicKeyFileName  = "public.key"

	SessionManager *scs.SessionManager
)
