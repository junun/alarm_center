package main

import (
	"alarm_center/internal/config"
	"alarm_center/internal/interfaces"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
)

var cfgFile string


func init() {
	flag.StringVar(&cfgFile, "config_file", "./app.yaml", "config file")
	flag.Parse()
}

func main() {
	config.CfgFile = "./app.yaml"
	viperEntry, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("read config error:%s\n", err.Error())
	}

	appConf, err := config.InitAppConf(viperEntry) // init app config
	if err != nil {
		log.Fatalf("init app config error:%s\n", err.Error())
	}

	r 		:= gin.Default()
	di 		:= interfaces.BuildContainer()
	server 	:= interfaces.NewServer(r, appConf, di)
	server.InitRouter()
	server.CronJob()
	server.Run()
}