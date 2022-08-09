package crud

import (
	"9minutes/db"
	"9minutes/np"
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

func DeleteUploadedFile(fileName, storeName string) error {
	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(db.Info.UploadTable)
	columnsFileName := np.CreateString(map[string]interface{}{"FILE_NAME": nil}, dbtype, "", false)
	columnsStorageName := np.CreateString(map[string]interface{}{"STORAGE_NAME": nil}, dbtype, "", false)

	sql := `
	DELETE
	FROM ` + tableName + `
	WHERE
		` + columnsFileName.Names + ` = '` + fileName + `' AND
		` + columnsStorageName.Names + ` = '` + storeName + `'`

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
