package service

import (
	"alarm_center/internal/config"
	"alarm_center/internal/domain/repo"
	"alarm_center/internal/infras/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

// DingtalkService
type DingtalkService struct {
	config  *config.AppConfig
	repo 	repo.DingTalkRepository
}

// NewUserService return user service
func NewDingtalkService(config *config.AppConfig, dingtalkRepo repo.DingTalkRepository) *DingtalkService {
	return &DingtalkService{
		config:     config,
		repo: 		dingtalkRepo,
	}
}

func (d *DingtalkService) SendDingTalk(c * gin.Context)  {
	var lists []string
	isAtAll 	:= false

	msgtype 	:= c.PostForm("msgtype")
	content 	:= c.PostForm("content")
	name	 	:= c.PostForm("name")

	all 		:= c.PostForm("all")
	if all == "true" {
		isAtAll = true
	}

	members 	:= c.PostForm("members")
	for _, item :=  range strings.Split(members,",") {
		lists = append(lists, item)
	}

	client, err := d.repo.GenerateClient(name)
	if err != nil {
		utils.JsonRespond(500, err.Error(), "", c)
		return
	}

	var robotSendRequest *repo.RobotSendRequest
	var e error
	switch msgtype {
	case "markdown":
		// markdown message
		title 	:= c.PostForm("title")
		robotSendRequest = &repo.RobotSendRequest{
			MsgType: "markdown",
			Markdown: repo.Markdown{
				Title: title,
				Text: content,
			},
			At: repo.At{
				AtMobiles: lists,
				IsAtAll:   isAtAll,
			},
		}
	case "link":
		// link message
		title 	:= c.PostForm("title")
		msgurl 	:= c.PostForm("msgurl")
		picurl 	:= c.PostForm("picurl")
		robotSendRequest = &repo.RobotSendRequest{
			MsgType: "link",
			Link: repo.Link{
				Title:      title,
				Text:       content,
				PicUrl:     picurl,
				MessageUrl: msgurl,
			},
		}
	default:
		// text message
		robotSendRequest = &repo.RobotSendRequest{
			MsgType: "text",
			Text:    repo.Text{Content: content},
			At:      repo.At{AtMobiles: lists, IsAtAll: isAtAll},
		}
	}

	_, e = d.repo.SendMessage(name, client, robotSendRequest)

	if e != nil {
		utils.JsonRespond(500, e.Error(), "", c)
		return
	}

	utils.JsonRespond(200, "发送成功", "", c)
}
