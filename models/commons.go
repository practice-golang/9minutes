package models

import (
	"gopkg.in/guregu/null.v4"
)

// Options - Paging options
type Options struct {
	Count null.Int    `json:"count"`
	Page  null.Int    `json:"page"`
	Order null.String `json:"order"`
}
