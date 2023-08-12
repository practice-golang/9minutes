package crud

import (
	"9minutes/internal/db"
	"9minutes/internal/np"
	"9minutes/model"
	"database/sql"
	"errors"
	"fmt"

	"github.com/blockloop/scan"
)

func GetUploadedFile(idx int) (model.StoredFileInfo, error) {
	finfo := model.StoredFileInfo{}
	finfos := []model.StoredFileInfo{}

	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(db.Info.UploadTable)
	column := np.CreateString(finfo, dbtype, "", false)
	where := np.CreateWhereString(map[string]interface{}{"IDX": idx}, dbtype, "=", "AND", "", false)

	sql := `
	SELECT
		` + column.Names + `
	FROM ` + tableName +
		where

	r, err := db.Con.Query(sql)
	if err != nil {
		return finfo, err
	}
	defer r.Close()

	err = scan.Rows(&finfos, r)
	if err != nil {
		return finfo, err
	}
	if len(finfos) == 0 {
		return finfo, errors.New("no file found")
	}

	finfo = finfos[0]

	return finfo, nil
}

func GetUploadedFiles(idxes []int) ([]model.StoredFileInfo, error) {
	finfo := model.StoredFileInfo{}
	finfos := []model.StoredFileInfo{}

	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(db.Info.UploadTable)
	column := np.CreateString(finfo, dbtype, "", false)

	var indices []map[string]interface{}
	for _, idx := range idxes {
		indices = append(indices, map[string]interface{}{"IDX": idx})
	}
	where := np.CreateWhereString(indices, dbtype, "=", "OR", "", false)

	sql := `
	SELECT
		` + column.Names + `
	FROM ` + tableName +
		where

	r, err := db.Con.Query(sql)
	if err != nil {
		return finfos, err
	}
	defer r.Close()

	err = scan.Rows(&finfos, r)
	if err != nil {
		return finfos, err
	}
	if len(finfos) == 0 {
		return finfos, errors.New("no file found")
	}

	return finfos, nil
}

func AddUploadedFile(fileName, storageName string) (sql.Result, error) {
	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(db.Info.UploadTable)
	columnsFileName := np.CreateString(map[string]interface{}{"FILE_NAME": nil}, dbtype, "", false)
	columnsStorageName := np.CreateString(map[string]interface{}{"STORAGE_NAME": nil}, dbtype, "", false)

	sql := `
	INSERT INTO ` + tableName + ` (
		` + columnsFileName.Names + `,` + columnsStorageName.Names + `
	) VALUES (
		'` + fileName + `', '` + storageName + `'
	)`

	result, err := db.Con.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateUploadedFile(boardIDX, postIDX int64, filename, storename string) error {
	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(db.Info.UploadTable)

	colFilename := np.CreateString(map[string]interface{}{"FILE_NAME": nil}, dbtype, "", false)
	colStorename := np.CreateString(map[string]interface{}{"STORAGE_NAME": nil}, dbtype, "", false)
	colBoardIDX := np.CreateString(map[string]interface{}{"BOARD_IDX": nil}, dbtype, "", false)
	colPostIDX := np.CreateString(map[string]interface{}{"POST_IDX": nil}, dbtype, "", false)

	columnsBoardIDX := np.CreateString(map[string]interface{}{"BOARD_IDX": nil}, dbtype, "", false)
	columnsPostIDX := np.CreateString(map[string]interface{}{"POST_IDX": nil}, dbtype, "", false)

	sql := `
	UPDATE ` + tableName + ` SET
		` + colBoardIDX.Names + ` = '` + fmt.Sprint(boardIDX) + `',
		` + colPostIDX.Names + ` = '` + fmt.Sprint(postIDX) + `'
	WHERE ` + colFilename.Names + ` = '` + filename + `'
		AND ` + colStorename.Names + ` = '` + storename + `'
		AND ` + columnsBoardIDX.Names + ` IS NULL
		AND ` + columnsPostIDX.Names + ` IS NULL`

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUploadedFile(idx int64) (err error) {
	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(db.Info.UploadTable)
	where := np.CreateWhereString(map[string]interface{}{"IDX": idx}, dbtype, "=", "AND", "", false)

	sql := `
	DELETE
	FROM ` + tableName +
		where

	_, err = db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
