package course_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *CourseController) FetchCourse(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	course, detail, err := c.CourseUseCase.FetchManyForEachCourse(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"detail": detail,
		"data":   course,
	})
}

func (c *CourseController) FetchCourseByID(ctx *gin.Context) {
	courseID := ctx.Query("_id")

	course, err := c.CourseUseCase.FetchByID(ctx, courseID)
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
