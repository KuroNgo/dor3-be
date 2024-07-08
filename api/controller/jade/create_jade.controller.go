package jade_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type JadeTest struct {
	Jade int64 `json:"jade" bson:"jade"`
}

func (j *JadeController) CreateJadeInUser(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not login!",
		})
		return
	}

	user, err := j.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not authorize in this action!",
		})
		return
	}

	var jadeTest JadeTest
	if err = ctx.ShouldBindJSON(&jadeTest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	err = j.JadeUseCase.CreateCurrencyInUser(ctx, user.ID, jadeTest.Jade)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
