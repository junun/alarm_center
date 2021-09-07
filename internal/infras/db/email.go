package db

import "gopkg.in/gomail.v2"

type EmailConfig struct {
	User 	string
	Pass 	string
	Host 	string
	Port	int
}

func  (e *EmailConfig) InitDialer() *gomail.Dialer {
	return gomail.NewDialer(e.Host, e.Port, e.User, e.Pass)
}



