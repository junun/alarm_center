package repo

import "gopkg.in/gomail.v2"

type EmailBody struct {
	From        string
	To			[]string
	Subject		string
	Body 		string
}

type EmailRepository interface {
	CreateMsg(b *EmailBody) *gomail.Message
	CreateAnnexMsg(b *EmailBody, annex string) *gomail.Message
	AddSendQueue(b *EmailBody) error
	SendMsg(gm *gomail.Message) error
	SendQueueMsg()
}

