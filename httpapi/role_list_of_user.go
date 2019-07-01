package httpapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type mUR struct {
	UserId     uint   `json:"userId" binding:"gt=0"`
	RoleIdList []uint `json:"roleIdList" binding:"gt=0,dive,gt=0"`
}

func AddRoleListOfUser(c *gin.Context) {
	var mr mUR
	err := c.BindJSON(&mr)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	err = addRoleListOfUser(mr.UserId, mr.RoleIdList)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	respOkEmpty(c)
}

func DeleteRoleListOfUser(c *gin.Context) {
	var mr mUR
	err := c.BindJSON(&mr)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	err = deleteRoleListOfUser(mr.UserId, mr.RoleIdList)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	respOkEmpty(c)
}

func deleteRoleListOfUser(userId uint, roleIdList []uint) error {
	rolesDelete := make([]Role, len(roleIdList))
	for i := range roleIdList {
		rolesDelete[i].ID = roleIdList[i]
	}

	user := User{Model: Model{ID: userId}}
	err := DB.Model(&user).Association("roles").Delete(&rolesDelete).Error
	if err != nil {
		return fmt.Errorf("删除'用户.角色'失败：%v", err)
	}
	fmt.Printf("删除用户.角色：user:%+v, roles:%+v\n", user, rolesDelete)
	loadUserOfRolePolicy(rolesDelete, user.Name, "del")
	return nil
}

func addRoleListOfUser(userId uint, roleIdList []uint) error {
	var rolesAppend []Role
	user := User{Model: Model{ID: userId}}
	err := DB.Model(&user).Where(`id in (?)`, roleIdList).Find(&rolesAppend).Error
	if err != nil {
		return fmt.Errorf("检查`用户.角色`：%v", err)
	}

	if len(rolesAppend) == 0 {
		return fmt.Errorf("添加`用户.角色`失败，要添加的角色不存在：%v", err)
	}

	//只添加数据库有的资源
	err = DB.Model(&user).Association("roles").Append(&rolesAppend).Error
	if err != nil {
		return fmt.Errorf("添加`角色.资源`失败：%v", err)
	}
	fmt.Printf("增加用户.角色：user:%+v, roles:%+v\n", user, rolesAppend)
	loadUserOfRolePolicy(rolesAppend, user.Name, "add")
	return nil
}
