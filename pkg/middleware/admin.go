package middleware

import (
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			c.JSON(403, gin.H{"error": "admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
