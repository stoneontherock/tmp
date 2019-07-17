package httpapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
)

type mRR struct {
	RoleName       string `json:"roleName" binding:"gt=0"`
	ResourceIdList []uint `json:"resourceIdList" binding:"gt=0,dive,gt=0"`
	Domain string `json:"domain"`
}

func addResourceListForRole(c *gin.Context) {
	var mr mRR
	err := c.ShouldBindJSON(&mr)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	dm, err := getDomain(c, ctxDomain(c))
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	code, err := verifyDomainOfRole(dm, mr.RoleName)
	if err != nil {
		respErr(c, code, err.Error())
		return
	}

	var rscAppend []Resource
	for _, rscID := range mr.ResourceIdList {
		rsc := Resource{ID: rscID}
		err := DB.First(&rsc).Error
		if err != nil {
			respErr(c, 500, err.Error())
			return
		}

		if rsc.Domain == "" {
			respErr(c, 400, fmt.Sprintf("资源id=%d不存在", rsc.ID))
			return
		}

		if strings.HasPrefix(ctxUser(c), "admin@") {
			if rsc.Domain != dm {
				respErr(c, 400, fmt.Sprintf("资源id=%d不在域内,要添加域外资源需要超管权限", rsc.ID))
				return
			}
		}

		rscAppend = append(rscAppend, rsc)
	}

	//只添加数据库有的资源
	role := Role{Name: mr.RoleName}
	err = DB.Model(&role).Association("resources").Append(&rscAppend).Error
	if err != nil {
		respErr(c, 500, err.Error())
		return
	}

	logrus.Debugf("增加角色.资源：role:%+v, rscs:%+v\n", role, rscAppend)
	loadResourceForRolePolicy(rscAppend, role.Name, "add")
	respOkEmpty(c)
}

func verifyDomainOfRole(domain, roleName string) (int, error) {
	var cnt int
	err := DB.Model(&Role{Name: roleName}).Where(`domain = ?`, domain).Count(&cnt).Error
	if err != nil {
		return 500, fmt.Errorf("查找`域-角色名`失败：%v", err)
	}

	if cnt == 0 {
		return 400, fmt.Errorf("角色不存在或角色不在域内：%v", err)
	}

	return 0, nil
}

func deleteResourceListForRole(c *gin.Context) {
	var mr mRR
	err := c.ShouldBindJSON(&mr)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	dm, err := getDomain(c, ctxDomain(c))
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	code, err := verifyDomainOfRole(dm, mr.RoleName)
	if err != nil {
		respErr(c, code, err.Error())
		return
	}

	err = deleteResourceListForRoleDo(mr.RoleName, mr.ResourceIdList)
	if err != nil {
		respErr(c, 500, err.Error())
		return
	}

	respOkEmpty(c)
}

func deleteResourceListForRoleDo(roleName string, resourceIdList []uint) error {
	role := Role{Name: roleName}
	rscDelete := make([]Resource, len(resourceIdList))
	for i := range resourceIdList {
		logrus.Debugf("rsc id=%d\n", resourceIdList[i])
		rscDelete[i].ID = resourceIdList[i]
	}

	err := DB.Model(&role).Association("resources").Delete(&rscDelete).Error
	if err != nil {
		return fmt.Errorf("删除`角色.资源`失败：%v", err)
	}

	logrus.Debugf("删除角色.资源：role:%+v, rscs:%+v\n", role, rscDelete)
	loadResourceForRolePolicy(rscDelete, role.Name, "del")
	return nil
}
