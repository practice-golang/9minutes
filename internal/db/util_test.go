package db

import (
	"9minutes/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTableName(t *testing.T) {
	tests := []struct {
		name string
		db   int
		want string
	}{
		{
			name: "SQLITE",
			db:   model.SQLITE,
			want: "test_table",
		},
		{
			name: "MYSQL",
			db:   model.MYSQL,
			want: "test_table",
		},
		{
			name: "POSTGRES",
			db:   model.POSTGRES,
			want: "test_table",
		},
	}
	for _, tt := range tests {
		Info = DBInfo{
			DatabaseType: tt.db,
			DatabaseName: "test_db",
			SchemaName:   "test_schema",
			TableName:    "test_table",
		}

		t.Run(tt.name, func(t *testing.T) {
			got := GetTableName()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetFullTableName(t *testing.T) {
	tests := []struct {
		name string
		db   int
		want string
	}{
		{
			name: "SQLITE",
			db:   model.SQLITE,
			want: `"test_table"`,
		},
		{
			name: "MYSQL",
			db:   model.MYSQL,
			want: "`test_db`.`test_table`",
		},
		{
			name: "POSTGRES",
			db:   model.POSTGRES,
			want: `"test_schema"."test_table"`,
		},
	}
	for _, tt := range tests {
		Info = DBInfo{
			DatabaseType: tt.db,
			DatabaseName: "test_db",
			SchemaName:   "test_schema",
			TableName:    "test_table",
		}

		t.Run(tt.name, func(t *testing.T) {
			got := GetFullTableName(Info.TableName)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetDatabaseTypeString(t *testing.T) {
	tests := []struct {
		name string
		db   int
		want string
	}{
		{
			name: "SQLITE",
			db:   model.SQLITE,
			want: "sqlite",
		},
		{
			name: "MYSQL",
			db:   model.MYSQL,
			want: "mysql",
		},
		{
			name: "POSTGRES",
			db:   model.POSTGRES,
			want: "postgres",
		},
	}
	for _, tt := range tests {
		Info = DBInfo{DatabaseType: tt.db}

		t.Run(tt.name, func(t *testing.T) {
			got := GetDatabaseTypeString()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestQuotesName(t *testing.T) {
	tests := []struct {
		name string
		db   int
		want string
	}{
		{
			name: "SQLITE",
			db:   model.SQLITE,
			want: `"field_name"`,
		},
		{
			name: "MYSQL",
			db:   model.MYSQL,
			want: `'field_name'`,
		},
		{
			name: "POSTGRES",
			db:   model.POSTGRES,
			want: `"field_name"`,
		},
	}
	for _, tt := range tests {
		Info = DBInfo{DatabaseType: tt.db}

		t.Run(tt.name, func(t *testing.T) {
			got := QuotesName("field_name")
			require.Equal(t, tt.want, got)
		})
	}
}

func TestQuotesValue(t *testing.T) {
	tests := []struct {
		name string
		db   int
		want string
	}{
		{
			name: "SQLITE",
			db:   model.SQLITE,
			want: "'field_value'",
		},
		{
			name: "MYSQL",
			db:   model.MYSQL,
			want: "'field_value'",
		},
		{
			name: "POSTGRES",
			db:   model.POSTGRES,
			want: "'field_value'",
		},
	}
	for _, tt := range tests {
		Info = DBInfo{DatabaseType: tt.db}

		t.Run(tt.name, func(t *testing.T) {
			got := QuotesValue("field_value")
			require.Equal(t, tt.want, got)
		})
	}
}
