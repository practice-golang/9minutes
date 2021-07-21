package models

import (
	"gopkg.in/guregu/null.v4"
)

// ContentsBasicBoardSET - Content fields of basic-board
type ContentsBasicBoardSET struct {
	Idx            null.String `json:"idx" db:"IDX" goqu:"skipinsert,skipupdate"`
	Title          null.String `json:"title" db:"TITLE"`
	Content        null.String `json:"content" db:"CONTENT"`
	IsMember       null.String `json:"is-member" db:"IS_MEMBER" goqu:"skipupdate"`
	WriterIdx      null.String `json:"writer-idx" db:"WRITER_IDX" goqu:"skipupdate"`
	WriterName     null.String `json:"writer-name" db:"WRITER_NAME" goqu:"skipupdate"`
	WriterPassword null.String `json:"writer-password" db:"WRITER_PASSWORD" goqu:"skipupdate"`
	Files          null.String `json:"files" db:"FILES"`
	RegDTTM        null.String `json:"reg-dttm" db:"REG_DTTM" goqu:"skipinsert,skipupdate"`
}

// ContentsBasicBoardGET - Content fields of basic-board
type ContentsBasicBoardGET struct {
	Idx        null.String `json:"idx" db:"IDX" goqu:"skipinsert,skipupdate"`
	Title      null.String `json:"title" db:"TITLE"`
	Content    null.String `json:"content" db:"CONTENT"`
	IsMember   null.String `json:"is-member" db:"IS_MEMBER"`
	WriterName null.String `json:"writer-name" db:"WRITER_NAME"`
	Files      null.String `json:"files" db:"FILES"`
	RegDTTM    null.String `json:"reg-dttm" db:"REG_DTTM" goqu:"skipinsert,skipupdate"`
}

// ContentsBasicBoardList - Content fields of basic-board
type ContentsBasicBoardList struct {
	Idx        null.String `json:"idx" db:"IDX" goqu:"skipinsert,skipupdate"`
	Title      null.String `json:"title" db:"TITLE"`
	IsMember   null.String `json:"is-member" db:"IS_MEMBER"`
	WriterName null.String `json:"writer-name" db:"WRITER_NAME"`
	RegDTTM    null.String `json:"reg-dttm" db:"REG_DTTM" goqu:"skipinsert,skipupdate"`
}

// ContentSearch - Search contents for basic-board
type ContentSearch struct {
	Keywords []ContentsBasicBoardSET `json:"keywords" db:"-"` // Search keywords
	Options  Options                 `json:"options" db:"-"`  // Paging, options for search (eg. and/or)
	Table    null.String             `json:"table" db:"-"`    // Target table name
}
