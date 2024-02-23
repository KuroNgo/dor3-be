package file

import (
	"errors"
	"github.com/dhowden/tag"
	"os"
)

// GetMetadataFileMP3 get metadata of file
func GetMetadataFileMP3(filepath string) (interface{}, error) {
	if !IsMP3(filepath) {
		return nil, errors.New("this file is not mp3")
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, errors.New("failed to open file")
	}
	defer file.Close()

	metadata, err := tag.ReadFrom(file)
	if err != nil {
		return nil, errors.New("failed to read metadata")
	}

	return metadata, nil
}
