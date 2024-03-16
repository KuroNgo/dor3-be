package audio_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *AudioController) FetchManyAudio(ctx *gin.Context) {
	audio, err := a.AudioUseCase.FetchMany(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   audio,
	})
}
