package np

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"gopkg.in/guregu/null.v4"
)

// ColumnStrings - strings for names and values
type ColumnStrings struct {
	Names  string
	Values string
}

var (
	TagName       = "db"
	TagNameNPSKIP = "npskip"
	Separator     = ","
)

func createString(o interface{}, dbtype, skipValue, separatorNames, separatorValues string, checkValid bool) ColumnStrings {
	names := ""
	values := ""

	ot := reflect.TypeOf(o)
	ov := reflect.ValueOf(o)

	switch ot.Kind() {
	case reflect.Struct:
		for i := 0; i < ov.NumField(); i++ {
			isValid := true
			skipTag := ov.Type().Field(i).Tag.Get(TagNameNPSKIP)
			// log.Println("skip::", ov.Field(i).Type().Name(), ov.Field(i).Type().PkgPath(), ov.Field(i).Type().Kind(), ov.Field(i).Type().String())
			if skipTag != "" && skipValue != "" && strings.Contains(skipTag, skipValue) {
				continue
			}

			switch true {
			case ov.Field(i).Type().PkgPath() == "database/sql":
				// sql.NullString, sql.NullInt..
				valid := ov.Field(i).Field(1).Bool()
				if valid {
					names += ov.Type().Field(i).Tag.Get(TagName)
					values += fmt.Sprint(ov.Field(i).Field(0))
				}

			case ov.Field(i).Kind() == reflect.Struct:
				// Maybe null.String, null.Int..
				valueStruct := createString(ov.Field(i).Interface(), dbtype, skipValue, separatorNames, separatorValues, checkValid)

				isValid = false
				if checkValid {
					ovInterface := ov.Field(i).Interface()

					switch ov.Field(i).Type() {
					case reflect.TypeOf(null.String{}):
						isValid = ovInterface.(null.String).Valid
					case reflect.TypeOf(null.Int{}):
						isValid = ovInterface.(null.Int).Valid
					case reflect.TypeOf(null.Float{}):
						isValid = ovInterface.(null.Float).Valid
					case reflect.TypeOf(null.Bool{}):
						isValid = ovInterface.(null.Bool).Valid
					case reflect.TypeOf(null.Time{}):
						isValid = ovInterface.(null.Time).Valid
					}
				} else {
					isValid = true
				}

				if isValid {
					if valueStruct.Names != "" {
						names += valueStruct.Names
					} else {
						names += ov.Type().Field(i).Tag.Get(TagName)
					}

					if valueStruct.Values != "" {
						values += valueStruct.Values
					}
				}

			case ov.Field(i).Kind() == reflect.Ptr:
				// Maybe pointer
				if !ov.Field(i).IsNil() {
					names += ov.Type().Field(i).Tag.Get(TagName)
					values += fmt.Sprint(ov.Field(i).Elem())
				}

			case ov.Field(i).Kind() == reflect.Interface:
				// Maybe interface
				// if ov.Type().Field(i).Tag.Get(TagName) != "" {
				if !ov.Field(i).IsNil() {
					names += ov.Type().Field(i).Tag.Get(TagName)
					values += fmt.Sprint(ov.Field(i).Elem())
				}

			default:
				// Maybe string, int.. How to do?
				unknownValue := fmt.Sprint(ov.Field(i).Field(0))
				log.Println("Maybe string, int..", unknownValue)
			}

			if isValid && i < ov.NumField()-1 {
				names += separatorNames
				values += separatorValues
			}
		}

	case reflect.Map:
		for i, k := range ov.MapKeys() {
			key := strings.ReplaceAll(k.String(), "-", "_")
			key = strings.ToUpper(key)
			names += key
			values += fmt.Sprint(ov.MapIndex(k))

			if i < ov.Len()-1 {
				names += separatorNames
				values += separatorValues
			}
		}

	}

	result := ColumnStrings{
		Names:  names,
		Values: values,
	}

	return result
}

