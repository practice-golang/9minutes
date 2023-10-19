package main // import "fake-posting"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/brianvoe/gofakeit/v6"
)

type LoginData struct {
	UserID   string `json:"userid,omitempty"`
	Password string `json:"password,omitempty"`
}

type PostingData struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	AuthorIdx  int    `json:"author-idx"`
	AuthorName string `json:"author-name"`
	Files      string `json:"files"`
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

func preparePostings(count int) (contents []PostingData) {
	for i := 0; i < count; i++ {
		content := PostingData{
			Title:      gofakeit.LetterN(10),
			Content:    gofakeit.LetterN(40),
			AuthorIdx:  1,
			AuthorName: "admin",
			Files:      "",
		}

		contents = append(contents, content)
	}

	return contents
}

func writePostings(uri, sess string, content PostingData) {
	contentsJSON, _ := json.Marshal(content)
	buf := bytes.NewBuffer(contentsJSON)

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

	count := 30
	contents := preparePostings(count)

	uriWriteContent := "http://localhost:5525/api/board/misc/posting"

	for _, c := range contents {
		writePostings(uriWriteContent, sess, c)
	}
}
