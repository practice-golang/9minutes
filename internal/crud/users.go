package crud

import (
	"9minutes/consts"
	"9minutes/internal/db"
	"9minutes/internal/np"
	"9minutes/model"
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/blockloop/scan"
	"golang.org/x/crypto/bcrypt"
)

func GetUserByUserIdAndPassword(name, password string) (interface{}, error) {
	var result interface{}

	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)

	// columns := ""
	wheres := np.CreateWhereString(map[string]interface{}{"USERID": name}, dbtype, "=", "AND", "", false)

	columnList, _ := GetUserColumnsList()
	columnNames := map[string]interface{}{}
	for _, column := range columnList {
		if column.ColumnName.Valid {
			// columns += column.ColumnName.String + ","
			columnNames[column.ColumnName.String] = nil
		}
	}
	// columns = strings.TrimSuffix(columns, ",")
	columns := np.CreateString(columnNames, dbtype, "", false)

	sql := `
	SELECT
		` + columns.Name + `
	FROM ` + tablename +
		wheres

	r, err := db.Con.Query(sql)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	users := []map[string]interface{}{}
	for r.Next() {
		scanedRow, _ := scanMap(r)
		users = append(users, scanedRow)
	}

	if len(users) == 0 {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(users[0]["password"].(string)), []byte(password))
	if err != nil {
		return nil, errors.New("invalid userid or password")
	}

	result = users[0]
	return result, nil
}

func GetUserByNameAndEmailMap(userid, email string) (interface{}, error) {
	var result interface{}

	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)

	whereUserId := np.CreateString(map[string]interface{}{"USERID": nil}, dbtype, "", false)
	whereEmail := np.CreateString(map[string]interface{}{"EMAIL": nil}, dbtype, "", false)
	whereGrade := np.CreateString(map[string]interface{}{"GRADE": nil}, dbtype, "", false)

	// Use map with default and user defined columns
	columnList, _ := GetUserColumnsList()
	columnNames := map[string]interface{}{}
	for _, column := range columnList {
		if column.ColumnName.Valid {
			// columns += column.ColumnName.String + ","
			columnNames[column.ColumnName.String] = nil
		}
	}

	// columns = strings.TrimSuffix(columns, ",")
	columns := np.CreateString(columnNames, dbtype, "", false)

	sql := `
	SELECT
		` + columns.Name + `
	FROM ` + tablename + `
	WHERE ` + whereUserId.Name + ` = '` + userid + `'
		AND ` + whereEmail.Name + ` = '` + email + `'
		AND ` + whereGrade.Name + ` != '` + "user_quit" + `'`

	r, err := db.Con.Query(sql)
	if err != nil {
		return result, err
	}
	defer r.Close()

	users := []map[string]interface{}{}
	for r.Next() {
		scanedRow, _ := scanMap(r)
		users = append(users, scanedRow)
	}

	if len(users) > 0 {
		result = users[0]
	}

	return result, nil
}

func GetUserByNameMap(userid string) (interface{}, error) {
	var result interface{}
	dbtype := db.GetDatabaseTypeString()

	tablename := db.GetFullTableName(consts.TableUsers)
	columnUserId := np.CreateString(map[string]interface{}{"USERID": nil}, db.GetDatabaseTypeString(), "", false)

	// Use map with default and user defined columns
	columnList, _ := GetUserColumnsList()
	columnNames := map[string]interface{}{}
	for _, column := range columnList {
		if column.ColumnName.Valid {
			// columns += column.ColumnName.String + ","
			columnNames[column.ColumnName.String] = nil
		}
	}

	// columns = strings.TrimSuffix(columns, ",")
	columns := np.CreateString(columnNames, dbtype, "", false)

	sql := `
	SELECT
		` + columns.Name + `
	FROM ` + tablename + `
	WHERE ` + columnUserId.Name + ` = '` + userid + `'`

	r, err := db.Con.Query(sql)
	if err != nil {
		return result, err
	}
	defer r.Close()

	users := []map[string]interface{}{}
	for r.Next() {
		scanedRow, _ := scanMap(r)
		users = append(users, scanedRow)
	}

	if len(users) > 0 {
		result = users[0]
	}

	return result, nil
}

