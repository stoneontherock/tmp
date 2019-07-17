package httpapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type getUsersForRoleIn struct {
	RoleName string `form:"roleName" binding:"required"`
}

//获取角色关联的所有用户
func getUsersForRole(c *gin.Context) {
	var rri getUsersForRoleIn
	err := c.ShouldBindQuery(&rri)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	role := Role{Name:rri.RoleName}
	err = DB.Find(&role).Error
	if err != nil {
		respErr(c, 500, "查找用户失败:"+err.Error())
		return
	}

	if ctxUser(c) != SA && role.Domain != ctxDomain(c) {
		respErr(c, 400, "域管理员只能查域内角色")
		return
	}

	userList, err := Enforcer.GetUsersForRole(rri.RoleName)
	if err != nil {
		respErr(c, 500, fmt.Sprintf("获取角色(%s)的用户失败:%v", rri.RoleName, err))
		return
	}

	c.JSON(200, gin.H{"code": 200, "result": userList, "totalCount": len(userList)})
}
