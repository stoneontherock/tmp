package httpapi

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func loadResourceForRolePolicy(resourceList []Resource, roleName, action string) {
	if action == "add" {
		for _, r := range resourceList {
			logrus.Debugf("增加角色.资源：role:%s, domain:%s rscname:%s, action%s\n", roleName, r.Domain, r.Name, r.Action)
			Enforcer.AddPermissionForUser(roleName, r.Domain+"@"+r.Name, r.Action)
		}
		return
	}

	for _, r := range resourceList {
		logrus.Debugf("删除角色.资源：role:%s, domain:%s rscname:%s, action%s\n", roleName, r.Domain, r.Name, r.Action)
		Enforcer.DeletePermissionForUser(roleName, r.Domain+"@"+r.Name, r.Action)
	}
	return
}

func loadAllRoleRescourcePolicy() error {
	var roleList []Role
	err := DB.Find(&roleList).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		} else {
			return fmt.Errorf("载入全局角色-资源策略失败:roleList; %v", err)
		}
	}

	for _, role := range roleList {
		var resources []Resource
		err := DB.Model(&role).Association("resources").Find(&resources).Error
		if err != nil {
			return fmt.Errorf("载入全局角色-资源策略失败:resources;%v", err)
		}
		for _, rsc := range resources {
			logrus.Debugf("载入全局角色资源策略：role=%s domain=%s rsc=%s action%s\n", role.Name, rsc.Domain, rsc.Name, rsc.Action)
			Enforcer.AddPermissionForUser(role.Name, rsc.Domain+"@"+rsc.Name, rsc.Action)
		}
	}

	return nil
}

func loadRoleForUserPolicy(roleList []Role, userName, action string) {
	if action == "add" {
		for _, role := range roleList {
			Enforcer.AddRoleForUser(userName, role.Name)
		}
		return
	}

	for _, role := range roleList {
		Enforcer.DeleteRoleForUser(userName, role.Name)
	}
}

func loadAllUserRolePolicy() error {
	var userList []User
	err := DB.Find(&userList).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		} else {
			return fmt.Errorf("载入全局角色-用户策略失败:userList; %v", err)
		}
	}
	logrus.Debugf("userList=%+v\n", userList)

	for i := range userList {
		var roles []Role
		err := DB.Model(&userList[i]).Association("roles").Find(&roles).Error
		if err != nil {
			return fmt.Errorf("载入全局角色-资源策略失败:roles;%v", err)
		}
		for _, role := range roles {
			logrus.Debugf("载入全局用户.角色：user=%s role=%s\n", userList[i].Name, role.Name)
			Enforcer.AddRoleForUser(userList[i].Name, role.Name)
		}
	}

	return nil
}
