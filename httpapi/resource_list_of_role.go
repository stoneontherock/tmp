package httpapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type mRR struct {
	RoleID         uint   `json:"roleID" binding:"gt=0"`
	ResourceIdList []uint `json:"resourceIdList" binding:"gt=0,dive,gt=0"`
}

func AddResourceListOfRole(c *gin.Context) {
	var mr mRR
	err := c.BindJSON(&mr)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	err = addResourceListOfRole(mr.RoleID, mr.ResourceIdList)
	if err != nil {
		respErr(c, 500, err.Error())
		return
	}

	respOkEmpty(c)
}

func addResourceListOfRole(roleId uint, resourceIdList []uint) error {
	var rscAppend []Resource
	role := Role{Model: Model{ID: roleId}}
	err := DB.Model(&role).Where(`id in (?)`, resourceIdList).Find(&rscAppend).Error
	if err != nil {
		return fmt.Errorf("检查`角色.资源`：%v", err)
	}

	if len(rscAppend) == 0 {
		return fmt.Errorf("添加`角色.资源`失败，要添加的资源不存在：%v", err)
	}

	//只添加数据库有的资源
	err = DB.Model(&role).Association("resources").Append(&rscAppend).Error
	if err != nil {
		return fmt.Errorf("添加`角色.资源`失败：%v", err)
	}
	fmt.Printf("增加角色.资源：role:%+v, rscs:%+v\n", role, rscAppend)
	loadResourceOfRolePolicy(rscAppend, role.Name, "add")
	return nil
}

func DeleteResourceListOfRole(c *gin.Context) {
	var mr mRR
	err := c.BindJSON(&mr)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	err = deleteResourceListOfRole(mr.RoleID, mr.ResourceIdList)
	if err != nil {
		respErr(c, 500, err.Error())
		return
	}

	respOkEmpty(c)
}

func deleteResourceListOfRole(roleId uint, resourceIdList []uint) error {
	role := Role{Model: Model{ID: roleId}}
	rscDelete := make([]Resource, len(resourceIdList))
	for i := range resourceIdList {
		rscDelete[i].ID = resourceIdList[i]
	}
	err := DB.Model(&role).Association("resources").Delete(&rscDelete).Error
	if err != nil {
		return fmt.Errorf("删除`角色.资源`失败：%v", err)
	}
	fmt.Printf("删除角色.资源：role:%+v, rscs:%+v\n", role, rscDelete)
	loadResourceOfRolePolicy(rscDelete, role.Name, "del")
	return nil
}
