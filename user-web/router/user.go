package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/api"
)

func InitUserRouter(Router *gin.RouterGroup) {
	ApiGroup := Router.Group("user")
	{
		ApiGroup.GET("/list", api.GetUserList)
		ApiGroup.POST("/login", api.PassWordLoginForms)
	}
}
