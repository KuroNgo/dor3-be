package file

import "strings"

// IsMP3 is used for checked format file mp3
func IsMP3(filename string) bool {
	return strings.ToLower(filename[len(filename)-4:]) == ".mp3"
}

// IsMP4 is used for checked format file mp4
func IsMP4(filename string) bool {
	return strings.ToLower(filename[len(filename)-4:]) == ".mp4"
}

// IsExcel is used for checked format file excel
func IsExcel(filename string) bool {
	ext := strings.ToLower(filename[len(filename)-4:])
	return ext == ".xls" || ext == ".xlsx"
}
