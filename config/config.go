package config

import (
	"aa/aalog"
	"aa/db"
	"fmt"
	"github.com/BurntSushi/toml"
	"gopkg.in/go-playground/validator.v9"
	"os"
	"path/filepath"
	"time"
)

type GRPCConf struct {
	Addr        string `toml:"listen_addr" validate:"required"`
	Certificate string `toml:"certificate" validate:"file"`
	Key         string `toml:"key" validate:"file"`
}

type Conf struct {
	HTTP struct {
		ListenAddr string `toml:"listen_addr" validate:"required"` //HTTP服务监听地址
	}

	GRPC GRPCConf

	JWT struct {
		Pstr   string        `toml:"pstr" validate:"gte=8"`
		Expire time.Duration `toml:"expire" validate:"gt=0"`
	}

	DB db.Conf

	LOG aalog.LogConf
}

var C Conf

func GetConf() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	_, err = toml.DecodeFile(dir+"/config/config.toml", &C)
	if err != nil {
		panic(err)
	}

	err = validator.New().Struct(&C)
	if err != nil {
		panic(fmt.Sprintf("%s配置有误,Conf结构体:%+v\nErr:%v\n", dir+"/config/config.toml", C, err))
	}
}

//func validateConf(c *Conf) {
//	if c.HTTP.ListenAddr == "" {
//		panic(fmt.Sprintf("conf file: one or more string value in [HTTP] are empty. %+v", c.HTTP))
//	}
//
//	if _, err := net.ResolveTCPAddr("tcp", c.GRPC.Addr); err != nil {
//		panic(fmt.Sprintf("conf file: invalid value %+v in [GRPC] block", c.GRPC))
//	}
//
//	if c.DB.Addr == "" || c.DB.User == "" || c.DB.DBName == "" {
//		panic(fmt.Sprintf("conf file: one or more string value in [DB] are empty. %+v", c.DB))
//	}
//
//	if c.LOG.MaxSize < 1 || c.LOG.MaxBackup <= 1 || c.LOG.Level > 6 {
//		panic("value of `max_size` or `max_backup` or `level` is invalid")
//	}
//}
