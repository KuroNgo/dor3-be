package unit_controller

import (
	unit_domain "clean-architecture/domain/unit"
	"clean-architecture/internal"
	file_internal "clean-architecture/internal/file"
	"clean-architecture/internal/file/excel"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"time"
)

func (u *UnitController) CreateOneUnit(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	user, err := u.UserUseCase.GetByID(ctx, fmt.Sprint(currentUser))

	var unitInput unit_domain.Input
	if err := ctx.ShouldBindJSON(&unitInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := internal.IsValidUnit(unitInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	unitRes := &unit_domain.Unit{
		ID:         primitive.NewObjectID(),
		LessonID:   unitInput.LessonID,
		Name:       unitInput.Name,
		Content:    unitInput.Content,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		WhoUpdates: user.FullName,
	}

	err = u.UnitUseCase.CreateOne(ctx, unitRes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (u *UnitController) CreateUnitWithFile(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(8 << 20) // 8MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("files")
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if !file_internal.IsExcel(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Not an Excel file",
		})
		return
	}

	err = ctx.SaveUploadedFile(file, file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		err := os.Remove(file.Filename)
		if err != nil {
			fmt.Printf("Failed to delete temporary file: %v\n", err)
		}
	}()

	result, err := excel.ReadFileForUnit(file.Filename)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var units []unit_domain.Unit

	for _, unit := range result {
		lessonID, err := u.UnitUseCase.FindLessonIDByLessonName(ctx, unit.LessonID)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		elUnit := unit_domain.Unit{
			ID:         primitive.NewObjectID(),
			LessonID:   lessonID,
			Name:       unit.Name,
			ImageURL:   "",
			Content:    "null",
			IsComplete: 0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			//WhoUpdates: user.FullName,
		}
		units = append(units, elUnit)
	}

	successCount := 0
	for _, unit := range units {
		err = u.UnitUseCase.CreateOneByNameLesson(ctx, &unit)
		if err != nil {
			continue
		}
		successCount++
	}

	if successCount == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create any unit",
			"message": "Any value have exist in database",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"success_count": successCount,
	})
}
