package file_internal

import (
	"errors"
	"log"
	"os"
	"strconv"
)

// GetDurationFileMP3 get time duration from file and convert it to string
func GetDurationFileMP3(filepath string) (string, error) {
	if !IsMP3(filepath) {
		return "", errors.New("this file is not mp3")
	}

	file, err := os.Open("./" + filepath)
	if err != nil {
		return "", errors.New("failed to open file")
	}
	defer file.Close()

	// Đọc phần header của file MP3 để xác định thời lượng
	header := make([]byte, 10)
	_, err = file.Read(header)
	if err != nil {
		log.Fatal(err)
	}

	// Tính thời lượng từ thông tin trong header
	// Ví dụ: Header byte 7-9 là thời lượng của file MP3
	duration := calculateDuration(header[6], header[7], header[8])
	timeDuration := strconv.Itoa(duration)

	return timeDuration, nil
}

// Hàm tính thời lượng từ thông tin trong header của file MP3
func calculateDuration(b1, b2, b3 byte) int {
	return int(b1)<<16 + int(b2)<<8 + int(b3)
}
