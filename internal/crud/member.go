package crud

import (
	"9minutes/internal/db"
	"9minutes/internal/np"
	"9minutes/model"
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

func AddMember(data model.MemberRequest) (count, idx int64, err error) {
	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(db.Info.MemberTable)

	columns := np.CreateString(data, dbtype, "insert", false)

	sql := `
	INSERT INTO ` + tablename + `(
		` + columns.Name + `
	) VALUES (
		` + columns.Value + `
	)`

	count, idx, err = db.Obj.Exec(sql, nil, "IDX")
	if err != nil {
		return -1, -1, err
	}

	return count, idx, err
}
