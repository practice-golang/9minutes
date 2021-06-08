package models

import (
	"gopkg.in/guregu/null.v4"
)

// ContentsBasicBoard - Post fields of basic-board
type ContentsBasicBoard struct {
	Idx            null.String `json:"idx" db:"IDX" goqu:"skipinsert,skipupdate"`
	Title          null.String `json:"title" db:"TITLE"`
	Content        null.String `json:"content" db:"CONTENT"`
	WriterIdx      null.String `json:"writer-idx" db:"WRITER_IDX"`
	WriterName     null.String `json:"writer-name" db:"WRITER_NAME"`
	WriterPassword null.String `json:"writer-password" db:"WRITER_PASSWORD"`
	RegDTTM        null.String `json:"reg-dttm" db:"REG_DTTM" goqu:"skipinsert,skipupdate"`
}

// ContentsListBasicBoard - Post fields of basic-board
type ContentsListBasicBoard struct {
	Idx        null.String `json:"idx" db:"IDX" goqu:"skipinsert,skipupdate"`
	Title      null.String `json:"title" db:"TITLE"`
	WriterIdx  null.String `json:"writer-idx" db:"WRITER_IDX"`
	WriterName null.String `json:"writer-name" db:"WRITER_NAME"`
	RegDTTM    null.String `json:"reg-dttm" db:"REG_DTTM" goqu:"skipinsert,skipupdate"`
}

// Comment - Comment fields
type Comment struct {
	Idx       null.String `json:"idx" db:"IDX" goqu:"skipinsert,skipupdate"`
	PostIdx   null.String `json:"post-idx" db:"POST_IDX" goqu:"skipupdate"`
	WriterIdx null.String `json:"writer-idx" db:"WRITER_IDX" goqu:"skipupdate"`
	Content   null.String `json:"content" db:"CONTENT"`
	RegDTTM   null.String `json:"reg-dttm" db:"REG_DTTM" goqu:"skipupdate"`
}

// ContentSearch - Search contents for basic-board
type ContentSearch struct {
	Keywords []ContentsBasicBoard `json:"keywords" db:"-"` // Search keywords
	Options  Options              `json:"options" db:"-"`  // Paging, options for search (eg. and/or)
	Table    null.String          `json:"table" db:"-"`    // Target table name
}
