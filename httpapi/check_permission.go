package httpapi

import "github.com/gin-gonic/gin"

type cpIn struct {
	Act          string `json:"act" binding:"required"`
	Domain       string `json:"domain"` //todo  空值处理
	ResourceName string `json:"resourceName" binding:"required"`
}

func CheckPermission(c *gin.Context) {
	var ci cpIn
	err := c.BindJSON(&ci)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	user := User{Name: ctxUser(c)}
	if ci.Domain == "" {
		err = DB.First(&user).Error
		if err != nil { //todo not found
			respErr(c, 500, err.Error())
			return
		}
	}

	if ci.Domain == "" {
		ci.Domain = user.DefaultDomain
	}

	ok := Enforcer.Enforce(user.Name, ci.ResourceName+"@"+ci.Domain, ci.Act)
	if ok {
		c.JSON(200, "true")
	} else {
		c.JSON(401, "false")
	}
}
