package user_controller

import (
	"clean-architecture/internal/cloud/google"
	subject_const "clean-architecture/internal/cloud/google/const"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (u *UserController) GetMail(ctx *gin.Context) {
	err := google.SendEmail("2011423@dlu.edu.vn", subject_const.SignInTheFirstTime, subject_const.ContentTitle4)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
