package board

import (
	"log"

	jsoniter "github.com/json-iterator/go"
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
		if result[i].Fields != nil && result[i].Fields.(string) != "" {
			err := json.Unmarshal([]byte(result[i].Fields.(string)), &rawJson)
			if err != nil {
				log.Println("GetData Unmarshal: ", err)
			}
			result[i].Fields = rawJson
		}
	}

	return result
}
