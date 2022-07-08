package model

import (
	"gopkg.in/guregu/null.v4"
)

// ContentListingOptions - Search, page
type ContentListingOptions struct {
	Search    null.String
	Page      null.Int
	ListCount null.Int
}

// ContentPageData - Contents list
type ContentPageData struct {
	ContentList []ContentList `json:"content-list"`
	PageList    []int         `json:"page-list"`
	TotalPage   int           `json:"total-page"`
	CurrentPage int           `json:"current-page"`
}

type ContentList struct {
	Idx          null.Int    `json:"idx"           db:"IDX"           npskip:"insert, update"`
	Title        null.String `json:"title"         db:"TITLE"`
	TitleImage   null.String `json:"title-image"   db:"TITLE_IMAGE"`
	AuthorIdx    null.Int    `json:"author-idx"    db:"AUTHOR_IDX"`
	AuthorName   null.String `json:"author-name"   db:"AUTHOR_NAME"   npskip:"insert, update, select, read"`
	CommentCount null.String `json:"comment-count" db:"COMMENT_COUNT" npskip:"insert, update, select, read"`
	Views        null.Int    `json:"views"         db:"VIEWS"`
	RegDTTM      null.String `json:"reg-dttm"      db:"REG_DTTM"`
}

type Content struct {
	Idx        null.Int    `json:"idx"         db:"IDX"         npskip:"insert, update, viewcount"`
	Title      null.String `json:"title"       db:"TITLE"       npskip:"viewcount"`
	TitleImage null.String `json:"title-image" db:"TITLE_IMAGE" npskip:"viewcount"`
	Content    null.String `json:"content"     db:"CONTENT"     npskip:"viewcount"`
	AuthorIdx  null.Int    `json:"author-idx"  db:"AUTHOR_IDX"  npskip:"viewcount"`
	AuthorName null.String `json:"author-name" db:"AUTHOR_NAME" npskip:"insert, update, select, read, viewcount"`
	Files      null.String `json:"files"       db:"FILES"       npskip:"viewcount"`
	Views      null.Int    `json:"views"       db:"VIEWS"       npskip:"update"`
	RegDTTM    null.String `json:"reg-dttm"    db:"REG_DTTM"    npskip:"update, viewcount"`
}

// CommentListingOptions - Search, page
type CommentListingOptions struct {
	Search    null.String
	Page      null.Int
	ListCount null.Int
}

// CommentPageData - Contents list
type CommentPageData struct {
	CommentList []Comment `json:"comment-list"`
	TotalPage   int       `json:"total-page"`
	CurrentPage int       `json:"current-page"`
	TotalCount  int       `json:"total-count"`
	ListCount   int       `json:"list-count"`
}

type Comment struct {
	Idx        null.Int    `json:"idx"         db:"IDX"         npskip:"insert, update"`
	BoardIdx   null.Int    `json:"board-idx"   db:"BOARD_IDX"`
	Content    null.String `json:"content"     db:"CONTENT"`
	AuthorIdx  null.Int    `json:"author-idx"  db:"AUTHOR_IDX"`
	AuthorName null.String `json:"author-name" db:"AUTHOR_NAME" npskip:"insert, update, select, read"`
	Files      null.String `json:"files"       db:"FILES"`
	RegDTTM    null.String `json:"reg-dttm"    db:"REG_DTTM"`
}
