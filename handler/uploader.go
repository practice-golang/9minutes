package handler

import (
	"9m/router"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
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

	tempFile, err := ioutil.TempFile(router.UploadPath, "upload-*-"+handler.Filename)
	if err != nil {
		log.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	tempFile.Write(fileBytes)

	resultMAP := map[string]string{
		"message":   "success",
		"filename":  handler.Filename,
		"storename": filepath.Base(tempFile.Name()),
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

	tempFile, err := ioutil.TempFile(router.UploadPath, "upload-*-"+handler.Filename)
	if err != nil {
		log.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	tempFile.Write(fileBytes)

	resultMAP := map[string]string{
		"message":   "success",
		"filename":  handler.Filename,
		"storename": filepath.Base(tempFile.Name()),
	}

	c.Json(http.StatusOK, resultMAP)
}
