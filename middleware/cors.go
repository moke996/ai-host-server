package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 允许CORS的中间件函数
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 对客户端支持的头部进行设置
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		// 处理实际的请求
		c.Next()
	}
}