// GetUsersMap - Get Users List map
func GetUsersMap(userListOption model.UserListingOption) (model.UserPageData, error) {
	result := model.UserPageData{}

	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(consts.TableUsers)
	// columns := ""
	columnUserId := np.CreateString(map[string]interface{}{"USERID": nil}, dbtype, "", false)
	columnEmail := np.CreateString(map[string]interface{}{"EMAIL": nil}, dbtype, "", false)
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)

	// Use map with default and user defined columns
	columnList, _ := GetUserColumnsList()
	columnsValid := map[string]interface{}{}
	for _, column := range columnList {
		if column.ColumnName.Valid {
			// columns += column.ColumnName.String + ","

			jsonName := strings.ToLower(column.ColumnName.String)
			jsonName = strings.ReplaceAll(jsonName, "_", "-")

			columnsValid[jsonName] = nil
		}
	}

	// columns = strings.TrimSuffix(columns, ",")
	columns := np.CreateString(columnsValid, db.GetDatabaseTypeString(), "", false)

	sqlSearch := ""

	if userListOption.Search.Valid && userListOption.Search.String != "" {
		sqlSearch = `
		WHERE ` + columnUserId.Name + ` LIKE '%` + userListOption.Search.String + `%'
			OR ` + columnEmail.Name + ` LIKE '%` + userListOption.Search.String + `%'`
	}

	paging := ``
	if userListOption.Page.Valid && userListOption.ListCount.Valid {
		paging = db.Obj.GetPagingQuery(int(userListOption.Page.Int64*userListOption.ListCount.Int64), int(userListOption.ListCount.Int64))
	}

	sql := `
	SELECT
		` + columns.Name + `
	FROM ` + tableName + `
	` + sqlSearch + `
	ORDER BY ` + columnIdx.Name + ` ASC
	` + paging

	r, err := db.Con.Query(sql)
	if err != nil {
		return result, err
	}
	defer r.Close()

	users := []map[string]interface{}{}
	for r.Next() {
		scanedRow, _ := scanMap(r)
		users = append(users, scanedRow)
	}

	var totalCount int64
	sql = `
	SELECT
		COUNT(` + columnIdx.Name + `)
	FROM ` + tableName + ` ` + sqlSearch

	r, err = db.Con.Query(sql)
	if err != nil {
		return result, err
	}
	defer r.Close()

	err = scan.Row(&totalCount, r)
	if err != nil {
		return result, err
	}

	totalPage := math.Ceil(float64(totalCount) / float64(userListOption.ListCount.Int64))

	result = model.UserPageData{
		UserList:    users,
		CurrentPage: int(userListOption.Page.Int64) + 1,
		TotalPage:   int(totalPage),
	}

	return result, nil
}

// GetUsersListMap - API Get Users List map
func GetUsersListMap(columnsMap map[string]interface{}, userListOption model.UserListingOption) (model.UserPageData, error) {
	result := model.UserPageData{}
	users := []map[string]interface{}{}

	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(consts.TableUsers)
	columns := np.CreateString(columnsMap, dbtype, "", false).Name

	sqlSearch := ""
	if userListOption.Search.Valid && userListOption.Search.String != "" {
		search := userListOption.Search.String
		sqlSearch += np.CreateWhereString(
			map[string]interface{}{"USERID": search, "EMAIL": search},
			db.GetDatabaseTypeString(), "LIKE", "OR", "", false,
		)
	}

	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)

	paging := db.Obj.GetPagingQuery(
		(int(userListOption.Page.Int64)-1)*int(userListOption.ListCount.Int64),
		int(userListOption.ListCount.Int64),
	)

	sql := `
	SELECT
		` + columns + `
	FROM ` + tableName + `
	` + sqlSearch + `
	ORDER BY ` + columnIdx.Name + ` ASC ` + `
	` + paging

	r, err := db.Con.Query(sql)
	if err != nil {
		return result, err
	}

	for r.Next() {
		scanedRow, _ := scanMap(r)
		users = append(users, scanedRow)
	}

	var totalCount int64
	sql = `
	SELECT
		COUNT(` + columnIdx.Name + `)
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

	totalPage := math.Ceil(float64(totalCount) / float64(userListOption.ListCount.Int64))

	result = model.UserPageData{
		UserList:    users,
		CurrentPage: int(userListOption.Page.Int64),
		TotalPage:   int(totalPage),
	}

	return result, nil
}

func AddUserMap(userMap map[string]interface{}) error {
	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)

	column := np.CreateString(userMap, dbtype, "", false)

	sql := `
	INSERT INTO ` + tablename + ` (
		` + column.Name + `
	) VALUES (
		` + column.Value + `
	)`

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func AddUserVerification(verificationData map[string]string) error {
	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers + "_verification")

	column := np.CreateString(verificationData, dbtype, "insert", false)

	sql := `
	INSERT INTO ` + tablename + ` (
		` + column.Name + `
	) VALUES (
		` + column.Value + `
	)`
	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserMap(dataMap map[string]interface{}) error {
	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)

	data := map[string]interface{}{}
	for k, v := range dataMap {
		data[k] = v
	}

	var idx int64
	switch idxUnknown := data["idx"].(type) {
	case int64:
		idx = idxUnknown
	case int:
		idx = int64(idxUnknown)
	case float64:
		idx = int64(idxUnknown)
	case string:
		idxSTR, _ := strconv.Atoi(idxUnknown)
		idx = int64(idxSTR)
	}
	delete(data, "idx")

	updateset := np.CreateUpdateString(data, dbtype, "", false)
	wheres := np.CreateWhereString(
		map[string]interface{}{"IDX": idx},
		dbtype, "=", "AND", "", false,
	)

	sql := `UPDATE ` + tablename + updateset + wheres

	r, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("affected nothing")
	}

	return nil
}

func QuitUser(idx int64) error {
	dbType := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)

	updateset := np.CreateUpdateString(
		map[string]interface{}{
			"GRADE":    "user_quit",
			"APPROVAL": "N",
		},
		dbType, "", false,
	)
	wheres := np.CreateWhereString(
		map[string]interface{}{"IDX": idx},
		dbType, "=", "AND", "", false,
	)

	sql := `UPDATE ` + tablename + updateset + wheres

	r, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("affected nothing")
	}

	return nil
}

func DeleteUser(idx int64) error {
	dbType := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)

	wheres := np.CreateWhereString(
		map[string]interface{}{"IDX": idx},
		dbType, "=", "AND", "", false,
	)

	sql := `DELETE FROM ` + tablename + wheres

	r, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("affected nothing")
	}

	return nil
}
