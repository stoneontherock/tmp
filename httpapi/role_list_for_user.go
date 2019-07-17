package httpapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
)

type rolesUserIn struct {
	UserName     string   `json:"userName" binding:"gt=0"`
	RoleNameList []string `json:"roleNameList" binding:"gt=0,dive,gt=0"`
	Domain       string   `json:"domain" binding:"omitempty,isDomain"`
}

func addOrDelRoleListForUserFunc(action string) func(c *gin.Context) {
	return func(c *gin.Context) {
		var rui rolesUserIn
		err := c.ShouldBindJSON(&rui)
		if err != nil {
			respErr(c, 400, err.Error())
			return
		}

		if strings.HasPrefix(rui.UserName, "admin@") {
			respErr(c, 400, "域管理员只负责域内权限管理,不参与具体其他具体业务,无须指派角色")
			return
		}

		dm, err := getDomain(c, rui.Domain)
		if err != nil {
			respErr(c, 400, err.Error())
			return
		}

		u := User{Name: rui.UserName}
		DB.Find(&u)
		if u.Domain != dm {
			respErr(c, 400, fmt.Sprintf("%s不在域(%s)内", u.Name, dm))
			return
		}

		err = addOrDelRoleListForUser(dm, rui.UserName, action, rui.RoleNameList)
		if err != nil {
			respErr(c, 400, err.Error())
			return
		}

		respOkEmpty(c)
	}
}

func addOrDelRoleListForUser(domain, userName, action string, roleNameList []string) error {
	for _, roleName := range roleNameList {
		cnt := 0
		DB.Table(`role`).Where(`name = ? AND domain = ?`, roleName, domain).Count(&cnt)
		if cnt == 0 {
			return fmt.Errorf("角色(%s)不在域内", roleName)
		}
	}

	roleList := make([]Role, len(roleNameList))
	for i := range roleNameList {
		roleList[i].Name = roleNameList[i]
	}

	user := User{Name: userName}
	actFunc := DB.Model(&user).Association("roles").Append
	if action == "del" {
		actFunc = DB.Model(&user).Association("roles").Delete
	}
	err := actFunc(&roleList).Error
	if err != nil {
		return fmt.Errorf("%s `角色.资源`失败：%v", action, err)
	}

	logrus.Debugf("addOrDelRoleListForUser: 用户.角色：user:%+v, roles:%+v\n", user, roleList)
	loadRoleForUserPolicy(roleList, user.Name, action)
	return nil
}
