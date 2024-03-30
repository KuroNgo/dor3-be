package course_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *CourseController) DeleteCourse(ctx *gin.Context) {
	courseID := ctx.Query("_id")
	err := c.CourseUseCase.DeleteOne(ctx, courseID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Trả về mảng dữ liệu dưới dạng JSON
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "the course is deleted!",
	})
}
