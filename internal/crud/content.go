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

func GetContentList(board model.Board, options model.ContentListingOptions) (model.ContentPageData, error) {
	result := model.ContentPageData{}

	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(board.BoardTable.String)
	userTableName := db.GetFullTableName(db.Info.UserTable)
	commentTableName := db.GetFullTableName(board.CommentTable.String)

	column := np.CreateString(model.ContentList{}, dbtype, "select", false)
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

	var list []model.ContentList
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

	result = model.ContentPageData{
		BoardCode:     board.BoardCode.String,
		SearchKeyword: options.Search.String,
		ContentList:   list,
		CurrentPage:   int(options.Page.Int64) + 1,
		TotalPage:     int(totalPage),
	}

	return result, err
}

func GetContent(board model.Board, idx string) (model.Content, error) {
	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(board.BoardTable.String)
	userTableName := db.GetFullTableName(db.Info.UserTable)

	column := np.CreateString(model.Content{}, dbtype, "select", false)
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
		return model.Content{}, err
	}
	defer r.Close()

	var content model.Content
	err = scan.Row(&content, r)
	if err != nil {
		return model.Content{}, err
	}

	return content, nil
}

// GetComment - get comment
func GetComment(board model.Board, boardIdx, commentIdx string) (model.Comment, error) {
	var comment model.Comment

	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(board.CommentTable.String)

	column := np.CreateString(model.Comment{}, dbtype, "select", false)
	columnBoardIdx := np.CreateString(map[string]interface{}{"BOARD_IDX": nil}, dbtype, "", false)
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)

	sql := `
	SELECT
		` + column.Names + `
	FROM ` + tableName + `
	WHERE ` + columnBoardIdx.Names + ` = ` + boardIdx + `
		AND ` + columnIdx.Names + ` = ` + commentIdx

	r, err := db.Con.Query(sql)
	if err != nil {
		return comment, err
	}
	defer r.Close()

	err = scan.Row(&comment, r)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

// GetComments - get comment list
func GetComments(board model.Board, contentIdx string, options model.CommentListingOptions) (model.CommentPageData, error) {
	result := model.CommentPageData{}

	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(board.CommentTable.String)
	userTableName := db.GetFullTableName(db.Info.UserTable)

	column := np.CreateString(model.Comment{}, dbtype, "select", false)

	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)
	columnUserId := np.CreateString(map[string]interface{}{"USERID": nil}, dbtype, "", false)
	columnBoardIdx := np.CreateString(map[string]interface{}{"BOARD_IDX": nil}, dbtype, "", false)
	columnAuthorIdx := np.CreateString(map[string]interface{}{"AUTHOR_IDX": nil}, dbtype, "", false)
	columnAuthorName := np.CreateString(map[string]interface{}{"AUTHOR_NAME": nil}, dbtype, "", false)

	var totalCount int64
	sql := `
	SELECT
		COUNT(` + columnIdx.Names + `)
	FROM ` + tableName + `
	WHERE ` + columnBoardIdx.Names + ` = ` + contentIdx

	r, err := db.Con.Query(sql)
	if err != nil {
		return result, err
	}
	defer r.Close()

	err = scan.Row(&totalCount, r)
	if err != nil {
		return result, err
	}

	totalPage := math.Ceil(float64(totalCount) / float64(options.ListCount.Int64))
	currentPage := int64(totalPage)

	offset := 0
	if totalCount > 0 {
		offset = int(int64(totalPage)*options.ListCount.Int64 - options.ListCount.Int64)

		if options.Page.Valid && options.Page.Int64 > -1 && options.ListCount.Valid {
			currentPage = options.Page.Int64 + 1
			offset = int(options.Page.Int64 * options.ListCount.Int64)
		}
	}

	paging := db.Obj.GetPagingQuery(offset, int(options.ListCount.Int64))

	sql = `
	SELECT
		` + column.Names + `,
		(
			SELECT
				` + columnUserId.Names + `
			FROM ` + userTableName + `
			WHERE ` + columnIdx.Names + ` = ` + tableName + `.` + columnAuthorIdx.Names + `
		) AS ` + columnAuthorName.Names + `
	FROM ` + tableName + `
	WHERE ` + columnBoardIdx.Names + ` = ` + contentIdx + `
	ORDER BY ` + columnIdx.Names + ` ASC
	` + paging

	r, err = db.Con.Query(sql)
	if err != nil {
		return result, err
	}
	defer r.Close()

	var comments []model.Comment
	err = scan.Rows(&comments, r)
	if err != nil {
		return result, err
	}

	result = model.CommentPageData{
		CommentList: comments,
		CurrentPage: int(currentPage),
		TotalPage:   int(totalPage),
	}

	return result, nil
}

func WriteContent(board model.Board, content model.Content) (sql.Result, error) {
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

func UpdateContent(board model.Board, content model.Content, skipTag string) error {
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

func DeleteContent(board model.Board, idx string) error {
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

func WriteComment(board model.Board, content model.Comment) error {
	tableName := db.GetFullTableName(board.CommentTable.String)

	column := np.CreateString(content, db.GetDatabaseTypeString(), "insert", false)

	sql := `
	INSERT INTO ` + tableName + ` (
		` + column.Names + `
	) VALUES (
		` + column.Values + `
	)`

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func DeleteComment(board model.Board, boardIdx, commentIdx string) error {
	tableName := db.GetFullTableName(board.CommentTable.String)

	columnBoardIdx := np.CreateString(map[string]interface{}{"BOARD_IDX": nil}, db.GetDatabaseTypeString(), "", false)
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, db.GetDatabaseTypeString(), "", false)

	sql := `
	DELETE
	FROM ` + tableName + `
	WHERE ` + columnBoardIdx.Names + ` = ` + boardIdx + `
		AND ` + columnIdx.Names + ` = ` + commentIdx

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func DeleteComments(board model.Board, idx string) error {
	tableName := db.GetFullTableName(board.CommentTable.String)

	columnBoardIdx := np.CreateString(map[string]interface{}{"BOARD_IDX": nil}, db.GetDatabaseTypeString(), "", false)

	sql := `
	DELETE FROM ` + tableName + `
	WHERE ` + columnBoardIdx.Names + ` = ` + idx

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
