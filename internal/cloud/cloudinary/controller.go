package cloudinary

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func FileUpload() gin.HandlerFunc {
	return func(c *gin.Context) {
		//upload
		formfile, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				MediaDto{
					StatusCode: http.StatusInternalServerError,
					Message:    "error",
					Data:       map[string]interface{}{"data": "Select a file to upload"},
				})
			return
		}

		uploadUrl, err := NewMediaUpload().FileUpload(File{File: formfile})
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				MediaDto{
					StatusCode: http.StatusInternalServerError,
					Message:    "error",
					Data:       map[string]interface{}{"data": "Error uploading file"},
				})
			return
		}

		c.JSON(
			http.StatusOK,
			MediaDto{
				StatusCode: http.StatusOK,
				Message:    "success",
				Data:       map[string]interface{}{"data": uploadUrl},
			})
	}
}

func RemoteUpload() gin.HandlerFunc {
	return func(c *gin.Context) {
		var url Url

		//validate the request body
		if err := c.BindJSON(&url); err != nil {
			c.JSON(
				http.StatusBadRequest,
				MediaDto{
					StatusCode: http.StatusBadRequest,
					Message:    "error",
					Data:       map[string]interface{}{"data": err.Error()},
				})
			return
		}

		uploadUrl, err := NewMediaUpload().RemoteUpload(url)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				MediaDto{
					StatusCode: http.StatusInternalServerError,
					Message:    "error",
					Data:       map[string]interface{}{"data": "Error uploading file"},
				})
			return
		}

		c.JSON(
			http.StatusOK,
			MediaDto{
				StatusCode: http.StatusOK,
				Message:    "success",
				Data:       map[string]interface{}{"data": uploadUrl},
			})
	}
}
