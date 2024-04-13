package exam_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExamsController) FetchManyExam(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")

	exam, err := e.ExamUseCase.FetchMany(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   exam,
	})
}
