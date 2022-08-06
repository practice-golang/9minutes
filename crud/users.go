package crud

import (
	"9minutes/consts"
	"9minutes/db"
	"9minutes/model"
	"9minutes/np"
	"math"
	"strconv"
	"strings"

	"github.com/blockloop/scan"
)

func GetUserByNameAndEmail(username, email string) (model.UserData, error) {
	result := model.UserData{}

	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)
	columns := np.CreateString(model.UserData{}, dbtype, "", false).Names
	whereUsername := np.CreateString(map[string]interface{}{"USERNAME": nil}, dbtype, "", false)
	whereEmail := np.CreateString(map[string]interface{}{"EMAIL": nil}, dbtype, "", false)

	sql := `
	SELECT
		` + columns + `
	FROM ` + tablename + `
	WHERE ` + whereUsername.Names + ` = '` + username + `'
		AND ` + whereEmail.Names + ` = '` + email + `'`

	r, err := db.Con.Query(sql)
	if err != nil {
		return result, err
	}
	defer r.Close()

	err = scan.Row(&result, r)
	if err != nil {
		return result, err
	}

	return result, nil
}

func GetUserByNameAndEmailMap(username, email string) (interface{}, error) {
	var result interface{}

	tablename := db.GetFullTableName(consts.TableUsers)
	columns := ""
	columnUsername := np.CreateString(map[string]interface{}{"USERNAME": nil}, db.GetDatabaseTypeString(), "", false)
	columnEmail := np.CreateString(map[string]interface{}{"EMAIL": nil}, db.GetDatabaseTypeString(), "", false)

	// Use map with default and user defined columns
	columnList, _ := GetUserColumnsList()
	for _, column := range columnList {
		if column.ColumnName.Valid {
			columns += column.ColumnName.String + ","

			jsonName := strings.ToLower(column.ColumnName.String)
			jsonName = strings.ReplaceAll(jsonName, "_", "-")
		}
	}

	columns = strings.TrimSuffix(columns, ",")

	sql := `
	SELECT
		` + columns + `
	FROM ` + tablename + `
	WHERE ` + columnUsername.Names + ` = '` + username + `'
		AND ` + columnEmail.Names + ` = '` + email + `'`

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

func GetUserByName(username string) (model.UserData, error) {
	result := model.UserData{}

	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)
	columns := np.CreateString(model.UserData{}, dbtype, "", false).Names
	where := np.CreateString(map[string]interface{}{"USERNAME": nil}, dbtype, "", false)

	sql := `SELECT ` + columns + ` FROM ` + tablename + ` WHERE ` + where.Names + ` = '` + username + `'`

	r, err := db.Con.Query(sql)
	if err != nil {
		return result, err
	}
	defer r.Close()

	err = scan.Row(&result, r)
	if err != nil {
		return result, err
	}

	return result, nil
}

