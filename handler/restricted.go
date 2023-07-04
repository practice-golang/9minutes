package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func RestrictedHello(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	name := sess.Get("name")
	if name == nil {
		return c.Status(http.StatusForbidden).Send([]byte("Unauthorized"))
	}

	return c.Status(http.StatusOK).Send([]byte("Hello " + name.(string)))
}

// func GetUserColumns(c *router.Context) {
// 	result, err := crud.GetUserColumnsList()
// 	if err != nil {
// 		c.Text(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	c.Json(http.StatusOK, result)
// }

// func AddUserColumn(c *router.Context) {
// 	var userColumn model.UserColumn

// 	err := json.NewDecoder(c.Body).Decode(&userColumn)
// 	if err != nil {
// 		c.Text(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	err = crud.AddUserColumn(userColumn)
// 	if err != nil {
// 		c.Text(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	result := map[string]string{
// 		"result": "ok",
// 	}

// 	c.Json(http.StatusOK, result)
// }

// func UpdateUserColumn(c *router.Context) {
// 	var userColumn model.UserColumn

// 	err := json.NewDecoder(c.Body).Decode(&userColumn)
// 	if err != nil {
// 		c.Text(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	err = crud.UpdateUserColumn(userColumn)
// 	if err != nil {
// 		c.Text(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	result := map[string]string{
// 		"result": "ok",
// 	}

// 	c.Json(http.StatusOK, result)
// }

// func DeleteUserColumn(c *router.Context) {
// 	var userColumn model.UserColumn

// 	uri := strings.Split(c.URL.Path, "/")
// 	idx, _ := strconv.Atoi(uri[len(uri)-1])

// 	userColumn.Idx = null.IntFrom(int64(idx))

// 	err := crud.DeleteUserColumn(userColumn)
// 	if err != nil {
// 		c.Text(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	result := map[string]string{
// 		"result": "ok",
// 	}

// 	c.Json(http.StatusOK, result)
// }

// func GetUsersList(c *router.Context) {
// 	// Use struct with default columns or map with default and user defined columns
// 	columnsCount, _ := crud.GetUserColumnsCount()
// 	// columnsCount, _ := db.Obj.GetColumnCount(db.Info.UserTable)

// 	queries := c.URL.Query()
// 	search := queries.Get("search")

// 	switch columnsCount {
// 	case model.UserDataFieldCount:
// 		result, err := crud.GetUsersList(search)
// 		if err != nil {
// 			c.Text(http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		c.Json(http.StatusOK, result)
// 	default:
// 		result, err := crud.GetUsersListMap(search)
// 		if err != nil {
// 			c.Text(http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		c.Json(http.StatusOK, result)
// 	}
// }

// func AddUser(c *router.Context) {
// 	var err error

// 	now := time.Now().Format("20060102150405")
// 	columnsCount, _ := crud.GetUserColumnsCount()

// 	switch columnsCount {
// 	case model.UserDataFieldCount:
// 		var userData model.UserData

// 		err = json.NewDecoder(c.Body).Decode(&userData)
// 		if err != nil {
// 			c.Text(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		password, err := bcrypt.GenerateFromPassword([]byte(userData.Password.String), consts.BcryptCost)
// 		if err != nil {
// 			c.Text(http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		userData.Password = null.StringFrom(string(password))

// 		userData.RegDTTM = null.StringFrom(now)

// 		err = crud.AddUser(userData)
// 		if err != nil {
// 			c.Text(http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 	default:
// 		userData := make(map[string]interface{})

// 		err = json.NewDecoder(c.Body).Decode(&userData)
// 		if err != nil {
// 			c.Text(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		password, err := bcrypt.GenerateFromPassword([]byte(userData["password"].(string)), consts.BcryptCost)
// 		if err != nil {
// 			c.Text(http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		userData["password"] = string(password)

// 		userData["reg-dttm"] = now

// 		err = crud.AddUserMap(userData)
// 		if err != nil {
// 			c.Text(http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 	}

// 	result := map[string]string{
// 		"result": "ok",
// 	}

// 	c.Json(http.StatusOK, result)
// }

// func UpdateUser(c *router.Context) {
// 	var err error
// 	columnsCount, _ := crud.GetUserColumnsCount()

// 	switch columnsCount {
// 	case model.UserDataFieldCount:
// 		var userData model.UserData

// 		err = json.NewDecoder(c.Body).Decode(&userData)
// 		if err != nil {
// 			c.Text(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		if userData.Password.Valid {
// 			password, err := bcrypt.GenerateFromPassword([]byte(userData.Password.String), consts.BcryptCost)
// 			if err != nil {
// 				c.Text(http.StatusInternalServerError, err.Error())
// 				return
// 			}
// 			userData.Password = null.StringFrom(string(password))
// 		}

// 		err = crud.UpdateUser(userData)
// 	default:
// 		userData := make(map[string]interface{})

// 		err = json.NewDecoder(c.Body).Decode(&userData)
// 		if err != nil {
// 			c.Text(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		if _, ok := userData["password"]; ok {
// 			password, err := bcrypt.GenerateFromPassword([]byte(userData["password"].(string)), consts.BcryptCost)
// 			if err != nil {
// 				c.Text(http.StatusInternalServerError, err.Error())
// 				return
// 			}
// 			userData["password"] = string(password)
// 		}

// 		err = crud.UpdateUserMap(userData)
// 	}

// 	if err != nil {
// 		c.Text(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	result := map[string]string{
// 		"result": "ok",
// 	}

// 	c.Json(http.StatusOK, result)
// }

// func DeleteUser(c *router.Context) {
// 	var userData model.UserData

// 	uri := strings.Split(c.URL.Path, "/")
// 	idx, _ := strconv.Atoi(uri[len(uri)-1])

// 	userData.Idx = null.IntFrom(int64(idx))

// 	err := crud.DeleteUser(userData)
// 	if err != nil {
// 		c.Text(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	result := map[string]string{
// 		"result": "ok",
// 	}

// 	c.Json(http.StatusOK, result)
// }
