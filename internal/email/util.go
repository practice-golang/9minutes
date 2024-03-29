package email

import (
	"net/mail"
	"strings"
)

func getAddress(addr string) string {
	e, err := mail.ParseAddress(addr)
	if err != nil {
		return addr
	}

	return e.String()
}

func getDomain(email string) (result string, err error) {
	parsedAddress, err := mail.ParseAddress(email)
	if err != nil {
		return "", err
	}

	result = strings.Split(parsedAddress.Address, "@")[1]

	return result, nil
}
