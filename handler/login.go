package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"9minutes/consts"
	"9minutes/crud"
	"9minutes/email"
	"9minutes/model"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

// LoginAPI - Login
func LoginAPI(c *fiber.Ctx) error {
	var err error

	signin := model.SignIn{}
	err = json.Unmarshal(c.Request().Body(), &signin)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	sess.Set("name", signin.Name.String)
	sess.Set("ip", c.IP())
	sess.Set("user-agent", c.Get("User-Agent"))
	sess.Set("duration", 60*60*24*7)

	err = sess.Save()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	result := map[string]string{
		"msg": "Signin success",
	}

	return c.Status(http.StatusOK).JSON(result)
}

// Signup - Create new user
func SignupAPI(c *fiber.Ctx) error {
	var err error

	now := time.Now().Format("20060102150405")
	columnsCount, _ := crud.GetUserColumnsCount()

	userIDX := ""
	username := ""
	useremail := ""

	rbody := c.Request().Body()

	switch columnsCount {
	case model.UserDataFieldCount:
		var userData model.UserData

		// err = json.NewDecoder(c.Body).Decode(&userData)
		err = json.Unmarshal(rbody, &userData)
		if err != nil {
			return c.Status(http.StatusBadRequest).Send([]byte(err.Error()))
		}

		password, err := bcrypt.GenerateFromPassword([]byte(userData.Password.String), consts.BcryptCost)
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
		}
		userData.Password = null.StringFrom(string(password))
		userData.RegDTTM = null.StringFrom(now)
		userData.Grade = null.StringFrom("pending_user")
		userData.Approval = null.StringFrom("N")

		switch true {
		case userData.UserName.String == "":
			return c.Status(http.StatusBadRequest).Send([]byte("Username is empty"))
		case userData.Email.String == "":
			return c.Status(http.StatusBadRequest).Send([]byte("Email is empty"))
		case userData.Password.String == "":
			return c.Status(http.StatusBadRequest).Send([]byte("Password is empty"))
		}

		err = crud.AddUser(userData)
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
		}

		username = userData.UserName.String
		useremail = userData.Email.String

		userInsertResult, err := crud.GetUserByNameAndEmail(username, useremail)
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
		}

		userIDX = fmt.Sprint(userInsertResult.Idx.Int64)

	default:
		userData := make(map[string]interface{})

		// err = json.NewDecoder(c.Body).Decode(&userData)
		err = json.Unmarshal(rbody, &userData)
		if err != nil {
			return c.Status(http.StatusBadRequest).Send([]byte(err.Error()))
		}

		password, err := bcrypt.GenerateFromPassword([]byte(userData["password"].(string)), consts.BcryptCost)
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
		}

		userData["password"] = string(password)
		userData["reg-dttm"] = now
		userData["grade"] = "pending_user"
		userData["approval"] = "N"

		switch true {
		case userData["username"].(string) == "":
			return c.Status(http.StatusBadRequest).Send([]byte("Username is empty"))
		case userData["email"].(string) == "":
			return c.Status(http.StatusBadRequest).Send([]byte("Email is empty"))
		case userData["password"].(string) == "":
			return c.Status(http.StatusBadRequest).Send([]byte("Password is empty"))
		}

		err = crud.AddUserMap(userData)
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
		}

		username = userData["username"].(string)
		useremail = userData["email"].(string)

		userInsertResult, err := crud.GetUserByNameAndEmailMap(username, useremail)
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
		}

		userIDX = userInsertResult.(map[string]interface{})["IDX"].(string)
	}

	verificationKEY := GetRandomString(32)
	verificationData := map[string]string{
		"USER_IDX": userIDX,
		"TOKEN":    verificationKEY,
	}

	err = crud.AddUserVerification(verificationData)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	// Send verification email
	domain := email.Info.Domain
	message := email.Message{
		Service:          email.Info.Service,
		AppendFromToName: false,
		From:             email.From{Email: email.Info.SenderInfo.Email, Name: email.Info.SenderInfo.Name},
		To:               email.To{Email: useremail, Name: username},
		Subject:          "EnjoyTools - Email Verification",
		Body: `
		Please click the link below to verify your email address.
		<br />
		<a href='` + domain + `/verify?username=` + username + `&email=` + useremail + `&token=` + verificationKEY + `'>Click here</a>`,
		BodyType: email.HTML,
	}

	if email.Info.UseEmail {
		err = email.SendVerificationEmail(message)
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
		}
	}

	result := map[string]string{
		"result": "ok",
	}

	return c.Status(http.StatusOK).JSON(result)
}
