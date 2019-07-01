package httpapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type addResourceIn struct {
	Act     string `json:"act" binding:"required"`
	Domain  string `json:"domain"` //Todo: 空表示使用默认域
	Name    string `json:"name" binding:"required"`
	Comment string `json:"comment" binding:"required"`
}

func AddResource(c *gin.Context) {
	var ari addResourceIn
	err := c.BindJSON(&ari)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	if ari.Domain == "" {
		ari.Domain = c.GetString("initDomain")
	}

	var cnt int
	err = DB.Model(&Resource{}).Where(&Resource{Act: ari.Act, Domain: ari.Domain, Name: ari.Name}).Count(&cnt).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("创建资源失败,Count：%v", err))
		return
	}

	if cnt > 0 {
		respErr(c, 500, "资源已经存在(act+domain+name)")
		return
	}

	r := Resource{Act: ari.Act, Domain: ari.Domain, Name: ari.Name}
	err = DB.Model(Resource{}).Create(&r).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("创建资源失败,Create：%v", err))
		return
	}

	respOkEmpty(c)
}

type deleteResourceIn struct {
	ResourceId uint `json:"resourceId" binding:"gt=0"`
}

func DeleteResource(c *gin.Context) {
	var dri deleteResourceIn
	err := c.BindJSON(&dri)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	var roleIdList []uint
	err = DB.Table(`role_resource`).Where(`resource_id = ?`, dri.ResourceId).Find(&roleIdList).Error
	if err == nil {
		for _, roleId := range roleIdList {
			err = deleteResourceListOfRole(roleId, []uint{dri.ResourceId})
			if err != nil {
				respErr(c, 500, fmt.Sprintf("删除角色.资源失败:%v", err))
				return
			}
		}
		deleteResourceByID(c, dri.ResourceId)
		return
	}

	if err != gorm.ErrRecordNotFound {
		deleteResourceByID(c, dri.ResourceId)
		return
	}

	respOkEmpty(c)

}

func deleteResourceByID(c *gin.Context, resourceId uint) {
	err := DB.Delete(&Resource{Model: Model{ID: resourceId}}).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("删除资源失败: ID=%d, %v", resourceId, err))
		return
	}
	respOkEmpty(c)
}
