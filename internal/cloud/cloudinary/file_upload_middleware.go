package cloudinary

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func FileUploadMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("files")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer file.Close() // close file properly

		c.Set("filePath", header.Filename)
		c.Set("file", file)

		c.Next()
	}
}
