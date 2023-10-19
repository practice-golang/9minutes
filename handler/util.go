package handler

import (
	"9minutes/internal/crud"
	"9minutes/model"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/guregu/null.v4"
)

// import (
// 	"9minutes/auth"
// 	"9minutes/config"
// 	"9minutes/internal/crud"
// 	"9minutes/internal/fd"
// 	"9minutes/model"
// 	"9minutes/router"
// 	"bytes"
// 	"errors"
// 	"math/rand"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"time"
// )

func LoadHTML(c *fiber.Ctx) ([]byte, error) {
	return nil, nil
}

// GetRandomString - Generate random string
func GetRandomString(length int) string {
	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	randomBytes := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		randomBytes[i] = charset[r.Intn(len(charset))]
	}

	return string(randomBytes)
}

func DeleteUploadFile(filepath string) {
	if _, err := os.Stat(filepath); err == nil {
		os.Remove(filepath)
	}
}

func SetBoardDataALL() map[string]model.Board {
	BoardListALL = map[string]model.Board{}

	boardListingOptions := model.BoardListingOptions{}
	boardListingOptions.Page = null.IntFrom(0)
	boardListingOptions.ListCount = null.IntFrom(99999)
	boardList, _ := crud.GetBoards(boardListingOptions)

	for _, b := range boardList.BoardList {
		BoardListALL[b.BoardCode.String] = b
	}

	return BoardListALL
}
