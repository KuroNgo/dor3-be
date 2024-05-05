package exam_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExamsController) FetchManyExam(ctx *gin.Context) {
	_, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}

	page := ctx.DefaultQuery("page", "1")

	exam, detail, err := e.ExamUseCase.FetchMany(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"detail": detail,
		"data":   exam,
	})
}

func (e *ExamsController) FetchExamByUnitID(ctx *gin.Context) {
	_, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in",
		})
		return
	}

	unitID := ctx.Query("unit_id")
	page := ctx.DefaultQuery("page", "1")

	exam, detail, err := e.ExamUseCase.FetchManyByUnitID(ctx, unitID, page)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"detail": detail,
		"data":   exam,
	})
}
