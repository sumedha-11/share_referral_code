package sql

import (
	"fmt"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	pkglog "github.com/sumedha-11/share_referral_code/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type DBConfig struct {
	Url      string `json:"url"`
	UserName string `json:"username"`
	Password string `json:"password"`
	DataBase string `json:"database"`
}

func DBConn(conf *DBConfig) (err error) {
	url := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.UserName, conf.Password, conf.Url, conf.DataBase)
	//DB, err = gorm.Open("mysql", url)
	DB, err = gorm.Open(mysql.Open(url), &gorm.Config{
		Logger: logger.New(pkglog.DBLogger, logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      true,
		}),
	})
	if err != nil {
		return
	}

	db, err := DB.DB()
	if err != nil {
		return
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(5)
	return
}
