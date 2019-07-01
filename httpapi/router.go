package httpapi

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
)

func regVfunc(tag string, vfunc validator.Func) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation(tag, vfunc)
	}
}

const URI_VER = "/api/v1"

func newEngine() *gin.Engine {
	e := gin.New()

	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	apiV := e.Group(URI_VER)

	apiV.POST("/login", Login)
	apiV.Use(JWTAuth())

	//注册参数验证：
	regVfunc("isName", isName)
	regVfunc("isPstr", isPstr)

	da := apiV.Group("/da")
	{
		//资源
		//da.GET("/menus", da.GetMenus)
		//da.GET("/menus/:id", da.GetMenu)
		da.POST("/resource", AddResource)
		da.DELETE("/resource", DeleteResource)

		//角色.资源
		da.POST("/role_resources", AddResourceListOfRole)
		da.DELETE("/role_resources", DeleteResourceListOfRole)

		//角色
		//da.PUT("/menus/:id", da.EditMenu)
		//da.DELETE("/menus/:id", da.DeleteMenu)
		//
		//da.GET("/roles", da.GetRoles)
		//da.GET("/roles/:id", da.GetRole)
		da.POST("/role", AddRole)

		//用户.角色
		da.POST("/user_roles", AddRoleListOfUser)
		da.DELETE("/user_roles", DeleteRoleListOfUser)

		//用户
		//da.DELETE("/roles/:id", da.DeleteRole)
		//da.GET("/users", api.GetUsers)
		da.GET("/user", ListUser)
		da.POST("/user", AddUser)

		//da.DELETE("/users/:id", api.DeleteUser)

	}

	//超级管理员
	apiV.POST("/sa", AddDomainAdmin)
	//鉴权
	apiV.GET("/permission", CheckPermission)

	return e
}

func Serve() {
	gin.SetMode(gin.DebugMode) //gin.ReleaseMode
	r := newEngine()
	err := r.Run(":8000")
	if err != nil {
		panic(err)
	}
}
