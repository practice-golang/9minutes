package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/practice-golang/9minutes/board"
)

var (
	JwtKey []byte
)

// CustomClaims - jwt custom claim
type CustomClaims struct {
	Idx          string `json:"idx"`
	UserName     string `json:"username"`
	Email        string `json:"email"`
	Admin        string `json:"admin"`
	RefreshUntil int64  `json:"refresh-until"`
	jwt.StandardClaims
}

// PrepareToken - Create and return token
func PrepareToken(data interface{}) (string, error) {
	d := data.(map[string]interface{})

	claims := &CustomClaims{
		fmt.Sprint(d["IDX"]),
		fmt.Sprint(d["USERNAME"]),
		fmt.Sprint(d["EMAIL"]),
		fmt.Sprint(d["ADMIN"]),
		time.Now().Add(time.Hour * 30 * 24).Unix(),
		jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 1 * 24).Unix()},
		// jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Second * 1 * 60).Unix()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	result, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}

	return result, nil
}

// CheckAuth - Check role is valid with board setting
func CheckAuth(c echo.Context) (isValid bool) {
	code := c.QueryParam("code")
	mode := c.QueryParam("mode") // read, write

	boardInfos := board.GetBoardByCode(code)

	if len(boardInfos) == 0 {
		isValid = false
		return
	}

	grantRead := boardInfos[0].GrantRead.String
	grantWrite := boardInfos[0].GrantWrite.String

	user := c.Get("user")

	if user == nil {
		switch true {
		case ((mode == "write" || mode == "edit" || mode == "delete") && grantWrite == "all") ||
			((mode != "write" && mode != "edit" && mode != "delete") && grantRead == "all"):
			isValid = true
		default:
			isValid = false
		}
	} else {
		claims := user.(*jwt.Token).Claims.(*CustomClaims)
		// log.Println("CheckAuth: ", claims)
		// log.Println("CheckAuth: ", claims.Admin, code, mode, boardInfos[0].GrantWrite.String, boardInfos[0].GrantRead.String)

		switch true {
		case ((mode == "write" || mode == "edit" || mode == "delete") && (grantWrite == "admin" && claims.Admin == "Y")) ||
			((mode != "write" && mode != "edit" && mode != "delete") && (grantRead == "admin" && claims.Admin == "Y")) ||
			((mode == "write" || mode == "edit" || mode == "delete") && grantWrite == "user") ||
			((mode != "write" && mode != "edit" && mode != "delete") && grantRead == "user") ||
			((mode == "write" || mode == "edit" || mode == "delete") && grantWrite == "all") ||
			((mode != "write" && mode != "edit" && mode != "delete") && grantRead == "all"):
			isValid = true
		default:
			isValid = false
		}
	}

	return
}

// IsAdmin - Check admin
func IsAdmin(c echo.Context) (result bool) {
	result = false

	user := c.Get("user")
	claims := user.(*jwt.Token).Claims.(*CustomClaims)

	if claims.Admin == "Y" {
		result = true
	}

	return
}
