package model

import (
	"gopkg.in/guregu/null.v4"
)

// PostingListingOptions - Search, page
type PostingListingOptions struct {
	Search    null.String
	Page      null.Int
	ListCount null.Int
}

// PostingPageData - Contents list
type PostingPageData struct {
	BoardCode     string        `json:"board-code"`
	SearchKeyword string        `json:"search-keyword"`
	PostingList   []PostingList `json:"posting-list"`
	PageList      []int         `json:"page-list"`
	TotalPage     int           `json:"total-page"`
	CurrentPage   int           `json:"current-page"`
	JumpPrev      int           `json:"jump-prev"`
	JumpNext      int           `json:"jump-next"`
}

type PostingList struct {
	Idx          null.Int    `json:"idx"           db:"IDX"           npskip:"insert, update"`
	Title        null.String `json:"title"         db:"TITLE"`
	TitleImage   null.String `json:"title-image"   db:"TITLE_IMAGE"`
	AuthorIdx    null.Int    `json:"author-idx"    db:"AUTHOR_IDX"`
	AuthorName   null.String `json:"author-name"   db:"AUTHOR_NAME"   npskip:"insert, update, select, read"`
	CommentCount null.String `json:"comment-count" db:"COMMENT_COUNT" npskip:"insert, update, select, read"`
	Views        null.Int    `json:"views"         db:"VIEWS"`
	RegDate      null.String `json:"regdate"       db:"REGDATE"`
}

type Posting struct {
	Idx        null.Int    `json:"idx"         db:"IDX"         npskip:"insert, update, viewcount"`
	Title      null.String `json:"title"       db:"TITLE"       npskip:"viewcount"`
	TitleImage null.String `json:"title-image" db:"TITLE_IMAGE" npskip:"viewcount"`
	Content    null.String `json:"content"     db:"CONTENT"     npskip:"viewcount"`
	AuthorIdx  null.Int    `json:"author-idx"  db:"AUTHOR_IDX"  npskip:"viewcount"`
	AuthorName null.String `json:"author-name" db:"AUTHOR_NAME" npskip:"insert, update, select, read, viewcount"`
	Files      null.String `json:"files"       db:"FILES"       npskip:"viewcount"`
	Images     null.String `json:"images"      db:"IMAGES"      npskip:"viewcount"`
	Views      null.Int    `json:"views"       db:"VIEWS"       npskip:"update"`
	RegDate    null.String `json:"regdate"     db:"REGDATE"     npskip:"update, viewcount"`
}

type FilesToDelete struct {
	DeleteFiles []File `json:"delete-files"`
}

type File struct {
	FileName  null.String `json:"filename"  db:"FILE_NAME"`
	StoreName null.String `json:"storename" db:"STORAGE_NAME"`
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
}

type Comment struct {
	Idx        null.Int    `json:"idx"         db:"IDX"         npskip:"insert, update"`
	BoardIdx   null.Int    `json:"board-idx"   db:"BOARD_IDX"`
	Content    null.String `json:"content"     db:"CONTENT"`
	AuthorIdx  null.Int    `json:"author-idx"  db:"AUTHOR_IDX"`
	AuthorName null.String `json:"author-name" db:"AUTHOR_NAME" npskip:"insert, update, select, read"`
	Files      null.String `json:"files"       db:"FILES"`
	Images     null.String `json:"images"      db:"IMAGES"`
	RegDate    null.String `json:"regdate"     db:"REGDATE"`
}
