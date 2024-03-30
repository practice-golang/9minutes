package crud

import (
	"9minutes/internal/db"
	"9minutes/internal/np"
	"9minutes/model"

	"github.com/blockloop/scan"
)

func GetMemberList(boardIDX int64) ([]model.Member, error) {
	var result []model.Member

	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(db.Info.MemberTable)

	columns := np.CreateString(model.Member{}, dbtype, "select", false)
	whereVar := map[string]interface{}{"BOARD_IDX": boardIDX}
	where := np.CreateWhereString(whereVar, dbtype, "=", "AND", "", false)

	sql := `
	SELECT
		` + columns.Name + `
	FROM ` + tablename + ` AS A
	` + where + ``

	r, err := db.Con.Query(sql)
	if err != nil {
		return result, err
	}
	defer r.Close()

	err = scan.Rows(&result, r)
	if err != nil {
		return result, err
	}

	return result, nil
}

func AddMember(data model.Member) (count, idx int64, err error) {
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
