package course_controller

import (
	course_domain "clean-architecture/domain/course"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// UpdateCourse in this method, system can not need to check valid
func (c *CourseController) UpdateCourse(ctx *gin.Context) {
	courseID := ctx.Query("_id")

	var courseInput course_domain.Input
	if err := ctx.ShouldBindJSON(&courseInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": "error",
		})
		return
	}

	updateCourse := course_domain.Course{
		Name:        courseInput.Name,
		Description: courseInput.Description,
		UpdatedAt:   time.Now(),
	}

	err := c.CourseUseCase.UpdateOne(ctx, courseID, updateCourse)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
