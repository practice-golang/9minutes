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

	columns := np.CreateString(model.UserColumn{}, dbtype, "", false).Names
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)

	sql := `
	SELECT
		` + columns + `
	FROM ` + tablename + `
	ORDER BY ` + columnIdx.Names + ` ASC`

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
		` + column.Names + `, ` + columnSortOrder.Names + `
	) VALUES (
		` + column.Values + `, (SELECT (MAX(B.` + columnSortOrder.Names + `) + 1) FROM ` + tablename + ` B)
	)`

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	err = db.Obj.AddTableColumn(db.Info.UserTable, userColumn)
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
		` + column.Names + `
	FROM ` + tablename + `
	WHERE ` + columnIdx.Names + ` = ` + idx

	r, err := db.Con.Query(sql)
	if err != nil {
		return err
	}

	err = scan.Rows(&userColumns, r)
	if err != nil {
		return err
	}

	userColumnOld := userColumns[0]

	colNames := strings.Split(column.Names, ",")
	colValues := strings.Split(column.Values, ",")
	holder := ""

	for i := 0; i < len(colNames); i++ {
		holder += colNames[i] + " = " + colValues[i] + ", "
	}
	holder = strings.TrimSuffix(holder, ", ")

	sql = `
	UPDATE ` + tablename + ` SET
		` + holder + `
	WHERE ` + columnIdx.Names + ` = ` + strconv.Itoa(int(userColumnNew.Idx.Int64))

	_, err = db.Con.Exec(sql)
	if err != nil {
		return err
	}

	err = db.Obj.EditTableColumn(db.Info.UserTable, userColumnOld, userColumnNew)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUserColumn(userColumn model.UserColumn) error {
	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUserColumns)

	columns := np.CreateString(model.UserColumn{}, dbtype, "", false).Names
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)
	idx := strconv.Itoa(int(userColumn.Idx.Int64))

	userColumns := []model.UserColumn{}
	sql := `
	SELECT
		` + columns + `
	FROM ` + tablename + `
	WHERE ` + columnIdx.Names + ` = ` + idx

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
	WHERE ` + columnIdx.Names + ` = ` + idx

	_, err = db.Con.Exec(sql)
	if err != nil {
		return err
	}

	err = db.Obj.DeleteTableColumn(db.Info.UserTable, userColumn)
	if err != nil {
		return err
	}

	return nil
}
