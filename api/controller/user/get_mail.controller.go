package user_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (u *UserController) GetMail(ctx *gin.Context) {
	//err := google.SendEmail("maiquangdinh.it.work@gmail.com", "")
	//
	//if err != nil {
	//	ctx.JSON(http.StatusForbidden, gin.H{
	//		"status":  "fail",
	//		"message": err.Error(),
	//	})
	//	return
	//}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
