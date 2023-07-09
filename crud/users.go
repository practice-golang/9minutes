package crud

import (
	"9minutes/consts"
	"9minutes/db"
	"9minutes/model"
	"9minutes/np"
	"errors"
	"math"
	"strings"

	"github.com/blockloop/scan"
	"golang.org/x/crypto/bcrypt"
)

func GetUserByUsernameAndPassword(name, password string) (interface{}, error) {
	var result interface{}

	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)

	columns := ""
	wheres := np.CreateWhereString(map[string]interface{}{"USERNAME": name}, dbtype, "=", "AND", "", false)

	columnList, _ := GetUserColumnsList()
	for _, column := range columnList {
		if column.ColumnName.Valid {
			columns += column.ColumnName.String + ","
		}
	}

	columns = strings.TrimSuffix(columns, ",")

	sql := `
	SELECT
		` + columns + `
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
		return nil, errors.New("invalid username or password")
	}

	result = users[0]
	return result, nil
}

func GetUserByNameAndEmailMap(username, email string) (interface{}, error) {
	var result interface{}

	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)

	columns := ""
	whereUsername := np.CreateString(map[string]interface{}{"USERNAME": nil}, dbtype, "", false)
	whereEmail := np.CreateString(map[string]interface{}{"EMAIL": nil}, dbtype, "", false)
	whereGrade := np.CreateString(map[string]interface{}{"GRADE": nil}, dbtype, "", false)

	// Use map with default and user defined columns
	columnList, _ := GetUserColumnsList()
	for _, column := range columnList {
		if column.ColumnName.Valid {
			columns += column.ColumnName.String + ","
			// jsonName := strings.ReplaceAll(strings.ToLower(column.ColumnName.String), "_", "-")
		}
	}

	columns = strings.TrimSuffix(columns, ",")

	sql := `
	SELECT
		` + columns + `
	FROM ` + tablename + `
	WHERE ` + whereUsername.Names + ` = '` + username + `'
		AND ` + whereEmail.Names + ` = '` + email + `'
		AND ` + whereGrade.Names + ` != '` + "resigned_user" + `'`

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
			// jsonName := strings.ReplaceAll(strings.ToLower(column.ColumnName.String), "_", "-")
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

// GetUsersListMap - API Get Users List map
func GetUsersListMap(search string, page int) ([]map[string]interface{}, error) {
	result := []map[string]interface{}{}

	tablename := db.GetFullTableName(consts.TableUsers)
	columns := np.CreateString(
		map[string]interface{}{
			"IDX":      nil,
			"USERNAME": nil,
			"EMAIL":    nil,
			"GRADE":    nil,
			"REG_DTTM": nil,
		},
		db.GetDatabaseTypeString(), "", false,
	).Names

	wheres := ""
	if search != "" {
		wheres += np.CreateWhereString(
			map[string]interface{}{"USERNAME": search, "EMAIL": search},
			db.GetDatabaseTypeString(), "LIKE", "OR", "", false,
		)
	}

	sql := `
	SELECT
		` + columns + `
	FROM ` + tablename + `
	` + wheres

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

func UpdateUserMap(userDataMap map[string]interface{}) error {
	dbtype := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)

	userData := map[string]interface{}{}
	for k, v := range userDataMap {
		userData[k] = v
	}

	idx := userData["idx"].(string)
	delete(userData, "idx")

	updateset := np.CreateUpdateString(userData, dbtype, "", false)
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

func ResignUser(idx int64) error {
	dbType := db.GetDatabaseTypeString()
	tablename := db.GetFullTableName(consts.TableUsers)

	updateset := np.CreateUpdateString(
		map[string]interface{}{
			"GRADE":    "resigned_user",
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
