package models

import (
	"gopkg.in/guregu/null.v4"
)

// Comments - Comment fields for both basic-board and custom-board
type Comments struct {
	Idx            null.String `json:"idx" db:"IDX" goqu:"skipinsert,skipupdate"`
	BIdx           null.String `json:"bidx" db:"BOARD_IDX" goqu:"skipupdate"`
	Content        null.String `json:"content" db:"CONTENT"`
	WriterIdx      null.String `json:"writer-idx" db:"WRITER_IDX"`
	WriterName     null.String `json:"writer-name" db:"WRITER_NAME"`
	WriterPassword null.String `json:"writer-password" db:"WRITER_PASSWORD"`
	RegDTTM        null.String `json:"reg-dttm" db:"REG_DTTM" goqu:"skipinsert,skipupdate"`
}

// CommentSearch - Search contents for basic-board
type CommentSearch struct {
	Keywords []Comments  `json:"keywords" db:"-"` // Search keywords
	Options  Options     `json:"options" db:"-"`  // Paging, options for search (eg. and/or)
	Table    null.String `json:"table" db:"-"`    // Target table name
}
