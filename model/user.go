package model

import (
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

type SignIn struct {
	UserID   null.String `json:"userid"   form:"userid"`
	Password null.String `json:"password" form:"password"`
}

type AuthInfo struct {
	Name       null.String `json:"name"        mapstructure:"name"`
	IpAddr     null.String `json:"ip-addr"     mapstructure:"ip-addr"`
	Device     null.String `json:"device"      mapstructure:"device"`
	DeviceType null.String `json:"device-type" mapstructure:"device-type"`
	Os         null.String `json:"os"          mapstructure:"os"`
	Browser    null.String `json:"browser"     mapstructure:"browser"`
	Duration   null.Int    `json:"duration"    mapstructure:"duration"`
}
