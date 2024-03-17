package handler

import (
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

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

func CheckFileExtensionIsImage(fname string) bool {
	result := false
	imageFileExtensions := []string{"jpg", "jpeg", "png", "gif", "bmp"}

	fnames := strings.Split(fname, ".")
	fext := fnames[len(fnames)-1]

	for _, imExt := range imageFileExtensions {
		if fext == imExt {
			result = true
			break
		}
	}

	return result
}

func CopyResizeImagePNG(filename string, resizeName string, width int, height int) error {
	var err error

	img, err := imgio.Open(filename)
	if err != nil {
		return err
	}

	wx := img.Bounds().Size().X
	wy := img.Bounds().Size().Y
	if wx > 800 || wy > 800 {
		resized := transform.Resize(img, width, height, transform.Linear)
		err = imgio.Save(resizeName, resized, imgio.PNGEncoder())
		if err != nil {
			return err
		}
	} else {
		// If image size is under 800x800, copy without resize
		imgio.Save(resizeName, img, imgio.PNGEncoder())
	}

	return nil
}
