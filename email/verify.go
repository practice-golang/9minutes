package email

import (
	"log"
	"strings"

	"github.com/emersion/go-msgauth/dkim"
)

func Verify(msg string) (verifications []*dkim.Verification, err error) {
	s := strings.NewReader(msg)

	verifications, err = dkim.Verify(s)
	if err != nil {
		return nil, err
	}

	for _, v := range verifications {
		if v.Err == nil {
			log.Println("Valid key for", v.Domain)
			break
		} else {
			// log.Println("Invalid key for", v.Domain, v.Err)
			err = v.Err
			break
		}
	}

	return verifications, err
}
