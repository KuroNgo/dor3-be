package course_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *CourseController) StatisticCourse(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	admin, err := c.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	count := c.CourseUseCase.StatisticCourse(ctx)

	ctx.JSON(http.StatusOK, gin.H{
		"total": count,
	})
}