// CreateString - create string from struct, map
func CreateString(o interface{}, dbtype, skipValue string, checkValid bool) ColumnStrings {
	quoteNames := ""
	quoteValues := ""
	separatorNames := Separator
	separatorValues := Separator
	switch dbtype {
	case "sqlite":
		quoteNames = `"`
		quoteValues = "'"
		separatorNames = `","`
		separatorValues = "','"
	case "mysql":
		quoteNames = "`"
		quoteValues = "'"
		separatorNames = "`,`"
		separatorValues = `','`
	case "postgres":
		quoteNames = `"`
		quoteValues = "'"
		separatorNames = `","`
		separatorValues = `','`
	case "sqlserver":
		quoteNames = `"`
		quoteValues = "'"
		separatorNames = `","`
		separatorValues = `','`
	case "oracle":
		quoteNames = `"`
		quoteValues = "'"
		separatorNames = `","`
		separatorValues = `','`
	}

	result := createString(o, dbtype, skipValue, separatorNames, separatorValues, checkValid)

	result.Names = strings.TrimSuffix(result.Names, separatorNames)
	result.Values = strings.TrimSuffix(result.Values, separatorValues)

	result.Names = quoteNames + result.Names + quoteNames
	result.Values = quoteValues + result.Values + quoteValues

	return result
}

// CreateWhereString - create string from struct, map
func CreateWhereString(o interface{}, dbtype, opValue, opCombine, skipValue string, checkValid bool) (result string) {
	if reflect.TypeOf(o).Kind() == reflect.Slice {
		switch objects := o.(type) {
		case []map[string]interface{}:
			for i, object := range objects {
				created := CreateString(object, dbtype, skipValue, checkValid)

				names := strings.Split(created.Names, ",")
				values := strings.Split(created.Values, ",")
				for j, name := range names {
					value := values[j]
					if opValue == "LIKE" {
						value = value[0:1] + "%" + value[1:len(value)-1] + "%" + value[len(value)-1:]
					}

					if i == 0 && j == 0 {
						result += " WHERE " + name + " " + opValue + " " + value
					} else {
						result += " " + opCombine + " " + name + " " + opValue + " " + value
					}
				}
			}
		default:
			// Todo - struct slice
			log.Println("CreateWhereString - Unknown o type", objects)
		}
	} else {
		created := CreateString(o, dbtype, skipValue, checkValid)

		names := strings.Split(created.Names, ",")
		values := strings.Split(created.Values, ",")
		for i, name := range names {
			value := values[i]
			if opValue == "LIKE" {
				value = value[0:1] + "%" + value[1:len(value)-1] + "%" + value[len(value)-1:]
			}

			if i == 0 {
				result += " WHERE " + name + " " + opValue + " " + value
			} else {
				result += " " + opCombine + " " + name + " " + opValue + " " + value
			}
		}
	}

	return result
}

func CreateUpdateString(o interface{}, dbtype, skipValue string, checkValid bool) (result string) {
	created := CreateString(o, dbtype, skipValue, checkValid)

	names := strings.Split(created.Names, ",")
	values := strings.Split(created.Values, ",")
	for i, name := range names {
		value := values[i]
		if i == 0 {
			result += " SET " + name + " = " + value
		} else {
			result += " , " + name + " = " + value
		}
	}

	return result
}

