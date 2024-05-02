package exam_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExamsController) FetchManyExam(ctx *gin.Context) {
	//currentUser, exists := ctx.Get("currentUser")
	//if !exists {
	//	ctx.JSON(http.StatusUnauthorized, gin.H{
	//		"status":  "fail",
	//		"message": "You are not logged in!",
	//	})
	//	return
	//}
	//user, err := e.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	//if err != nil || user == nil {
	//	ctx.JSON(http.StatusUnauthorized, gin.H{
	//		"status":  "Unauthorized",
	//		"message": "You are not authorized to perform this action!",
	//	})
	//	return
	//}
	//
	//admin, err := e.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	//if err != nil || admin == nil {
	//	ctx.JSON(http.StatusUnauthorized, gin.H{
	//		"status":  "Unauthorized",
	//		"message": "You are not authorized to perform this action!",
	//	})
	//	return
	//}

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

	exam, err := e.ExamUseCase.FetchManyByUnitID(ctx, unitID)
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
