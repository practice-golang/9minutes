package model

import (
	"gopkg.in/guregu/null.v4"
)

type StoredFileInfo struct {
	Idx         null.Int    `json:"idx" db:"IDX"`
	FileName    null.String `json:"filename" db:"FILE_NAME"`
	StorageName null.String `json:"storage_name" db:"STORAGE_NAME"`
}

type FilePath struct {
	Path  null.String `json:"path"`
	Sort  null.String `json:"sort"`
	Order null.String `json:"order"`
}

type FileInfo struct {
	Name     null.String `json:"name"`
	Size     null.Int    `json:"size"`
	DateTime null.String `json:"datetime"`
	IsDir    null.Bool   `json:"isdir"`
}

type FileList struct {
	Path     null.String `json:"path"`
	FullPath null.String `json:"full-path"`
	Files    []FileInfo  `json:"files"`
}