// CreateMapSlice
// create map which contains slice interface of names and values
// from struct or map
func CreateMapSlice(o interface{}, skipValue string) map[string][]interface{} {
	names := []interface{}{}
	values := []interface{}{}

	ot := reflect.TypeOf(o)
	ov := reflect.ValueOf(o)

	switch ot.Kind() {
	case reflect.Struct:
		for i := 0; i < ov.NumField(); i++ {
			skipTag := ov.Type().Field(i).Tag.Get(TagNameNPSKIP)
			if skipTag != "" && strings.Contains(skipTag, skipValue) {
				continue
			}

			switch true {
			case ov.Field(i).Type().PkgPath() == "database/sql":
				// sql.NullString, swl.NullInt..
				valid := ov.Field(i).Field(1).Bool()
				if valid {
					names = append(names, ov.Type().Field(i).Tag.Get(TagName))
					values = append(values, fmt.Sprint(ov.Field(i).Field(0)))
				}

			case ov.Field(i).Kind() == reflect.Struct:
				// Maybe null.String, null.Int..
				valueStruct := CreateMapSlice(ov.Field(i).Interface(), skipValue)
				if len(valueStruct["names"]) > 0 && valueStruct["names"][0] != "" {
					names = append(names, valueStruct["names"]...)
				} else {
					names = append(names, ov.Type().Field(i).Tag.Get(TagName))
				}
				if len(valueStruct["values"]) > 0 {
					values = append(values, valueStruct["values"]...)
				}

			case ov.Field(i).Kind() == reflect.Ptr:
				// Maybe pointer
				if !ov.Field(i).IsNil() {
					names = append(names, ov.Type().Field(i).Tag.Get(TagName))
					values = append(values, fmt.Sprint(ov.Field(i).Elem()))
				}

			case ov.Field(i).Kind() == reflect.Interface:
				// Maybe interface
				// if ov.Type().Field(i).Tag.Get(TagName) != "" {
				if !ov.Field(i).IsNil() {
					names = append(names, ov.Type().Field(i).Tag.Get(TagName))
					values = append(values, fmt.Sprint(ov.Field(i).Elem()))
				}

			default:
				// Maybe string, int.. How to do?
				unknownValue := fmt.Sprint(ov.Field(i).Field(0))
				log.Println("Maybe string, int..", unknownValue)
			}
		}

	case reflect.Map:
		for _, k := range ov.MapKeys() {
			key := strings.ReplaceAll(k.String(), "-", "_")
			key = strings.ToUpper(key)
			names = append(names, key)
			values = append(values, fmt.Sprint(ov.MapIndex(k)))
		}

	}

	result := map[string][]interface{}{
		"names":  names,
		"values": values,
	}

	return result
}

// CreateMap - create map from struct, map
func CreateMap(o interface{}, skipValue string) map[string]string {
	pairs := map[string]string{}
	name := ""
	value := ""

	ot := reflect.TypeOf(o)
	ov := reflect.ValueOf(o)

	switch ot.Kind() {
	case reflect.Struct:
		for i := 0; i < ov.NumField(); i++ {
			skipTag := ov.Type().Field(i).Tag.Get(TagNameNPSKIP)
			if skipTag != "" && strings.Contains(skipTag, skipValue) {
				continue
			}

			switch true {
			case ov.Field(i).Type().PkgPath() == "database/sql":
				// sql.NullString, swl.NullInt..
				valid := ov.Field(i).Field(1).Bool()
				if valid {
					name = ov.Type().Field(i).Tag.Get(TagName)
					value = fmt.Sprint(ov.Field(i).Field(0))

					pairs[name] = value
				}

			case ov.Field(i).Kind() == reflect.Struct:
				// Maybe null.String, null.Int..
				valueStruct := CreateMap(ov.Field(i).Interface(), skipValue)
				if len(valueStruct) > 0 {
					for _, v := range valueStruct {
						name := ov.Type().Field(i).Tag.Get(TagName)
						pairs[name] = v
					}
				}

			case ov.Field(i).Kind() == reflect.Ptr:
				// Maybe pointer
				if !ov.Field(i).IsNil() {
					name = ov.Type().Field(i).Tag.Get(TagName)
					value = fmt.Sprint(ov.Field(i).Elem())

					pairs[name] = value
				}

			case ov.Field(i).Kind() == reflect.Interface:
				// Maybe interface
				if ov.Type().Field(i).Tag.Get(TagName) != "" {
					name = ov.Type().Field(i).Tag.Get(TagName)
					value = fmt.Sprint(ov.Field(i).Elem())

					pairs[name] = value
				}

			default:
				// Maybe string, int.. How to do?
				unknownValue := fmt.Sprint(ov.Field(i).Field(0))
				log.Println("Maybe string, int..", unknownValue)
			}
		}

	case reflect.Map:
		for _, k := range ov.MapKeys() {
			key := strings.ReplaceAll(k.String(), "-", "_")
			key = strings.ToUpper(key)

			pairs[key] = fmt.Sprint(ov.MapIndex(k))
		}

	}

	return pairs
}
