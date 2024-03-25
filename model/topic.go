package model

import (
	"gopkg.in/guregu/null.v4"
)

// TopicListingOption - Search, page
type TopicListingOption struct {
	Search    null.String
	Page      null.Int
	ListCount null.Int
}

// TopicPageData - Contents list
type TopicPageData struct {
	BoardCode     string      `json:"board-code"`
	SearchKeyword string      `json:"search-keyword"`
	TopicList     []TopicList `json:"topic-list"`
	PageList      []int       `json:"page-list"`
	TotalPage     int         `json:"total-page"`
	CurrentPage   int         `json:"current-page"`
	JumpPrev      int         `json:"jump-prev"`
	JumpNext      int         `json:"jump-next"`
	ListCount     int         `json:"list-count"`
}

type TopicList struct {
	Idx          null.Int    `json:"idx"            db:"IDX"            npskip:"insert, update"`
	Title        null.String `json:"title"          db:"TITLE"`
	TitleImage   null.String `json:"title-image"    db:"TITLE_IMAGE"`
	AuthorIdx    null.Int    `json:"author-idx"     db:"AUTHOR_IDX"     npskip:"insert, update"`
	AuthorName   null.String `json:"author-name"    db:"AUTHOR_NAME"    npskip:"insert, update"`
	AuthorIpFull null.String `json:"-"              db:"AUTHOR_IP"`
	AuthorIP     null.String `json:"author-ip"      db:"AUTHOR_IP_CUT"  npskip:"update, viewcount"`
	CommentCount null.String `json:"comment-count"  db:"COMMENT_COUNT"  npskip:"insert, update, select, read"`
	Views        null.Int    `json:"views"          db:"VIEWS"`
	RegDate      null.String `json:"regdate"        db:"REGDATE"`
}

type Topic struct {
	Idx          null.Int    `json:"idx"            db:"IDX"            npskip:"insert, update, viewcount"`
	Title        null.String `json:"title"          db:"TITLE"          npskip:"viewcount"`
	TitleImage   null.String `json:"title-image"    db:"TITLE_IMAGE"    npskip:"viewcount"`
	Content      null.String `json:"content"        db:"CONTENT"        npskip:"viewcount"`
	AuthorIdx    null.Int    `json:"author-idx"     db:"AUTHOR_IDX"     npskip:"update, viewcount"`
	AuthorName   null.String `json:"author-name"    db:"AUTHOR_NAME"    npskip:"update, viewcount"`
	AuthorIpFull null.String `json:"-"              db:"AUTHOR_IP"      npskip:"select, update, read, viewcount"`
	AuthorIP     null.String `json:"author-ip"      db:"AUTHOR_IP_CUT"  npskip:"update, viewcount"`
	EditPassword null.String `json:"edit-password"  db:"EDIT_PASSWORD"  npskip:"update, viewcount"`
	Files        null.String `json:"files"          db:"FILES"          npskip:"viewcount"`
	DeleteFiles  null.String `json:"delete-files"   db:"-"              npskip:"insert, select, update, read, viewcount"`
	Views        null.Int    `json:"views"          db:"VIEWS"          npskip:"update"`
	RegDate      null.String `json:"regdate"        db:"REGDATE"        npskip:"update, viewcount"`
}
