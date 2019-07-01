package httpapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type addUserIn struct {
	Name string `json:"name" binding:"isName"`
	Pstr string `json:"pstr" binding:"isPstr"`
}

func AddUser(c *gin.Context) {
	var aui addUserIn
	err := c.BindJSON(&aui)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}


	defaultDomain := c.GetString("initDomain")

	cnt := 0
	//不允许重名
	err = DB.Model(&User{}).Where(&User{Name: aui.Name/*, DefaultDomain: defaultDomain*/}).Count(&cnt).Error
	if err != nil {
		respErr(c, 500, err.Error())
		return
	}

	if cnt > 0 {
		respErr(c, 500, fmt.Sprintf("用户已经存在(name+defaultdomain)"))
		return
	}

	user := User{
		Name:          aui.Name,
		Pstr:          aui.Pstr,
		DefaultDomain: defaultDomain,
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

type luIn struct {
	Name   string `form:"name"`
	Offset int    `form:"offset" binding:"omitempty,gte=0"`
	Limit  int    `form:"limit" binding:"omitempty,gt=0"`
	Sort   string `form:"sort" binding:"omitempty,gt=0"`
	Order  string `form:"order"`
}

func ListUser(c *gin.Context) {
	var li luIn
	err := c.ShouldBindQuery(&li)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	defaultDomain := c.GetString("initDomain")

	total := 0
	query := DB.Model(&User{}).Where(&User{Name: li.Name, DefaultDomain: defaultDomain})
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

	var userList []User
	err = query.Offset(li.Offset).Find(&userList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		respErr(c, 500, err.Error())
		return
	}

	c.JSON(200, gin.H{"code": 200, "result": userList, "totalCount": total})
}
