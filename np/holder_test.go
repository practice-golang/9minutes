package np

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateHolders(t *testing.T) {
	type args struct {
		dbtype   string
		colNames string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "POSTGRES",
			args: args{
				dbtype:   "postgres",
				colNames: "col1,col2,col3",
			},
			want: "$1,$2,$3",
		},
		{
			name: "SQLITE",
			args: args{
				dbtype:   "sqlite",
				colNames: "col1,col2,col3",
			},
			want: "?,?,?",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateHolders(tt.args.dbtype, tt.args.colNames)
			if err != nil {
				t.Error(err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCreateUpdateHolders(t *testing.T) {
	type args struct {
		dbtype  string
		columns interface{}
		offset  int
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 int
	}{
		{
			name: "POSTGRES",
			args: args{
				dbtype:  "postgres",
				columns: "col1,col2,col3",
				offset:  0,
			},
			want:  "col1=$1,col2=$2,col3=$3",
			want1: 3,
		},
		{
			name: "SQLITE",
			args: args{
				dbtype:  "sqlite",
				columns: "col1,col2,col3",
				offset:  0,
			},
			want:  "col1=?,col2=?,col3=?",
			want1: 3,
		},
		{
			name: "string slice",
			args: args{
				dbtype:  "sqlite",
				columns: []string{"col1", "col2", "col3"},
				offset:  0,
			},
			want:  "col1=?,col2=?,col3=?",
			want1: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			holders, count, err := CreateAssignHolders(tt.args.dbtype, tt.args.columns, tt.args.offset)
			if err != nil {
				t.Error(err)
				return
			}

			require.Equal(t, tt.want, holders)
			require.Equal(t, tt.want1, count)
		})
	}
}