func GetUserByNameMap(username string) (interface{}, error) {
	var result interface{}

	tablename := db.GetFullTableName(consts.TableUsers)
	columns := ""
	columnUsername := np.CreateString(map[string]interface{}{"USERNAME": nil}, db.GetDatabaseTypeString(), "", false)

	// Use map with default and user defined columns
	columnList, _ := GetUserColumnsList()
	for _, column := range columnList {
		if column.ColumnName.Valid {
			columns += column.ColumnName.String + ","

			jsonName := strings.ToLower(column.ColumnName.String)
			jsonName = strings.ReplaceAll(jsonName, "_", "-")
		}
	}

	columns = strings.TrimSuffix(columns, ",")

	sql := `
	SELECT
		` + columns + `
	FROM ` + tablename + `
	WHERE ` + columnUsername.Names + ` = '` + username + `'`

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

// GetUsers - Get Users List
func GetUsers(options model.UserListingOptions) (model.UserPageData, error) {
	result := model.UserPageData{}

	dbtype := db.GetDatabaseTypeString()
	tableName := db.GetFullTableName(consts.TableUsers)

	// Use struct with default column
	column := np.CreateString(model.UserData{}, dbtype, "", false)
	columnUsername := np.CreateString(map[string]interface{}{"USERNAME": nil}, db.GetDatabaseTypeString(), "", false)
	columnEmail := np.CreateString(map[string]interface{}{"EMAIL": nil}, db.GetDatabaseTypeString(), "", false)
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, db.GetDatabaseTypeString(), "", false)

	sqlSearch := ""

	if options.Search.Valid && options.Search.String != "" {
		sqlSearch = `
		WHERE ` + columnUsername.Names + ` LIKE '%` + options.Search.String + `%'
			OR ` + columnEmail.Names + ` LIKE '%` + options.Search.String + `%'`
	}

	paging := ``
	if options.Page.Valid && options.ListCount.Valid {
		paging = db.Obj.GetPagingQuery(int(options.Page.Int64*options.ListCount.Int64), int(options.ListCount.Int64))
	}

	sql := `
	SELECT
		` + column.Names + `
	FROM ` + tableName + `
	` + sqlSearch + `
	ORDER BY ` + columnIdx.Names + ` ASC
	` + paging

	r, err := db.Con.Query(sql)
	if err != nil {
		return result, err
	}
	defer r.Close()

	var users []model.UserData
	err = scan.Rows(&users, r)
	if err != nil {
		return result, err
	}

	var totalCount int64
	sql = `
	SELECT
		COUNT(` + columnIdx.Names + `)
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

	totalPage := math.Ceil(float64(totalCount) / float64(options.ListCount.Int64))

	result = model.UserPageData{
		UserList:    users,
		CurrentPage: int(options.Page.Int64) + 1,
		TotalPage:   int(totalPage),
	}

	return result, nil
}

// GetUsersMap - Get Users List map
func GetUsersMap(options model.UserListingOptions) (model.UserPageData, error) {
	result := model.UserPageData{}

	tableName := db.GetFullTableName(consts.TableUsers)
	// columns := ""
	columnUsername := np.CreateString(map[string]interface{}{"USERNAME": nil}, db.GetDatabaseTypeString(), "", false)
	columnEmail := np.CreateString(map[string]interface{}{"EMAIL": nil}, db.GetDatabaseTypeString(), "", false)
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, db.GetDatabaseTypeString(), "", false)

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

	if options.Search.Valid && options.Search.String != "" {
		sqlSearch = `
		WHERE ` + columnUsername.Names + ` LIKE '%` + options.Search.String + `%'
			OR ` + columnEmail.Names + ` LIKE '%` + options.Search.String + `%'`
	}

	paging := ``
	if options.Page.Valid && options.ListCount.Valid {
		paging = db.Obj.GetPagingQuery(int(options.Page.Int64*options.ListCount.Int64), int(options.ListCount.Int64))
	}

	sql := `
	SELECT
		` + columns.Names + `
	FROM ` + tableName + `
	` + sqlSearch + `
	ORDER BY ` + columnIdx.Names + ` ASC
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
		COUNT(` + columnIdx.Names + `)
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

	totalPage := math.Ceil(float64(totalCount) / float64(options.ListCount.Int64))

	result = model.UserPageData{
		UserList:    users,
		CurrentPage: int(options.Page.Int64) + 1,
		TotalPage:   int(totalPage),
	}

	return result, nil
}

// GetUsersList - API Get Users List
func GetUsersList(search string) ([]model.UserData, error) {
	result := []model.UserData{}

	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)

	// Use struct with default columns
	columns := np.CreateString(model.UserData{}, dbtype, "", false).Names
	sqlSearch := ""
	columnUsername := np.CreateString(map[string]interface{}{"USERNAME": nil}, db.GetDatabaseTypeString(), "", false)
	columnEmail := np.CreateString(map[string]interface{}{"EMAIL": nil}, db.GetDatabaseTypeString(), "", false)

	if search != "" {
		sqlSearch = `
		WHERE ` + columnUsername.Names + ` LIKE '%` + search + `%'
			OR ` + columnEmail.Names + ` LIKE '%` + search + `%'`
	}

	sql := `
	SELECT
		` + columns + `
	FROM ` + tablename + `
	` + sqlSearch

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

