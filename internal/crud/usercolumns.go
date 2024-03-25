package crud

import (
	"9minutes/consts"
	"9minutes/internal/db"
	"9minutes/internal/np"
	"9minutes/model"
	"strconv"
	"strings"

	"github.com/blockloop/scan"
)

func GetUserColumnsCount() (int, error) {
	result := 0

	tablename := db.GetFullTableName(consts.TableUserColumns)

	sql := `SELECT COUNT(*) AS CNT FROM ` + tablename

	r, err := db.Con.Query(sql)
	if err != nil {
		return -1, err
	}

	err = scan.Row(&result, r)
	if err != nil {
		return -1, err
	}

	return result, nil
}

func GetUserColumnsList() ([]model.UserColumn, error) {
	result := []model.UserColumn{}

	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUserColumns)

	columns := np.CreateString(model.UserColumn{}, dbtype, "", false).Name
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)

	sql := `
	SELECT
		` + columns + `
	FROM ` + tablename + `
	ORDER BY ` + columnIdx.Name + ` ASC`

	// where := []interface{}{}
	// r, err := db.Con.Query(sql, where...)
	r, err := db.Con.Query(sql)
	if err != nil {
		return nil, err
	}

	err = scan.Rows(&result, r)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func AddUserColumn(userColumn model.UserColumn) error {
	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUserColumns)

	column := np.CreateString(userColumn, dbtype, "insert", false)
	columnSortOrder := np.CreateString(map[string]interface{}{"SORT_ORDER": nil}, dbtype, "", false)

	// sql := `INSERT INTO ` + tablename + ` (` + columns + `) VALUES (` + values + `)`
	sql := `
	INSERT INTO ` + tablename + ` (
		` + column.Name + `, ` + columnSortOrder.Name + `
	) VALUES (
		` + column.Value + `, (SELECT (MAX(B.` + columnSortOrder.Name + `) + 1) FROM ` + tablename + ` B)
	)`

	// _, err := db.Con.Exec(sql)
	_, _, err := db.Obj.Exec(sql, []interface{}{}, "IDX")
	if err != nil {
		return err
	}

	usertable := db.GetFullTableName(db.Info.UserTable)
	err = db.Obj.AddTableColumn(usertable, userColumn)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserColumn(userColumnNew model.UserColumn) error {
	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUserColumns)

	column := np.CreateString(userColumnNew, dbtype, "update", false)
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)
	idx := strconv.Itoa(int(userColumnNew.Idx.Int64))

	userColumns := []model.UserColumn{}

	sql := `
	SELECT
		` + column.Name + `
	FROM ` + tablename + `
	WHERE ` + columnIdx.Name + ` = ` + idx

	r, err := db.Con.Query(sql)
	if err != nil {
		return err
	}

	err = scan.Rows(&userColumns, r)
	if err != nil {
		return err
	}

	userColumnOld := userColumns[0]

	colNames := strings.Split(column.Name, ",")
	colValues := strings.Split(column.Value, ",")
	holder := ""

	for i := 0; i < len(colNames); i++ {
		holder += colNames[i] + " = " + colValues[i] + ", "
	}
	holder = strings.TrimSuffix(holder, ", ")

	sql = `
	UPDATE ` + tablename + ` SET
		` + holder + `
	WHERE ` + columnIdx.Name + ` = ` + strconv.Itoa(int(userColumnNew.Idx.Int64))

	_, err = db.Con.Exec(sql)
	if err != nil {
		return err
	}

	usertable := db.GetFullTableName(db.Info.UserTable)
	err = db.Obj.EditTableColumn(usertable, userColumnOld, userColumnNew)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUserColumn(userColumn model.UserColumn) error {
	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUserColumns)

	columns := np.CreateString(model.UserColumn{}, dbtype, "", false).Name
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)
	idx := strconv.Itoa(int(userColumn.Idx.Int64))

	userColumns := []model.UserColumn{}
	sql := `
	SELECT
		` + columns + `
	FROM ` + tablename + `
	WHERE ` + columnIdx.Name + ` = ` + idx

	r, err := db.Con.Query(sql)
	if err != nil {
		return err
	}

	err = scan.Rows(&userColumns, r)
	if err != nil {
		return err
	}

	userColumn = userColumns[0]

	sql = `
	DELETE
	FROM ` + tablename + `
	WHERE ` + columnIdx.Name + ` = ` + idx

	_, err = db.Con.Exec(sql)
	if err != nil {
		return err
	}

	usertable := db.GetFullTableName(db.Info.UserTable)
	err = db.Obj.DeleteTableColumn(usertable, userColumn)
	if err != nil {
		return err
	}

	return nil
}
