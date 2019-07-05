package httpapi

import (
	"aa/config"
	"fmt"
	"github.com/casbin/casbin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var jwtSecret []byte

func Init(db *gorm.DB) {
	DB = db
	DB.AutoMigrate(&Resource{}, &Role{}, &User{})
	createSuperAdmin()
	initEnforcer()
	jwtSecret = []byte(config.C.JWT.Pstr)
}

func createSuperAdmin() {
	pstr := os.Getenv("sa_pstr")
	if pstr == "" {
		pstr = "SuperAdmin@123"
	}
	salt := getRandomStr(8)
	hash := md5sum(pstr + salt)

	sa := User{Name: SA}
	err := DB.First(&sa).Error
	if err == nil {
		logrus.Debugf("超管账户已经存在: %v\n", sa)
		DB.Model(&sa).Updates(map[string]interface{}{"pstr": hash, "salt": salt})
		return
	}

	if err != gorm.ErrRecordNotFound {
		panic(fmt.Sprintf("查询超管记录时出现错误:%v", err))
	}

	sa.Salt = getRandomStr(8)
	sa.Pstr = md5sum(pstr + sa.Salt)
	sa.Domain = "root"

	err = DB.Create(&sa).Error
	if err != nil {
		panic(fmt.Sprintf("创建超级管理员账户失败:%v", err))
	}
}

var Enforcer *casbin.Enforcer

func initEnforcer() {
	binDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	Enforcer = casbin.NewEnforcer(binDir+"/config/model.conf", false)
	err := loadAllRoleRescourcePolicy()
	if err != nil {
		panic(err)
	}
	err = loadAllUserRolePolicy()
	if err != nil {
		panic(err)
	}
}
