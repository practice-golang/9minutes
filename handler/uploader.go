package handler

import (
	"9minutes/config"
	"9minutes/internal/crud"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/guregu/null.v4"
)

// https://tutorialedge.net/golang/go-file-upload-tutorial
func UploadFile(c *fiber.Ctx) (err error) {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	resultMAP := map[string]interface{}{
		"message": "success",
		"files":   []map[string]string{},
	}

	fdatas := form.File["upload-files"]
	for _, fdata := range fdatas {
		fname := fdata.Filename

		sha := sha256.New()
		sha.Write([]byte(filepath.Base(fname) + time.Now().String()))
		sha.Write([]byte(filepath.Ext(fname) + time.Now().String()))
		storageName := base64.StdEncoding.EncodeToString(sha.Sum(nil))
		storageName = strings.NewReplacer("=", "", "+", "", "/", "").Replace(storageName)
		storageName = storageName + GetRandomString(16) + "_" + time.Now().Format("20060102150405") + filepath.Ext(fname)

		err := c.SaveFile(fdata, config.UploadPath+"/"+storageName)
		if err != nil {
			return err
		}

		r, err := crud.AddUploadedFile(fname, storageName)
		if err != nil {
			return err
		}

		fidx, err := r.LastInsertId()
		if err != nil {
			return err
		}

		files := map[string]string{
			"idx":         strconv.FormatInt(fidx, 10),
			"filename":    fname,
			"storagename": storageName,
		}

		resultMAP["files"] = append(resultMAP["files"].([]map[string]string), files)
	}

	return c.Status(http.StatusOK).JSON(resultMAP)
}

func DeleteFiles(c *fiber.Ctx) (err error) {
	type uploadIdx struct {
		Idx null.Int `json:"idx" db:"IDX"`
	}
	var uploadIndices []uploadIdx

	err = c.BodyParser(&uploadIndices)
	if err != nil {
		return c.Status(http.StatusBadRequest).Send([]byte(err.Error()))
	}

	if len(uploadIndices) == 0 {
		return c.Status(http.StatusBadRequest).Send([]byte("no files to delete"))
	}

	for _, f := range uploadIndices {

		fdata, err := crud.GetUploadedFile(int(f.Idx.Int64))
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
		}

		err = crud.DeleteUploadedFile(f.Idx.Int64)
		if err != nil {
			return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
		}
		filepath := config.UploadPath + "/" + fdata.StorageName.String
		DeleteUploadFile(filepath)
	}

	return c.Status(http.StatusOK).Send([]byte("success"))
}

// func UploadImage(c *router.Context) {
// 	// w := c.ResponseWriter
// 	r := c.Request

// 	// w.Header().Set("Access-Control-Allow-Origin", "*")
// 	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 	// w.Header().Set("Access-Control-Allow-Credentials", "true")

// 	r.ParseMultipartForm(10 << 20) // 10 << 20 specifies a maximum upload of 10 MB files

// 	file, handler, err := r.FormFile("file")
// 	if err != nil {
// 		log.Println("Error Retrieving the File")
// 		log.Println(err)
// 		return
// 	}
// 	defer file.Close()

// 	sha := sha256.New()
// 	sha.Write([]byte(filepath.Base(handler.Filename) + time.Now().String()))
// 	sha.Write([]byte(filepath.Ext(handler.Filename) + time.Now().String()))
// 	storageName := base64.StdEncoding.EncodeToString(sha.Sum(nil))
// 	storageName = strings.NewReplacer("=", "", "+", "", "/", "").Replace(storageName)
// 	storageName = storageName + "_" + time.Now().Format("20060102150405") + filepath.Ext(handler.Filename)
// 	// storageName := GetRandomString(128) + time.Now().Format("20060102150405") + "." + filepath.Ext(handler.Filename)

// 	// // tempFile, err := ioutil.TempFile(router.UploadPath, "upload-*-"+handler.Filename)
// 	// tempFile, err := os.CreateTemp(router.UploadPath, "upload-*-"+handler.Filename)
// 	tempFile, err := os.CreateTemp(router.UploadPath, "*"+storageName)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer tempFile.Close()

// 	fileBytes, err := io.ReadAll(file)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	tempFile.Write(fileBytes)
// 	storageName = filepath.Base(tempFile.Name())

// 	err = crud.AddUploadedFile(handler.Filename, storageName)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	resultMAP := map[string]string{
// 		"message":   "success",
// 		"filename":  handler.Filename,
// 		"storename": storageName,
// 	}

// 	c.Json(http.StatusOK, resultMAP)
// }

// func DeleteFiles(c *router.Context) {
// 	type deleteFiles struct {
// 		BoardIdx null.Int     `json:"board-idx"`
// 		PostIdx  null.Int     `json:"post-idx"`
// 		Files    []model.File `json:"delete-files"`
// 	}
// 	var requestDelete deleteFiles

// 	err := json.NewDecoder(c.Body).Decode(&requestDelete)
// 	if err != nil {
// 		log.Println("Cancel process:", err)
// 		return
// 	}

// 	if len(requestDelete.Files) == 0 {
// 		return
// 	}

// 	for _, f := range requestDelete.Files {
// 		if f.FileName.Valid && f.StoreName.Valid {
// 			err = crud.DeleteUploadedFile(requestDelete.BoardIdx.Int64, requestDelete.PostIdx.Int64, f.FileName.String, f.StoreName.String)
// 			if err != nil {
// 				log.Println(err)
// 				return
// 			}

// 			filepath := router.UploadPath + "/" + f.StoreName.String
// 			DeleteUploadFile(filepath)
// 		}
// 	}

// 	// Because of browser have already gone so, response nothing
// 	// c.Json(http.StatusOK, "")
// }
