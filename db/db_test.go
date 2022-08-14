package db

import (
	"9minutes/model"
	"os"
	"testing"
)

func TestSetupDB(t *testing.T) {
	var err error
	tests := []struct {
		name string
		info DBInfo
	}{
		{
			name: "SQLITE",
			info: DBInfo{
				DatabaseType: model.SQLITE,
				DatabaseName: "books",
				TableName:    "books",
				FilePath:     "../books.db",
			},
		},
		{
			name: "MYSQL",
			info: DBInfo{
				DatabaseType:  model.MYSQL,
				Protocol:      "tcp",
				Addr:          "localhost",
				Port:          "13306",
				DatabaseName:  "9minutes",
				SchemaName:    "",
				TableName:     "users",
				GrantID:       "root",
				GrantPassword: "",
			},
		},
		{
			name: "POSTGRES",
			info: DBInfo{
				DatabaseType:  model.POSTGRES,
				Protocol:      "tcp",
				Addr:          "localhost",
				Port:          "5432",
				DatabaseName:  "postgres",
				SchemaName:    "public",
				TableName:     "books",
				GrantID:       "root",
				GrantPassword: "pgsql",
			},
		},
		{
			name: "SQLSERVER",
			info: DBInfo{
				DatabaseType:  model.SQLSERVER,
				Protocol:      "tcp",
				Addr:          "localhost",
				Port:          "1433",
				DatabaseName:  "mysitedb",
				SchemaName:    "dbo",
				TableName:     "books",
				GrantID:       "sa",
				GrantPassword: "SQLServer1433",
			},
		},
		{
			name: "NOTHING",
			info: DBInfo{
				DatabaseType: 999,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Info = tt.info
			err = SetupDB()
			if tt.name == "NOTHING" {
				// It's OK
				if err.Error() == "database type not supported" {
					return
				}
				t.Error(err)
			}

			if err != nil {
				t.Fatal(err)
			}

			Con.Close()
			if tt.info.DatabaseType == model.SQLITE {
				os.Remove(tt.info.FilePath)
			}
		})
	}
}
