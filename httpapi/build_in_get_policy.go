package httpapi

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strings"
)

type PermissonOut struct {
	RoleName     string `json:"roleName"`
	Action       string `json:"action"`
	Domain       string `json:"domain"`
	ResourceName string `json:"resourceName"`
	ResourceID   uint   `json:"resourceID"`
}

type permissonIn struct {
	RoleName string `json:"roleName"`
	Domain   string `json:"domain"`
}

//获取角色的权限
func getPermissionsForRole(c *gin.Context) {
	var pi permissonIn
	err := c.ShouldBindQuery(&pi)
	dm, err := getDomain(c, pi.Domain)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	var roleList []Role
	err = DB.Model(&Role{}).Where(&Role{Name: pi.RoleName, Domain: dm}).Find(&roleList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		respErr(c, 500, "查询域中的角色列表失败:"+err.Error())
		return
	}

	plist := Enforcer.GetPolicy()

	var permList []PermissonOut
	//p类似: ["dom2_admin_0","dom2@/a/data1","POST"]
	for _, p := range plist {
		for i := range roleList {
			if roleList[i].Name != p[0] {
				continue
			}
			domAtRsc := strings.Split(p[1], "@")
			action := p[2]
			id, err := getResourceId(action, dm, domAtRsc[1])
			if err != nil {
				respErr(c, 500, "查询资源ID失败(by action,domain,name): "+err.Error())
				return
			}
			permList = append(permList, PermissonOut{RoleName: p[0], Action: action, Domain: dm, ResourceName: domAtRsc[1], ResourceID: id})
		}
	}

	//批量查
	if pi.RoleName == "" {
		c.JSON(200, gin.H{"code": 200, "result": permList, "totalCount": len(permList)})
		return
	}

	//单查
	var po PermissonOut
	if len(permList) > 0 {
		po = permList[0]
	}

	//单查(空返回)
	c.JSON(200, struct {
		Code int `json:"code"`
		PermissonOut
	}{200, po})
}

func getResourceId(action, domain, resourceName string) (uint, error) {
	var rsc Resource
	err := DB.First(&rsc, `action = ? AND domain = ? AND name = ?`, action, domain, resourceName).Error
	if err != nil {
		return 0, err
	}

	return rsc.ID, nil
}
