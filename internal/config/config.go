package config

import (
	"alarm_center/internal/infras/db"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"log"
	"path/filepath"
	"time"
)

type AppConfig struct {
	Port      int
	DBDsn     string
	AppEnv    string
	AppDebug  bool
	GraceWait time.Duration
}

type DingTalk struct {
	Webhook   string
	Secret    string
	Keyword   string
}


var (
	// CfgFile 配置文件路径
	CfgFile string
	DingTalkConfig = make(map[string]*DingTalk)
)

// LoadConfig read config from file.
func LoadConfig() (*viper.Viper, error) {
	viperEntry := viper.New()
	if CfgFile != "" {
		// Use config file from the flag.
		viperEntry.SetConfigFile(CfgFile)
	} else {
		configDir, err := filepath.Abs("./")
		if err != nil {
			return nil, err
		}

		log.Println("config_dir: ", configDir)
		// Search config in ./ directory with name "app.yaml"
		viperEntry.AddConfigPath(configDir)
		viperEntry.SetConfigType("yaml")
		viperEntry.SetConfigName("app")
	}

	viperEntry.AutomaticEnv()
	if err := viperEntry.ReadInConfig(); err != nil {
		return nil, err
	}
	log.Println("using config file:", viperEntry.ConfigFileUsed())

	for k,v := range viperEntry.GetStringMap("dingtalk") {
		var d *DingTalk
		mapstructure.Decode(v, &d)
		DingTalkConfig[k] = d
	}

	//fmt.Println(fmt.Println(DingTalkConfig["default"].Keyword))

	return viperEntry, nil
}

// InitAppConf 初始化配置
func InitAppConf(viperEntry *viper.Viper) (*AppConfig, error) {
	conf := &AppConfig{}
	if err := viperEntry.UnmarshalKey("AppConfig", conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func InitEmailConf(viperEntry *viper.Viper) (*gomail.Dialer, error) {
	conf := &db.EmailConfig{}
	if err := viperEntry.UnmarshalKey("Email", conf); err != nil {
		return nil, err
	}

	return conf.InitDialer(), nil
}

// InitDBConf init db
func InitDBConf(viperEntry *viper.Viper) (*gorm.DB, error) {
	conf := &db.DBConf{}
	if err := viperEntry.UnmarshalKey("DBConfig", conf); err != nil {
		return nil, err
	}

	return conf.ConnectDB()
}

// InitDBConf init redis
func InitRedisConn(viperEntry *viper.Viper) (*redis.Client, error) {
	conf := &db.RedisConn{}
	if err := viperEntry.UnmarshalKey("RedisConfig", conf); err != nil {
		return nil, err
	}

	return conf.ConnectDB(), nil
}
