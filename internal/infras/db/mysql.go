package db

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

// ConnectDB create db entry
func (dbConf *DBConf) ConnectDB() (*gorm.DB, error) {
	log.Println("mysql orm of dbr run...")
	dsn, err := dbConf.DSN()
	if err != nil {
		log.Fatalf("format dsn err: %s\n", err.Error())
	}
	var db *gorm.DB
	//db, err = sql.Open("mysql", dsn)
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("open mysql err: %s\n", err.Error())
	}

	// 设置连接池
	if dbConf.MaxLifetime == 0 {
		dbConf.MaxLifetime = 20 * time.Minute
	}
	if dbConf.MaxIdleTime == 0 {
		dbConf.MaxIdleTime = 10 * time.Minute
	}
	if dbConf.MaxOpenConns == 0 {
		dbConf.MaxOpenConns = 300
	}
	if dbConf.MaxIdleConns == 0 {
		dbConf.MaxIdleConns = 10
	}

	db.SingularTable(true)
	db.DB().SetConnMaxIdleTime(dbConf.MaxIdleTime)
	db.DB().SetConnMaxLifetime(dbConf.MaxLifetime)
	db.DB().SetMaxOpenConns(dbConf.MaxIdleConns)
	db.DB().SetMaxIdleConns(dbConf.MaxOpenConns)

	//db.SetConnMaxIdleTime(dbConf.MaxIdleTime)
	//db.SetConnMaxLifetime(dbConf.MaxLifetime)
	//db.SetMaxOpenConns(dbConf.MaxOpenConns)
	//db.SetMaxIdleConns(dbConf.MaxIdleConns)

	return db, nil
}

// DBConf DB config
type DBConf struct {
	Type 	  string
	Ip        string
	Port      int // 默认3306
	User      string
	Password  string
	Database  string
	Charset   string // 字符集 utf8mb4 支持表情符号
	Collation string // 整理字符集 utf8mb4_unicode_ci

	MaxIdleConns int // 空闲pool个数
	MaxOpenConns int // 最大open connection个数

	// sets the maximum amount of time a connection may be reused.
	// 设置连接可以重用的最大时间
	// 给db设置一个超时时间，时间小于数据库的超时时间
	MaxLifetime time.Duration // 数据库超时时间
	MaxIdleTime time.Duration // 最大空闲时间

	// 连接超时/读取超时/写入超时设置
	Timeout      time.Duration // Dial timeout
	ReadTimeout  time.Duration // I/O read timeout
	WriteTimeout time.Duration // I/O write timeout

	ParseTime bool   // 格式化时间类型
	Loc       string // 时区字符串 Local,PRC
}

func (conf *DBConf) DSN() (string, error) {
	if conf.Type == "" {
		conf.Type = "mysql"
	}

	if conf.Ip == "" {
		conf.Ip = "127.0.0.1"
	}

	if conf.Port == 0 {
		conf.Port = 3306
	}

	if conf.Charset == "" {
		conf.Charset = "utf8mb4"
	}

	// 默认字符序，定义了字符的比较规则
	if conf.Collation == "" {
		conf.Collation = "utf8mb4_general_ci"
	}

	if conf.Loc == "" {
		conf.Loc = "Local"
	}

	if conf.Timeout == 0 {
		conf.Timeout = 10 * time.Second
	}

	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = 5 * time.Second
	}

	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = 5 * time.Second
	}

	// mysql connection time loc.
	loc, err := time.LoadLocation(conf.Loc)
	if err != nil {
		return "", err
	}

	// mysql config
	mysqlConf := mysql.Config{
		User:   conf.User,
		Passwd: conf.Password,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%d", conf.Ip, conf.Port),
		DBName: conf.Database,
		// Connection parameters
		Params: map[string]string{
			"charset": conf.Charset,
		},
		Collation:            conf.Collation,
		Loc:                  loc,               // Location for time.Time values
		Timeout:              conf.Timeout,      // Dial timeout
		ReadTimeout:          conf.ReadTimeout,  // I/O read timeout
		WriteTimeout:         conf.WriteTimeout, // I/O write timeout
		AllowNativePasswords: true,              // Allows the native password authentication method
		ParseTime:            conf.ParseTime,    // Parse time values to time.Time
	}

	dsn	:= fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		mysqlConf.User,
		mysqlConf.Passwd,
		mysqlConf.Addr,
		mysqlConf.DBName)
	//return mysqlConf.FormatDSN(), nil
	return dsn, nil
}

