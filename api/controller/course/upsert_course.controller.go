package course_controller

import (
	course_domain "clean-architecture/domain/course"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// UpsertOneQuiz in this method, system can not need to check valid
func (c *CourseController) UpsertOneQuiz(ctx *gin.Context) {
	courseID := ctx.Query("_id")

	var course course_domain.Input
	if err := ctx.ShouldBindJSON(&course); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	upsertCourse := &course_domain.Course{
		Name:        course.Name,
		Description: course.Description,
		Level:       course.Level,
		UpdatedAt:   time.Now(),
	}

	courseRes, err := c.CourseUseCase.UpsertOne(ctx, courseID, upsertCourse)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   courseRes,
	})
}
