package email

import (
	"net"
	"strings"
)

func GetMX(to string) (mx string, err error) {
	domain, err := getDomain(to)
	if err != nil {
		return "", err
	}

	var mxs []*net.MX
	mxs, err = net.LookupMX(domain)
	if err != nil {
		return "", err
	}

	for _, x := range mxs {
		mx = strings.TrimSuffix(x.Host, ".")
		return
	}

	return "", nil
}
