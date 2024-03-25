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
func GetComment(board model.Board, topicIdx, commentIdx string) (model.Comment, error) {
	var comment model.Comment

	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(board.CommentTable.String)

	column := np.CreateString(model.Comment{}, dbtype, "select", false)
	columnTopicIdx := np.CreateString(map[string]interface{}{"TOPIC_IDX": nil}, dbtype, "", false)
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)

	sql := `
	SELECT
		` + column.Name + `
	FROM ` + tableName + `
	WHERE ` + columnTopicIdx.Name + ` = ` + topicIdx + `
		AND ` + columnIdx.Name + ` = ` + commentIdx

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
func GetComments(board model.Board, contentIdx string, commentListOption model.CommentListingOptions) (model.CommentPageData, error) {
	result := model.CommentPageData{}

	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(board.CommentTable.String)
	// userTableName := db.GetFullTableName(db.Info.UserTable)

	column := np.CreateString(model.Comment{}, dbtype, "select", false)

	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)
	columnTopicIdx := np.CreateString(map[string]interface{}{"TOPIC_IDX": nil}, dbtype, "", false)

	var totalCount int64
	sql := `
	SELECT
		COUNT(` + columnIdx.Name + `)
	FROM ` + tableName + `
	WHERE ` + columnTopicIdx.Name + ` = ` + contentIdx

	r, err := db.Con.Query(sql)
	if err != nil {
		return result, err
	}
	defer r.Close()

	err = scan.Row(&totalCount, r)
	if err != nil {
		return result, err
	}

	totalPage := math.Ceil(float64(totalCount) / float64(commentListOption.ListCount.Int64))
	currentPage := int64(totalPage)

	offset := 0
	if totalCount > 0 {
		offset = int(int64(totalPage)*commentListOption.ListCount.Int64 - commentListOption.ListCount.Int64)

		if commentListOption.Page.Valid && commentListOption.Page.Int64 > -1 && commentListOption.ListCount.Valid {
			currentPage = commentListOption.Page.Int64 + 1
			offset = int(commentListOption.Page.Int64 * commentListOption.ListCount.Int64)
		}
	}

	paging := db.Obj.GetPagingQuery(offset, int(commentListOption.ListCount.Int64))

	sql = `
	SELECT
		` + column.Name + `
	FROM ` + tableName + `
	WHERE ` + columnTopicIdx.Name + ` = ` + contentIdx + `
	ORDER BY ` + columnIdx.Name + ` ASC
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

// func WriteComment(board model.Board, content model.Comment) error {
func WriteComment(board model.Board, content model.Comment) (int64, int64, error) {
	tableName := db.GetFullTableName(board.CommentTable.String)

	content.Content.String = EscapeString(content.Content.String)

	column := np.CreateString(content, db.GetDatabaseTypeString(), "insert", false)

	sql := `
	INSERT INTO ` + tableName + ` (
		` + column.Name + `
	) VALUES (
		` + column.Value + `
	)`

	count, idx, err := db.Obj.Exec(sql, []interface{}{}, "IDX")
	if err != nil {
		// return err
		return 0, -1, err
	}

	// return nil
	return count, idx, nil
}

func UpdateComment(board model.Board, comment model.Comment, topicIdx string) error {
	tableName := db.GetFullTableName(board.CommentTable.String)
	commentIDX := strconv.Itoa(int(comment.Idx.Int64))

	comment.Content.String = EscapeString(comment.Content.String)

	column := np.CreateString(comment, db.GetDatabaseTypeString(), "update", false)
	colNames := strings.Split(column.Name, ",")
	colValues := strings.Split(column.Value, ",")

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
	WHERE ` + columnIdx.Name + ` = ` + commentIDX

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func DeleteComment(board model.Board, topicIdx, commentIdx string) error {
	tableName := db.GetFullTableName(board.CommentTable.String)

	columnTopicIdx := np.CreateString(map[string]interface{}{"TOPIC_IDX": nil}, db.GetDatabaseTypeString(), "", false)
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, db.GetDatabaseTypeString(), "", false)

	sql := `
	DELETE
	FROM ` + tableName + `
	WHERE ` + columnTopicIdx.Name + ` = ` + topicIdx + `
		AND ` + columnIdx.Name + ` = ` + commentIdx

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func DeleteComments(board model.Board, idx string) error {
	tableName := db.GetFullTableName(board.CommentTable.String)

	columnTopicIdx := np.CreateString(map[string]interface{}{"TOPIC_IDX": nil}, db.GetDatabaseTypeString(), "", false)

	sql := `
	DELETE FROM ` + tableName + `
	WHERE ` + columnTopicIdx.Name + ` = ` + idx

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
