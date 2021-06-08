package models

import (
	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/guregu/null.v4"
)

// User - User(basic-user) data
type User struct {
	Idx      null.String `json:"idx" db:"IDX" goqu:"skipinsert,skipupdate"`
	Id       null.String `json:"id" db:"ID"`
	Password null.String `json:"Password" db:"PASSWORD"`
	Name     null.String `json:"name" db:"NAME"`
	Email    null.String `json:"email" db:"EMAIL"`
}

// Token - JWT token
type Token struct {
	Id    string
	Name  string
	Admin string
	jwt.StandardClaims
}
