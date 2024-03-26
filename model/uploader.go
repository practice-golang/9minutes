package model

import (
	"gopkg.in/guregu/null.v4"
)

type StoredFileInfo struct {
	Idx         null.Int    `json:"idx"           db:"IDX"`
	TopicIdx    null.String `json:"topic-idx"     db:"TOPIC_IDX"`
	CommentIdx  null.String `json:"comment-idx"   db:"COMMENT_IDX"`
	FileName    null.String `json:"filename"      db:"FILE_NAME"`
	StorageName null.String `json:"storage_name"  db:"STORAGE_NAME"`
	RegDate     null.String `json:"regdate"       db:"REGDATE"        npskip:"update"`
}

// Not use
type FilePath struct {
	Path  null.String `json:"path"`
	Sort  null.String `json:"sort"`
	Order null.String `json:"order"`
}

// Not use
type FileList struct {
	Path     null.String `json:"path"`
	FullPath null.String `json:"full-path"`
	Files    []FileInfo  `json:"files"`
}

// Not use
type FilesToDelete struct {
	DeleteFiles []File `json:"delete-files"`
}

// Not use
type FileInfo struct {
	Name     null.String `json:"name"`
	Size     null.Int    `json:"size"`
	DateTime null.String `json:"datetime"`
	IsDir    null.Bool   `json:"isdir"`
}

// Not use
type File struct {
	FileName   null.String `json:"filename"     db:"FILE_NAME"`
	TopicIdx   null.String `json:"topic-idx"    db:"TOPIC_IDX"`
	CommentIdx null.String `json:"comment-idx"  db:"COMMENT_IDX"`
	StoreName  null.String `json:"storename"    db:"STORAGE_NAME"`
}
