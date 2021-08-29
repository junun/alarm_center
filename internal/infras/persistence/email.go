package persistence

import (
	"alarm_center/internal/domain/repo"
	"alarm_center/internal/infras/db"
	"encoding/json"
	"github.com/go-redis/redis"
	"gopkg.in/gomail.v2"
	"sync"
)

type emailRepo struct {
	rdb *redis.Client
	gd *gomail.Dialer
}

func NewEmailRepository(rdb *redis.Client, gd *gomail.Dialer) repo.EmailRepository {
	return &emailRepo {
		rdb: rdb,
		gd: gd,
	}
}

// 生成消息体
func (e *emailRepo) CreateMsg(b *repo.EmailBody) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From","Monitor" + "<" + b.From + ">")
	m.SetHeader("To", b.To...)  //发送给多个用户
	m.SetHeader("Subject", b.Subject)  //设置邮件主题
	m.SetBody("text/plain", b.Body)  //设置邮件正文

	return m
}

// 生成带附件的消息体， 不支持非实时发送。
func (e *emailRepo) CreateAnnexMsg(b *repo.EmailBody, annex string) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From","Monitor" + "<" + b.From + ">")
	m.SetHeader("To", b.To...)  //发送给多个用户
	m.SetHeader("Subject", b.Subject)  //设置邮件主题
	m.SetBody("text/html", b.Body)  //设置邮件正文
	m.Attach(annex)							// 设置附件

	return m
}

func (e *emailRepo) AddSendQueue(b *repo.EmailBody) error {
	item, err := json.Marshal(b)
	if err != nil {
		return err
	}

	e.rdb.LPush(db.EmailSendQueue, string(item))
	return nil
}

func (e *emailRepo) SendMsg(msg *gomail.Message) error {
	return e.gd.DialAndSend(msg)
}

func (e *emailRepo) SendQueueMsg()  {
	var wg sync.WaitGroup
	for {
		length := db.GetQueueLength(e.rdb, db.EmailSendQueue)
		if length == 0 {
			break
		}

		msg, _  := e.rdb.RPop(db.EmailSendQueue).Result()
		var data *repo.EmailBody
		json.Unmarshal([]byte(msg), &data)
		//if err = json.Unmarshal([]byte(msg), &data); err != nil {
		//	return err
		//}

		//协程执行该任务
		wg.Add(1)
		go func(data *repo.EmailBody) {
			m := e.CreateMsg(data)
			err := e.gd.DialAndSend(m)
			// 发送失败 重新加入队列 也可以同时记录到日志
			if err != nil {
				e.AddSendQueue(data)
			}
			defer wg.Done()
		}(data)
	}
}




