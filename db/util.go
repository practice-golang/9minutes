package db

import "strings"

func GetFullTableName(table string) string {
	tablename := ""
	table = strings.ToLower(table)

	switch Info.DatabaseType {
	case SQLITE:
		tablename = `"` + table + `"`
	case MYSQL:
		tablename = "`" + Info.DatabaseName + "`.`" + table + "`"
	case POSTGRES:
		tablename = `"` + Info.SchemaName + `"."` + table + `"`
	case SQLSERVER:
		tablename = `"` + Info.DatabaseName + `"."` + Info.SchemaName + `"."` + table + `"`
	}

	return tablename
}

func GetTableName() string {
	tablename := Info.TableName

	GetFullTableName(tablename)

	return tablename
}

func GetDatabaseTypeString() string {
	dbtype := ""

	switch Info.DatabaseType {
	case SQLITE:
		dbtype = "sqlite"
	case MYSQL:
		dbtype = "mysql"
	case POSTGRES:
		dbtype = "postgres"
	case SQLSERVER:
		dbtype = "sqlserver"
	}

	return dbtype
}

func QuotesName(data string) string {
	result := ""

	switch Info.DatabaseType {
	case SQLITE:
		data = strings.ReplaceAll(data, `"`, `""`)
		result = `"` + data + `"`
	case MYSQL:
		data = strings.ReplaceAll(data, "`", "``")
		result = "'" + data + "'"
	case POSTGRES:
		data = strings.ReplaceAll(data, `"`, `""`)
		result = `"` + data + `"`
	case SQLSERVER:
		data = strings.ReplaceAll(data, `"`, `""`)
		result = `"` + data + `"`
	}

	return result
}

func QuotesValue(data string) string {
	result := ""

	switch Info.DatabaseType {
	case SQLITE:
		data = strings.ReplaceAll(data, "'", "''")
		result = "'" + data + "'"
	case MYSQL:
		data = strings.ReplaceAll(data, "'", "\\'")
		result = "'" + data + "'"
	case POSTGRES:
		data = strings.ReplaceAll(data, "'", "''")
		result = "'" + data + "'"
	case SQLSERVER:
		data = strings.ReplaceAll(data, "'", "''")
		result = "'" + data + "'"
	}

	return result
}
