package config

import (
	"aa/aalog"
	"aa/db"
	"aa/panicerr"
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
    panicerr.PE(err)

	_, err = toml.DecodeFile(dir+"/config/config.toml", &C)
	panicerr.PE(err)

	err = validator.New().Struct(&C)
	panicerr.PE(err,fmt.Sprintf("%s配置有误,Conf结构体:%+v\n", dir+"/config/config.toml", C))
}

