package audio_route

import (
	audio_controller "clean-architecture/api/controller/audio"
	"clean-architecture/bootstrap"
	audio_domain "clean-architecture/domain/audio"
	audio_repository "clean-architecture/repository/audio"
	audio_usecase "clean-architecture/usecase/audio"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AudioRoute(env *bootstrap.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	au := audio_repository.NewAudioRepository(db, audio_domain.CollectionAudio)
	audio := &audio_controller.AudioController{
		AudioUseCase: audio_usecase.NewAudioUseCase(au, timeout),
		Database:     env,
	}

	router := group.Group("/audio")
	router.GET("/fetch", audio.FetchManyAudio)
}
