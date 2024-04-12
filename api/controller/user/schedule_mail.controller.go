package user_controller

import (
	"github.com/gin-gonic/gin"
)

type Time struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

func (u *UserController) AddScheduleReminderNotification(ctx *gin.Context) {
	
}
