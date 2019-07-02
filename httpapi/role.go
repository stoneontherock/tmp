package httpapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strings"
)

type roleIn struct {
	Name   string `json:"name" binding:"required"`
	Domain string `json:"domain"`
}

func AddRole(c *gin.Context) {
	var ri roleIn
	err := c.BindJSON(&ri)
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

	r := Role{Name: ri.Name, DefaultDomain: dm}
	err = DB.Create(&r).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("创建角色失败,Create：%v", err))
		return
	}

	respOkEmpty(c)
}

func DelRole(c *gin.Context) {
	var ri roleIn
	err := c.BindJSON(&ri)
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
	DB.Where(`name = ? AND default_domain = ?`, ri.Name, dm).First(&role)
	if role.Name == "" {
		respErr(c, 500, "角色找不到")
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

	err = tx.Table(`user_role`).Where(`role_name = ?`, ri.Name).Delete(userRole{}).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("删除用户-角色失败：%v", err))
		return
	}

	Enforcer.DeleteRole(ri.Name)

	respOkEmpty(c)
}

type listRoleIn struct {
	Domain string `json:"domain"`
}

func ListRole(c *gin.Context) {
	var lri listRoleIn
	dm, err := getDomain(c, lri.Domain)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	var jdom string
	if ctxUser(c) == SA {
		jdom,err =getJoinedDomainByInitialDomain(dm)
		if err != nil {
			respErr(c, 400, err.Error())
			return
		}
	}
	domList := strings.Fields(strings.Replace(dm+","+jdom, ",", " ", -1))

	roleList, err := listRole(domList)
	if err != nil {
		respErr(c, 500, err.Error())
		return
	}

	c.JSON(200, gin.H{"code": 200, "roleList": roleList, "total": len(roleList)})
}

func listRole(domList []string) ([]Role, error) {
	var roleList []Role
	err := DB.Where(`default_domain IN (?)`, domList).Find(&roleList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return roleList, nil
}

