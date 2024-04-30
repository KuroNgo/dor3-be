package user_controller

import (
	"clean-architecture/internal/cloud/google"
	"clean-architecture/internal/cloud/google/mail"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (u *UserController) GetMailTest(ctx *gin.Context) {
	err := google.Cron.AddFunc("@every 0h0m1s", func() {
		err := mail.SendMailTest()
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}
	})

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
