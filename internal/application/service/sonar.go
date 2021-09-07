package service

import (
	"alarm_center/internal/config"
	"alarm_center/internal/domain/repo"
	"alarm_center/internal/infras/db"
	"alarm_center/internal/infras/utils"
	"github.com/gin-gonic/gin"
)


type SonarService struct {
	config  *config.AppConfig
	repo 	repo.SonarRepository
	emailRepo 	repo.EmailRepository
}

func NewSonarService(config *config.AppConfig, sonarRepo repo.SonarRepository, emailRepo repo.EmailRepository) *SonarService {
	return &SonarService{
		config: config,
		repo: 	sonarRepo,
		emailRepo: emailRepo,
	}
}

func (s SonarService) GetMeasuresComponent(c * gin.Context) {
	name 	:= c.Param("name")
	email 	:=  c.Param("email")

	res, err := s.repo.GetMeasuresComponent(db.SonarMeasuresComponent, name)
	if err != nil {
		utils.JsonRespond(500, err.Error(), "", c)
		return
	}
	var to []string
	to = append(to, email)

	msg := "<!DOCTYPE html>" +
		"<html lang=\"en\">\n" +
		"<head>\n" +
		"    <title></title>\n" +
		"   <meta charset=\"utf-8\">\n" +
		"</head>\n" +
		"<body>\n" +
		"<div class=\"page\" style=\"margin-left: 30px\">\n" +
		"<h3>"+ email +", 你好</h3>\n" +
		"<h3> 本次提交代码检查结果如下</h3>\n" +
		"<h3> 项目名称：" + name + " </h3>\n" +
		"<h4>一、总体情况</h4>\n" +
		"<ul>\n" +
		"<li style=\"font-weight:bold;\">\n" +
		"检查结果: &nbsp;" + res.AlertStatus + "\n" +
		"</li>\n" +
		"<li style=\"font-weight:bold;\">\n" +
		"新增bug: &nbsp;" + res.NewBugs + "\n" +
		"</li>\n" +
		"<li style=\"font-weight:bold;\">\n" +
		"bug总数: &nbsp;" + res.Bugs + "\n" +
		"</li>\n" +
		"<li style=\"font-weight:bold;\">\n" +
		"新增漏洞: &nbsp;" + res.NewVulnerabilities + "\n" +
		"</li>\n" +
		"<li style=\"font-weight:bold;\">\n" +
		"漏洞总数: &nbsp;" + res.Vulnerabilities + "\n" +
		"</li>\n" +
		"<li style=\"font-weight:bold;\">\n" +
		"新增安全热点: &nbsp;" + res.NewSecurityHotspots + "\n" +
		"</li>\n" +
		"<li style=\"font-weight:bold;\">\n" +
		"安全热点总数: &nbsp;" + res.SecurityHotspots + "\n" +
		"</li>\n" +
		"<li style=\"font-weight:bold;\">\n" +
		"新增不规范代码: &nbsp;" + res.NewCodeSmells + "\n" +
		"</li>\n" +
		"<li style=\"font-weight:bold;\">\n" +
		"不规范代码总数: &nbsp;" + res.CodeSmells + "\n" +
		"</li>\n" +
		"</ul>\n" +
		"<h4>二、信息详情</h4>\n" +
		"<ul>\n" +
		"<li style=\"font-weight:bold;margin-top: 10px;\">\n" +
		"URL地址：&nbsp;\n" +
		"<a style=\"font-weight:bold;\" href=\"" + config.SonarConfig.Url  + "/dashboard?id=" + name + "\">" + name + "\n" +
		"</a>\n" +
		"</li>\n" +
		"</ul>\n" +
		"</div>\n</body>\n</html>"

	emailBody := &repo.EmailBody{
		From: "470499989@qq.com",
		To: to,
		Subject: name,
		Body: msg,
	}

	m 	:= s.emailRepo.CreateMsg(emailBody)
	m.SetBody("text/html", msg)
	// 实时发送
	err = s.emailRepo.SendMsg(m)
	if err != nil {
		utils.JsonRespond(500, err.Error(), "", c)
		return
	}

	utils.JsonRespond(200, "sonar", "", c)
}
