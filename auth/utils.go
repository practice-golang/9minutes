package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	JwtKey []byte
)

// CustomClaims - jwt custom claim
type CustomClaims struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Admin        string `json:"admin"`
	RefreshUntil int64  `json:"refresh-until"`
	jwt.StandardClaims
}

// PrepareToken - Create and return token
func PrepareToken(data interface{}) (string, error) {
	d := data.(map[string]interface{})

	claims := &CustomClaims{
		fmt.Sprint(d["NAME"]),
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
