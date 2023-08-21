package model

import (
	"gopkg.in/guregu/null.v4"
)

// CommentListingOptions - Search, page
type CommentListingOptions struct {
	Search    null.String
	Page      null.Int
	ListCount null.Int
}

// CommentPageData - Contents list
type CommentPageData struct {
	CommentList []Comment `json:"comment-list"`
	PageList    []int     `json:"page-list"`
	TotalPage   int       `json:"total-page"`
	CurrentPage int       `json:"current-page"`
	JumpPrev    int       `json:"jump-prev"`
	JumpNext    int       `json:"jump-next"`
}

type Comment struct {
	Idx        null.Int    `json:"idx"         db:"IDX"         npskip:"insert, update"`
	BoardIdx   null.Int    `json:"board-idx"   db:"BOARD_IDX"` // TODO - remove
	Content    null.String `json:"content"     db:"CONTENT"`
	AuthorIdx  null.Int    `json:"author-idx"  db:"AUTHOR_IDX"`
	AuthorName null.String `json:"author-name" db:"AUTHOR_NAME" npskip:"update"`
	Files      null.String `json:"files"       db:"FILES"`
	Images     null.String `json:"images"      db:"IMAGES"` // TODO - remove
	RegDate    null.String `json:"regdate"     db:"REGDATE"`
}
