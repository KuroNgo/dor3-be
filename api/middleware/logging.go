package middleware

import (
	activity_controller "clean-architecture/api/controller/activity"
	activity_log_domain "clean-architecture/domain/activity_log"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func StructuredLogger(logger *zerolog.Logger, activity *activity_controller.ActivityControllerV2) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		//raw := ctx.

		ctx.Next()

		param := gin.LogFormatterParams{
			TimeStamp: time.Now(),
			Path:      path,
			ClientIP:  ctx.ClientIP(),
			Method:    ctx.Request.Method,
		}

		if ctx.Writer.Status() >= 500 || ctx.Errors != nil || param.Method == "DELETE" {
			param.Latency = time.Since(start).Truncate(time.Millisecond)
			param.StatusCode = ctx.Writer.Status()
			param.ErrorMessage = ctx.Errors.ByType(gin.ErrorTypePrivate).String()

			logger.Error().
				Str("client_id", param.ClientIP).
				Str("method", param.Method).
				Int("status_code", param.StatusCode).
				Int("body_size", ctx.Writer.Size()).
				Str("path", param.Path).
				Str("latency", param.Latency.String()).
				Msg(param.ErrorMessage)

			newLog := activity_log_domain.ActivityLog{
				LogID:        primitive.NewObjectID(),
				ClientIP:     param.ClientIP,
				Method:       param.Method,
				StatusCode:   param.StatusCode,
				BodySize:     ctx.Writer.Size(),
				Path:         path,
				Latency:      param.Latency.String(),
				Error:        param.ErrorMessage,
				ActivityTime: param.TimeStamp,
			}

			err := activity.ActivityUseCase.CreateOne(ctx, newLog)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to create activity log")
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status": "error",
					"error":  "Failed to create activity log",
				})
				return
			}
		}
	}
}
