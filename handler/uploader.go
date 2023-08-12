package handler

import (
	"9minutes/config"
	"9minutes/internal/crud"
	"crypto/sha256"
	"encoding/base64"
	"log"
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

	for fiels, headers := range form.File {
		log.Println("uploading file", fiels, headers)
	}

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

func FilesInfo(c *fiber.Ctx) (err error) {
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

	var indices []int
	for _, u := range uploadIndices {
		indices = append(indices, int(u.Idx.Int64))
	}

	var result = []fiber.Map{}
	fdatas, err := crud.GetUploadedFiles(indices)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Send([]byte(err.Error()))
	}

	for _, fdata := range fdatas {
		result = append(result, fiber.Map{
			"idx":         strconv.FormatInt(fdata.Idx.Int64, 10),
			"filename":    fdata.FileName.String,
			"storagename": fdata.StorageName.String,
		})
	}

	return c.Status(http.StatusOK).JSON(result)
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
