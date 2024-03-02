package file_internal

import "strings"

func HaveVocabulary(cell string) bool {
	if strings.ToLower(cell) != "từ vựng" || strings.ToLower(cell) != "vocabulary" {
		return false
	}
	return true
}

func HaveWordType(cell string) bool {
	if strings.ToLower(cell) == "word type" || strings.ToLower(cell) == "loại từ" {
		return true
	}
	return false
}

func HaveSpelling(cell string) bool {
	if strings.ToLower(cell) == "spelling" || strings.ToLower(cell) == "phiên âm" {
		return true
	}
	return false
}

func HaveVieMean(cell string) bool {
	if strings.ToLower(cell) == "vietnamese explanation" || strings.ToLower(cell) == "giải thích tiếng việt" {
		return true
	}
	return false
}
