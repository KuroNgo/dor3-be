package audio_controller

import (
	"clean-architecture/bootstrap"
	audio_domain "clean-architecture/domain/audio"
)

type AudioController struct {
	AudioUseCase audio_domain.IAudioUseCase
	Database     *bootstrap.Database
}
