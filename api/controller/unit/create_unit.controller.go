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
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := u.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

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
		ID:        primitive.NewObjectID(),
		LessonID:  unitInput.LessonID,
		Name:      unitInput.Name,
		Level:     unitInput.Level,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		WhoCreate: admin.FullName,
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
	currentUser, exists := ctx.Get("currentUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "You are not logged in!",
		})
		return
	}
	admin, err := u.AdminUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
	if err != nil || admin == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Unauthorized",
			"message": "You are not authorized to perform this action!",
		})
		return
	}

	// Parse form
	err = ctx.Request.ParseMultipartForm(8 << 20) // 8MB max size
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Error parsing form",
			"message": err.Error(),
		})
		return
	}

	// Get uploaded file
	file, err := ctx.FormFile("files")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kiểm tra định dạng file
	if !file_internal.IsExcel(file.Filename) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Not an Excel file",
		})
		return
	}

	// Lưu file vào thư mục tạm
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

	// Đọc dữ liệu từ file Excel
	result, err := excel.ReadFileForUnit(file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Tạo kênh để gửi và nhận lỗi từ goroutine
	errChan := make(chan error)
	defer close(errChan)

	// Tạo các đơn vị từ dữ liệu bằng cách sử dụng goroutine
	successCount := 0
	for _, unit := range result {
		go func(unit file_internal.Unit) {
			// Tìm ID của bài học từ tên bài học
			lessonID, err := u.LessonUseCase.FindLessonIDByLessonNameInAdmin(ctx, unit.LessonID)
			if err != nil {
				errChan <- fmt.Errorf("failed to find lesson ID for lesson '%s': %v", unit.LessonID, err)
				return
			}

			// Tạo đơn vị
			elUnit := unit_domain.Unit{
				ID:        primitive.NewObjectID(),
				LessonID:  lessonID,
				Name:      unit.Name,
				ImageURL:  "",
				Level:     unit.Level,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				WhoCreate: admin.FullName,
			}

			// Tạo đơn vị trong cơ sở dữ liệu
			err = u.UnitUseCase.CreateOneByNameLesson(ctx, &elUnit)
			if err != nil {
				errChan <- fmt.Errorf("failed to create unit '%s': %v", elUnit.Name, err)
				return
			}

			// Gửi thông báo thành công
			errChan <- nil
		}(unit)
	}

	// Đợi goroutine kết thúc và xử lý lỗi nếu có
	for range result {
		if err := <-errChan; err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			successCount++
		}
	}

	// Trả về kết quả
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
