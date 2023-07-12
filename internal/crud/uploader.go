package crud

import (
	"9minutes/internal/db"
	"9minutes/internal/np"
	"fmt"
)

func AddUploadedFile(fileName, storageName string) error {
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

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
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

func DeleteUploadedFile(boardIDX, postIDX int64, fileName, storeName string) error {
	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(db.Info.UploadTable)
	columnsFileName := np.CreateString(map[string]interface{}{"FILE_NAME": nil}, dbtype, "", false)
	columnsStorageName := np.CreateString(map[string]interface{}{"STORAGE_NAME": nil}, dbtype, "", false)

	columnsBoardIDX := np.CreateString(map[string]interface{}{"BOARD_IDX": nil}, dbtype, "", false)
	columnsPostIDX := np.CreateString(map[string]interface{}{"POST_IDX": nil}, dbtype, "", false)

	sql := `
	DELETE
	FROM ` + tableName + `
	WHERE ` + columnsFileName.Names + ` = '` + fileName + `'
		AND ` + columnsStorageName.Names + ` = '` + storeName + `'`

	if boardIDX > 0 && postIDX > 0 {
		sql += `
			AND ` + columnsBoardIDX.Names + ` = '` + fmt.Sprint(boardIDX) + `'
			AND ` + columnsPostIDX.Names + ` = '` + fmt.Sprint(postIDX) + `'`
	} else {
		sql += `
			AND ` + columnsBoardIDX.Names + ` IS NULL
			AND ` + columnsPostIDX.Names + ` IS NULL`
	}

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
