package aalog

import (
	"aa/panicerr"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path/filepath"
)

type LogConf struct {
	Level     uint32 `toml:"level" validate:"lte=6"`      //单位MB
	MaxSize   int    `toml:"max_size" validate:"gte=1"`   //单位MB
	MaxBackup int    `toml:"max_backup" validate:"gte=1"` //最多备份数
}

func InitLog(c *LogConf) {
	fmtr := new(logrus.TextFormatter)
	fmtr.FullTimestamp = true                      // 显示完整时间
	fmtr.TimestampFormat = "06-01-02 15:04:05.000" // 时间格式
	fmtr.DisableTimestamp = false                  // 禁止显示时间
	fmtr.DisableColors = true                      // 禁止颜色显示

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	panicerr.PE(err)

	dir += "/log"
	err = os.MkdirAll(dir, 0700)
	panicerr.PE(err,"创建目录")

	f := filepath.Join(dir, filepath.Base(os.Args[0])+".log")

	log.Printf("log file: %s", f)

	jack := &lumberjack.Logger{
		Filename: f,         //如果没目录，它会自己建立
		MaxSize:  c.MaxSize, //MBytes
		//MaxAge: 1, //day
		MaxBackups: c.MaxBackup,
		LocalTime:  true,
		Compress:   true,
	}
	logrus.SetOutput(jack)

	logrus.SetFormatter(fmtr)

	logrus.SetLevel(logrus.DebugLevel)
	return
}
