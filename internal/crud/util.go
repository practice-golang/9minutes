package crud

import (
	"9minutes/internal/db"
	"9minutes/model"
	"database/sql"
	"strings"
)

func scanMap(r *sql.Rows) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	columns, err := r.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	err = r.Scan(scanArgs...)
	if err != nil {
		return nil, err
	}

	for i, col := range columns {
		jsonColumnName := strings.ToLower(col)
		jsonColumnName = strings.ReplaceAll(jsonColumnName, "_", "-")

		switch db.Info.DatabaseType {
		case model.MYSQL:
			if values[i] != nil {
				result[jsonColumnName] = string(values[i].([]byte))
			} else {
				result[jsonColumnName] = ""
			}
		default:
			result[jsonColumnName] = values[i]
		}
	}

	return result, nil
}
