package crud

import (
	"9minutes/internal/db"
	"9minutes/internal/np"
	"9minutes/model"
	"math"
	"strconv"
	"strings"

	"github.com/blockloop/scan"
)

// GetComment - get comment
func GetComment(board model.Board, postingIdx, commentIdx string) (model.Comment, error) {
	var comment model.Comment

	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(board.CommentTable.String)

	column := np.CreateString(model.Comment{}, dbtype, "select", false)
	columnPostingIdx := np.CreateString(map[string]interface{}{"POSTING_IDX": nil}, dbtype, "", false)
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)

	sql := `
	SELECT
		` + column.Names + `
	FROM ` + tableName + `
	WHERE ` + columnPostingIdx.Names + ` = ` + postingIdx + `
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
	// userTableName := db.GetFullTableName(db.Info.UserTable)

	column := np.CreateString(model.Comment{}, dbtype, "select", false)

	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)
	columnPostingIdx := np.CreateString(map[string]interface{}{"POSTING_IDX": nil}, dbtype, "", false)

	var totalCount int64
	sql := `
	SELECT
		COUNT(` + columnIdx.Names + `)
	FROM ` + tableName + `
	WHERE ` + columnPostingIdx.Names + ` = ` + contentIdx

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
		` + column.Names + `
	FROM ` + tableName + `
	WHERE ` + columnPostingIdx.Names + ` = ` + contentIdx + `
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

func WriteComment(board model.Board, content model.Comment) error {
	tableName := db.GetFullTableName(board.CommentTable.String)

	content.Content.String = EscapeString(content.Content.String)

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

func UpdateComment(board model.Board, comment model.Comment, postingIdx string) error {
	tableName := db.GetFullTableName(board.CommentTable.String)
	commentIDX := strconv.Itoa(int(comment.Idx.Int64))

	comment.Content.String = EscapeString(comment.Content.String)

	column := np.CreateString(comment, db.GetDatabaseTypeString(), "update", false)
	colNames := strings.Split(column.Names, ",")
	colValues := strings.Split(column.Values, ",")

	holder := ""
	for i := 0; i < len(colNames); i++ {
		if colNames[i] != `"FILES"` && colValues[i] == "''" {
			continue
		}
		holder += colNames[i] + " = " + colValues[i] + ", "
	}
	holder = strings.TrimSuffix(holder, ", ")

	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, db.GetDatabaseTypeString(), "", false)

	sql := `
	UPDATE ` + tableName + ` SET
		` + holder + `
	WHERE ` + columnIdx.Names + ` = ` + commentIDX

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func DeleteComment(board model.Board, postingIdx, commentIdx string) error {
	tableName := db.GetFullTableName(board.CommentTable.String)

	columnPostingIdx := np.CreateString(map[string]interface{}{"POSTING_IDX": nil}, db.GetDatabaseTypeString(), "", false)
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, db.GetDatabaseTypeString(), "", false)

	sql := `
	DELETE
	FROM ` + tableName + `
	WHERE ` + columnPostingIdx.Names + ` = ` + postingIdx + `
		AND ` + columnIdx.Names + ` = ` + commentIdx

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func DeleteComments(board model.Board, idx string) error {
	tableName := db.GetFullTableName(board.CommentTable.String)

	columnPostingIdx := np.CreateString(map[string]interface{}{"POSTING_IDX": nil}, db.GetDatabaseTypeString(), "", false)

	sql := `
	DELETE FROM ` + tableName + `
	WHERE ` + columnPostingIdx.Names + ` = ` + idx

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
