package crud

import (
	"9minutes/internal/db"
	"9minutes/internal/np"
	"9minutes/model"
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/blockloop/scan"
	"gopkg.in/guregu/null.v4"
)

func GetTopicList(board model.Board, topicListOption model.TopicListingOption) (model.TopicPageData, error) {
	result := model.TopicPageData{}

	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(board.BoardTable.String)
	commentTableName := db.GetFullTableName(board.CommentTable.String)

	column := np.CreateString(model.TopicList{}, dbtype, "select", false)

	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)
	columnTitle := np.CreateString(map[string]interface{}{"TITLE": nil}, dbtype, "", false)
	columnContent := np.CreateString(map[string]interface{}{"CONTENT": nil}, dbtype, "", false)
	columnTopicIdx := np.CreateString(map[string]interface{}{"TOPIC_IDX": nil}, dbtype, "", false)
	columnCommentCount := np.CreateString(map[string]interface{}{"COMMENT_COUNT": nil}, dbtype, "", false)

	sqlSearch := ""
	if topicListOption.Search.Valid && topicListOption.Search.String != "" {
		sqlSearch = `
		WHERE LOWER(` + columnTitle.Names + `) LIKE LOWER('%` + topicListOption.Search.String + `%')
			OR LOWER(` + columnContent.Names + `) LIKE LOWER('%` + topicListOption.Search.String + `%')`
	}

	paging := ``
	if topicListOption.Page.Valid && topicListOption.ListCount.Valid {
		paging = db.Obj.GetPagingQuery(int(topicListOption.Page.Int64*topicListOption.ListCount.Int64), int(topicListOption.ListCount.Int64))
	}

	sql := `
	SELECT
		` + column.Names + `,
		(
			SELECT
				COUNT(` + columnIdx.Names + `)
			FROM ` + commentTableName + `
			WHERE ` + columnTopicIdx.Names + ` = A.` + columnIdx.Names + `
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

	var list []model.TopicList
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

	totalPage := math.Ceil(float64(totalCount) / float64(topicListOption.ListCount.Int64))

	result = model.TopicPageData{
		BoardCode:     board.BoardCode.String,
		SearchKeyword: topicListOption.Search.String,
		TopicList:     list,
		CurrentPage:   int(topicListOption.Page.Int64) + 1,
		TotalPage:     int(totalPage),
	}

	return result, err
}

func GetTopic(board model.Board, idx string) (model.Topic, error) {
	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(board.BoardTable.String)
	// userTableName := db.GetFullTableName(db.Info.UserTable)

	column := np.CreateString(model.Topic{}, dbtype, "select", false)
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)
	// columnUserId := np.CreateString(map[string]interface{}{"USERID": nil}, dbtype, "", false)
	// columnAuthorIdx := np.CreateString(map[string]interface{}{"AUTHOR_IDX": nil}, dbtype, "", false)
	// columnAuthorName := np.CreateString(map[string]interface{}{"AUTHOR_NAME": nil}, dbtype, "", false)

	/*
		(
		SELECT
			` + columnUserId.Names + `
		FROM ` + userTableName + `
		WHERE ` + columnIdx.Names + ` = ` + tableName + `.` + columnAuthorIdx.Names + `
		) AS ` + columnAuthorName.Names + `
	*/
	sql := `
	SELECT
		` + column.Names + `
	FROM ` + tableName + `
	WHERE ` + columnIdx.Names + ` = ` + idx

	r, err := db.Con.Query(sql)
	if err != nil {
		return model.Topic{}, err
	}
	defer r.Close()

	var content model.Topic
	err = scan.Row(&content, r)
	if err != nil {
		return model.Topic{}, err
	}

	return content, nil
}

func WriteTopic(board model.Board, content model.Topic) (sql.Result, error) {
	tableName := db.GetFullTableName(board.BoardTable.String)

	content.Title = null.StringFrom(content.Title.String)
	content.TitleImage = null.StringFrom(content.TitleImage.String)
	content.Content = null.StringFrom(content.Content.String)

	content.Title.String = EscapeString(content.Title.String)
	content.Content.String = EscapeString(content.Content.String)

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

func UpdateTopic(board model.Board, content model.Topic, skipTag string) error {
	tableName := db.GetFullTableName(board.BoardTable.String)
	idx := fmt.Sprint(content.Idx.Int64)

	content.Title = null.StringFrom(content.Title.String)
	content.Content = null.StringFrom(content.Content.String)

	content.Title.String = EscapeString(content.Title.String)
	content.Content.String = EscapeString(content.Content.String)

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

func DeleteTopic(board model.Board, idx string) error {
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
