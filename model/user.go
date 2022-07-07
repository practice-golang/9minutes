package model

import (
	"reflect"

	"gopkg.in/guregu/null.v4"
)

// UserListingOptions - Search, page
type UserListingOptions struct {
	Search    null.String
	Page      null.Int
	ListCount null.Int
}

// UserPageData - User list
type UserPageData struct {
	// UserList    []UserData `json:"user-list"`
	UserList    interface{} `json:"user-list"`
	TotalPage   int         `json:"total-page"`
	CurrentPage int         `json:"current-page"`
}

type UserInfo struct {
	Name null.String `json:"name"`
	Age  null.Int    `json:"age"`
}

// UserColumn - Columns for user table
type UserColumn struct {
	Idx         null.Int    `json:"idx"          db:"IDX"          npskip:"insert, update"`
	DisplayName null.String `json:"display-name" db:"DISPLAY_NAME"`
	ColumnCode  null.String `json:"column-code"  db:"COLUMN_CODE"`
	ColumnType  null.String `json:"column-type"  db:"COLUMN_TYPE"`
	ColumnName  null.String `json:"column-name"  db:"COLUMN_NAME"`
	SortOrder   null.Int    `json:"sort-order"   db:"SORT_ORDER"   npskip:"insert"`
}

// UserData - User data
type UserData struct {
	Idx      null.Int    `json:"idx"      db:"IDX"      npskip:"insert, update"`
	UserName null.String `json:"username" db:"USERNAME"`
	Password null.String `json:"password" db:"PASSWORD"`
	Email    null.String `json:"email"    db:"EMAIL"`
	Grade    null.String `json:"grade"    db:"GRADE"`
	Approval null.String `json:"approval" db:"APPROVAL"`
	RegDTTM  null.String `json:"reg-dttm" db:"REG_DTTM" npskip:"update"`
}

var UserDataFieldCount = reflect.TypeOf(UserData{}).NumField()

type SignIn struct {
	Name     null.String `json:"name"     form:"username"`
	Password null.String `json:"password" form:"password"`
}

type AuthInfo struct {
	Name     null.String `json:"name"     mapstructure:"name"`
	IpAddr   null.String `json:"ip-addr"  mapstructure:"ip-addr"`
	Platform null.String `json:"platform" mapstructure:"platform"`
	Duration null.Int    `json:"duration" mapstructure:"duration"`
}
