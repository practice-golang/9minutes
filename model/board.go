package model

import (
	"gopkg.in/guregu/null.v4"
)

// BoardListingOption - Search, page
type BoardListingOption struct {
	Search    null.String
	Page      null.Int
	ListCount null.Int
}

// BoardPageData - Board list
type BoardPageData struct {
	BoardList   []Board `json:"board-list"`
	TotalPage   int     `json:"total-page"`
	CurrentPage int     `json:"current-page"`
}

// Board - Board data
type Board struct {
	Idx          null.Int    `json:"idx"           db:"IDX"           npskip:"insert,update,delete"`
	BoardName    null.String `json:"board-name"    db:"BOARD_NAME"`
	BoardCode    null.String `json:"board-code"    db:"BOARD_CODE"`
	BoardType    null.String `json:"board-type"    db:"BOARD_TYPE"`
	BoardTable   null.String `json:"board-table"   db:"BOARD_TABLE"`
	CommentTable null.String `json:"comment-table" db:"COMMENT_TABLE"`
	GrantRead    null.String `json:"grant-read"    db:"GRANT_READ"`
	GrantWrite   null.String `json:"grant-write"   db:"GRANT_WRITE"`
	GrantComment null.String `json:"grant-comment" db:"GRANT_COMMENT"`
	GrantUpload  null.String `json:"grant-upload"  db:"GRANT_UPLOAD"`
	Fields       interface{} `json:"fields"        db:"FIELDS"`
	// Fields []Field `json:"fields" db:"FIELDS"`
}

// Field - Fields for custom board
type Field struct {
	Idx         null.Int    `json:"idx"`
	Name        null.String `json:"name"`
	ColumnName  null.String `json:"column"`
	Type        null.String `json:"type"`
	UseCommment null.Int    `json:"use-comment"` // 0: Not use comment, 1: Use comment
	Json        null.String `json:"json-name"`
	Order       null.Int    `json:"order"`
}
