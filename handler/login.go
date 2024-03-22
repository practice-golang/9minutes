package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"9minutes/consts"
	"9minutes/internal/crud"
	"9minutes/internal/email"
	"9minutes/model"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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

	user, err := crud.GetUserByUserIdAndPassword(signin.UserID.String, signin.Password.String)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	sess.Set("idx", user.(map[string]interface{})["idx"])
	sess.Set("userid", signin.UserID.String)
	sess.Set("grade", user.(map[string]interface{})["grade"])
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

func LogoutAPI(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	userid := sess.Get("userid")
	if userid == nil {
		return c.Status(http.StatusForbidden).Send([]byte("Unauthorized"))
	}

	sess.Destroy()

	result := map[string]string{
		"msg": "Signout success",
	}

	return c.Status(http.StatusOK).JSON(result)
}

// Signup - Create new user
func SignupAPI(c *fiber.Ctx) error {
	var err error

	now := time.Now().Format("20060102150405")

	userIDX := ""
	userid := ""
	useremail := ""

	rbody := c.Request().Body()

	data := make(map[string]interface{})

	err = json.Unmarshal(rbody, &data)
	if err != nil {
		return c.Status(http.StatusBadRequest).Send([]byte("Unmarshal:" + err.Error()))
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data["password"].(string)), consts.BcryptCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte("GenerateFromPassword:" + err.Error()))
	}

	data["password"] = string(password)
	data["regdate"] = now
	data["grade"] = "user_hold"
	data["approval"] = "N"

	switch true {
	case data["userid"].(string) == "":
		return c.Status(http.StatusBadRequest).Send([]byte("User id is empty"))
	case data["email"].(string) == "":
		return c.Status(http.StatusBadRequest).Send([]byte("Email is empty"))
	case data["password"].(string) == "":
		return c.Status(http.StatusBadRequest).Send([]byte("Password is empty"))
	}

	err = crud.AddUserMap(data)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte("AddUserMap:" + err.Error()))
	}

	userid = data["userid"].(string)
	useremail = data["email"].(string)

	userInsertResult, err := crud.GetUserByNameAndEmailMap(userid, useremail)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte("GetUserByNameAndEmailMap:" + err.Error()))
	}

	userIDX = strconv.Itoa(int(userInsertResult.(map[string]interface{})["idx"].(int64)))

	verificationKEY := GetRandomString(32)
	verificationData := map[string]string{
		"USER_IDX": userIDX,
		"TOKEN":    verificationKEY,
	}

	err = crud.AddUserVerification(verificationData)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte("AddUserVerification:" + err.Error()))
	}

	// Send verification email
	domain := email.Info.Domain
	message := email.Message{
		Service:          email.Info.Service,
		AppendFromToName: false,
		From:             email.From{Email: email.Info.SenderInfo.Email, Name: email.Info.SenderInfo.Name},
		To:               email.To{Email: useremail, Name: userid},
		Subject:          "Sign up Verification",
		Body: `
		Please click the link below to verify your email address.
		<br />
		<a href='` + domain + `/verify?userid=` + userid + `&email=` + useremail + `&token=` + verificationKEY + `'>Click here</a>`,
		BodyType: email.HTML,
	}

	if email.Info.UseEmail {
		err = email.SendVerificationEmail(message)
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte("SendVerificationEmail:" + err.Error()))
		}
	}

	result := map[string]string{
		"result": "ok",
	}

	return c.Status(http.StatusOK).JSON(result)
}
