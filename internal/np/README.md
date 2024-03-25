# NP

Create comma separated names and values from `struct` or `map` for raw SQL query.

## Which types of
* `database/sql` Nullnnnn
* `github.com/guregu/null`
* Map - Convert kebab case keys to upper case with underscore
* Pointer

## Example

See `np_test.go`

```go
package main

import (
	"github.com/practice-golang/9minutes/np"
)

// Human
type Human struct {
	Name         null.String `json:"name" db:"NAME"`
	Age          null.Int    `json:"age"  db:"AGE"`
	EmailAddress null.String `json:"email-address,omitempty" db:"EMAIL_ADDRESS"`
}

func main() {
	john := Human{
		null.NewString("John", true),
		null.NewInt(777, true),
		null.NewString("john@human.io", true),
	}

	np.TagName = "db"
	np.Separator = ","

	colString := np.CreateString(john, "mysql", "", false)

	selectQuery := "SELECT " + colString.Names + " FROM table_name;"
	insertQuery := "INSERT INTO table_name (" + colString.Names + ") VALUES(" + colString.Values + ");"

	fmt.Println(selectQuery)
	fmt.Println(insertQuery)

	whereString := np.CreateWhereString(john, "mysql", "", false)
	whereQuery := whereString

	fmt.Println(whereQuery)
}

// Result:
// SELECT NAME,AGE,EMAIL_ADDRESS FROM table_name;
// INSERT INTO table_name (NAME,AGE,EMAIL_ADDRESS) VALUES(John,777,john@human.io)
// WHERE `NAME`='John' AND `AGE`='777' AND `EMAIL_ADDRESS`='john@human.io'
```
