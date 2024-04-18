package course_controller

import (
	course_domain "clean-architecture/domain/course"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// UpdateCourse in this method, system can not need to check valid
func (c *CourseController) UpdateCourse(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	admin, err := c.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	var courseInput course_domain.Input
	if err := ctx.ShouldBindJSON(&courseInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": "error",
		})
		return
	}

	updateCourse := course_domain.Course{
		Id:          courseInput.Id,
		Name:        courseInput.Name,
		Description: courseInput.Description,
		UpdatedAt:   time.Now(),
		WhoUpdated:  admin.FullName,
	}

	data, err := c.CourseUseCase.UpdateOne(ctx, updateCourse)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})

}
