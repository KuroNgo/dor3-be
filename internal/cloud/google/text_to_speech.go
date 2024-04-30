package google

import (
	htgotts "github.com/hegedustibor/htgo-tts"
	"os"
	"strings"
)

func CreateTextToSpeech(word string) error {
	speech := htgotts.Speech{Folder: "audio", Language: "en"}
	wordRemoveSpace := strings.ReplaceAll(word, " ", "")
	_, err := speech.CreateSpeechFile(word, wordRemoveSpace)
	if err != nil {
		return err
	}
	return nil
}

type FileInfo struct {
	OriginalName string
	TrimmedName  string
}

func ListFilesInDirectory(dir string) ([]string, error) {
	var filenames []string

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			filenames = append(filenames, file.Name())
		}
	}

	return filenames, nil
}

func DeleteAllFilesInDirectory(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return nil
}
