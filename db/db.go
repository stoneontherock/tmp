package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

type Conf struct {
	User   string `toml:"user"` //数据库用户
	Addr   string `toml:"addr"` //数据库监听地址
	PStr   string `toml:"pstr"` //数据库密码
	DBName string `toml:"db_name"`
}

var DB *gorm.DB //standard db


func InitDataBase(conf *Conf) {
	var err error
	source := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=Local&parseTime=true", conf.User, conf.PStr, conf.Addr, conf.DBName)
	DB, err = gorm.Open("mysql", source)
	if err != nil {
		panic(err)
	}

	DB.DB().SetMaxOpenConns(50)
	DB.DB().SetMaxIdleConns(30)
	DB.DB().SetConnMaxLifetime(20 * time.Second)
	DB.SingularTable(true) //表名非复数形式

	if err = DB.DB().Ping(); err != nil {
		log.Printf("DB.DB().Ping():%v", err.Error())
		panic(err)
	}
}
