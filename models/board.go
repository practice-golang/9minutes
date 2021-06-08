package models

import (
	"gopkg.in/guregu/null.v4"
)

// Board - Board data
type Board struct {
	Idx    null.String `json:"idx" db:"IDX" goqu:"skipinsert,skipupdate"`
	Name   null.String `json:"name" db:"NAME"`
	Code   null.String `json:"code" db:"CODE"`
	Type   null.String `json:"type" db:"TYPE"`
	Table  null.String `json:"table" db:"TABLE"`
	Fields interface{} `json:"fields" db:"FIELDS"`
	// Fields []Field `json:"fields" db:"FIELDS"`
}

// Field - Fields for custom board
type Field struct {
	Idx         null.Int    `json:"idx"`
	Name        null.String `json:"name"`
	ColumnName  null.String `json:"column"`
	Type        null.String `json:"type"`
	UseCommment null.Int    `json:"use-comment"` // 0: Not use comment, 1: Use comment
	Json        null.String `json:"json-name"`
	Order       null.Int    `json:"order"`
}

// BoardSearch - Search
type BoardSearch struct {
	Keywords []Board `json:"keywords" db:"-"` // Search keywords
	Options  Options `json:"options" db:"-"`  // Paging, options for search (eg. and/or)
}
