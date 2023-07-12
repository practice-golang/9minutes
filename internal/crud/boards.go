package crud

import (
	"9minutes/consts"
	"9minutes/internal/db"
	"9minutes/internal/np"
	"9minutes/model"
	"encoding/json"
	"math"
	"strconv"
	"strings"

	"github.com/blockloop/scan"
)

func GetBoardByIdx(board model.Board) (model.Board, error) {
	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(consts.TableBoards)

	column := np.CreateString(board, dbtype, "", false)
	where := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)

	sql := `
	SELECT
		` + column.Names + `
	FROM ` + tableName + `
	WHERE ` + where.Names + ` = ` + strconv.Itoa(int(board.Idx.Int64))

	r, err := db.Con.Query(sql)
	if err != nil {
		return board, err
	}
	defer r.Close()

	scan.Row(&board, r)
	if err != nil {
		return board, err
	}

	// Ignore until board type sheet
	if board.Fields != nil {
		var fields model.Field
		err = json.Unmarshal(board.Fields.([]byte), &fields)
		if err != nil {
			return board, err
		}
		board.Fields = fields
	}

	return board, nil
}

func GetBoardByCode(board model.Board) (model.Board, error) {
	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(consts.TableBoards)

	column := np.CreateString(board, dbtype, "", false)
	where := np.CreateString(map[string]interface{}{"BOARD_CODE": nil}, dbtype, "", false)

	sql := `
	SELECT
		` + column.Names + `
	FROM ` + tableName + `
	WHERE ` + where.Names + `='` + board.BoardCode.String + `'`

	r, err := db.Con.Query(sql)
	if err != nil {
		return board, err
	}
	defer r.Close()

	scan.Row(&board, r)
	if err != nil {
		return board, err
	}

	return board, nil
}

func GetBoards(options model.BoardListingOptions) (model.BoardPageData, error) {
	result := model.BoardPageData{}

	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(consts.TableBoards)

	column := np.CreateString(model.Board{}, dbtype, "", false)
	whereBoardName := np.CreateString(map[string]interface{}{"BOARD_NAME": nil}, dbtype, "", false)
	whereBoardCode := np.CreateString(map[string]interface{}{"BOARD_CODE": nil}, dbtype, "", false)
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)

	sqlSearch := ""

	if options.Search.Valid && options.Search.String != "" {
		sqlSearch = `
		WHERE ` + whereBoardName.Names + ` LIKE '%` + options.Search.String + `%'
			OR ` + whereBoardCode.Names + ` LIKE '%` + options.Search.String + `%'`
	}

	paging := ``
	if options.Page.Valid && options.ListCount.Valid {
		paging = db.Obj.GetPagingQuery(int(options.Page.Int64*options.ListCount.Int64), int(options.ListCount.Int64))
	}

	sql := `
	SELECT
		` + column.Names + `
	FROM ` + tableName + `
	` + sqlSearch + `
	ORDER BY ` + columnIdx.Names + ` ASC
	` + paging

	r, err := db.Con.Query(sql)
	if err != nil {
		return result, err
	}
	defer r.Close()

	var boards []model.Board
	err = scan.Rows(&boards, r)
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

	result = model.BoardPageData{
		BoardList:   boards,
		CurrentPage: int(options.Page.Int64) + 1,
		TotalPage:   int(totalPage),
	}

	return result, nil
}

func AddBoard(board model.Board) error {
	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(consts.TableBoards)

	column := np.CreateString(board, dbtype, "insert", false)
	sql := `
	INSERT INTO ` + tableName + ` (
		` + column.Names + `
	) VALUES (
		` + column.Values + `
	)`

	_, err := db.Con.Query(sql)
	if err != nil {
		return err
	}

	// 테이블 생성 - 작업전까지 일단 무시

	return nil
}

func UpdateBoard(board model.Board) error {
	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(consts.TableBoards)

	column := np.CreateString(board, dbtype, "update", false)
	idx := strconv.Itoa(int(board.Idx.Int64))

	colNames := strings.Split(column.Names, ",")
	colValues := strings.Split(column.Values, ",")
	holder := ""

	for i := 0; i < len(colNames); i++ {
		holder += colNames[i] + " = " + colValues[i] + ", "
	}
	holder = strings.TrimSuffix(holder, ", ")

	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)

	sql := `
	UPDATE ` + tableName + ` SET
		` + holder + `
	WHERE ` + columnIdx.Names + ` = ` + idx

	_, err := db.Con.Query(sql)
	if err != nil {
		return err
	}

	return nil
}

func DeleteBoard(board model.Board) error {
	tableName := db.GetFullTableName(consts.TableBoards)

	idx := strconv.Itoa(int(board.Idx.Int64))
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, db.GetDatabaseTypeString(), "", false)

	sql := `
	DELETE
	FROM ` + tableName + `
	WHERE ` + columnIdx.Names + ` = ` + idx

	_, err := db.Con.Query(sql)
	if err != nil {
		return err
	}

	return nil
}
