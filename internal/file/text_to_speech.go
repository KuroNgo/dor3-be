package file_internal

import (
	htgotts "github.com/hegedustibor/htgo-tts"
)

func CreateTextToSpeech(word string) error {
	speech := htgotts.Speech{Folder: "audio", Language: "en"}
	_, err := speech.CreateSpeechFile(word, word)
	if err != nil {
		return err
	}
	return nil
}
