package email

import (
	"net/mail"
	"strings"
)

func getDomain(email string) (result string, err error) {
	parsedAddress, err := mail.ParseAddress(email)
	if err != nil {
		return "", err
	}

	result = strings.Split(parsedAddress.Address, "@")[1]

	return result, nil
}
