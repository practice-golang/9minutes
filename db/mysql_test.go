package db

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestMysql_Exec(t *testing.T) {
	var err error
	Info = DBInfo{
		DatabaseType:  MYSQL,
		Protocol:      "tcp",
		Addr:          "localhost",
		Port:          "13306",
		DatabaseName:  "9minutestest",
		SchemaName:    "",
		TableName:     "users",
		UserTable:     "users",
		GrantID:       "root",
		GrantPassword: "",
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
			name: "MYSQL_INSERT",
			args: args{
				sql:       "INSERT INTO " + GetFullTableName(Info.UserTable) + " (USERNAME,PASSWORD) VALUES (?,?)",
				colValues: []interface{}{"test2", "test3"},
				options:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := Info.GrantID + ":" + Info.GrantPassword + "@" + Info.Protocol + "(" + Info.Addr + ":" + Info.Port + ")/"
			Obj = &Mysql{dsn: dsn}
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

			_, err = Con.Exec(`DROP TABLE IF EXISTS ` + GetFullTableName(Info.UserTable) + `;`)

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

			_, err = Con.Exec(`DROP DATABASE IF EXISTS ` + Info.DatabaseName + `;`)

			require.Equal(t, int64(1), count)
		})
	}
}
