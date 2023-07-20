package main // import "fakegen-boards"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/brianvoe/gofakeit/v6"
)

type LoginData struct {
	UserID   string `json:"userid,omitempty"`
	Password string `json:"password,omitempty"`
}

type boardData struct {
	BoardName    string `json:"board-name,omitempty"`
	BoardCode    string `json:"board-code,omitempty"`
	BoardType    string `json:"board-type,omitempty"`
	BoardTable   string `json:"board-table,omitempty"`
	CommentTable string `json:"comment-table,omitempty"`
	GrantRead    string `json:"grant-read,omitempty"`
	GrantWrite   string `json:"grant-write,omitempty"`
	GrantComment string `json:"grant-comment,omitempty"`
	GrantUpload  string `json:"grant-upload,omitempty"`
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

func getBoardList(uri, sess string) {
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

func prepareBoardData(count int) (boards []boardData) {
	boardcodes := []string{}

	for i := 0; i < count; i++ {
		sameExists := false

		randomBoard := gofakeit.LetterN(10)
		boardName := strings.ToLower(randomBoard)
		boardNameRune := []rune(boardName)
		boardNameRune[0] = unicode.ToUpper(boardNameRune[0])
		boardName = string(boardNameRune)

		board := boardData{
			BoardName:    boardName,
			BoardCode:    strings.ToLower(randomBoard),
			BoardType:    "board",
			BoardTable:   strings.ToUpper("BOARD_" + randomBoard),
			CommentTable: strings.ToUpper("COMMENT_" + randomBoard),
			GrantRead:    "guest",
			GrantWrite:   "regular_user",
			GrantComment: "regular_user",
			GrantUpload:  "regular_user",
		}

		for i := 0; i < len(boardcodes); i++ {
			if boardcodes[i] == board.BoardCode {
				i--
				sameExists = true
				break
			}
		}

		if sameExists {
			continue
		}

		boardcodes = append(boardcodes, board.BoardCode)
		boards = append(boards, board)
	}

	return
}

func addBoardData(uri, sess string, board boardData) {
	boardJSON, _ := json.Marshal(board)
	buf := bytes.NewBuffer(boardJSON)

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

	uriBoardList := "http://localhost:5525/api/admin/board"
	getBoardList(uriBoardList, sess)

	count := 30
	boards := prepareBoardData(count)

	uriAddBoard := "http://localhost:5525/api/admin/board"

	for _, b := range boards {
		addBoardData(uriAddBoard, sess, b)
	}
}
