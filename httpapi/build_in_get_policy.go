package httpapi

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type policyOut struct {
	RoleName     string `json:"roleName"`
	Action       string `json:"action"`
	Domain       string `json:"domain"`
	ResourceName string `json:"resourceName"`
	ResourceID   uint   `json:"resourceID"`
}

//获取角色的权限
func getPermissionsForRole(c *gin.Context) {
	dom := c.Query("domain")
	dm, err := getDomain(c, dom)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	plist := Enforcer.GetPolicy()

	var policies []policyOut
	//p类似: ["dom2_admin_0","dom2@/a/data1","POST"]
	for _, p := range plist {
		if strings.HasPrefix(p[1], dm+"@") {
			domAtRsc := strings.Split(p[1], "@")
			action := p[2]
			id, err := getResourceId(action, dm, domAtRsc[1])
			if err != nil {
				respErr(c, 500, "查询资源ID失败(by action,domain,name): "+err.Error())
				return
			}
			policies = append(policies, policyOut{RoleName: p[0], Action: action, Domain: dm, ResourceName: domAtRsc[1], ResourceID: id})
		}
	}

	c.JSON(200, gin.H{"code": 200, "result": policies, "totalCount": len(policies)})
}

func getResourceId(action, domain, resourceName string) (uint, error) {
	var rsc Resource
	err := DB.First(&rsc, `action = ? AND domain = ? AND name = ?`, action, domain, resourceName).Error
	if err != nil {
		return 0, err
	}

	return rsc.ID, nil
}
