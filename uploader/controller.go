package uploader

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type (
	FileData struct {
		Name    null.String `json:"name"`
		TmpName null.String `json:"tmp_name,omitempty"`
		DbName  null.String `json:"dbname,omitempty"`
		Type    null.String `json:"type"`
		Content null.String `json:"content,omitempty"`
		Success null.String `json:"success,omitempty"`
	}
)

var (
	pathTMP  = "data_tmp"
	pathData = "data"
)

// UniqueID - Get Unique ID
func UniqueID() string {
	t := time.Unix(time.Now().Unix(), 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)

	return fmt.Sprint(ulid.MustNew(ulid.Timestamp(t), entropy))
}

// UploadTMP - Upload file as temporary
func UploadTMP(c echo.Context) (err error) {
	var file FileData

	// 에러 확인
	// return c.JSON(http.StatusBadRequest, map[string]string{"msg": "You cannot upload"})

	dataROOT, err := os.Getwd()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}
	dataPathTMP := dataROOT + "/../" + pathTMP + "/"

	if _, err := os.Stat(dataPathTMP); os.IsNotExist(err) {
		err := os.Mkdir(dataPathTMP, 0644)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
		}
	}

	if err = c.Bind(&file); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	fName := strings.TrimSuffix(file.Name.String, filepath.Ext(file.Name.String))
	fExt := filepath.Ext(file.Name.String)

	tmpName := fName + "_" + UniqueID() + fExt

	content, _ := b64.StdEncoding.DecodeString(file.Content.String)
	err = ioutil.WriteFile(dataPathTMP+tmpName, content, 0644)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	result := map[string]string{
		"result":   "done",
		"name":     file.Name.String,
		"type":     file.Type.String,
		"tmp_name": tmpName,
	}

	return c.JSON(http.StatusOK, result)
}

// UploadFINISH - Move file from temporary to data dir
func UploadFINISH(c echo.Context) (err error) {
	var files []FileData

	dataROOT, err := os.Getwd()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}
	dataPathTMP := dataROOT + "/../" + pathTMP + "/"
	finPath := dataROOT + "/../" + pathData + "/"

	if _, err := os.Stat(finPath); os.IsNotExist(err) {
		err := os.Mkdir(finPath, 0644)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
		}
	}

	if err = c.Bind(&files); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg": err.Error()})
	}

	for i, f := range files {
		// fName := strings.TrimSuffix(f.Name.String, filepath.Ext(f.Name.String))
		// fExt := filepath.Ext(f.Name.String)

		// dbName := fName + "_" + UniqueID() + fExt
		// err := os.Rename(dataPathTMP+f.TmpName.String, finPath+dbName)
		err := os.Rename(dataPathTMP+f.TmpName.String, finPath+f.TmpName.String)

		files[i].Success = null.NewString("success", true)
		files[i].DbName = null.NewString(f.TmpName.String, true)
		if err != nil {
			files[i].Success = null.NewString(err.Error(), true)
		}
	}

	// filesJSON, _ := json.Marshal(files)
	result := map[string]interface{}{
		"result": "done",
		// "files":  string(filesJSON),
	}

	return c.JSON(http.StatusOK, result)
}

// DeleteTMP - Delete temporary file
func DeleteTMP(c echo.Context) (err error) {
	var file FileData

	dataROOT, err := os.Getwd()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg1": err.Error()})
	}
	dataPathTMP := dataROOT + "/../" + pathTMP + "/"

	if err = c.Bind(&file); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg2": err.Error()})
	}

	tmpName := file.TmpName.String

	err = os.Remove(dataPathTMP + tmpName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"msg3": err.Error()})
	}

	result := map[string]string{
		"result":   "done",
		"name":     file.Name.String,
		"type":     file.Type.String,
		"tmp_name": tmpName,
	}

	return c.JSON(http.StatusOK, result)
}