// GetUsersListMap - API Get Users List map
func GetUsersListMap(search string) ([]map[string]interface{}, error) {
	result := []map[string]interface{}{}

	tablename := db.GetFullTableName(consts.TableUsers)
	columns := ""
	columnUsername := np.CreateString(map[string]interface{}{"USERNAME": nil}, db.GetDatabaseTypeString(), "", false)
	columnEmail := np.CreateString(map[string]interface{}{"EMAIL": nil}, db.GetDatabaseTypeString(), "", false)

	// Use map with default and user defined columns
	columnList, _ := GetUserColumnsList()
	for _, column := range columnList {
		if column.ColumnName.Valid {
			columns += column.ColumnName.String + ","

			jsonName := strings.ToLower(column.ColumnName.String)
			jsonName = strings.ReplaceAll(jsonName, "_", "-")
		}
	}

	columns = strings.TrimSuffix(columns, ",")
	sqlSearch := ""

	if search != "" {
		sqlSearch = `
		WHERE ` + columnUsername.Names + ` LIKE '%` + search + `%'
			OR ` + columnEmail.Names + ` LIKE '%` + search + `%'`
	}

	sql := `
	SELECT
		` + columns + `
	FROM ` + tablename + `
	` + sqlSearch

	r, err := db.Con.Query(sql)
	if err != nil {
		return nil, err
	}

	for r.Next() {
		scanedRow, _ := scanMap(r)
		result = append(result, scanedRow)
	}

	return result, nil
}

func AddUser(userColumn model.UserData) error {
	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)

	column := np.CreateString(userColumn, dbtype, "insert", false)

	// sql := `INSERT INTO ` + tablename + ` (` + columns + `) VALUES (` + values + `)`
	sql := `
	INSERT INTO ` + tablename + ` (
		` + column.Names + `
	) VALUES (
		` + column.Values + `
	)`

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func AddUserMap(userMap map[string]interface{}) error {
	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)

	column := np.CreateString(userMap, dbtype, "", false)

	sql := `
	INSERT INTO ` + tablename + ` (
		` + column.Names + `
	) VALUES (
		` + column.Values + `
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
		` + column.Names + `
	) VALUES (
		` + column.Values + `
	)`
	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUser(userColumn model.UserData) error {
	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)

	column := np.CreateString(userColumn, dbtype, "update", true)
	idx := strconv.Itoa(int(userColumn.Idx.Int64))

	colNames := strings.Split(column.Names, ",")
	colValues := strings.Split(column.Values, ",")
	holder := ""

	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)

	for i := 0; i < len(colNames); i++ {
		holder += colNames[i] + " = " + colValues[i] + ", "
	}
	holder = strings.TrimSuffix(holder, ", ")

	sql := `
	UPDATE ` + tablename + ` SET
		` + holder + `
	WHERE ` + columnIdx.Names + ` = ` + idx

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserMap(userMap map[string]interface{}) error {
	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)

	idx := userMap["idx"].(string)
	delete(userMap, "idx")

	column := np.CreateString(userMap, dbtype, "", false)
	directive, _, _ := np.CreateAssignHolders(dbtype, column.Names, 0)

	values := strings.Split(column.Values, ",")
	// values = append(values, idx)

	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, dbtype, "", false)

	valuesi := make([]interface{}, len(values))
	for i, v := range values {
		switch len(v) {
		case 1:
			valuesi[i] = v
		default:
			valuesi[i] = v[1 : len(v)-1]
		}
	}

	sql := `
	UPDATE ` + tablename + ` SET
		` + directive + `
	WHERE ` + columnIdx.Names + ` = ` + idx

	_, err := db.Con.Exec(sql, valuesi...)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(userColumn model.UserData) error {
	tablename := db.GetFullTableName(consts.TableUsers)

	idx := strconv.Itoa(int(userColumn.Idx.Int64))
	columnIdx := np.CreateString(map[string]interface{}{"IDX": nil}, db.GetDatabaseTypeString(), "", false)

	sql := `
	DELETE FROM ` + tablename + `
	WHERE ` + columnIdx.Names + ` = ` + idx

	_, err := db.Con.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
