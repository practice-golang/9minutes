package models

import (
	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/guregu/null.v4"
)

// User - (Minimal) User data
type User struct {
	Idx      null.String `json:"idx" db:"IDX" goqu:"skipinsert,skipupdate"`
	UserID   null.String `json:"user-id" db:"USER_ID"`
	Password null.String `json:"password" db:"PASSWORD"`
	Name     null.String `json:"name" db:"NAME"`
	Email    null.String `json:"email" db:"EMAIL"`
}

// Token - JWT token
type Token struct {
	Id      string
	Name    string
	IsAdmin string
	jwt.StandardClaims
}
