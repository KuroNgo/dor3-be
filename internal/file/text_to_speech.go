package file_internal

import (
	"fmt"
	htgotts "github.com/hegedustibor/htgo-tts"
	"mime/multipart"
	"os"
)

func CreateTextToSpeech(word string) error {
	speech := htgotts.Speech{Folder: "audio", Language: "en"}
	_, err := speech.CreateSpeechFile(word, word)
	if err != nil {
		return err
	}
	return nil
}

func ListFilesInDirectory(dir string) ([]*multipart.FileHeader, error) {
	var fileHeaders []*multipart.FileHeader

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			fileHeader := &multipart.FileHeader{
				Filename: file.Name(),
			}
			fileHeaders = append(fileHeaders, fileHeader)
		}
	}

	return fileHeaders, nil
}

func DeleteFile(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		fmt.Printf("Failed to delete temporary file: %v\n", err)
	}
}
