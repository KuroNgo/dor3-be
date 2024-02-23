package file_internal

import (
	"errors"
	"github.com/dhowden/tag"
	"github.com/go-audio/wav"
	"os"
)

// GetNameFileMP3 get filename from file
func GetNameFileMP3(filepath string) (string, error) {
	if !IsMP3(filepath) {
		return "", errors.New("this file is not mp3")
	}

	file, err := os.Open(filepath)
	if err != nil {
		return "", errors.New("failed to open file")
	}
	defer file.Close()

	res, err := tag.ReadFrom(file)
	if err != nil {
		return "", errors.New("failed to read metadata")
	}

	filename := res.Title()
	return filename, nil
}

// GetDurationFileMP3 get time duration from file and convert it to string
func GetDurationFileMP3(filepath string) (string, error) {
	if !IsMP3(filepath) {
		return "", errors.New("this file is not mp3")
	}

	file, err := os.Open(filepath)
	if err != nil {
		return "", errors.New("failed to open file")
	}
	defer file.Close()

	decoder := wav.NewDecoder(file)
	if decoder == nil {
		errors.New("failed to create decoder")
	}

	duration, err := decoder.Duration()
	if err != nil {
		return "", errors.New("failed to open file")
	}

	timeDuration := duration.String()

	return timeDuration, nil
}
