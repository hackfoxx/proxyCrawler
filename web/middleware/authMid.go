package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"proxyCrawler/config"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != config.GetConfig().Web.Authorization { // 简单的 token 验证
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
