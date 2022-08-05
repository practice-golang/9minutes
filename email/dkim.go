package email

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"log"
	"strings"

	"github.com/emersion/go-msgauth/dkim"
)

func GetSignedMessage(msg []byte, domain, pemKEY string) ([]byte, error) {
	r := strings.NewReader(string(msg))
	block, _ := pem.Decode([]byte(pemKEY))

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	options := &dkim.SignOptions{
		Domain:   domain,
		Selector: "mail",
		Signer:   privateKey,
	}

	var b bytes.Buffer
	if err := dkim.Sign(&b, r, options); err != nil {
		log.Fatal(err)
	}

	resultBytes := b.Bytes()
	// resultString := b.String()
	// resultString := string(resultBytes)

	return resultBytes, nil
}
