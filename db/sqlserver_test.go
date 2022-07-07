package db

import (
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/stretchr/testify/require"
)

func TestSqlServer_Exec(t *testing.T) {
	var err error
	Info = DBInfo{
		DatabaseType:  SQLSERVER,
		Protocol:      "tcp",
		Addr:          "localhost",
		Port:          "1433",
		DatabaseName:  "mysitedb",
		SchemaName:    "dbo",
		TableName:     "books",
		GrantID:       "sa",
		GrantPassword: "SQLServer1433",
	}

	type args struct {
		sql       string
		colValues []interface{}
		options   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "SQLSERVER",
			args: args{
				sql:       "INSERT INTO " + GetFullTableName(Info.TableName) + " (TITLE,AUTHOR) VALUES (@p1,@p2)",
				colValues: []interface{}{"test2", "test3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := "sqlserver://" + Info.GrantID + ":" + Info.GrantPassword + "@" + Info.Addr + ":" + Info.Port + "?" + Info.DatabaseName + "&connction+timeout=30"
			Obj = &SqlServer{dsn: dsn}
			Con, err = Obj.connect()
			if err != nil {
				t.Error(err)
				return
			}
			defer Con.Close()

			err = Obj.CreateDB()
			if err != nil {
				t.Error(err)
				return
			}

			// err = Obj.CreateTable()
			// if err != nil {
			// 	t.Error(err)
			// 	return
			// }

			count, _, err := Obj.Exec(tt.args.sql, tt.args.colValues, tt.args.options)
			if err != nil {
				t.Error(err)
				return
			}

			require.Equal(t, int64(1), count)
		})
	}
}
