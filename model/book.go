package model

import "gopkg.in/guregu/null.v4"

// Book
type Book struct {
	Idx    null.Int    `json:"idx"    db:"IDX"    npskip:"insert,update,delete"` // IDX
	Title  null.String `json:"title"  db:"TITLE"`                                // TITLE
	Author null.String `json:"author" db:"AUTHOR"`                               // AUTHOR
}
