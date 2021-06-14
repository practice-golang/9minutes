package models

import (
	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/guregu/null.v4"
)

// UserColumn - Fields for user table
type UserColumn struct {
	Idx        null.Int    `json:"idx,omitempty" db:"IDX" goqu:"skipinsert,skipupdate"`
	Name       null.String `json:"name" db:"NAME"`
	Code       null.String `json:"code" db:"CODE"`
	Type       null.String `json:"type" db:"TYPE"`
	ColumnName null.String `json:"column" db:"COLUMN_NAME"`
	Order      null.Int    `json:"order" db:"ORDER"`
}

// Token - JWT token
type Token struct {
	Id      string
	Name    string
	IsAdmin string
	jwt.StandardClaims
}
