package httpapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"regexp"
	"strings"
)

type adaIn struct {
	InitialDomain string   `json:"initialDomain" binding:"omitempty,isDomain"` //数字或字母或点号
	JoinedDomain  []string `json:"joinedDomain" binding:"omitempty,gt=0,dive,isDomain"`
	Pstr          string   `json:"pstr" binding:"required"`
}

func addDomain(c *gin.Context) {
	var ai adaIn
	err := c.BindJSON(&ai)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	//验证domain存在
	//for _, jd := range ai.JoinedDomain {
	//	var cnt int
	//	DB.Model(&DomainAdmin{}).Where(`joined_domain = ?`,jd).Count(&cnt)
	//	if cnt == 0 {
	//		respErr(c, 500, err.Error())
	//		return
	//	}
	//}

	salt := getRandomStr(8)
	var da = DomainAdmin{
		Name:          "admin@" + ai.InitialDomain,
		Pstr:          md5sum(ai.Pstr + salt),
		Salt:          salt,
		InitialDomain: ai.InitialDomain,
		JoinedDomain:  strings.Join(ai.JoinedDomain, ","),
	}

	err = DB.Create(&da).Error
	if err != nil {
		respErr(c, 500, "新增域管理者失败:"+err.Error())
		return
	}

	respOkEmpty(c)
}

type delDomainIn struct {
	Domain string `json:"domain" binding:"isDomain"`
}

func delDomain(c *gin.Context) {
	var ddi delDomainIn
	err := c.BindJSON(&ddi)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	commit := false
	tx, cf := txCommit(DB, &commit)
	defer cf()

	var roleList []Role
	err = tx.Where(`default_domain = ?`, ddi.Domain).Find(&roleList).Error
	if err != nil {
		respErr(c, 500, "查找角色失败:"+err.Error())
		return
	}

	roleNameList := make([]string, len(roleList))
	for i := range roleNameList {
		roleNameList[i] = roleList[i].Name
	}

	var userList []User
	err = tx.Where(`default_domain = ?`, ddi.Domain).Find(&userList).Error
	if err != nil {
		respErr(c, 500, "查找用户失败:"+err.Error())
		return
	}

	userNameList := make([]string, len(userList))
	for i := range userList {
		userNameList[i] = userList[i].Name
	}

	var resourceList []Resource
	err = tx.Where(`domain = ?`, ddi.Domain).Find(&resourceList).Error
	if err != nil {
		respErr(c, 500, "查找用户失败:"+err.Error())
		return
	}

	log.Printf("*********** 0")
	//删 user_role (by role_name)
	if len(roleList) > 0 {
		for i := range userList {
			asso := tx.Model(&userList[i]).Association("roles").Delete(roleList)
			if asso.Error != nil {
				respErr(c, 500, "删除用户.角色(by roleName)失败:"+asso.Error.Error())
				return
			}
		}
	}
	log.Printf("*********** 1")

	//删 user_role (by user_name)
	if len(userList) > 0 {
		err = tx.Delete(&userRole{}, `user_name IN (?)`, userNameList).Error
		if err != nil {
			respErr(c, 500, "查找用户.角色(by userName)失败:"+err.Error())
			return
		}
	}

	log.Printf("*********** 2")

	//删 role_resource (by resourceID)
	if len(resourceList) > 0 {
		for i := range roleList {
			asso := tx.Model(&roleList[i]).Association("resources").Delete(resourceList)
			if asso.Error != nil {
				respErr(c, 500, "删除角色.资源(by resourceID)失败:"+asso.Error.Error())
				return
			}
		}
	}

	log.Printf("*********** 3")
	//删 role_resource (by roleName)
	if len(userList) > 0 {
		err = tx.Delete(&roleResource{}, `role_name IN (?)`, roleNameList).Error
		if err != nil {
			respErr(c, 500, "查找用户.角色(by userName)失败:"+err.Error())
			return
		}
	}

	log.Printf("*********** 4")
	err = tx.Delete(&User{}, `default_domain = ?`, ddi.Domain).Error
	if err != nil {
		respErr(c, 500, "删除用户失败:"+err.Error())
		return
	}

	err = tx.Delete(&Role{}, `default_domain = ?`, ddi.Domain).Error
	if err != nil {
		respErr(c, 500, "删除角色失败:"+err.Error())
		return
	}

	err = tx.Delete(&Resource{}, `domain = ?`, ddi.Domain).Error
	if err != nil {
		respErr(c, 500, "删除资源失败:"+err.Error())
		return
	}

	err = tx.Delete(&DomainAdmin{}, `initial_domain = ?`, ddi.Domain).Error
	if err != nil {
		respErr(c, 500, "删除域管理记录失败:"+err.Error())
		return
	}

	var daList []DomainAdmin
	err = tx.Where(`joined_domain = ?`, ddi.Domain).Find(&daList).Error
	if err != nil {
		respErr(c, 500, "查询关联扩展域记录失败:"+err.Error())
		return
	}

	for i := range daList {
		delRegex := regexp.MustCompile(fmt.Sprintf(`(%s,)|(,%[1]s$)`, regexp.QuoteMeta(ddi.Domain)))
		daList[i].JoinedDomain = delRegex.ReplaceAllString(daList[i].JoinedDomain, "")
		err = tx.Save(&daList[i]).Error
		if err != nil {
			respErr(c, 500, "更新记录的扩展域失败:"+err.Error())
			return
		}
	}

	commit = true
	respOkEmpty(c)
}

func listDomain(c *gin.Context) {
	var domList []struct {
		InitialDomain string
	}

	err := DB.Table(`domain_admin`).Find(&domList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		respErr(c, 500, "查询域失败:"+err.Error())
		return
	}

	result := make([]string, len(domList))
	for i := range result {
		result[i] = domList[i].InitialDomain
	}

	c.JSON(200, gin.H{"code": 200, "result": result, "totalCount": len(result)})
}
