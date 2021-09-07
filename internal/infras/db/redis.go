package db

import (
	"github.com/go-redis/redis"
)

const (
	EmailSendQueue  	= "email_send_queue"
	DingTalkSendQueue  	= "ding_talk_send_queue"
	CronNameEntryId 	= "cron_name_entry_id"
	RedisScanCount		= 100
)
type RedisConn struct {
	Addr        string
	Password    string
	DB			int
}

type SentinelConn struct {
	MasterName 	string
	SentinelNodes []string
	Password 	string
	DB 			int
	Client 		*redis.Client
}

type ClusterConn struct {
	StartNodes []string
	Password 	string
	DB 			int
}

func  (r *RedisConn) ConnectDB() *redis.Client {
	var rdb *redis.Client
	rdb = redis.NewClient(&redis.Options{
		Addr: r.Addr,
		Password: r.Password,
		DB: r.DB,
	})

	return rdb
}

func GetQueueLength(rdb *redis.Client, name string) int64{
	res, _ 	:= rdb.Exists(name).Result()
	if res == 0 {
		return 0
	}
	len,e:= rdb.LLen(name).Result()
	if e != nil{
		return 0
	}

	return len
}

