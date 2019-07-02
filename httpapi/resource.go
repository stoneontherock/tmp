package httpapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type addResourceIn struct {
	Act     string `json:"act" binding:"required"`
	Domain  string `json:"domain"` //空则使用默认域
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

	dm, err := getDomain(c, ari.Domain)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	//if ari.Domain != "" {
	//	dm = ari.Domain
	//}
	//
	//var da DomainAdmin
	//if ctxUser(c) != SA {
	//	DB.Where(`name = ?`, "admin@"+ctxIDom(c)).Find(&da)
	//	var ok bool
	//	for _, d := range strings.Split(da.JoinedDomain, ",") {
	//		if dm == d {
	//			ok = true
	//			break
	//		}
	//	}
	//
	//	if !(ok || dm == ctxIDom(c)) {
	//		respErr(c, 400, "域管理员只能添加初始域或扩展域的资源")
	//		return
	//	}
	//}

	var cnt int
	err = DB.Model(&Resource{}).Where(&Resource{Act: ari.Act, Domain: dm, Name: ari.Name}).Count(&cnt).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("创建资源失败,Count：%v", err))
		return
	}

	if cnt > 0 {
		respErr(c, 500, "资源已经存在(act+domain+name)")
		return
	}

	r := Resource{Act: ari.Act, Domain: dm, Name: ari.Name}
	err = DB.Model(Resource{}).Create(&r).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("创建资源失败,Create：%v", err))
		return
	}

	respOkEmpty(c)
}

type deleteResourceIn struct {
	ResourceId uint `json:"resourceID" binding:"gt=0"`
}

type roleResource struct {
	RoleName string
}

func DeleteResource(c *gin.Context) {
	var dri deleteResourceIn
	err := c.BindJSON(&dri)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	var resource Resource
	if ctxUser(c) == SA {
		DB.First(&resource, "id = ?", dri.ResourceId)
	} else {
		DB.First(&resource, "id = ? AND domain = ?", dri.ResourceId, ctxIDom(c))
	}

	if resource.Name == "" {
		respErr(c, 500, "资源不存在或不属于域内资源")
		return
	}

	var rrList []roleResource
	err = DB.Table(`role_resource`).Select(`role_name`).Where(`resource_id = ?`, dri.ResourceId).Find(&rrList).Error
	if err == nil {
		for _, rr := range rrList {
			err = deleteResourceListOfRole(rr.RoleName, []uint{dri.ResourceId})
			if err != nil {
				respErr(c, 500, fmt.Sprintf("删除角色.资源失败:%v", err))
				return
			}
		}
		deleteResourceByName(c, dri.ResourceId)
		return
	}

	if err != gorm.ErrRecordNotFound {
		deleteResourceByName(c, dri.ResourceId)
		return
	}

	respOkEmpty(c)

}

func deleteResourceByName(c *gin.Context, resourceId uint) {
	err := DB.Delete(&Resource{ID: resourceId}).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("删除资源失败: ID=%s, %v", resourceId, err))
		return
	}
	respOkEmpty(c)
}
