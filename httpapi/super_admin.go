package httpapi

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type adaIn struct {
	Domain string `json:"domain" binding:"omitempty,isDomain"` //数字或字母或点号
	Pstr   string `json:"password" binding:"required"`
}

func addDomain(c *gin.Context) {
	var ai adaIn
	err := c.ShouldBindJSON(&ai)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	salt := getRandomStr(8)
	var da = User{
		Name:   "admin@" + ai.Domain,
		Pstr:   md5sum(ai.Pstr + salt),
		Salt:   salt,
		Domain: ai.Domain,
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
	err := c.ShouldBindJSON(&ddi)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	commit := false
	tx, cf := txCommit(DB, &commit)
	defer cf()

	//查角色列表
	var roleList []Role
	err = tx.Where(`domain = ?`, ddi.Domain).Find(&roleList).Error
	if err != nil {
		respErr(c, 500, "查找角色失败:"+err.Error())
		return
	}
	roleNameList := make([]string, len(roleList))
	for i := range roleNameList {
		roleNameList[i] = roleList[i].Name
	}

	//查用戶列表
	var userList []User
	err = tx.Where(`domain = ?`, ddi.Domain).Find(&userList).Error
	if err != nil {
		respErr(c, 500, "查找用户失败:"+err.Error())
		return
	}
	userNameList := make([]string, len(userList))
	for i := range userList {
		userNameList[i] = userList[i].Name
	}

	//查資源列表
	var resourceList []Resource
	err = tx.Where(`domain = ?`, ddi.Domain).Find(&resourceList).Error
	if err != nil {
		respErr(c, 500, "查找資源失败:"+err.Error())
		return
	}

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

	//删 user_role (by user_name)
	if len(userList) > 0 {
		err = tx.Delete(&userRole{}, `user_name IN (?)`, userNameList).Error
		if err != nil {
			respErr(c, 500, "查找用户.角色(by userName)失败:"+err.Error())
			return
		}
	}

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

	//删 role_resource (by roleName)
	if len(userList) > 0 {
		err = tx.Delete(&roleResource{}, `role_name IN (?)`, roleNameList).Error
		if err != nil {
			respErr(c, 500, "查找用户.角色(by userName)失败:"+err.Error())
			return
		}
	}

	err = tx.Delete(&User{}, `domain = ?`, ddi.Domain).Error
	if err != nil {
		respErr(c, 500, "删除用户失败:"+err.Error())
		return
	}

	err = tx.Delete(&Role{}, `domain = ?`, ddi.Domain).Error
	if err != nil {
		respErr(c, 500, "删除角色失败:"+err.Error())
		return
	}

	err = tx.Delete(&Resource{}, `domain = ?`, ddi.Domain).Error
	if err != nil {
		respErr(c, 500, "删除资源失败:"+err.Error())
		return
	}

	commit = true
	respOkEmpty(c)
}

func listDomain(c *gin.Context) {
	var userList []User
	err := DB.Find(&userList, `name LIKE ?`, "admin@%").Error
	if err != nil && err != gorm.ErrRecordNotFound {
		respErr(c, 500, "查询失败:"+err.Error())
		return
	}

	result := make([]string, len(userList))
	for i := range result {
		result[i] = userList[i].Domain
	}

	c.JSON(200, gin.H{"code": 200, "result": result, "totalCount": len(result)})
}
