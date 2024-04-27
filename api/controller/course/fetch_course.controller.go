package course_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *CourseController) FetchCourse(ctx *gin.Context) {
	course, err := c.CourseUseCase.FetchManyForEachCourse(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   course,
	})
}
