package main // import "fakegen-users"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
)

type LoginData struct {
	UserID   string `json:"userid,omitempty"`
	Password string `json:"password,omitempty"`
}

type userData struct {
	UserID   string `json:"userid,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Grade    string `json:"grade,omitempty"`
	Approval string `json:"approval,omitempty"`
}

func getSession(uri string) (cookieSession string) {
	admin := LoginData{UserID: "admin", Password: "admin"}
	// adminJSON, _ := json.MarshalIndent(admin, "", " ")
	adminJSON, _ := json.Marshal(admin)
	buf := bytes.NewBuffer(adminJSON)
	response, err := http.Post(uri, "application/json", buf)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	cookieSession = response.Cookies()[0].Name + "=" + response.Cookies()[0].Value

	body, _ := io.ReadAll(response.Body)
	fmt.Println("getSession:", string(body))

	return
}

func getUserList(uri, sess string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Set("Cookie", sess)
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	fmt.Println(string(body))
}

func prepareUserData(count int) (users []userData) {
	userids := []string{}
	emails := []string{}

	for i := 0; i < count; i++ {
		sameExists := false

		user := userData{
			UserID:   strings.ToLower(gofakeit.Username()),
			Password: "1234",
			Email:    gofakeit.Email(),
			Grade:    "user_hold",
			Approval: "N",
		}

		for i := 0; i < len(userids); i++ {
			if userids[i] == user.UserID {
				i--
				sameExists = true
				break
			}
			if emails[i] == user.Email {
				i--
				sameExists = true
				break
			}
		}

		if sameExists {
			continue
		}

		userids = append(userids, user.UserID)
		emails = append(emails, user.Email)

		users = append(users, user)
	}

	return
}

func addUserData(uri, sess string, user userData) {
	userJSON, _ := json.Marshal(user)
	buf := bytes.NewBuffer(userJSON)

	req, _ := http.NewRequest("POST", uri, buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Cookie", sess)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Print("body buf:")
		panic(err)
	}
	defer response.Body.Close()

	responseBody, _ := io.ReadAll(response.Body)
	fmt.Println(string(responseBody))
}

func main() {
	uriLogin := "http://localhost:5525/api/login"
	sess := getSession(uriLogin)

	uriUserList := "http://localhost:5525/api/admin/user"
	getUserList(uriUserList, sess)

	personnel := 30
	users := prepareUserData(personnel)

	uriAddUser := "http://localhost:5525/api/admin/user"

	for _, u := range users {
		addUserData(uriAddUser, sess, u)
	}
}
