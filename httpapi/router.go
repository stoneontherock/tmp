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

	//注册参数验证：
	regVfunc("isName", isName)
	regVfunc("isPstr", isPstr)
	regVfunc("isDomain", isDomain)

	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	apiV := e.Group(URI_VER)

	apiV.POST("/login", login)
	apiV.Use(jwtAuth())

	da := apiV.Group("/da")
	{
		//资源
		da.POST("/resource", addResource)
		da.DELETE("/resource", deleteResource)
		da.GET("/resource", listResource)

		//角色.资源
		da.POST("/role_resources", addResourceListOfRole)
		da.DELETE("/role_resources", deleteResourceListOfRole)

		//角色
		da.POST("/role", addRole)
		da.DELETE("/role", delRole)
		da.GET("/role", listRole)

		//用户.角色
		da.POST("/user_roles", addOrDelRoleListOfUserFunc("add"))
		da.DELETE("/user_roles", addOrDelRoleListOfUserFunc("del"))

		//用户
		da.POST("/user", addUser)
		da.DELETE("/user", delUser)
		da.GET("/user", listUser)

	}

	//超级管理员
	apiV.POST("/sa", addDomain)
	apiV.DELETE("/sa", delDomain)
	apiV.GET("/sa", listDomain)
	//鉴权
	apiV.GET("/permission", checkPermission)
	apiV.GET("/test", test)

	return e
}

func Serve(addr string) {
	gin.SetMode(gin.DebugMode) //gin.ReleaseMode
	r := newEngine()
	err := r.Run(addr)
	if err != nil {
		panic(err)
	}
}
