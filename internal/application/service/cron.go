package service

import (
	"alarm_center/internal/config"
	"alarm_center/internal/domain/repo"
	"alarm_center/internal/infras/db"
	"github.com/go-redis/redis"
)

type CronService struct {
	config  *config.AppConfig
	rdb 	*redis.Client
	repo 	repo.CronRepository
	eRepo 	repo.EmailRepository
	dRepo 	repo.DingTalkRepository
}

func NewCronService(
	config *config.AppConfig,
	rdb *redis.Client,
	cronRepo repo.CronRepository,
	eRepo repo.EmailRepository,
	dRepo repo.DingTalkRepository) *CronService {
	return &CronService{
		config: config,
		rdb: rdb,
		repo: 	cronRepo,
		eRepo: 	eRepo,
		dRepo: dRepo,
	}
}

func (cron *CronService) JobRun()  {
	cron.repo.RunJob()
}

// 系统任务， 重启后要重新加载
func (cron *CronService) StartJobOnBoot()  {
	cron.rdb.Del(db.CronNameEntryId)
	emailQueueMsgSend := &repo.Job{
		Name: "emailQueueMsgSend",
		Spec: "* * * * *", // every min
		Cmd: cron.eRepo.SendQueueMsg,
	}
	err := cron.repo.NewLocalJob(emailQueueMsgSend)
	if err != nil {
		cron.ErrorJobSendDingTalkMsg("添加定时任务 emailQueueMsgSend 失败：" + err.Error())
	}

	dingTalkMsgSend := &repo.Job{
		Name: "dingTalkMsgSend",
		Spec: "@every 30s", // every 30s
		Cmd: cron.dRepo.SendQueueDingTalkMsg,
	}
	err = cron.repo.NewLocalJob(dingTalkMsgSend)
	if err != nil {
		cron.ErrorJobSendDingTalkMsg("添加定时任务 dingTalkMsgSend 失败：" + err.Error())
	}
}

func (cron *CronService) ErrorJobSendDingTalkMsg(context string)  {
	dingTalkClient,_ := cron.dRepo.GenerateClient("default")
	robotSendRequest := &repo.RobotSendRequest{
		MsgType: "text",
		Text:    repo.Text{Content: context},
		At:      repo.At{AtMobiles: []string{}, IsAtAll: false},
	}

	cron.dRepo.SendMessage("default", dingTalkClient, robotSendRequest)
}

