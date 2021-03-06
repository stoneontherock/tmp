package httpapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type addUserIn struct {
	Name   string `json:"name" binding:"isName"`
	Pstr   string `json:"password" binding:"isPstr"`
	Domain string `json:"domain" binding:"omitempty,isDomain"`
}

func addUser(c *gin.Context) {
	var aui addUserIn
	err := c.ShouldBindJSON(&aui)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	dm, err := getDomain(c, aui.Domain)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	cnt := 0
	//不允许重名
	err = DB.Model(&User{}).Where(&User{Name: aui.Name}).Count(&cnt).Error
	if err != nil {
		respErr(c, 500, err.Error())
		return
	}

	if cnt > 0 {
		respErr(c, 500, fmt.Sprintf("用户已经存在"))
		return
	}

	user := User{
		Name:   aui.Name,
		Pstr:   aui.Pstr,
		Domain: dm,
	}

	user.Salt = getRandomStr(8)
	user.Pstr = md5sum(user.Pstr + user.Salt)
	err = DB.Create(&user).Error
	if err != nil {
		respErr(c, 500, err.Error())
		return
	}

	respOkEmpty(c)
}

type delUserIn struct {
	Name   string `json:"name" binding:"required"`
	Domain string `json:"domain" binding:"omitempty,isDomain"`
}

func delUser(c *gin.Context) {
	var dui delUserIn
	err := c.ShouldBindJSON(&dui)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	dm, err := getDomain(c, dui.Domain)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	var u User
	DB.Where(`name = ? AND domain = ?`, dui.Name, dm).First(&u)
	if u.Name == "" {
		respErr(c, 500, "域内找不到用户:"+dui.Name)
		return
	}

	commit := false
	tx, commitFunc := txCommit(DB, &commit)
	defer commitFunc()

	user := User{Name: dui.Name}
	err = tx.Table(`user_role`).Where("user_name = ?", dui.Name).Delete(&userRole{}).Error
	if err != nil {
		respErr(c, 500, "删除用户失败:user+role:"+err.Error())
		return
	}

	if err := tx.Delete(&user).Error; err != nil {
		respErr(c, 500, "删除用户失败:user:"+err.Error())
		return
	}

	Enforcer.DeleteUser(dui.Name)
	commit = true
	respOkEmpty(c)
}

type luIn struct {
	Name   string `form:"name"`
	Domain string `form:"domain"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
	Order  string `form:"order"`
}

func listUser(c *gin.Context) {
	var li luIn
	err := c.ShouldBindQuery(&li)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	dm, err := getDomain(c, li.Domain)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	total := 0
	query := DB.Model(&User{}).Where(&User{Name: li.Name, Domain: dm})
	err = query.Count(&total).Error
	if err != nil {
		respErr(c, 500, err.Error())
		return
	}

	if total == 0 {
		c.JSON(200, gin.H{"code": 200, "result": nil, "totalCount": total})
		return
	}

	if li.Limit > 0 {
		query = query.Limit(li.Limit)
	}

	if li.Order != "" {
		query = query.Order("name " + li.Order)
	}

	var userList []User
	err = query.Offset(li.Offset).Find(&userList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		respErr(c, 500, err.Error())
		return
	}

	us := make([]string, len(userList))
	for i := range userList {
		us[i] = userList[i].Name
	}
	c.JSON(200, gin.H{"code": 200, "result": us, "totalCount": total})
	return
}

type chgPwdIn struct {
	UserName  string `json:"userName"`
	OldPasswd string `json:"oldPasswd"`
	NewPasswd string `json:"newPasswd" binding:"isPstr"`
}

func changePassword(c *gin.Context) {
	var cpi chgPwdIn
	err := c.ShouldBindJSON(&cpi)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	var userName string
	if (strings.HasPrefix(ctxUser(c), "admin@") || ctxUser(c) == SA) && cpi.UserName != "" {
		userName = cpi.UserName
	} else {
		userName = ctxUser(c)
	}

	u := User{Name: userName}
	err = DB.First(&u).Error
	if err != nil {
		respErr(c, 500, err.Error())
		return
	}

	if ctxUser(c) != SA {
		if u.Domain != ctxDomain(c) {
			respErr(c, 400, "不可以修改域外用户的密码")
			return
		}
	}

	//普通用户需要校验密码
	if ctxUser(c) != SA && !strings.HasPrefix(ctxUser(c), "admin@") {
		if u.Pstr != md5sum(cpi.OldPasswd+u.Salt) {
			respErr(c, 400, "原密码错误")
			return
		}
	}

	err = DB.Model(&User{Name: userName}).Update("pstr", md5sum(cpi.NewPasswd+u.Salt)).Error
	if err != nil {
		respErr(c, 500, err.Error())
		return
	}

	respOkEmpty(c)
}
