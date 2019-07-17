package httpapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type getRolesForUserIn struct {
	UserName string `form:"userName" binding:"required"`
}

//获取用户的所有角色
func getRolesForUser(c *gin.Context) {
	var rri getRolesForUserIn
	err := c.ShouldBindQuery(&rri)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	user := User{Name:rri.UserName}
	err = DB.Find(&user).Error
	if err != nil {
		respErr(c, 500, "查找用户失败:"+err.Error())
		return
	}

	if ctxUser(c) != SA && user.Domain != ctxDomain(c) {
		respErr(c, 400, "域管理员只能查域内用户")
		return
	}

	roleList, err := Enforcer.GetRolesForUser(rri.UserName)
	if err != nil {
		respErr(c, 500, fmt.Sprintf("获取用户(%s)的角色失败:%v", rri.UserName, err))
		return
	}

	c.JSON(200, gin.H{"code": 200, "result": roleList, "totalCount": len(roleList)})
}
