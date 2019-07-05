package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

type Conf struct {
	User   string `toml:"user" validate:"required"` //数据库用户
	Addr   string `toml:"addr" validate:"required"` //数据库监听地址
	PStr   string `toml:"pstr" validate:"required"` //数据库密码
	DBName string `toml:"db_name" validate:"required"`
}

//var DB *gorm.DB //standard db

func InitDataBase(conf *Conf) *gorm.DB {
	var err error
	source := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=Local&parseTime=true", conf.User, conf.PStr, conf.Addr, conf.DBName)
	db, err := gorm.Open("mysql", source)
	if err != nil {
		panic(err)
	}

	db.DB().SetMaxOpenConns(50)
	db.DB().SetMaxIdleConns(30)
	db.DB().SetConnMaxLifetime(20 * time.Second)
	db.SingularTable(true) //表名非复数形式

	if err = db.DB().Ping(); err != nil {
		log.Printf("DB.DB().Ping():%v", err.Error())
		panic(err)
	}

	return db
}
