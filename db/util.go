package db

import (
	"9minutes/model"
	"strings"
)

func GetFullTableName(table string) string {
	tablename := ""
	table = strings.ToLower(table)

	switch Info.DatabaseType {
	case model.SQLITE:
		tablename = `"` + table + `"`
	case model.MYSQL:
		tablename = "`" + Info.DatabaseName + "`.`" + table + "`"
	case model.POSTGRES:
		tablename = `"` + Info.SchemaName + `"."` + table + `"`
	case model.SQLSERVER:
		tablename = `"` + Info.DatabaseName + `"."` + Info.SchemaName + `"."` + table + `"`
	case model.ORACLE:
		tablename = `"` + strings.ToUpper(Info.GrantID) + `"."` + strings.ToUpper(Info.TableName) + `"`
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
	case model.SQLITE:
		dbtype = "sqlite"
	case model.MYSQL:
		dbtype = "mysql"
	case model.POSTGRES:
		dbtype = "postgres"
	case model.SQLSERVER:
		dbtype = "sqlserver"
	case model.ORACLE:
		dbtype = "oracle"
	}

	return dbtype
}

func QuotesName(data string) string {
	result := ""

	switch Info.DatabaseType {
	case model.SQLITE:
		data = strings.ReplaceAll(data, `"`, `""`)
		result = `"` + data + `"`
	case model.MYSQL:
		data = strings.ReplaceAll(data, "`", "``")
		result = "'" + data + "'"
	case model.POSTGRES:
		data = strings.ReplaceAll(data, `"`, `""`)
		result = `"` + data + `"`
	case model.SQLSERVER:
		data = strings.ReplaceAll(data, `"`, `""`)
		result = `"` + data + `"`
	case model.ORACLE:
		data = strings.ReplaceAll(data, `"`, `""`)
		result = `"` + data + `"`
	}

	return result
}

func QuotesValue(data string) string {
	result := ""

	switch Info.DatabaseType {
	case model.SQLITE:
		data = strings.ReplaceAll(data, "'", "''")
		result = "'" + data + "'"
	case model.MYSQL:
		data = strings.ReplaceAll(data, "'", "\\'")
		result = "'" + data + "'"
	case model.POSTGRES:
		data = strings.ReplaceAll(data, "'", "''")
		result = "'" + data + "'"
	case model.SQLSERVER:
		data = strings.ReplaceAll(data, "'", "''")
		result = "'" + data + "'"
	case model.ORACLE:
		data = strings.ReplaceAll(data, "'", "''")
		result = "'" + data + "'"
	}

	return result
}
