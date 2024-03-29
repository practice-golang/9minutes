package crud

import (
	"9minutes/internal/db"
	"9minutes/internal/np"
	"log"
)

func GetMemberList(boardIDX int64) {
	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(db.Info.UserTable)

	columnList, _ := GetUserColumnsList()
	columnNames := map[string]interface{}{}
	for _, column := range columnList {
		if column.ColumnName.Valid {
			columnNames["A."+column.ColumnName.String] = nil
		}
	}
	columns := np.CreateString(columnNames, dbtype, "", false)

	sql := `
	SELECT
		` + columns.Name + `
	FROM ` + tablename + ` AS A
	`

	log.Println(sql)
}
