package course_controller

import (
	course_domain "clean-architecture/domain/course"
	"clean-architecture/internal"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (c *CourseController) CreateOneCourse(ctx *gin.Context) {
	var courseInput course_domain.Input
	if err := ctx.ShouldBindJSON(&courseInput); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := internal.IsValidCourse(courseInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	course := &course_domain.Course{
		Name:        courseInput.Name,
		Description: courseInput.Description,
		Level:       courseInput.Level,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := c.CourseUseCase.CreateOne(ctx, course)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
