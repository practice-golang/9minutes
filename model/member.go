package model

import (
	"gopkg.in/guregu/null.v4"
)

type MemberRequest struct {
	Idx      null.Int    `json:"idx"          db:"IDX"          npskip:"insert, update"`
	BoardIdx null.Int    `json:"board-idx"    db:"BOARD_IDX"`
	UserIdx  null.Int    `json:"user-idx"     db:"USER_IDX"`
	Grade    null.String `json:"grade"        db:"GRADE"`
	RegDate  null.String `json:"regdate"      db:"REGDATE"      npskip:"update"`
}
