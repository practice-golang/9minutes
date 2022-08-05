package email

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"strings"
)

func encodeRFC2047(str string) string {
	addr := mail.Address{Address: str}
	return strings.Trim(addr.String(), " <>@")
}

func ComposeMimeMailTEXT(to string, from string, subject string, body string) []byte {
	header := make(map[string]string)
	header["From"] = getAddress(from)
	header["To"] = getAddress(to)
	header["Subject"] = encodeRFC2047(subject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	return []byte(message)
}

func ComposeMimeMailHTML(to string, from string, subject string, body string) []byte {
	header := make(map[string]string)
	header["From"] = getAddress(from)
	header["To"] = getAddress(to)
	header["Subject"] = encodeRFC2047(subject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	return []byte(message)
}
