package httpapi

import (
	"aa/config"
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

	//注册参数验证：
	regVfunc("isName", isName)
	regVfunc("isPstr", isPstr)
	regVfunc("isDomain", isDomain)

	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	apiV := e.Group(URI_VER)

	apiV.POST("/login", login)
	apiV.POST("/passwd",changePassword)
	apiV.Use(jwtAuth())

	da := apiV.Group("/da")
	{
		//资源
		da.POST("/resource", addResource)
		da.DELETE("/resource", deleteResource)
		da.PUT("/resource",editResource)
		da.GET("/resource", listResource) //单查
		da.GET("/resources", listResource)  //批量查
		//da.GET("/resource/:resourceID", getResource)

		//角色.资源
		da.POST("/role_resources", addResourceListForRole)
		da.DELETE("/role_resources", deleteResourceListForRole)
		da.GET("/role_resources", getPermissionsForRole)

		//角色
		da.POST("/role", addRole)
		da.DELETE("/role", delRole)
		da.GET("/role", listRole)

		//用户.角色
		da.POST("/user_roles", addOrDelRoleListForUserFunc("add"))
		da.DELETE("/user_roles", addOrDelRoleListForUserFunc("del"))
		da.GET("/user_roles", getRolesForUser)

		//查询角色下有几个用户
		da.GET("/role_users", getUsersForRole)

		//用户
		da.POST("/user", addUser)
		da.DELETE("/user", delUser)
		da.GET("/user", listUser)
	}

	//超级管理员
	apiV.POST("/sa/domain", addDomain)
	apiV.DELETE("/sa/domain", delDomain)
	apiV.GET("/sa/domain", listDomain)

	return e
}

func Serve() {
	r := newEngine()
	gin.SetMode(gin.DebugMode) //gin.ReleaseMode
	err := r.Run(config.C.HTTP.ListenAddr)
	if err != nil {
		panic(err)
	}
}
