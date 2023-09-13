package handler

import (
	"9minutes/consts"
	"9minutes/internal/crud"
	"9minutes/internal/email"
	"9minutes/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
)

func GetUserListAPI(c *fiber.Ctx) error {
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

func AddUserAPI(c *fiber.Ctx) error {
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

func UpdateUserAPI(c *fiber.Ctx) error {
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

func DeleteUserAPI(c *fiber.Ctx) error {
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
		// idx, err := strconv.Atoi(userData["idx"].(string))
		// if err != nil {
		// 	responseData := map[string]interface{}{"data": userData, "error": err.Error()}
		// 	userDatasFailed = append(userDatasFailed, responseData)
		// 	continue
		// }
		idx := int64(userData["idx"].(float64))

		if isDelete {
			err = crud.DeleteUser(idx)
		} else {
			err = crud.ResignUser(idx)
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

func GetMyInfo(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	userid := sess.Get("userid")
	if userid == nil {
		return c.Status(http.StatusForbidden).Send([]byte("Unauthorized"))
	}

	user, err := crud.GetUserByNameMap(userid.(string))
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	delete(user.(map[string]interface{}), "password")

	return c.Status(http.StatusOK).JSON(user)
}

func UpdateMyInfo(c *fiber.Ctx) error {
	var err error

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	userid := sess.Get("userid")
	if userid == nil {
		return c.Status(http.StatusForbidden).Send([]byte("Unauthorized"))
	}

	userDataOldRaw, err := crud.GetUserByNameMap(userid.(string))
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}
	userDataOld := userDataOldRaw.(map[string]interface{})

	userDataNEW := make(map[string]interface{})
	err = json.Unmarshal(c.Body(), &userDataNEW)
	if err != nil {
		return c.Status(http.StatusBadRequest).Send([]byte(err.Error()))
	}

	userDataNEW["idx"] = fmt.Sprint(userDataOld["idx"].(int64))

	if _, ok := userDataNEW["password"]; ok {
		err = bcrypt.CompareHashAndPassword([]byte(userDataOld["password"].(string)), []byte(userDataNEW["old-password"].(string)))
		if err != nil {
			return c.Status(http.StatusBadRequest).Send([]byte("wrong password"))
		}

		password, err := bcrypt.GenerateFromPassword([]byte(userDataNEW["password"].(string)), consts.BcryptCost)
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
		}

		userDataNEW["password"] = string(password)
		delete(userDataNEW, "old-password")
	}

	err = crud.UpdateUserMap(userDataNEW)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	result := map[string]string{
		"result": "ok",
	}

	return c.Status(http.StatusOK).JSON(result)
}

func ResignUser(c *fiber.Ctx) error {
	var err error

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	userid := sess.Get("userid")
	if userid == nil {
		return c.Status(http.StatusForbidden).Send([]byte("Unauthorized"))
	}

	userDataRaw, err := crud.GetUserByNameMap(userid.(string))
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}
	userData := userDataRaw.(map[string]interface{})

	err = crud.ResignUser(userData["idx"].(int64))
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	result := map[string]string{
		"result": "ok",
	}

	return c.Status(http.StatusOK).JSON(result)
}

func ResetPasswordAPI(c *fiber.Ctx) error {
	var err error

	columnsCount, _ := crud.GetUserColumnsCount()

	userid := c.FormValue("userid")
	useremail := c.FormValue("email")

	if userid == "" {
		return c.Status(http.StatusBadRequest).Send([]byte("userid is empty"))
	}
	if useremail == "" {
		return c.Status(http.StatusBadRequest).Send([]byte("Email is empty"))
	}

	password := GetRandomString(16)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	switch columnsCount {
	// case model.UserDataFieldCount:
	// 	// Waive struct at this time
	// 	user, err := crud.GetUserByNameAndEmail(userid, useremail)
	// 	if err != nil {
	// 		return c.Status(http.StatusOK).Send([]byte(consts.MsgPasswordResetUserNotFound))
	// 	}

	// 	user.Password = null.StringFrom(string(passwordHash))
	// 	crud.UpdateUser(user)
	default:
		user, err := crud.GetUserByNameAndEmailMap(userid, useremail)
		if err != nil {
			return c.Status(http.StatusOK).Send([]byte(consts.MsgPasswordResetUserNotFound))
		}

		user.(map[string]interface{})["password"] = string(passwordHash)
		crud.UpdateUserMap(user.(map[string]interface{}))
	}

	// Send password reset email
	message := email.Message{
		Service:          email.Info.Service,
		AppendFromToName: false,
		From:             email.From{Email: email.Info.SenderInfo.Email, Name: email.Info.SenderInfo.Name},
		To:               email.To{Email: useremail, Name: userid},
		Subject:          "EnjoyTools - Password changed",
		Body: `
		The password for your account was changed on ` + time.Now().UTC().Format("2006-01-02 15:04:05 UTC") + `
		<br /><br />
		` + password,
		BodyType: email.HTML,
	}

	if email.Info.UseEmail {
		err = email.SendVerificationEmail(message)
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
		}
	}

	return c.Status(http.StatusOK).Send([]byte(consts.MsgPasswordResetEmail))
}

func Get2FaQR(c *fiber.Ctx) error {
	// randomSecret := gotp.RandomSecret(16)
	randomSecret := "ILOYEUDHGQJUSG7WP4RRP3RLT4"
	fmt.Println(randomSecret)

	totp := gotp.NewDefaultTOTP(randomSecret)
	fmt.Println("current one-time password is:", totp.Now())

	uri := totp.ProvisioningUri("user@email.com", consts.SiteName)
	fmt.Println(uri)

	// err := qrcode.WriteFile(uri, qrcode.Medium, 256, "qr.png")
	qrb64, err := qrcode.Encode(uri, qrcode.Medium, 256)
	if err != nil {
		return err
	}

	c.Set("Content-type", "image/png")
	return c.Status(http.StatusOK).Send(qrb64)
}

// func Verify2FA(randomSecret string) {
func Verify2FA(c *fiber.Ctx) error {
	randomSecret := "ILOYEUDHGQJUSG7WP4RRP3RLT4"

	totp := gotp.NewDefaultTOTP(randomSecret)
	otpValue := totp.Now()
	fmt.Println("current one-time password is:", otpValue)

	ok := totp.Verify(otpValue, time.Now().Unix())
	fmt.Println("verify OTP success:", ok)

	return c.Status(http.StatusOK).Send([]byte(otpValue))
}
