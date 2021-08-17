package db

import (
	"errors"
	"log"
	"reflect"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/go-cmp/cmp"
	"github.com/practice-golang/9minutes/models"
	"github.com/thoas/go-funk"
	"gopkg.in/guregu/null.v4"
)

// PrepareWhere - Prepare where condition for search
func PrepareWhere(data interface{}) goqu.Ex {
	result := goqu.Ex{}

	values := reflect.ValueOf(data)
	dataReflect := values.Type()

	for i := 0; i < values.NumField(); i++ {
		f := values.Field(i).Interface()
		b := dataReflect.Field(i)
		// dbName := b.Tag.Get("json")
		dbName := b.Tag.Get("db")
		switch f := f.(type) {
		case null.String:
			if f.Valid {
				result[dbName], _ = f.Value()
			}
		case null.Int:
			if f.Valid {
				result[dbName], _ = f.Value()
			}
		case null.Float:
			if f.Valid {
				result[dbName], _ = f.Value()
			}
		}
	}

	log.Println("Prepare Where: ", result)

	return result
}

// CheckValidAndPrepareWhere - Prepare where and check missing values for update
func CheckValidAndPrepareWhere(data interface{}) (goqu.Ex, error) {
	result := goqu.Ex{}

	values := reflect.ValueOf(data)
	dataReflect := values.Type()

	for i := 0; i < values.NumField(); i++ {
		f := values.Field(i).Interface()
		b := dataReflect.Field(i)
		// log.Println("Check where struct: ", b.Name)
		// dbName := b.Tag.Get("json")
		dbName := b.Tag.Get("db")
		switch f := f.(type) {
		case null.String:
			if !f.Valid {
				return nil, errors.New("`" + dbName + "` must have a value")
			}
			if funk.Contains(UpdateScope, dbName) {
				result[dbName], _ = f.Value()
			}
		case null.Int:
			if !f.Valid {
				return nil, errors.New("`" + dbName + "` must have a value")
			}
			if funk.Contains(UpdateScope, dbName) {
				result[dbName], _ = f.Value()
			}
		case null.Float:
			if !f.Valid {
				return nil, errors.New("`" + dbName + "` must have a value")
			}
			if funk.Contains(UpdateScope, dbName) {
				result[dbName], _ = f.Value()
			}
		}
	}

	return result, nil
}

func diffCustomBoardFields(old, new []map[string]interface{}) (add, remove, modify []map[string]interface{}) {
	var diff []map[string]interface{}

	for i := 0; i < 2; i++ {
		for _, s1 := range old {
			found := false
			for _, s2 := range new {
				if s1["idx"] == s2["idx"] {
					if i == 0 && !cmp.Equal(s1, s2) {
						modify = append(modify, s2)
					}

					found = true
					break
				}
			}

			if !found {
				diff = append(diff, s1)
			}
		}

		if i == 0 {
			remove = diff
			old, new = new, old
		} else {
			add = diff
		}

		diff = []map[string]interface{}{}
	}

	return
}

func diffUserTableFields(fieldsInfoOld, fieldsInfoNew []models.UserColumn) (add, remove, modify []models.UserColumn) {
	var diff []models.UserColumn

	old := fieldsInfoOld
	new := fieldsInfoNew

	for i := 0; i < 2; i++ {
		for _, s1 := range old {
			found := false
			for _, s2 := range new {
				if s1.Idx == s2.Idx {
					if i == 0 && !cmp.Equal(s1, s2) {
						modify = append(modify, s2)
					}

					found = true
					break
				}
			}

			if !found {
				diff = append(diff, s1)
			}
		}

		if i == 0 {
			remove = diff
			old, new = new, old
		} else {
			add = diff
		}

		diff = []models.UserColumn{}
	}

	return
}
