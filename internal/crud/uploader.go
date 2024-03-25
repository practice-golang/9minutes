package crud

import (
	"9minutes/internal/db"
	"9minutes/internal/np"
	"9minutes/model"
	"strconv"
	"strings"

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
		` + column.Name + `
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
		// return finfo, errors.New("no file found")
		return finfo, nil
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
		` + column.Name + `
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
	// if len(finfos) == 0 {
	// 	// return finfos, errors.New("no file found")
	// 	return finfos, nil
	// }

	return finfos, nil
}

func AddUploadedFile(fileName, storageName string) (int64, int64, error) {
	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(db.Info.UploadTable)

	data := map[string]interface{}{
		"FILE_NAME":    fileName,
		"TOPIC_IDX":    -1,
		"COMMENT_IDX":  -1,
		"STORAGE_NAME": storageName,
	}
	columns := np.CreateString(data, dbtype, "", false)

	sql := `
	INSERT INTO ` + tableName + ` (
		` + columns.Name + `
	) VALUES (
		` + columns.Value + `
	)`

	count, idx, err := db.Obj.Exec(sql, []interface{}{}, "IDX")
	if err != nil {
		return -1, -1, err
	}

	return count, idx, nil
}

func DeleteUploadedFile(idx int64) (err error) {
	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(db.Info.UploadTable)
	where := np.CreateWhereString(map[string]interface{}{"IDX": idx}, dbtype, "=", "AND", "", false)

	sql := `
	DELETE
	FROM ` + tableName + `
	` + where

	_, err = db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func SetUploadedFileIndex(idx, topicIDX, commentIDX int64) (err error) {
	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(db.Info.UploadTable)

	dataWhere := map[string]interface{}{"IDX": idx}
	where := np.CreateWhereString(dataWhere, dbtype, "=", "AND", "", false)

	data := map[string]string{}
	if topicIDX > 0 {
		data["TOPIC_IDX"] = strconv.FormatInt(topicIDX, 10)
	}
	if commentIDX > 0 {
		data["COMMENT_IDX"] = strconv.FormatInt(commentIDX, 10)
	}
	columns := np.CreateString(data, dbtype, "update", false)
	colNames := strings.Split(columns.Name, ",")
	colValues := strings.Split(columns.Value, ",")
	holder := ""

	for i := 0; i < len(colNames); i++ {
		holder += colNames[i] + " = " + colValues[i] + ", "
	}
	holder = strings.TrimSuffix(holder, ", ")

	sql := `
	UPDATE ` + tableName + ` SET
		` + holder + `
	` + where

	_, err = db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
