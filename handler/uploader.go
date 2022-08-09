package handler

import (
	"9minutes/crud"
	"9minutes/router"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// https://tutorialedge.net/golang/go-file-upload-tutorial
func UploadFile(c *router.Context) {
	// w := c.ResponseWriter
	r := c.Request

	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	// w.Header().Set("Access-Control-Allow-Credentials", "true")

	r.ParseMultipartForm(1 << 20) // 10 << 20 specifies a maximum upload of 10 MB files

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}
	defer file.Close()

	sha := sha256.New()
	sha.Write([]byte(filepath.Base(handler.Filename) + time.Now().String()))
	sha.Write([]byte(filepath.Ext(handler.Filename) + time.Now().String()))
	storageName := base64.StdEncoding.EncodeToString(sha.Sum(nil))
	storageName = strings.NewReplacer("=", "", "+", "", "/", "").Replace(storageName)
	storageName = storageName + "_" + time.Now().Format("20060102150405") + filepath.Ext(handler.Filename)
	// storageName := GetRandomString(128) + time.Now().Format("20060102150405") + "." + filepath.Ext(handler.Filename)

	// // tempFile, err := ioutil.TempFile(router.UploadPath, "upload-*-"+handler.Filename)
	// tempFile, err := os.CreateTemp(router.UploadPath, "upload-*-"+handler.Filename)
	tempFile, err := os.CreateTemp(router.UploadPath, "*"+storageName)
	if err != nil {
		log.Println(err)
		return
	}
	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}
	tempFile.Write(fileBytes)
	storageName = filepath.Base(tempFile.Name())

	err = crud.AddUploadedFile(handler.Filename, storageName)
	if err != nil {
		log.Println(err)
		return
	}

	resultMAP := map[string]string{
		"message":   "success",
		"filename":  handler.Filename,
		"storename": storageName,
	}

	c.Json(http.StatusOK, resultMAP)
}

func UploadImage(c *router.Context) {
	// w := c.ResponseWriter
	r := c.Request

	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	// w.Header().Set("Access-Control-Allow-Credentials", "true")

	r.ParseMultipartForm(10 << 20) // 10 << 20 specifies a maximum upload of 10 MB files

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}
	defer file.Close()

	sha := sha256.New()
	sha.Write([]byte(filepath.Base(handler.Filename) + time.Now().String()))
	sha.Write([]byte(filepath.Ext(handler.Filename) + time.Now().String()))
	storageName := base64.StdEncoding.EncodeToString(sha.Sum(nil))
	storageName = strings.NewReplacer("=", "", "+", "", "/", "").Replace(storageName)
	storageName = storageName + "_" + time.Now().Format("20060102150405") + filepath.Ext(handler.Filename)
	// storageName := GetRandomString(128) + time.Now().Format("20060102150405") + "." + filepath.Ext(handler.Filename)

	// // tempFile, err := ioutil.TempFile(router.UploadPath, "upload-*-"+handler.Filename)
	// tempFile, err := os.CreateTemp(router.UploadPath, "upload-*-"+handler.Filename)
	tempFile, err := os.CreateTemp(router.UploadPath, "*"+storageName)
	if err != nil {
		log.Println(err)
		return
	}
	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}
	tempFile.Write(fileBytes)
	storageName = filepath.Base(tempFile.Name())

	err = crud.AddUploadedFile(handler.Filename, storageName)
	if err != nil {
		log.Println(err)
		return
	}

	resultMAP := map[string]string{
		"message":   "success",
		"filename":  handler.Filename,
		"storename": storageName,
	}

	c.Json(http.StatusOK, resultMAP)
}
