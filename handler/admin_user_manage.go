package handler

import (
	"9minutes/consts"
	"9minutes/internal/crud"
	"9minutes/model"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

func GetUserList(c *fiber.Ctx) error {
	queries := c.Queries()

	search := strings.TrimSpace(queries["search"])
	page := 1
	if queries["page"] != "" {
		page, _ = strconv.Atoi(queries["page"])
		if page <= 0 {
			page = 1
		}
	}
	listCount := 10
	if queries["list-count"] != "" {
		listCount, _ = strconv.Atoi(queries["list-count"])
	}

	listingOption := model.UserListingOptions{
		Search:    null.StringFrom(search),
		Page:      null.IntFrom(int64(page)),
		ListCount: null.IntFrom(int64(listCount)),
	}

	/* Todo: Move to setup */
	columnNames, _ := crud.GetUserColumnsList()
	selectUserColumnsMap := map[string]interface{}{}
	for _, c := range columnNames {
		if c.ColumnName.Valid && c.ColumnName.String != "PASSWORD" {
			selectUserColumnsMap[c.ColumnName.String] = nil
		}
	}
	/* Todo: Move to setup */

	result, err := crud.GetUsersListMap(selectUserColumnsMap, listingOption)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	for i := range result.UserList {
		result.UserList[i]["password"] = ""
	}

	return c.Status(http.StatusOK).JSON(result)
}

func AddUser(c *fiber.Ctx) error {
	var err error

	now := time.Now().Format("20060102150405")

	userData := make(map[string]interface{})

	err = c.BodyParser(&userData)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	password, err := bcrypt.GenerateFromPassword([]byte(userData["password"].(string)), consts.BcryptCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	userData["password"] = string(password)
	userData["regdate"] = now

	err = crud.AddUserMap(userData)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	result := map[string]string{
		"result": "ok",
	}

	return c.Status(http.StatusOK).JSON(result)
}

func UpdateUser(c *fiber.Ctx) error {
	var err error

	userDatas := []map[string]interface{}{}
	userDatasSuccess := []map[string]interface{}{}
	userDatasFailed := []map[string]interface{}{}

	err = c.BodyParser(&userDatas)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	for _, userData := range userDatas {
		if _, ok := userData["password"]; ok {
			password, err := bcrypt.GenerateFromPassword([]byte(userData["password"].(string)), consts.BcryptCost)
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}
			userData["password"] = string(password)
		}

		err = crud.UpdateUserMap(userData)
		if err != nil {
			responseData := map[string]interface{}{"data": userData, "error": err.Error()}
			userDatasFailed = append(userDatasFailed, responseData)
			continue
		}
		responseData := map[string]interface{}{"data": userData, "error": ""}
		userDatasSuccess = append(userDatasSuccess, responseData)
	}

	result := map[string]interface{}{"result": "ok"}
	if len(userDatasFailed) > 0 {
		result["result"] = "failed"
		result["failed"] = userDatasFailed
		result["success"] = userDatasSuccess
	}

	return c.Status(http.StatusOK).JSON(result)
}

func DeleteUser(c *fiber.Ctx) error {
	isDelete := false
	queries := c.Queries()
	if strings.TrimSpace(queries["mode"]) == "delete" {
		isDelete = true
	}

	userDatas := []map[string]interface{}{}
	userDatasSuccess := []map[string]interface{}{}
	userDatasFailed := []map[string]interface{}{}

	err := c.BodyParser(&userDatas)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	for _, userData := range userDatas {
		idx, err := strconv.Atoi(userData["idx"].(string))
		if err != nil {
			responseData := map[string]interface{}{"data": userData, "error": err.Error()}
			userDatasFailed = append(userDatasFailed, responseData)
			continue
		}

		if isDelete {
			err = crud.DeleteUser(int64(idx))
		} else {
			err = crud.ResignUser(int64(idx))
		}

		if err != nil {
			responseData := map[string]interface{}{"data": userData, "error": err.Error()}
			userDatasFailed = append(userDatasFailed, responseData)
			continue
		}
		responseData := map[string]interface{}{"data": userData, "error": ""}
		userDatasSuccess = append(userDatasSuccess, responseData)
	}

	result := map[string]interface{}{"result": "ok"}
	if len(userDatasFailed) > 0 {
		result["result"] = "failed"
		result["failed"] = userDatasFailed
		result["success"] = userDatasSuccess
	}

	return c.Status(http.StatusOK).JSON(result)
}
