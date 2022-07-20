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
		DatabaseName:  "9minutestest",
		SchemaName:    "dbo",
		TableName:     "users",
		UserTable:     "users",
		GrantID:       "sa",
		GrantPassword: "SQLServer1433",
	}

	var TableUserColumns string = "user_fields"

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
				sql:       "INSERT INTO " + GetFullTableName(Info.UserTable) + " (USERNAME,PASSWORD) VALUES (@p1,@p2)",
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

			dropUserTableSQL := `
			USE "` + Info.DatabaseName + `"
			IF OBJECT_ID('` + Info.UserTable + `','U') IS NOT NULL
			DROP TABLE "` + Info.UserTable + `"
			IF OBJECT_ID('` + TableUserColumns + `','U') IS NOT NULL
			DROP TABLE "` + TableUserColumns + `"
			-- GO`

			_, err = Con.Exec(dropUserTableSQL)

			err = Obj.CreateDB()
			if err != nil {
				t.Error(err)
				return
			}

			err = Obj.CreateUserTable()
			if err != nil {
				t.Error(err)
				return
			}

			count, _, err := Obj.Exec(tt.args.sql, tt.args.colValues, tt.args.options)
			if err != nil {
				t.Error(err)
				return
			}

			_, err = Con.Exec(`DROP DATABASE "` + Info.DatabaseName + `";`)
			if err != nil {
				t.Error(err)
				return
			}

			require.Equal(t, int64(1), count)
		})
	}
}
