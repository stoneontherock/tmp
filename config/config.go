package config

import (
	"aa/aalog"
	"aa/db"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Conf struct {
	HTTP struct {
		ListenAddr string `toml:"listen_addr"` //HTTP服务监听地址
	}

	DB db.Conf

	LOG aalog.LogConf
}

func GetConf() *Conf {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	var c Conf
	_, err = toml.DecodeFile(dir+"/config/config.toml", &c)
	if err != nil {
		panic(err)
	}

	//todo 完成校验
	validateConf(&c)
	return &c
}

func validateConf(c *Conf) {
	if c.HTTP.ListenAddr == "" {
		panic(fmt.Sprintf("conf file: one or more string value in [HTTP] are empty. %+v", c.HTTP))
	}

	if c.DB.Addr == "" || c.DB.User == "" || c.DB.DBName == "" {
		panic(fmt.Sprintf("conf file: one or more string value in [DB] are empty. %+v", c.DB))
	}

	if c.LOG.MaxSize < 1 || c.LOG.MaxBackup <= 1 || c.LOG.Level > 6 {
		panic("value of `max_size` or `max_backup` or `level` is invalid")
	}
}
