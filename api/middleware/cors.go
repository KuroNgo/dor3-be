package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORSPublic() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization,accept,origin,Cache-Control,X-Requested-With,Access-Control-Allow-Origin,access-control-allow-methods")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func CORSPrivate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization,accept,origin,Cache-Control,X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Max-Age", "86400") // 24 giờ

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func OptionMessage(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	c.Header("Access-Control-Allow-Methods", "GET, PUT, PATCH, POST, DELETE, OPTIONS")
	c.Header("Access-Control-Max-Age", "86400") // 24 giờ

	if c.Request.Method == "OPTIONS" {
		c.JSON(http.StatusFound, 204)
	}

	c.Next()
}
