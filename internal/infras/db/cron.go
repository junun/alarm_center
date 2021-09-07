package db

import (
	"github.com/robfig/cron/v3"
)

func  NewCron() *cron.Cron {
	return cron.New()
}