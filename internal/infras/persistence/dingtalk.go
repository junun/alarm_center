package persistence

import (
	"alarm_center/internal/config"
	"alarm_center/internal/domain/repo"
	"alarm_center/internal/infras/db"
	"alarm_center/internal/infras/utils"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// dingtalk repo
type dingtalkRepo struct {
	rdb *redis.Client
}

// NewDingtalkRepository new dingtalk repo
func NewDingtalkRepository(rdb *redis.Client) repo.DingTalkRepository {
	return &dingtalkRepo {
		rdb: rdb,
	}
}

func (d *dingtalkRepo) GenerateClient(name string) (*repo.DingTalkClient, error) {
	var client *repo.DingTalkClient
	if name != ""  {
		client = &repo.DingTalkClient{
			Webhook: config.DingTalkConfig[name].Webhook,
			Secret: config.DingTalkConfig[name].Secret,
			Keyword: config.DingTalkConfig[name].Keyword,
		}
	} else {
		client = &repo.DingTalkClient{
			Webhook: config.DingTalkConfig["default"].Webhook,
			Secret: config.DingTalkConfig["default"].Secret,
			Keyword: config.DingTalkConfig["default"].Keyword,
		}
	}
	return client, nil
}

func (d *dingtalkRepo) Sign(timestamp string, c *repo.DingTalkClient) (string, error) {
	stringToSign := fmt.Sprintf("%s\n%s", timestamp, c.Secret)
	h := hmac.New(sha256.New, []byte(c.Secret))
	if _, err := io.WriteString(h, stringToSign); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func (d dingtalkRepo) GetPushUrl(c *repo.DingTalkClient) (string, error) {
	if c.Secret == "" {
		return c.Webhook, nil
	}
	timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)
	sign, err := d.Sign(timestamp, c)
	if err != nil {
		return c.Webhook, err
	}

	query := url.Values{}
	query.Set("timestamp", timestamp)
	query.Set("sign", sign)
	return c.Webhook + "&" + query.Encode(), nil
}

func (e *dingtalkRepo) AddSendQueue(name string, b *repo.RobotSendRequest)  {
	item, _ := json.Marshal(b)
	key 	:= db.DingTalkSendQueue+"_"+name
	e.rdb.LPush(key, string(item))
}

func (d *dingtalkRepo) SendMessage(name string, c *repo.DingTalkClient, request *repo.RobotSendRequest) (*repo.RobotSendResponse, error) {
	if c.Keyword != "" {
		request.Text.Content = "【" + c.Keyword + "】" + request.Text.Content
	}
	key 	:= utils.SHA1HashString(c.Webhook)
	num,_  	:= d.rdb.Get(key).Int()
	// 一分钟发送不能超过18, 超过放入redis发送队列
	if num > 18 {
		d.AddSendQueue(name, request)
		return nil, errors.New(fmt.Sprintf("已经到达单分钟最大的发送能力：%d", num))
	}

	d.rdb.Set(key, num+1, 60*time.Second)
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return  nil, err
	}

	pushUrl, err := d.GetPushUrl(c)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(pushUrl, "application/json", bytes.NewReader(reqBytes))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("http response status code is: %d", resp.StatusCode))
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var oResponse repo.RobotSendResponse
	if err = json.Unmarshal(responseBody, &oResponse); err != nil {
		return nil, err
	}
	if oResponse.ErrCode != 0 {
		return &oResponse, errors.New(fmt.Sprintf("response: %s", responseBody))
	}
	return &oResponse, nil
}

func (d *dingtalkRepo) SendQueueDingTalkMsg()  {
	var wg sync.WaitGroup
	for {
		keyIndex := db.DingTalkSendQueue + "_*"
		keys, _ := d.rdb.Scan(0, keyIndex, db.RedisScanCount).Val()
		if len(keys) == 0 {
			break
		}
		for _,v := range keys {
			length := db.GetQueueLength(d.rdb, v)
			if length == 0 {
				break
			}

			name 	:= utils.LastString(strings.Split(v, "_"))
			client,_:= d.GenerateClient(name)
			key 	:= utils.SHA1HashString(client.Webhook)
			num,_  	:=d.rdb.Get(key).Int()
			// 一分钟发送不能超过18, 超过放入redis发送队列
			if num > 18 {
				return
			}

			msg, _  := d.rdb.RPop(v).Result()
			var data *repo.RobotSendRequest
			json.Unmarshal([]byte(msg), &data)

			//协程执行该任务
			wg.Add(1)
			go func(name string, client *repo.DingTalkClient, data *repo.RobotSendRequest) {
				_, err 	:= d.SendMessage(name, client, data)
				// 发送失败 重新加入队列 也可以同时记录到日志
				if err 	!= nil {
					d.AddSendQueue(name, data)
				}
				defer wg.Done()
			}(name, client, data)
		}
	}
}