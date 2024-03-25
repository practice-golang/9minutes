package db

import (
	"9minutes/model"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestPostgres_Exec(t *testing.T) {
	var err error
	Info = DBInfo{
		DatabaseType:  model.POSTGRES,
		Protocol:      "tcp",
		Addr:          "localhost",
		Port:          "5432",
		DatabaseName:  "postgres",
		SchemaName:    "public",
		TableName:     "users",
		UserTable:     "users",
		GrantID:       "root",
		GrantPassword: "pgsql",
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
			name: "POSTGRES",
			args: args{
				sql:       `INSERT INTO ` + GetFullTableName(Info.TableName) + ` ("USERID","PASSWORD") VALUES ($1,$2)`,
				colValues: []interface{}{"test2", "test3"},
				options:   "IDX",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := `host=` + Info.Addr + ` port=` + Info.Port + ` user=` + Info.GrantID + ` password=` + Info.GrantPassword + ` dbname=` + Info.DatabaseName + ` sslmode=disable`
			Obj = &Postgres{dsn: dsn}
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
			_, err = Con.Exec(`DROP TABLE IF EXISTS ` + GetFullTableName(TableUserColumns) + `;`)

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

			_, err = Con.Exec(`DROP TABLE IF EXISTS ` + GetFullTableName(Info.UserTable) + `;`)
			if err != nil {
				t.Error(err)
				return
			}

			_, err = Con.Exec(`DROP TABLE IF EXISTS ` + GetFullTableName(TableUserColumns) + `;`)
			if err != nil {
				t.Error(err)
				return
			}

			require.Equal(t, int64(1), count)
		})
	}
}
