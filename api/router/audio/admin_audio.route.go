package audio_route

import (
	audio_controller "clean-architecture/api/controller/audio"
	"clean-architecture/bootstrap"
	audio_domain "clean-architecture/domain/audio"
	"clean-architecture/infrastructor/mongo"
	audio_repository "clean-architecture/repository/audio"
	audio_usecase "clean-architecture/usecase/audio"
	"github.com/gin-gonic/gin"
	"time"
)

func AdminAudioRouter(env *bootstrap.Database, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	au := audio_repository.NewAudioRepository(db, audio_domain.CollectionAudio)
	audio := &audio_controller.AudioController{
		AudioUseCase: audio_usecase.NewAudioUseCase(au, timeout),
		Database:     env,
	}

	router := group.Group("/audio")
	router.POST("/create", audio.CreateAudioInFireBaseAndSaveMetaDataInDatabase)
}
