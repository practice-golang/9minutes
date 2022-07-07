package auth

import (
	"log"
	"reflect"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

// ConvertToNullTypeHookFunc - https://github.com/mitchellh/mapstructure/issues/164
func ConvertToNullTypeHookFunc(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	nullTypes := []reflect.Kind{reflect.String, reflect.Int64}

	isNullTypes := false
	for _, v := range nullTypes {
		if f.Kind() != v {
			isNullTypes = true
			break
		}
	}

	if !isNullTypes {
		return data, nil
	}

	switch t {
	case reflect.TypeOf(null.String{}):
		d := null.NewString(data.(string), true)
		return d, nil
	case reflect.TypeOf(null.Int{}):
		switch data.(type) {
		case int:
			d := null.NewInt(int64(data.(int)), true)
			return d, nil
		case int64:
			d := null.NewInt(data.(int64), true)
			return d, nil
		case float64:
			d := null.NewInt(int64(data.(float64)), true)
			return d, nil
		}
	case reflect.TypeOf(null.Float{}):
		d := null.NewFloat(data.(float64), true)
		return d, nil
	}

	return data, nil
}

func ComparePassword(password, hashedPassword string) bool {
	result := false

	encryptionErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if encryptionErr == nil {
		result = true
	}

	if encryptionErr != nil {
		log.Println("auth/util ComparePassword: ", encryptionErr)
	}

	return result
}
