package handler

import (
	"9minutes/consts"
	"9minutes/internal/crud"
	"9minutes/model"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/guregu/null.v4"
)

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

func LoadBoardListData() map[string]model.Board {
	BoardListData = map[string]model.Board{}

	listingOptions := model.BoardListingOptions{}
	listingOptions.Page = null.IntFrom(0)
	listingOptions.ListCount = null.IntFrom(99999)
	boardList, _ := crud.GetBoards(listingOptions)

	for _, b := range boardList.BoardList {
		BoardListData[b.BoardCode.String] = b
	}

	return BoardListData
}

func LoadUserColumnDatas() {
	UserColumnsData, _ = crud.GetUserColumnsList()
}

func checkBoardActionExist(action string) bool {
	result := false
	for _, a := range boardActions {
		if a == action {
			result = true
			break
		}
	}

	return result
}

func checkBoardAccessible(boardGradeKey string, userGradeKey string) bool {
	boardRank := consts.BoardGrades[boardGradeKey].Rank
	userRank := consts.UserGrades[userGradeKey].Rank

	result := false
	if userRank <= boardRank {
		result = true
	}

	return result
}
