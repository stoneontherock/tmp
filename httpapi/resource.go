package httpapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strings"
)

type addResourceIn struct {
	Action  string `json:"action" binding:"required"`
	Domain  string `json:"domain" binding:"omitempty,isDomain"`
	Name    string `json:"name" binding:"required"`
	Comment string `json:"comment" binding:"required"`
}

func addResource(c *gin.Context) {
	var ari addResourceIn
	err := c.ShouldBindJSON(&ari)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}
	//fmt.Printf("action=%s\n", ari.Action) //todo DEL

	dm, err := getDomain(c, ari.Domain)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	var cnt int
	err = DB.Model(&Resource{}).Where(&Resource{Action: ari.Action, Domain: dm, Name: ari.Name}).Count(&cnt).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("创建资源失败,Count：%v", err))
		return
	}

	if cnt > 0 {
		respErr(c, 500, "资源已经存在(action+domain+name)")
		return
	}

	r := Resource{Action: ari.Action, Domain: dm, Name: ari.Name, Comment: ari.Comment}
	err = DB.Model(Resource{}).Create(&r).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("创建资源失败,Create：%v", err))
		return
	}

	respOkEmpty(c)
}

type deleteResourceIn struct {
	ResourceId uint `json:"ID" binding:"gt=0"`
}

func deleteResource(c *gin.Context) {
	var dri deleteResourceIn
	err := c.ShouldBindJSON(&dri)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	var resource Resource
	if ctxUser(c) == SA {
		DB.First(&resource, "id = ?", dri.ResourceId)
	} else {
		DB.First(&resource, "id = ? AND domain = ?", dri.ResourceId, ctxDomain(c))
	}

	if resource.Name == "" {
		respErr(c, 500, "资源不存在或不属于域内资源")
		return
	}

	var rrList []roleResource
	err = DB.Table(`role_resource`).Select(`role_name`).Where(`resource_id = ?`, dri.ResourceId).Find(&rrList).Error
	if err == nil {
		for _, rr := range rrList {
			err = deleteResourceListForRoleDo(rr.RoleName, []uint{dri.ResourceId})
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
	err := DB.Delete(&Resource{ID: resourceId}).Error
	if err != nil {
		respErr(c, 500, fmt.Sprintf("删除资源失败: ID=%d, %v", resourceId, err))
		return
	}
	respOkEmpty(c)
}

type listResourceIn struct {
	ID     uint   `form:"ID"`
	Domain string `form:"domain" binding:"omitempty,isDomain"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
	Sort   string `form:"sort"`
	Order  string `form:"order"`
}

func listResource(c *gin.Context) {
	var lri listResourceIn
	err := c.ShouldBindQuery(&lri)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	dm, err := getDomain(c, lri.Domain)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	//
	total := 0
	query := DB.Model(&Resource{}).Where(&Resource{ID: lri.ID, Domain: dm})
	err = query.Count(&total).Error
	if err != nil {
		respErr(c, 500, err.Error())
		return
	}

	if total == 0 {
		c.JSON(200, gin.H{"code": 200, "result": nil, "totalCount": total})
		return
	}

	if lri.Limit > 0 {
		query = query.Limit(lri.Limit)
	}

	if lri.Sort != "" {
		query = query.Order(lri.Sort + " " + lri.Order)
	}

	var rscList []Resource
	err = query.Offset(lri.Offset).Find(&rscList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		respErr(c, 500, err.Error())
		return
	}

	//批量查
	if strings.HasSuffix(c.Request.URL.Path, "resources") {
		c.JSON(200, gin.H{"code": 200, "result": rscList, "totalCount": len(rscList)})
		return
	}

	//单查
	var r Resource
	if len(rscList) > 0 {
		r = rscList[0]
	}

	c.JSON(200, struct {
		Code int `json:"code"`
		Resource
	}{
		200,
		r,
	})
}

type editResourceIn struct {
	Domain     string `json:"domain"`
	ResourceId uint   `json:"ID" binding:"gt=0"`
	Name       string `json:"name" binding:"required"`
	Action     string `json:"action" binding:"required"`
	Comment    string `json:"comment" binding:"required"`
}

func editResource(c *gin.Context) {
	var eri editResourceIn
	err := c.ShouldBindJSON(&eri)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	dm, err := getDomain(c, eri.Domain)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	var cnt int
	var rsc0 Resource
	err = DB.Model(&rsc0).Where(`Name = ? AND action = ? AND domain = ?`, eri.Name, eri.Action, dm).Count(&cnt).Error
	if err != nil {
		respErr(c, 500, err.Error())
		return
	}

	if cnt != 0 {
		respErr(c, 400, fmt.Sprintf("name=%s,action=%s,domain=%s的资源已存在", eri.Name, eri.Action, dm))
		return
	}

	rsc := Resource{ID: eri.ResourceId}
	if ctxUser(c) == SA {
		DB.First(&rsc)
	} else {
		DB.First(&rsc, "domain = ?", ctxDomain(c))
	}

	if rsc.Name == "" {
		respErr(c, 500, "资源不存在或不属于域内资源")
		return
	}

	var rrList []roleResource
	err = DB.Table(`role_resource`).Where(`resource_id = ?`, rsc.ID).Find(&rrList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		respErr(c, 500, fmt.Sprintf("查找角色-资源失败: %v", err))
		return
	}
	for _, rr := range rrList {
		loadResourceForRolePolicy([]Resource{rsc}, rr.RoleName, "del")
	}

	err = DB.Model(&rsc).Updates(Resource{Name: eri.Name, Action: eri.Action, Comment: eri.Comment}).Error
	if err != nil {
		respErr(c, 500, err.Error())
		return
	}

	for _, rr := range rrList {
		loadResourceForRolePolicy([]Resource{rsc}, rr.RoleName, "add")
	}

	respOkEmpty(c)
}
