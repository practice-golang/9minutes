package models

import (
	"gopkg.in/guregu/null.v4"
)

// CommentSET - Comment fields for both basic-board and custom-board
type CommentSET struct {
	// WriterIdx      null.String `json:"writer-idx" db:"WRITER_IDX"`
	Idx            null.String `json:"idx" db:"IDX" goqu:"skipinsert,skipupdate"`
	BIdx           null.String `json:"bidx" db:"BOARD_IDX" goqu:"skipupdate"`
	Content        null.String `json:"content" db:"CONTENT"`
	IsMember       null.String `json:"is-member" db:"IS_MEMBER" goqu:"skipupdate"`
	WriterName     null.String `json:"writer-name" db:"WRITER_NAME" goqu:"skipupdate"`
	WriterPassword null.String `json:"writer-password" db:"WRITER_PASSWORD" goqu:"skipupdate"`
	RegDTTM        null.String `json:"reg-dttm" db:"REG_DTTM" goqu:"skipinsert,skipupdate"`
}

// CommentList - Comment fields for both basic-board and custom-board
type CommentList struct {
	// WriterIdx  null.String `json:"writer-idx" db:"WRITER_IDX"`
	Idx        null.String `json:"idx" db:"IDX" goqu:"skipinsert,skipupdate"`
	BIdx       null.String `json:"bidx" db:"BOARD_IDX" goqu:"skipupdate"`
	Content    null.String `json:"content" db:"CONTENT"`
	IsMember   null.String `json:"is-member" db:"IS_MEMBER" goqu:"skipupdate"`
	WriterName null.String `json:"writer-name" db:"WRITER_NAME"`
	RegDTTM    null.String `json:"reg-dttm" db:"REG_DTTM" goqu:"skipinsert,skipupdate"`
}

// CommentSearch - Search contents for basic-board
type CommentSearch struct {
	Keywords []CommentSET `json:"keywords" db:"-"` // Search keywords
	Options  Options      `json:"options" db:"-"`  // Paging, options for search (eg. and/or)
	Table    null.String  `json:"table" db:"-"`    // Target table name
}
