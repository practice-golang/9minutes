package board

import (
	"fmt"
	"log"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/practice-golang/9minutes/models"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// var json = jsoniter.ConfigFastest

func prepareInsertData(data []byte) ([]models.Board, [][]models.Field) {
	var result []models.Board
	var resultFields [][]models.Field

	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Println("AddBoards Unmarshal: ", err)
	}

	for i := range result {
		fields := []models.Field{}

		// log.Println("prepareInsertData: ", i, result[i].Table.Valid)

		if result[i].Fields != nil {
			rawJson, err := json.MarshalToString(result[i].Fields)
			if err != nil {
				log.Println("GetData Unmarshal: ", err)
			}
			result[i].Fields = rawJson

			jsoniter.UnmarshalFromString(rawJson, &fields)
		}

		resultFields = append(resultFields, fields)
	}

	return result, resultFields
}

func prepareUpdateData(data []byte) models.Board {
	var result models.Board

	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Println("AddBoards Unmarshal: ", err)
	}

	if result.Fields != nil {
		rawJson, err := json.MarshalToString(result.Fields)
		if err != nil {
			log.Println("GetData Unmarshal: ", err)
		}
		result.Fields = rawJson
	}

	return result
}

func prepareSelectData(data interface{}) []models.Board {
	var result []models.Board

	if data != nil {
		result = data.([]models.Board)
	} else {
		log.Println("Null data")
	}

	for i := range result {
		var rawJson interface{}
		if result[i].Fields != nil {
			switch fields := result[i].Fields.(type) {
			case string:
				if fields != "" {
					err := json.Unmarshal([]byte(fields), &rawJson)
					if err != nil {
						log.Println("prepareSelectData Unmarshal: ", err)
					}
					result[i].Fields = rawJson
				}
			case []byte:
				if string(fields) != "" {
					err := json.Unmarshal(fields, &rawJson)
					if err != nil {
						log.Println("prepareSelectData Unmarshal: ", err)
					}
					result[i].Fields = rawJson
				}
			default:
				t := fmt.Sprintf("Unknown type %T", result[i].Fields)
				log.Println("prepareSelectData: Unknown type " + t)
			}
		}
	}

	return result
}

// CheckPermission - Check permission
func CheckUpload(c echo.Context) bool {
	var isFileUpload bool

	code := c.QueryParam("code")
	boardInfos := GetBoardByCode(code)

	if len(boardInfos) == 0 {
		return false
	}

	if boardInfos[0].Type.String == "custom-tablelist" {
		return false
	}

	fileUploadSET := boardInfos[0].FileUpload.String

	isFileUpload = false
	if fileUploadSET == "Y" {
		isFileUpload = true
	}

	return isFileUpload
}
