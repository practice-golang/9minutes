package crud

import (
	"9minutes/internal/db"
	"9minutes/internal/np"
	"9minutes/model"
	"database/sql"
	"fmt"
	"html"
	"math"
	"strings"

	"github.com/blockloop/scan"
	"gopkg.in/guregu/null.v4"
)

func GetPostingList(board model.Board, options model.PostingListingOptions) (model.PostingPageData, error) {
	result := model.PostingPageData{}

	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(board.BoardTable.String)
	userTableName := db.GetFullTableName(db.Info.UserTable)
	commentTableName := db.GetFullTableName(board.CommentTable.String)

	column := np.CreateString(model.PostingList{}, dbtype, "select", false)
	sqlSearch := ""

	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)
	columnTitle := np.CreateString(map[string]interface{}{"TITLE": nil}, dbtype, "", false)
	columnContent := np.CreateString(map[string]interface{}{"CONTENT": nil}, dbtype, "", false)
	columnUserId := np.CreateString(map[string]interface{}{"USERID": nil}, dbtype, "", false)
	columnBoardIdx := np.CreateString(map[string]interface{}{"BOARD_IDX": nil}, dbtype, "", false)
	columnAuthorIdx := np.CreateString(map[string]interface{}{"AUTHOR_IDX": nil}, dbtype, "", false)
	columnAuthorName := np.CreateString(map[string]interface{}{"AUTHOR_NAME": nil}, dbtype, "", false)
	columnCommentCount := np.CreateString(map[string]interface{}{"COMMENT_COUNT": nil}, dbtype, "", false)

	if options.Search.Valid && options.Search.String != "" {
		sqlSearch = `
		WHERE LOWER(` + columnTitle.Names + `) LIKE LOWER('%` + options.Search.String + `%')
			OR LOWER(` + columnContent.Names + `) LIKE LOWER('%` + options.Search.String + `%')`
	}

	paging := ``
	if options.Page.Valid && options.ListCount.Valid {
		paging = db.Obj.GetPagingQuery(int(options.Page.Int64*options.ListCount.Int64), int(options.ListCount.Int64))
	}

	sql := `
	SELECT
		` + column.Names + `,
		(
			SELECT
				` + columnUserId.Names + `
			FROM ` + userTableName + `
			WHERE ` + columnIdx.Names + ` = A.` + columnAuthorIdx.Names + `
		) AS ` + columnAuthorName.Names + `,
		(
			SELECT
				COUNT(` + columnIdx.Names + `)
			FROM ` + commentTableName + `
			WHERE ` + columnBoardIdx.Names + ` = A.` + columnIdx.Names + `
		) AS ` + columnCommentCount.Names + `
	FROM ` + tableName + ` A
	` + sqlSearch + `
	ORDER BY ` + columnIdx.Names + ` DESC
	` + paging

	r, err := db.Con.Query(sql)
	if err != nil {
		return result, err
	}
	defer r.Close()

	var list []model.PostingList
	err = scan.Rows(&list, r)
	if err != nil {
		return result, err
	}

	var totalCount int64
	sql = `
	SELECT
		COUNT(` + columnIdx.Names + `)
	FROM ` + tableName + `
	` + sqlSearch

	r, err = db.Con.Query(sql)
	if err != nil {
		return result, err
	}
	defer r.Close()

	err = scan.Row(&totalCount, r)
	if err != nil {
		return result, err
	}

	totalPage := math.Ceil(float64(totalCount) / float64(options.ListCount.Int64))

	result = model.PostingPageData{
		BoardCode:     board.BoardCode.String,
		SearchKeyword: options.Search.String,
		PostingList:   list,
		CurrentPage:   int(options.Page.Int64) + 1,
		TotalPage:     int(totalPage),
	}

	return result, err
}

func GetPosting(board model.Board, idx string) (model.Posting, error) {
	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(board.BoardTable.String)
	userTableName := db.GetFullTableName(db.Info.UserTable)

	column := np.CreateString(model.Posting{}, dbtype, "select", false)
	columnUserId := np.CreateString(map[string]interface{}{"USERID": nil}, dbtype, "", false)
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)
	columnAuthorIdx := np.CreateString(map[string]interface{}{"AUTHOR_IDX": nil}, dbtype, "", false)
	columnAuthorName := np.CreateString(map[string]interface{}{"AUTHOR_NAME": nil}, dbtype, "", false)

	sql := `
	SELECT
		` + column.Names + `,
		(
			SELECT
				` + columnUserId.Names + `
			FROM ` + userTableName + `
			WHERE ` + columnIdx.Names + ` = ` + tableName + `.` + columnAuthorIdx.Names + `
		) AS ` + columnAuthorName.Names + `
	FROM ` + tableName + `
	WHERE ` + columnIdx.Names + ` = ` + idx

	r, err := db.Con.Query(sql)
	if err != nil {
		return model.Posting{}, err
	}
	defer r.Close()

	var content model.Posting
	err = scan.Row(&content, r)
	if err != nil {
		return model.Posting{}, err
	}

	return content, nil
}

func WritePosting(board model.Board, content model.Posting) (sql.Result, error) {
	tableName := db.GetFullTableName(board.BoardTable.String)

	content.Title = null.StringFrom(html.EscapeString(content.Title.String))
	content.Content = null.StringFrom(html.EscapeString(content.Content.String))

	column := np.CreateString(content, db.GetDatabaseTypeString(), "insert", false)

	sql := `
	INSERT INTO ` + tableName + ` (
		` + column.Names + `
	) VALUES (
		` + column.Values + `
	)`

	result, err := db.Con.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func UpdatePosting(board model.Board, content model.Posting, skipTag string) error {
	tableName := db.GetFullTableName(board.BoardTable.String)
	idx := fmt.Sprint(content.Idx.Int64)

	content.Title = null.StringFrom(html.EscapeString(content.Title.String))
	content.Content = null.StringFrom(html.EscapeString(content.Content.String))

	column := np.CreateString(content, db.GetDatabaseTypeString(), skipTag, false)
	colNames := strings.Split(column.Names, ",")
	colValues := strings.Split(column.Values, ",")

	holder := ""
	for i := 0; i < len(colNames); i++ {
		if colValues[i] == "''" {
			continue
		}
		holder += colNames[i] + " = " + colValues[i] + ", "
	}
	holder = strings.TrimSuffix(holder, ", ")

	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, db.GetDatabaseTypeString(), "", false)

	sql := `
	UPDATE ` + tableName + ` SET
		` + holder + `
	WHERE ` + columnIdx.Names + ` = ` + idx

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func DeletePosting(board model.Board, idx string) error {
	tableName := db.GetFullTableName(board.BoardTable.String)

	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, db.GetDatabaseTypeString(), "", false)

	sql := `
	DELETE FROM ` + tableName + `
	WHERE ` + columnIdx.Names + ` = ` + idx

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
