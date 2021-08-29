package service

import (
	"alarm_center/internal/config"
	"alarm_center/internal/domain/repo"
	"alarm_center/internal/infras/utils"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)


type EmailService struct {
	config  *config.AppConfig
	repo 	repo.EmailRepository
}

func NewEmailService(config *config.AppConfig, emailRepo repo.EmailRepository) *EmailService {
	return &EmailService{
		config:     config,
		repo: 		emailRepo,
	}
}

func (e *EmailService) SendEmail(c * gin.Context)  {
	var to []string

	mailFrom 	:= c.PostForm("mail_from")
	mailTo 		:= c.PostForm("mail_to")
	typeValue 	:= c.PostForm("mail_type")
	if typeValue == " " {
		utils.JsonRespond(500, "邮件发送类型设置错误", "", c)
	}
	mailType,_ 	:= strconv.Atoi(typeValue)

	for _, item :=  range strings.Split(mailTo,",") {
		to = append(to, item)
	}

	subject := c.PostForm("subject")
	content := c.PostForm("content")

	emailBody := &repo.EmailBody{
		From: mailFrom,
		To: to,
		Subject: subject,
		Body: content,
	}

	m 	:= e.repo.CreateMsg(emailBody)
	if  mailType != 0 && mailType != 1 {
		utils.JsonRespond(500, "邮件发送类型设置错误", "", c)
	}

	if mailType == 1  {
		// 实时发送
		err := e.repo.SendMsg(m)
		if err != nil {
			utils.JsonRespond(500, err.Error(), "", c)
			return
		}
	} else  {
		// 发送到消息队列
		err := e.repo.AddSendQueue(emailBody)
		if err != nil {
			utils.JsonRespond(500, err.Error(), "", c)
			return
		}
	}

	utils.JsonRespond(200, "邮件发送成功", "", c)
}
