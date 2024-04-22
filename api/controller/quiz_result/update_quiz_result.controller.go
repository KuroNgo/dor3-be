package quiz_result_controller

//func (q *QuizResultController) UpdateComplete(ctx *gin.Context) {
//	currentUser, exists := ctx.Get("currentUser")
//	if !exists {
//		ctx.JSON(http.StatusUnauthorized, gin.H{
//			"status":  "fail",
//			"message": "You are not logged in!",
//		})
//		return
//	}
//	user, err := q.UserUseCase.GetByID(ctx, fmt.Sprintf("%s", currentUser))
//	if err != nil || user == nil {
//		ctx.JSON(http.StatusUnauthorized, gin.H{
//			"status":  "Unauthorized",
//			"message": "You are not authorized to perform this action!",
//		})
//		return
//	}
//
//	exerciseID := ctx.Query("exercise_id")
//
//	_, err = q.QuizResultUseCase(ctx, exerciseID, 1)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{
//			"status":  "error",
//			"message": err.Error(),
//		})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, gin.H{
//		"status": "success",
//	})
//}
