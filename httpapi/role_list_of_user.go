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

func addOrDelRoleListOfUserFunc(act string) func(c *gin.Context) {
	return func(c *gin.Context) {
		var rui rolesUserIn
		err := c.BindJSON(&rui)
		if err != nil {
			respErr(c, 400, err.Error())
			return
		}

		dm, err := getDomain(c, rui.Domain)
		if err != nil {
			respErr(c, 400, err.Error())
			return
		}

		u := User{Name: rui.UserName}
		DB.Find(&u)
		if u.DefaultDomain != dm {
			respErr(c, 400, fmt.Sprintf("%s不在域(%s)内", u.Name, dm))
			return
		}

		var jdom string
		if ctxUser(c) == SA {
			logrus.Debugf("rui.Domain=%s\n", rui.Domain)
			jdom, err = getJoinedDomainByInitialDomain(rui.Domain)
			if err != nil {
				respErr(c, 400, err.Error())
				return
			}
		}
		domList := strings.Fields(strings.Replace(dm+","+jdom, ",", " ", -1))
		//logrus.Debugf("domList=  %s\n", domList)

		err = addOrDelRoleListOfUser(act, rui.UserName, rui.RoleNameList, domList)
		if err != nil {
			respErr(c, 400, err.Error())
			return
		}

		respOkEmpty(c)
	}
}

//func deleteRoleListOfUser(c *gin.Context) {
//	var rui rolesUserIn
//	err := c.BindJSON(&rui)
//	if err != nil {
//		respErr(c, 400, err.Error())
//		return
//	}
//
//	err = deleteRoleListOfUser(rui.UserName, rui.RoleNameList)
//	if err != nil {
//		respErr(c, 400, err.Error())
//		return
//	}
//
//	rolesDelete := make([]Role, len(roleNameList))
//	for i := range roleNameList {
//		rolesDelete[i].Name = roleNameList[i]
//	}
//
//	user := User{Name: userName}
//	err := DB.Model(&user).Association("roles").Delete(&rolesDelete).Error
//	if err != nil {
//		return fmt.Errorf("删除'用户.角色'失败：%v", err)
//	}
//	fmt.Printf("删除用户.角色：user:%+v, roles:%+v\n", user, rolesDelete)
//	loadRoleOfUserPolicy(rolesDelete, user.Name, "del")
//	return nil
//
//
//	respOkEmpty(c)
//}

//func deleteRoleListOfUser(userName string, roleNameList []string) error {
//	rolesDelete := make([]Role, len(roleNameList))
//	for i := range roleNameList {
//		rolesDelete[i].Name = roleNameList[i]
//	}
//
//	user := User{Name: userName}
//	err := DB.Model(&user).Association("roles").Delete(&rolesDelete).Error
//	if err != nil {
//		return fmt.Errorf("删除'用户.角色'失败：%v", err)
//	}
//	fmt.Printf("删除用户.角色：user:%+v, roles:%+v\n", user, rolesDelete)
//	loadRoleOfUserPolicy(rolesDelete, user.Name, "del")
//	return nil
//}

func addOrDelRoleListOfUser(action, userName string, roleNameList, domList []string) error {
	for _, roleName := range roleNameList {
		cnt := 0
		DB.Table(`role`).Where(`name = ? AND default_domain IN (?)`, roleName, domList).Count(&cnt)
		if cnt == 0 {
			return fmt.Errorf("角色(%s)不属于初始域或扩展域", roleName)
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

	logrus.Debugf("addOrDelRoleListOfUser: 用户.角色：user:%+v, roles:%+v\n", user, roleList)
	loadRoleOfUserPolicy(roleList, user.Name, action)
	return nil
}
