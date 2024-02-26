package smtpurl

import (
	"fmt"
	"net/smtp"
	"net/url"
	"strings"
)

const (
	defaultPort = "25"
	authToken   = ";AUTH="
)

// Parse STMP URL
func Parse(raw string) (string, smtp.Auth, error) {
	if !strings.HasPrefix(raw, "smtp://") {
		return "", nil, fmt.Errorf("invalid url: %s", raw)
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", nil, err
	}
	sh := strings.Split(u.Host, ":")
	var hostname, port string
	switch len(sh) {
	case 1:
		hostname = sh[0]
		port = defaultPort
	case 2:
		hostname = sh[0]
		port = sh[1]
	default:
		return "", nil, fmt.Errorf("invalid url: %s", raw)
	}
	host := fmt.Sprintf("%s:%s", hostname, port)
	var auth smtp.Auth
	if strings.Contains(strings.ToUpper(u.User.Username()), authToken) {
		// ref: https://datatracker.ietf.org/doc/html/draft-earhart-url-smtp-00
		su := strings.Split(u.User.Username(), authToken)
		switch strings.ToUpper(su[1]) {
		case "PLAIN":
			pass, _ := u.User.Password()
			auth = smtp.PlainAuth("", su[0], pass, hostname)
		case "CRAM-MD5":
			pass, _ := u.User.Password()
			auth = smtp.CRAMMD5Auth(su[0], pass)
		default:
			return "", nil, fmt.Errorf("unsupported auth method: %s", strings.ToUpper(su[1]))
		}
	} else if u.User != nil{
		// PLAIN
		pass, _ := u.User.Password()
		auth = smtp.PlainAuth("", u.User.Username(), pass, hostname)
	} else {
		return host, auth, nil
	}
	return host, auth, nil
}
