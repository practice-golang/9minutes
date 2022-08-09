package handler

import (
	"9minutes/auth"
	"9minutes/config"
	"9minutes/crud"
	"9minutes/fd"
	"9minutes/model"
	"9minutes/router"
	"bytes"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func LoadFile(c *router.Context) ([]byte, error) {
	var h []byte
	var err error

	// If the file exists in the real storage, read real instead of embed.
	storePath := StoreRoot + c.URL.Path // Real storage
	embedPath := EmbedRoot + c.URL.Path // Embed storage
	switch true {
	case fd.CheckFileExists(storePath, false):
		h, err = os.ReadFile(storePath)
	case fd.CheckFileExists(embedPath, true):
		h, err = router.Content.ReadFile(embedPath)
	default:
		return nil, errors.New("file not found")
	}

	if err != nil {
		return nil, err
	}

	return h, err
}

func LoadHTML(c *router.Context) ([]byte, error) {
	h, err := LoadFile(c)
	if err != nil {
		return nil, err
	}

	if filepath.Base(c.URL.Path) == "login.html" {
		return h, err
	}

	isLoggedIn := ""
	name := ""
	var userInfo model.UserData
	userGrade := 999

	claim, _ := auth.GetClaim(*c.Request, "cookie")
	if claim.Name.String != "" {
		name = claim.Name.String
		isLoggedIn = "true"
	}

	if auth.SessionManager != nil && auth.SessionManager.Exists(c.Context(), "userid") {
		ses, err := auth.GetCookieSession(c)
		if err != nil {
			auth.ExpireCookie(c.ResponseWriter)
		}
		c.AuthInfo = ses
		if c.AuthInfo != nil {
			authinfo := c.AuthInfo.(model.AuthInfo)
			if authinfo.Name.Valid {
				name = authinfo.Name.String
			}
		}
	}

	if name == "" {
		name = "Guest"
		isLoggedIn = ""

		userGrade = config.UserGrades.IndexOf("guest")
	} else {
		userInfo, err = crud.GetUserByName(name)
		if err != nil {
			auth.ExpireCookie(c.ResponseWriter)
			return nil, err
		}

		userGrade = config.UserGrades.IndexOf(userInfo.Grade.String)
	}

	m := reIncludes.FindAllSubmatch(h, -1)
	for _, v := range m {
		includeFileName := string(v[1])
		includeDirective := bytes.TrimSpace(v[0])
		includeStoreFilePath := strings.TrimSpace(StoreRoot + "/" + includeFileName)
		includeEmbedFilePath := strings.TrimSpace(EmbedRoot + "/" + includeFileName)

		include := []byte{}
		switch true {
		case fd.CheckFileExists(includeStoreFilePath, false):
			include, _ = os.ReadFile(includeStoreFilePath)
		case fd.CheckFileExists(includeEmbedFilePath, true):
			include, _ = router.Content.ReadFile(includeEmbedFilePath)
		}

		h = bytes.ReplaceAll(h, includeDirective, include)
	}

	if name != "Guest" {
		h = reLogin.ReplaceAll(h, []byte(""))

		if userGrade > config.UserGrades.IndexOf("manager") {
			h = reAdmin.ReplaceAll(h, []byte(""))
		}

		if userInfo.Grade.String != "pending_user" {
			h = reYouArePending.ReplaceAll(h, []byte(""))
		}
	} else {
		h = reLogout.ReplaceAll(h, []byte(""))
		h = reAdmin.ReplaceAll(h, []byte(""))
		h = reMyPage.ReplaceAll(h, []byte(""))
		h = reYouArePending.ReplaceAll(h, []byte(""))
	}

	h = bytes.ReplaceAll(h, []byte("$USERNAME$"), []byte(name))
	h = bytes.ReplaceAll(h, []byte("$LinkLogin$"), []byte(""))
	h = bytes.ReplaceAll(h, []byte("$LinkLogout$"), []byte(""))
	h = bytes.ReplaceAll(h, []byte("$LinkAdmin$"), []byte(""))
	h = bytes.ReplaceAll(h, []byte("$LinkMyPage$"), []byte(""))
	h = bytes.ReplaceAll(h, []byte("$YouArePending$"), []byte(""))

	h = bytes.ReplaceAll(h, []byte("$ISLOGGEDIN$"), []byte(isLoggedIn))

	return h, err
}

// GetRandomString - Generate random string
func GetRandomString(length int) string {
	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	rand.Seed(time.Now().UnixNano())

	randomBytes := make([]byte, length)

	for i := 0; i < length; i++ {
		randomBytes[i] = charset[rand.Intn(len(charset))]
	}

	return string(randomBytes)
}

func DeleteUploadFile(filepath string) {
	if fd.CheckFileExists(filepath, false) {
		os.Remove(filepath)
	}
}
