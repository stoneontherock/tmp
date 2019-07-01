package httpapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type addRoleIn struct {
	Name string `json:"name" binding:"required"`
}

func AddRole(c *gin.Context) {
	var ari addRoleIn
	err := c.BindJSON(&ari)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	var cnt int
	err = DB.Model(&Role{}).Where(`name = ?`, ari.Name).Count(&cnt).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("创建角色失败,Count：%v", err))
		return
	}

	if cnt > 0 {
		respErr(c, 400, "角色已经存在")
		return
	}

	r := Role{Name: ari.Name}
	err = DB.Create(&r).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("创建角色失败,Create：%v", err))
		return
	}

	respOkEmpty(c)
}


