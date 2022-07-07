package np

import (
	"database/sql"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

// Human - guregu/null
type Human struct {
	Name         null.String `json:"name" db:"NAME"`
	Age          null.Int    `json:"age" db:"AGE"`
	EmailAddress null.String `json:"email-address,omitempty" db:"EMAIL_ADDRESS"`
}

// Agent - sql.Null..
type Agent struct {
	Name         sql.NullString `json:"name" db:"NAME"`
	Age          sql.NullInt64  `json:"age" db:"AGE"`
	EmailAddress sql.NullString `json:"email-address,omitempty" db:"EMAIL_ADDRESS"`
}

// Ander - pointer
type Ander struct {
	Name         *string `json:"name" db:"NAME"`
	Age          *int    `json:"age" db:"AGE"`
	EmailAddress *string `json:"email-address,omitempty" db:"EMAIL_ADDRESS"`
}

func TestCreateColString_null_struct(t *testing.T) {
	john := Human{
		null.NewString("John", true),
		null.NewInt(777, true),
		null.NewString("john@human.io", true),
	}

	colString := CreateString(john, "sqlite", "", false)

	names := strings.Split(colString.Names, ",")
	values := strings.Split(colString.Values, ",")

	nameSample := []string{`"NAME"`, `"AGE"`, `"EMAIL_ADDRESS"`}
	valueSample := []string{"'John'", "'777'", "'john@human.io'"}

	for i, name := range names {
		isExist := false
		for j, n := range nameSample {
			if n == name {
				if values[i] != valueSample[j] {
					t.Fatal("\nExpect:", strings.Join(valueSample, ","), "\nResult:", colString.Values)
				}

				isExist = true
				break
			}
		}
		if !isExist {
			t.Fatal("\nExpect:", strings.Join(nameSample, ","), "\nResult:", colString.Names)
		}
	}

	// log.Println()
	// log.Println("INSERT INTO table_name (" + colString.Names + ") VALUES(" + colString.Values + ")")
}

func TestCreateColString_json_struct(t *testing.T) {
	jane := Human{}
	janeJSON := []byte("{ \"name\": \"Jane\", \"age\": 999, \"email-address\": \"jane@human.io\" }")
	json.Unmarshal(janeJSON, &jane)

	colString := CreateString(jane, "sqlite", "", false)

	names := strings.Split(colString.Names, ",")
	values := strings.Split(colString.Values, ",")

	nameSample := []string{`"NAME"`, `"AGE"`, `"EMAIL_ADDRESS"`}
	valueSample := []string{"'Jane'", "'999'", "'jane@human.io'"}

	for i, name := range names {
		isExist := false
		for j, n := range nameSample {
			if n == name {
				if values[i] != valueSample[j] {
					t.Fatal("\nExpect:", strings.Join(valueSample, ","), "\nResult:", colString.Values)
				}
				isExist = true

				break
			}
		}
		if !isExist {
			t.Fatal("\nExpect:", strings.Join(nameSample, ","), "\nResult:", colString.Names)
		}
	}
}

func TestCreateColString_json_map(t *testing.T) {
	james := map[string]interface{}{}
	jamesJSON := []byte("{ \"name\": \"James\", \"age\": 888, \"email-address\": \"james@human.io\" }")
	json.Unmarshal(jamesJSON, &james)

	colString := CreateString(james, "sqlite", "", false)

	names := strings.Split(colString.Names, ",")
	values := strings.Split(colString.Values, ",")

	nameSample := []string{`"NAME"`, `"AGE"`, `"EMAIL_ADDRESS"`}
	valueSample := []string{"'James'", "'888'", "'james@human.io'"}

	for i, name := range names {
		isExist := false
		for j, n := range nameSample {
			if n == name {
				if values[i] != valueSample[j] {
					t.Fatal("\nExpect:", strings.Join(valueSample, ","), "\nResult:", colString.Values)
				}

				isExist = true
				break
			}
		}
		if !isExist {
			t.Fatal("\nExpect:", strings.Join(nameSample, ","), "\nResult:", colString.Names)
		}
	}
}

func Test_createColString_sql_null(t *testing.T) {
	smith := Agent{
		Name:         sql.NullString{String: "Smith", Valid: true},
		Age:          sql.NullInt64{Int64: int64(222), Valid: true},
		EmailAddress: sql.NullString{String: "smith@machine.io", Valid: true},
	}

	colString := CreateString(smith, "sqlite", "", false)

	names := strings.Split(colString.Names, ",")
	values := strings.Split(colString.Values, ",")

	nameSample := []string{`"NAME"`, `"AGE"`, `"EMAIL_ADDRESS"`}
	valueSample := []string{"'Smith'", "'222'", "'smith@machine.io'"}

	for i, name := range names {
		isExist := false
		for j, n := range nameSample {
			if n == name {
				if values[i] != valueSample[j] {
					t.Fatal("\nExpect:", strings.Join(valueSample, ","), "\nResult:", colString.Values)
				}

				isExist = true
				break
			}
		}
		if !isExist {
			t.Fatal("\nExpect:", strings.Join(nameSample, ","), "\nResult:", colString.Names)
		}
	}
}

func Test_createColString_pointer(t *testing.T) {
	thomas := Ander{}
	thomasJSON := []byte("{ \"name\": \"Thomas\", \"age\": 444, \"email-address\": \"thomas@son.io\" }")
	json.Unmarshal(thomasJSON, &thomas)

	colString := CreateString(thomas, "sqlite", "", false)

	names := strings.Split(colString.Names, ",")
	values := strings.Split(colString.Values, ",")

	nameSample := []string{`"NAME"`, `"AGE"`, `"EMAIL_ADDRESS"`}
	valueSample := []string{"'Thomas'", "'444'", "'thomas@son.io'"}

	for i, name := range names {
		isExist := false
		for j, n := range nameSample {
			if n == name {
				if values[i] != valueSample[j] {
					t.Fatal("\nExpect:", strings.Join(valueSample, ","), "\nResult:", colString.Values)
				}

				isExist = true
				break
			}
		}
		if !isExist {
			t.Fatal("\nExpect:", strings.Join(nameSample, ","), "\nResult:", colString.Names)
		}
	}
}

func TestCreateMapSlice(t *testing.T) {
	ptrName := "Smith"
	ptrAge := int(222)
	ptrEmailAddress := "smith@machine.io"

	type args struct {
		o         interface{}
		skipValue string
	}
	tests := []struct {
		name string
		args args
		want map[string][]interface{}
	}{
		{
			name: "sql null",
			args: args{
				o: Agent{
					Name:         sql.NullString{String: "Smith", Valid: true},
					Age:          sql.NullInt64{Int64: int64(222), Valid: true},
					EmailAddress: sql.NullString{String: "smith@machine.io", Valid: true},
				},
			},
			want: map[string][]interface{}{
				"names":  {"NAME", "AGE", "EMAIL_ADDRESS"},
				"values": {"Smith", "222", "smith@machine.io"},
			},
		},
		{
			name: "guregu null",
			args: args{
				o: Human{
					Name:         null.StringFrom("Smith"),
					Age:          null.IntFrom(int64(222)),
					EmailAddress: null.StringFrom("smith@machine.io"),
				},
			},
			want: map[string][]interface{}{
				"names":  {"NAME", "AGE", "EMAIL_ADDRESS"},
				"values": {"Smith", "222", "smith@machine.io"},
			},
		},
		{
			name: "pointer",
			args: args{
				o: Ander{
					Name:         &ptrName,
					Age:          &ptrAge,
					EmailAddress: &ptrEmailAddress,
				},
			},
			want: map[string][]interface{}{
				"names":  {"NAME", "AGE", "EMAIL_ADDRESS"},
				"values": {"Smith", "222", "smith@machine.io"},
			},
		},
		{
			name: "map",
			args: args{
				o: map[string]interface{}{
					"name":          "Smith",
					"age":           222,
					"email-address": "smith@machine.io",
				},
			},
			want: map[string][]interface{}{
				"names":  {"NAME", "AGE", "EMAIL_ADDRESS"},
				"values": {"Smith", "222", "smith@machine.io"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateMapSlice(tt.args.o, tt.args.skipValue)

			for i, n := range tt.want["names"] {
				catched := false
				for j, v := range got["names"] {
					if n == v {
						require.Equal(t, tt.want["values"][i], got["values"][j])
						catched = true
						break
					}
				}
				if !catched {
					require.Equal(t, tt.want, got)
				}
			}
		})
	}
}

func TestCreateMap(t *testing.T) {
	ptrName := "Smith"
	ptrAge := int(222)
	ptrEmailAddress := "smith@machine.io"

	type args struct {
		o         interface{}
		skipValue string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "sql null",
			args: args{
				o: Agent{
					Name:         sql.NullString{String: "Smith", Valid: true},
					Age:          sql.NullInt64{Int64: int64(222), Valid: true},
					EmailAddress: sql.NullString{String: "smith@machine.io", Valid: true},
				},
			},
			want: map[string]string{
				"NAME":          "Smith",
				"AGE":           "222",
				"EMAIL_ADDRESS": "smith@machine.io",
			},
		},
		{
			name: "guregu null",
			args: args{
				o: Human{
					Name:         null.StringFrom("Smith"),
					Age:          null.IntFrom(int64(222)),
					EmailAddress: null.StringFrom("smith@machine.io"),
				},
			},
			want: map[string]string{
				"NAME":          "Smith",
				"AGE":           "222",
				"EMAIL_ADDRESS": "smith@machine.io",
			},
		},
		{
			name: "pointer",
			args: args{
				o: Ander{
					Name:         &ptrName,
					Age:          &ptrAge,
					EmailAddress: &ptrEmailAddress,
				},
			},
			want: map[string]string{
				"NAME":          "Smith",
				"AGE":           "222",
				"EMAIL_ADDRESS": "smith@machine.io",
			},
		},
		{
			name: "map",
			args: args{
				o: map[string]interface{}{
					"name":          "Smith",
					"age":           222,
					"email-address": "smith@machine.io",
				},
			},
			want: map[string]string{
				"NAME":          "Smith",
				"AGE":           "222",
				"EMAIL_ADDRESS": "smith@machine.io",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateMap(tt.args.o, tt.args.skipValue)
			for k, v := range tt.want {
				require.Equal(t, v, got[k])
			}
		})
	}
}
