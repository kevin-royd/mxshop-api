package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/models"
)

// IsAdminAuth 中间件，检查用户是否为管理员
func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从上下文中获取 claims
		claims, exists := ctx.Get("claims")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"msg": "未授权",
			})
			ctx.Abort()
			return
		}

		// 类型断言，将 claims 转换为 *CustomClaims 类型
		currentUser, ok := claims.(*models.CustomClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"msg": "权限解析失败",
			})
			ctx.Abort()
			return
		}

		// 检查用户权限
		if currentUser.AuthorityId != 1 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"msg": "权限不足",
			})
			ctx.Abort()
			return
		}

		// 权限通过，继续处理请求
		ctx.Next()
	}
}
