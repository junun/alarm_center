package persistence

import (
	"alarm_center/internal/domain/repo"
	"alarm_center/internal/infras/db"
	"github.com/go-redis/redis"
	"github.com/robfig/cron/v3"
	"strconv"
)

type cronRepo struct {
	rdb 	*redis.Client
	cron	*cron.Cron
}

func NewCronRepository(rdb *redis.Client, cron *cron.Cron) repo.CronRepository {
	return &cronRepo {
		rdb:  rdb,
		cron: cron,
	}
}

func (c *cronRepo) JobList() map[string]string {
	return  c.rdb.HGetAll(db.CronNameEntryId).Val()
}

func (c *cronRepo) NewLocalJob(job *repo.Job) error {
	entryID, err :=  c.cron.AddFunc(job.Spec, job.Cmd)
	if err != nil {
		return err
	}
	// 用redis 存储job name和cron id对应关系
	return c.rdb.HSet(db.CronNameEntryId, job.Name, strconv.Itoa(int(entryID))).Err()
}

func (c *cronRepo) RunJob() {
	c.cron.Start()
}

func (c *cronRepo) RemoveJob(name string, id int) error {
	c.cron.Remove(cron.EntryID(id))
	return c.rdb.HDel(db.CronNameEntryId, name).Err()
}