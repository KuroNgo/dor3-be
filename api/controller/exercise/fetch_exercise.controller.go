package exercise_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (e *ExerciseController) FetchManyExercise(ctx *gin.Context) {
	_, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not login!",
		})
		return
	}

	page := ctx.DefaultQuery("page", "1")

	exercise, count, err := e.ExerciseUseCase.FetchMany(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   count,
		"data":   exercise,
	})
}

func (e *ExerciseController) FetchManyByUnitID(ctx *gin.Context) {
	_, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not login!",
		})
		return
	}

	unitID := ctx.Query("unit_id")
	if unitID == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
		})
		return
	}

	page := ctx.DefaultQuery("page", "1")

	exercise, detail, err := e.ExerciseUseCase.FetchManyByUnitID(ctx, unitID, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"detail": detail,
		"data":   exercise,
	})
}

func (e *ExerciseController) FetchOneByUnitID(ctx *gin.Context) {
	_, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not login!",
		})
		return
	}

	unitID := ctx.Query("unit_id")
	if unitID == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
		})
		return
	}

	ex, err := e.ExerciseUseCase.FetchOneByUnitID(ctx, unitID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   ex,
	})
}

func (e *ExerciseController) FetchOneExerciseByID(ctx *gin.Context) {
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	user, err := e.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	id := ctx.Query("_id")
	exercise, err := e.ExerciseUseCase.FetchByID(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   exercise,
	})
}

func (e *ExerciseController) FetchManyExerciseInAdmin(ctx *gin.Context) {
	_, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not login!",
		})
		return
	}

	page := ctx.DefaultQuery("page", "1")

	exercise, count, err := e.ExerciseUseCase.FetchMany(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"page":   count,
		"data":   exercise,
	})
}

func (e *ExerciseController) FetchManyByUnitIDInAdmin(ctx *gin.Context) {
	_, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not login!",
		})
		return
	}

	unitID := ctx.Query("unit_id")
	if unitID == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
		})
		return
	}

	page := ctx.DefaultQuery("page", "1")

	exercise, detail, err := e.ExerciseUseCase.FetchManyByUnitID(ctx, unitID, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"detail": detail,
		"data":   exercise,
	})
}

func (e *ExerciseController) FetchOneByUnitIDInAdmin(ctx *gin.Context) {
	_, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not login!",
		})
		return
	}

	unitID := ctx.Query("unit_id")
	if unitID == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
		})
		return
	}

	ex, err := e.ExerciseUseCase.FetchOneByUnitID(ctx, unitID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   ex,
	})
}
