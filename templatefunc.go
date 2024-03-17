package server

import (
	"fmt"
	"html"
	"html/template"
	"strconv"
	"strings"
	"time"
)

func unEscape(s string) template.HTML {
	return template.HTML(html.UnescapeString(s))
}

func formatDate(s string) string {
	t, err := time.Parse("20060102150405", s)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func jsArray(s []int) string {
	return strings.Trim(strings.Replace(fmt.Sprint(s), " ", ",", -1), "[]")
}

func moreEQ(a, b int) bool {
	result := false
	if a >= b {
		result = true
	}

	return result
}

func lessEQ(a, b int) bool {
	result := false
	if a <= b {
		result = true
	}

	return result
}

func sliceValue(s []int, idx string) int {
	i, _ := strconv.Atoi(idx)
	return s[i]
}
