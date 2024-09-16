package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup) {
	userGroup := Router.Group("user")
	{
		userGroup.GET("/list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		userGroup.POST("/login", api.PassWordLoginForms)
		userGroup.POST("/register", api.RegisterUser)
	}
}
