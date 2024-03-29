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
	Idx          null.Int    `json:"idx"            db:"IDX"            npskip:"insert, update"`
	TopicIdx     null.Int    `json:"topic-idx"      db:"TOPIC_IDX"      npskip:"update"`
	Content      null.String `json:"content"        db:"CONTENT"`
	AuthorIdx    null.Int    `json:"author-idx"     db:"AUTHOR_IDX"     npskip:"update"`
	AuthorName   null.String `json:"author-name"    db:"AUTHOR_NAME"    npskip:"update"`
	AuthorIpFull null.String `json:"-"              db:"AUTHOR_IP"      npskip:"select, update, read, viewcount"`
	AuthorIP     null.String `json:"author-ip"      db:"AUTHOR_IP_CUT"  npskip:"update, viewcount"`
	EditPassword null.String `json:"edit-password"  db:"EDIT_PASSWORD"  npskip:"update, viewcount"`
	Files        null.String `json:"files"          db:"FILES"          npskip:"viewcount"`
	DeleteFiles  null.String `json:"delete-files"   db:"-"              npskip:"insert, select, update, read, viewcount"`
	RegDate      null.String `json:"regdate"        db:"REGDATE"        npskip:"update"`
}
