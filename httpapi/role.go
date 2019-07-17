package httpapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type roleIn struct {
	Name   string `json:"name" binding:"required"`
	Domain string `json:"domain" binding:"omitempty,isDomain"`
}

func addRole(c *gin.Context) {
	var ri roleIn
	err := c.ShouldBindJSON(&ri)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	dm, err := getDomain(c, ri.Domain)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	var cnt int
	err = DB.Model(&Role{}).Where(`name = ?`, ri.Name).Count(&cnt).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("创建角色失败,Count：%v", err))
		return
	}

	if cnt > 0 {
		respErr(c, 400, "角色已经存在")
		return
	}

	r := Role{Name: ri.Name, Domain: dm}
	err = DB.Create(&r).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("创建角色失败,Create：%v", err))
		return
	}

	respOkEmpty(c)
}

func delRole(c *gin.Context) {
	var ri roleIn
	err := c.ShouldBindJSON(&ri)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	dm, err := getDomain(c, ri.Domain)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	var role Role
	DB.Where(`name = ? AND domain = ?`, ri.Name, dm).First(&role)
	if role.Name == "" {
		respErr(c, 500, "域内找不到该角色")
		return
	}

	commit := false
	tx, cf := txCommit(DB, &commit)
	defer cf()

	err = tx.Delete(&Role{Name: ri.Name}).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("删除角色失败：%v", err))
		return
	}

	//删除该角色和用户的绑定关系
	err = tx.Table(`user_role`).Where(`role_name = ?`, ri.Name).Delete(userRole{}).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("删除用户角色绑定关系失败：%v", err))
		return
	}

	//删除该角色和资源的绑定关系
	err = tx.Table(`role_resource`).Where(`role_name = ?`, ri.Name).Delete(roleResource{}).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("删除角色资源绑定关系失败：%v", err))
		return
	}

	Enforcer.DeleteRole(ri.Name)

	respOkEmpty(c)
}

type listRoleIn struct {
	Domain string `form:"domain" binding:"omitempty,isDomain"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
	Order  string `form:"order"`
}

func listRole(c *gin.Context) {
	var lri listRoleIn
	err := c.ShouldBindQuery(&lri)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	dm, err := getDomain(c, lri.Domain)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	query := DB.Model(&Role{}).Where(&Role{Domain: dm})

	total := 0
	err = query.Count(&total).Error
	if err != nil {
		respErr(c, 500, err.Error())
		return
	}

	if total == 0 {
		c.JSON(200, gin.H{"code": 200, "result": nil, "totalCount": total})
		return
	}

	if lri.Limit > 0 {
		query = query.Limit(lri.Limit)
	}

	if lri.Order != "" {
		query = query.Order("name " + lri.Order)
	}

	var roleList []Role
	err = query.Offset(lri.Offset).Find(&roleList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		respErr(c, 500, err.Error())
		return
	}

	rs := make([]string, len(roleList))
	for i := range roleList {
		rs[i] = roleList[i].Name
	}

	c.JSON(200, gin.H{"code": 200, "result": rs, "totalCount": total})
	return
}
