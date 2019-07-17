package db

import (
	"aa/panicerr"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
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
    panicerr.PE(err,"连接Mysql")

	db.DB().SetMaxOpenConns(50)
	db.DB().SetMaxIdleConns(30)
	db.DB().SetConnMaxLifetime(20 * time.Second)
	db.SingularTable(true) //表名非复数形式

	//todo :Debug, 发布后注释掉
	db.SetLogger(log.New(os.Stdout, "", log.LstdFlags))
	db.LogMode(true)

	err = db.DB().Ping()
	panicerr.PE(err,"检查数据库连接")

	return db
}
